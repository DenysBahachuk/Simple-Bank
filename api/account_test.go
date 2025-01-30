package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/DenysBahachuk/Simple_Bank/db/mock"
	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	"github.com/DenysBahachuk/Simple_Bank/token"
	"github.com/DenysBahachuk/Simple_Bank/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	user, _ := randomUser()
	account := randomAccount(user.Username)

	testCase := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		buildStubs    func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: -1,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, request, maker, authTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStore := mockdb.NewMockStore(ctrl)

			tc.buildStubs(mockStore)

			//start test server and send request
			server, err := newTestServer(mockStore)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  utils.RandomAmount(),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account

	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}
