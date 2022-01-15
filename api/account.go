package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/tpmdigital/simplebank/db/sqlc"
)

//// Single

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"` 
}

func (server *Server) getAccount(ctx *gin.Context) {
	// Bind request to getAccountRequest struct
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// no errors return to the client
	ctx.JSON(http.StatusOK, account)
}

//// Many

type listAccountRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"` 
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"` 
}

func (server *Server) listAccount(ctx *gin.Context) {
	// Bind request to getAccountRequest struct
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:    req.PageSize,
		Offset: (req.PageID -1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// no errors return to the client
	ctx.JSON(http.StatusOK, accounts)
}

////

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR GBP"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	// Bind request to createAccountRequest struct
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Call into the db to create the account
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// no errors return to the client
	ctx.JSON(http.StatusOK, account)
}

////

type updateAccountRequest struct {
	ID    int64 `json:"id" binding:"required"`
	Balance int64 `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {

	// Bind request to createAccountRequest struct
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Call into the db to create the account
	arg := db.UpdateAccountParams{
		ID:    req.ID,
		Balance:  req.Balance,
	}
	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// no errors return to the client
	ctx.JSON(http.StatusOK, account)
}

///

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"` 
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	// Bind request to deleteAccountRequest struct
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

