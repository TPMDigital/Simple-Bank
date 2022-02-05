package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/tpmdigital/simplebank/db/sqlc"
	"github.com/tpmdigital/simplebank/token"
)

//// TransferTx
type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {

	// Bind request to transferRequest struct
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check that both accounts have the same currency as the one in the transfer request
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}	

	toAccount, valid := server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}	

	// Check that we are transferring to the same account
	if fromAccount.ID == toAccount.ID {
		err := errors.New("cannot transfer from and to the same account")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check the transfer is ONLY from the logged in user 
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload) // cast to token.Payload as return is general interface
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Call into the db to transfer the amount between accounts
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// No errors return to the client
	ctx.JSON(http.StatusOK, result)
}

// Check the account with the passed in accountID has the same currency as the one supplied
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {

	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
