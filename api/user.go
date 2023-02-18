package api

import (
	"database/sql"
	"fmt"
	db "lesson/simple-bank/db/sqlc"
	"lesson/simple-bank/utils"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"user_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=5"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username         string    `json:"user_name" binding:"required,alphanum"`
	FullName         string    `json:"full_name" binding:"required"`
	Email            string    `json:"email" binding:"required,email"`
	CreatedAt        time.Time `json:"create_at"`
	PasswordChangeAt time.Time `json:"password_change_at"`
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
	rsp := createUserResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		CreatedAt:        user.CreatedAt,
		PasswordChangeAt: user.PasswordChangeAt,
	}
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

	rsp := createUserResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		CreatedAt:        user.CreatedAt,
		PasswordChangeAt: user.PasswordChangeAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
