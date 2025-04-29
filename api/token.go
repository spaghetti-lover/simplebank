package api

import (
<<<<<<< HEAD
	"database/sql"
=======
	"errors"
>>>>>>> d4d0e58 (refactor)
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
<<<<<<< HEAD
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct {
	AccessToken           string    `json:"access_token"`
	AccesssTokenExpiresAt time.Time `json:"access_token_expires_at"`
=======
	db "github.com/spaghetti-lover/simplebank/db/sqlc"
	"github.com/spaghetti-lover/simplebank/token"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
>>>>>>> d4d0e58 (refactor)
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

<<<<<<< HEAD
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
=======
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken, token.TokenTypeRefreshToken)
>>>>>>> d4d0e58 (refactor)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
<<<<<<< HEAD
		if err == sql.ErrNoRows {
=======
		if errors.Is(err, db.ErrRecordNotFound) {
>>>>>>> d4d0e58 (refactor)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
<<<<<<< HEAD
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
=======
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
>>>>>>> d4d0e58 (refactor)
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
<<<<<<< HEAD
=======
		return
>>>>>>> d4d0e58 (refactor)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
<<<<<<< HEAD
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	rsp := renewAccessTokenResponse{
		AccessToken:           accessToken,
		AccesssTokenExpiresAt: accessPayload.ExpiredAt,
=======
		refreshPayload.Role,
		server.config.AccessTokenDuration,
		token.TokenTypeAccessToken,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
>>>>>>> d4d0e58 (refactor)
	}
	ctx.JSON(http.StatusOK, rsp)
}
