package main

import (
	"bytes"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSignInOK(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", signinUrl(),
		func(req *http.Request) (*http.Response, error) {

			r := SignInResponse{
				Token: "jwt_token",
			}

			resp, err := httpmock.NewJsonResponse(200, r)
			return resp, err
		},
	)
	var s SignInRequest
	s.Email = "a@a.com"
	s.Password = "1234"
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(s)
	body, _ := json.Marshal(s)
	b := bytes.NewBuffer(body)

	r, _ := http.NewRequest("POST", "/signin", b)
	w := httptest.NewRecorder()

	SignIn(w, r)

	var responseObject SignInResponse
	json.Unmarshal(w.Body.Bytes(), &responseObject)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "jwt_token", responseObject.Token)
}

func TestGetVirtualCardsUnauthorized(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", virtualCardsUrl(),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(401, nil)
			return resp, err
		},
	)

	r, _ := http.NewRequest("GET", "/virtualcards", nil)
	w := httptest.NewRecorder()

	GetVirtualCards(w, r)

	var responseObject []VirtualCard
	json.Unmarshal(w.Body.Bytes(), &responseObject)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetVirtualCardsOK(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCardId1 := "testcard1"
	testCardId2 := "testcard2"
	displayName1 := "name 1"
	displayName2 := "name 2"

	httpmock.RegisterResponder("GET", virtualCardsUrl(),
		func(req *http.Request) (*http.Response, error) {
			vc := VirtualCardsResponse{
				VirtualCards: []VirtualCard{
					{
						Id:          testCardId1,
						DisplayName: displayName1,
					},
					{
						Id:          testCardId2,
						DisplayName: displayName2,
					},
				},
			}
			resp, err := httpmock.NewJsonResponse(200, vc)
			return resp, err
		},
	)

	r, _ := http.NewRequest("GET", "/virtualcards", nil)
	w := httptest.NewRecorder()

	GetVirtualCards(w, r)

	var responseObject []VirtualCard
	json.Unmarshal(w.Body.Bytes(), &responseObject)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, responseObject[0].Id, testCardId1)
	assert.Equal(t, responseObject[1].Id, testCardId2)
	assert.Equal(t, responseObject[0].DisplayName, displayName1)
	assert.Equal(t, responseObject[1].DisplayName, displayName2)
}

func TestGetVirtualCardTransactionsUnauthorized(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", virtualCardsTransactionsUrl("vc1"),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(401, nil)
			return resp, err
		},
	)

	r, _ := http.NewRequest("GET", "/virtualcards/vc1/transactions", nil)
	r = mux.SetURLVars(r, map[string]string{
		"id": "vc1",
	})
	w := httptest.NewRecorder()

	GetVirtualCardTransactions(w, r)

	var responseObject []VirtualCard
	json.Unmarshal(w.Body.Bytes(), &responseObject)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetVirtualCardTransactionsOK(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	txnId1 := "txn1"
	txnId2 := "txn2"
	merch1 := "merch 1"
	merch2 := "merch 2"
	httpmock.RegisterResponder("GET", virtualCardsTransactionsUrl("vc1"),
		func(req *http.Request) (*http.Response, error) {
			t := TransactionsResponse{
				Transactions: []Transaction{
					{
						Id:           txnId1,
						MerchantName: merch1,
					},
					{
						Id:           txnId2,
						MerchantName: merch2,
					},
				},
			}
			resp, err := httpmock.NewJsonResponse(200, t)
			return resp, err
		},
	)

	r, _ := http.NewRequest("GET", "/virtualcards/vc1/transactions", nil)
	r = mux.SetURLVars(r, map[string]string{
		"id": "vc1",
	})

	w := httptest.NewRecorder()

	GetVirtualCardTransactions(w, r)

	var responseObject []VirtualCard
	json.Unmarshal(w.Body.Bytes(), &responseObject)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, responseObject[0].Id, txnId1)
	assert.Equal(t, responseObject[1].Id, txnId2)
	assert.Equal(t, responseObject[0].Id, txnId1)
	assert.Equal(t, responseObject[1].Id, txnId2)
}

func TestGetTransactionDetailsUnauthorized(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", transactionsUrl("txn1"),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(401, nil)
			return resp, err
		},
	)

	r, _ := http.NewRequest("GET", "/virtualcards/transactions/txn1", nil)
	r = mux.SetURLVars(r, map[string]string{
		"id": "txn1",
	})
	w := httptest.NewRecorder()

	GetTransactionDetails(w, r)

	var responseObject Transaction
	json.Unmarshal(w.Body.Bytes(), &responseObject)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetTransactionDetailsOK(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	txnId := "txn1"
	merch := "merch1"
	httpmock.RegisterResponder("GET", transactionsUrl("txn1"),
		func(req *http.Request) (*http.Response, error) {
			t := Transaction{
				Id:           txnId,
				MerchantName: merch,
			}
			resp, err := httpmock.NewJsonResponse(200, t)
			return resp, err
		},
	)

	r, _ := http.NewRequest("GET", "/virtualcards/txn1/transactions", nil)
	r = mux.SetURLVars(r, map[string]string{
		"id": "txn1",
	})

	w := httptest.NewRecorder()

	GetTransactionDetails(w, r)

	var responseObject Transaction
	json.Unmarshal(w.Body.Bytes(), &responseObject)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, responseObject.Id, txnId)
	assert.Equal(t, responseObject.MerchantName, merch)
}
