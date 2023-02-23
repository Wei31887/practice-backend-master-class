package api

import (
	"database/sql"
	"errors"
	db "lesson/simple-bank/db/sqlc"
	"lesson/simple-bank/token"
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
	fromAccount, valid := server.vaildAccount(ctx, req.FromAccountID, req.Currency) 
	if !valid {
		return
	}

	_, valid = server.vaildAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	userPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != userPayload.Username {
		err := errors.New("from account does not belong to you")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
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

func (server *Server) vaildAccount(ctx *gin.Context, accountId int64, currency string ) (db.Account, bool) { 
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
        return account, false
	}

	if account.Currency != currency { 
		ctx.JSON(http.StatusBadRequest, errResponse(err))
        return account, false
	}

	return account, true
}