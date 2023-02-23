package api

import (
	"errors"
	"fmt"
	"lesson/simple-bank/token"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
)

var (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	// to get the payload at next router
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header not found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		field := strings.Fields(authorizationHeader)
		if len(field) < 2 {
			err := errors.New("invalid authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		authorizationType := strings.ToLower(field[0])
		if authorizationType != authorizationTypeBearer{
			err := fmt.Errorf("invalid authorization type %s", authorizationType)
            c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
            return
		}

		token := field[1]
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}
