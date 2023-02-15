package api

import (
	"database/sql"
	db "lesson/simple-bank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
	Currency	string `json:"currency" binding:"required"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	err := ctx.ShouldBindJSON(&req)
	if err!= nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	if !server.vaildAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
	if !server.vaildAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}
	result, err := server.store.TranserTx(ctx, arg)
    if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
        return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) vaildAccount(ctx *gin.Context, accountId int64, currency string ) bool { 
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
        return false
	}

	if account.Currency != currency { 
		ctx.JSON(http.StatusBadRequest, errResponse(err))
        return false
	}

	return true
}