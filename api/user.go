package api

import (
	"database/sql"
	"fmt"
	db "lesson/simple-bank/db/sqlc"
	"lesson/simple-bank/utils"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"user_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=5"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username         string    `json:"user_name" binding:"required,alphanum"`
	FullName         string    `json:"full_name" binding:"required"`
	Email            string    `json:"email" binding:"required,email"`
	CreatedAt        time.Time `json:"create_at"`
	PasswordChangeAt time.Time `json:"password_change_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		CreatedAt:        user.CreatedAt,
		PasswordChangeAt: user.PasswordChangeAt,
	}
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPd, err := utils.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPd,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type getUserRequest struct {
	Username string `json:"user_name" binding:"required"`
}

func (server *Server) GetUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required,min=1"`
}

type loginUserResponse struct {
	SessionId        uuid.UUID    `json:"session_id"`
	AccessToken      string       `json:"access_token"`
	AccessExpiredAt  time.Time    `json:"access_expired_at"`
	RefreshToken     string       `json:"refresh_token"`
	RefreshExpiredAt time.Time    `json:"refresh_expired_at"`
	User             userResponse `json:"user"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = utils.ComparePassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	// create the access token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	// create the refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	// store refresh token into session
	arg := db.CreateSessionParams{
		ID:           refreshPayload.Id,
		Username:     req.Username,
		RefreshToken: refreshToken,
		ClientIp:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	}

	session, err := server.store.CreateSession(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	// send the token to the client
	rsp := loginUserResponse{
		SessionId: session.ID,
		AccessToken:      accessToken,
		AccessExpiredAt:  accessPayload.ExpiresAt,
		RefreshToken:     refreshToken,
		RefreshExpiredAt: refreshPayload.ExpiresAt,
		User:             newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
