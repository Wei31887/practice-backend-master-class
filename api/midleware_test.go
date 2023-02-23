package api

import (
	"fmt"
	"lesson/simple-bank/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// function to setup the authorization request for testing purposes
func authorizationSetup(
	t *testing.T,
    request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	user string,
    duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(user, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	requestHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, requestHeader)
}


func TestAuthorization(t *testing.T) {
	testCase := []struct {
		name      string
		setupFunc func(t *testing.T, reqest *http.Request, tokenMaker token.Maker)
		authFunc  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
				authorizationSetup(t, reqest, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			authFunc:  func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},
		{
			name: "NoAuthorization",
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
			},
			authFunc:  func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "InvalidAuthorization",
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
				authorizationSetup(t, reqest, tokenMaker, "", "user", time.Minute)
			},
			authFunc:  func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "UnsportedType",
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
				authorizationSetup(t, reqest, tokenMaker, "Unsurported", "user", time.Minute)
			},
			authFunc:  func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "ExpiredToken",
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
				authorizationSetup(t, reqest, tokenMaker, authorizationTypeBearer, "user", -time.Minute)
			},
			authFunc:  func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// build the server and the router
			server := newTestServer(t, nil)
			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker), // middleware fot test
				func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{})
				})
			
				// set the request and record
				request, err := http.NewRequest(http.MethodGet, authPath, nil)
				recorder := httptest.NewRecorder()
				require.NoError(t, err)
				
				// set up the authorization before send the request
				tc.setupFunc(t, request, server.tokenMaker)
				// send the request
				server.router.ServeHTTP(recorder, request)
				// check the response
				tc.authFunc(t, recorder)
			
			})
	}
}
