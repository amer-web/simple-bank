package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/amer-web/simple-bank/db/mock"
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/helper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccount(t *testing.T) {
	acc := createRandomAccount()
	testCases := []struct {
		name      string
		builsStup func(store *mockdb.MockStore)
		checkRes  func(t *testing.T, res *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			builsStup: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), acc.ID).
					Times(1).
					Return(acc, nil)
			},
			checkRes: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, res.Code)
				matchAccount(t, res, acc)
			},
		},
		{
			name: "notFound",
			builsStup: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), acc.ID).
					Times(1).
					Return(db.Account{}, db.ErrorRecordNotFound)
			},
			checkRes: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, res.Code)
			},
		},
		{
			name: "serverInternalError",
			builsStup: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), acc.ID).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkRes: func(t *testing.T, res *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, res.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	server := NewServer(store)
	server.Run()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.builsStup(store)
			recorder := httptest.NewRecorder()
			req := createRequest(t, http.MethodGet, fmt.Sprintf("/account/%d", acc.ID), nil)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(t, recorder)
		})

	}

}
func createRandomAccount() db.Account {
	return db.Account{
		ID:       helper.RandomInt(1, 10),
		Currency: helper.RandomCurrency(),
		Owner:    "mohamed amer",
		Balance:  helper.RandomInt(100, 600),
	}
}

func createRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)
	return req
}

func matchAccount(t *testing.T, recorder *httptest.ResponseRecorder, acc db.Account) {
	var gotAccount db.Account
	err := json.NewDecoder(recorder.Body).Decode(&gotAccount)
	require.NoError(t, err)
	require.Equal(t, acc, gotAccount)
}
