package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authTypeBearer         = "bearer"
	authPayloadKey         = "authPayload"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		header := ctx.GetHeader(authorizationHeaderKey)
		if len(header) == 0 {
			err := fmt.Errorf("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(header)

		if len(fields) < 2 {
			err := fmt.Errorf("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])

		if authType != authTypeBearer {
			err := fmt.Errorf("unsupported authorization type: %s", authType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]

		payload, err := s.tokenMaker.ValidateToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authPayloadKey, payload)

		ctx.Next()

	})
}
