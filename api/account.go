package api

import (
	"database/sql"
	"fmt"
	db "lesson/simple-bank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=TWD USD EUR"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err!= nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
        Currency: req.Currency,
	}
	account, err := server.store.CreateAccount(ctx, arg)
    if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
				case "unique_violation", "foreign_key_violation":
					ctx.JSON(http.StatusForbidden, errResponse(err))
					return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err!= nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
        return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err!= nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
        return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageId int32 `json:"page_id" binding:"required,min=1"`
	PageSize int32 	`json:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) ListAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindJSON(&req); err!= nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
        return
	}

	arg := db.ListAccountParams {
		Limit: req.PageSize,
		Offset: (req.PageId-1) * req.PageSize, 
	}
	accounts, err := server.store.ListAccount(ctx, arg)
	if err!= nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
        return
	}
	ctx.JSON(http.StatusOK, accounts)
}