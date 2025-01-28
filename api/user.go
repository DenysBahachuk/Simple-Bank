package api

import (
	"net/http"
	"time"

	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	"github.com/DenysBahachuk/Simple_Bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// CreateUser godoc
//
//	@Summary	Creates a user
//	@Schemes
//	@Description	Creates a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createUserRequest	true	"User payload"
//	@Success		200		{object}	db.User
//	@Failure		400		{string}	error	"Bad request"
//	@Failure		403		{string}	error	"Forbidden"
//	@Failure		500		{string}	error	"Internal server error"
//	@Router			/users [post]
func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
