package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// CreateTransfer godoc
//
//	@Summary	Creates a transfer
//	@Schemes
//	@Description	Creates a money transfer between two accounts
//	@Tags			transfer
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		transferRequest	true	"Post payload"
//	@Success		200		{object}	db.TransferTxresult
//	@Failure		400		{string}	error	"Bad request"
//	@Failure		500		{string}	error	"Internal server error"
//	@Failure		404		{string}	error	"Account not found"
//	@Router			/transfers [post]
func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !s.validateAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !s.validateAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	args := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := s.store.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
