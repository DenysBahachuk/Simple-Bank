package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// RenewAccessToken godoc
//
//	@Summary	Renew of the access token
//	@Schemes
//	@Description	Renew of the access token
//	@Tags			tokens
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		renewAccessTokenRequest	true	"Renew token payload"
//	@Success		200		{object}	renewAccessTokenResponse
//	@Failure		400		{string}	error	"Bad request"
//	@Failure		401		{string}	error	"Unauthorized"
//	@Failure		401		{string}	error	"blocked session"
//	@Failure		401		{string}	error	"incorrect session user"
//	@Failure		401		{string}	error	"mismatched session token"
//	@Failure		401		{string}	error	"expired session"
//	@Failure		500		{string}	error	"Internal server error"
//	@Router			/tokens/renew_access [post]
func (s *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := s.tokenMaker.ValidateToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := s.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(refreshPayload.Username, s.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
