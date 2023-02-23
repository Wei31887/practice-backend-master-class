package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type refreshTokenResponse struct {
	AccessToken      string       `json:"access_token"`
	AccessExpiredAt  time.Time    `json:"access_expired_at"`
}

func (server *Server) RefreshToken(ctx *gin.Context) {
	var req refreshTokenRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
        ctx.JSON(http.StatusUnauthorized, errResponse(err))
        return
    }

	session, err := server.store.GetSession(ctx, refreshPayload.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	    return
	}

	if session.Username != refreshPayload.Username {
		err := errors.New("Session username not match")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
	    return
	}

	if time.Now().After(session.ExpiresAt) {
		err = errors.New("Refresh token has expired")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
	    return	
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(session.Username, server.config.AccessTokenDuration)
	if err != nil {
        ctx.JSON(http.StatusUnauthorized, errResponse(err))
        return
    }

	rsp := refreshTokenResponse{
		AccessToken:      accessToken,
        AccessExpiredAt:  accessPayload.ExpiresAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}
