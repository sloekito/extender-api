package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/sloekito/extend-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

const jsonContentType string = "application/json"
const extendApiVersion string = "application/vnd.paywithextend.v2021-03-12+json"
const extendUrlBase string = "api.paywithextend.com"

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SignInResponse struct {
	Token string `json:"token"`
}

type VirtualCardsResponse struct {
	VirtualCards []VirtualCard `json:"virtualCards"`
}

type VirtualCard struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Id           string `json:"id"`
	MerchantName string `json:"merchantName"`
}

// SignIn godoc
// @Summary Sign in
// @Description Sign in with username and password
// @Accept  json
// @Produce  json
// @Param order body SignInRequest true "SignIn"
// @Success 200 {object} SignInResponse
// @Router /signin [post]
func SignIn(w http.ResponseWriter, r *http.Request) {
	var signin SignInRequest
	err := json.NewDecoder(r.Body).Decode(&signin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(signin)

	req, err := http.NewRequest(http.MethodPost, signinUrl(), payloadBuf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", jsonContentType)
	req.Header.Set("Accept", extendApiVersion)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(response.StatusCode)
		return
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject SignInResponse
	json.Unmarshal(responseData, &responseObject)

	jsonResponse, jsonError := json.Marshal(responseObject)
	if jsonError != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// VirtualCards godoc
// @Summary Get Virtual Cards for the user
// @Description Get Virtual Cards for the user
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} VirtualCardsResponse
// @Router /virtualcards [get]
func GetVirtualCards(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	req, err := http.NewRequest(http.MethodGet, virtualCardsUrl(), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", jsonContentType)
	req.Header.Set("Accept", extendApiVersion)
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(response.StatusCode)
		return
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var responseObject VirtualCardsResponse
	json.Unmarshal(responseData, &responseObject)

	jsonResponse, jsonError := json.Marshal(responseObject.VirtualCards)
	if jsonError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// VirtualCardsTransactions godoc
// @Summary Get Transactions for a virtual card
// @Description Get Transactions for a virtual card
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param        id   path      string  true  "Virtual Card ID"
// @Success 200 {object} TransactionsResponse
// @Router /virtualcards/{id}/transactions [get]
func GetVirtualCardTransactions(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	token := r.Header.Get("Authorization")
	req, err := http.NewRequest(http.MethodGet, virtualCardsTransactionsUrl(id), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", jsonContentType)
	req.Header.Set("Accept", extendApiVersion)
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(response.StatusCode)
		return
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var responseObject TransactionsResponse
	json.Unmarshal(responseData, &responseObject)

	jsonResponse, jsonError := json.Marshal(responseObject.Transactions)
	if jsonError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Get Transaction Details godoc
// @Summary  Get transaction details
// @Description Get Transaction Details for a transaction ID
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param        id   path      string  true  "Transaction ID"
// @Success 200 {object} Transaction
// @Router /transactions/{id} [get]
func GetTransactionDetails(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	token := r.Header.Get("Authorization")
	req, err := http.NewRequest(http.MethodGet, transactionsUrl(id), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", jsonContentType)
	req.Header.Set("Accept", extendApiVersion)
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(response.StatusCode)
		return
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var responseObject Transaction
	json.Unmarshal(responseData, &responseObject)
	jsonResponse, jsonError := json.Marshal(responseObject)
	if jsonError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func virtualCardsUrl() string {
	u := url.URL{
		Scheme: "https",
		Host:   extendUrlBase,
		Path:   "virtualcards",
	}
	return u.String()
}

func virtualCardsTransactionsUrl(virtualcardsId string) string {
	p := path.Join("virtualcards", virtualcardsId, "transactions")
	u := url.URL{
		Scheme: "https",
		Host:   extendUrlBase,
		Path:   p,
	}
	return u.String()
}

func transactionsUrl(transactionId string) string {
	p := path.Join("transactions", transactionId)
	u := url.URL{
		Scheme: "https",
		Host:   extendUrlBase,
		Path:   p,
	}
	return u.String()
}

func signinUrl() string {
	u := url.URL{
		Scheme: "https",
		Host:   extendUrlBase,
		Path:   "signin",
	}
	return u.String()
}

// @title Extend API
// @version 1.0
// @description This is a service that calls Extend API t
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email sheila.loekito@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signin", SignIn).Methods("POST")
	r.HandleFunc("/virtualcards", GetVirtualCards).Methods("GET")
	r.HandleFunc("/virtualcards/{id}/transactions", GetVirtualCardTransactions).Methods("GET")
	r.HandleFunc("/transactions/{id}", GetTransactionDetails).Methods("GET")

	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
