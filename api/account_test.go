package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	mockdb "lesson/simple-bank/db/mock"
	db "lesson/simple-bank/db/sqlc"
	"lesson/simple-bank/token"
	"lesson/simple-bank/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	user := utils.RandomOwner()
	account := randomAccount(user)

	// test case table
	testCases := []struct {
		name string
		accountId int64
		setupFunc func(t *testing.T, reqest *http.Request, tokenMaker token.Maker)
		buidStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account)
	}{
		{
			name: "OK",
			accountId: account.ID,
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
				authorizationSetup(t, reqest, tokenMaker, authorizationTypeBearer, user, time.Minute)
			},
			buidStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
				require.Equal(t, http.StatusOK, recorder.Code)
				responseBodyMatch(t, recorder.Body, account)	
			},
		},
		{
			name: "UnauthorizedUser",
			accountId: account.ID,
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
				authorizationSetup(t, reqest, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buidStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			accountId: account.ID,
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
			},
			buidStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFound",
			accountId: account.ID,
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
				authorizationSetup(t, reqest, tokenMaker, authorizationTypeBearer, user, time.Minute)
			},
			buidStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			accountId: account.ID,
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
				authorizationSetup(t, reqest, tokenMaker, authorizationTypeBearer, user, time.Minute)
			},
			buidStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidId",
			accountId: 0,
			setupFunc: func(t *testing.T, reqest *http.Request, tokenMaker token.Maker) {
				authorizationSetup(t, reqest, tokenMaker, authorizationTypeBearer, user, time.Minute)
			},
			buidStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// mock store obj
			store := mockdb.NewMockStore(ctrl)
			// bulid the stubs
			tc.buidStubs(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountId)
			request, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			// set up authorizations
			tc.setupFunc(t, request, server.tokenMaker)

			// make the request
			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder, account)
		})
	}
}

func randomAccount(user string) db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    user,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func responseBodyMatch(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var reqAccount db.Account
	err = json.Unmarshal(data, &reqAccount)
	require.NoError(t, err)
	require.Equal(t, reqAccount, account)
}
