package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
)

type CreateRequest struct {
	URL string `json:"url"`
}

type CreateResponse struct {
	URL string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type API struct {
	urlShortener *UrlShortener
}

func NewAPI(urlShortener *UrlShortener) *API {
	return &API{
		urlShortener: urlShortener,
	}
}

func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	var body CreateRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, ErrorResponse{Error: "invalid request body"})
		return
	}

	_, err = url.ParseRequestURI(body.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, ErrorResponse{"invalid request url"})
		return
	}

	location, err := a.urlShortener.Create(r.Context(), body.URL)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, ErrorResponse{Error: "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeResponse(w, CreateResponse{URL: location})
}

func (a *API) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	redirectUrl, err := a.urlShortener.Get(r.Context(), slug)
	if err != nil {
		slog.Error(err.Error())
		writeResponse(w, ErrorResponse{Error: "internal server error"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if redirectUrl == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", redirectUrl)
	w.WriteHeader(http.StatusFound)
}

func writeResponse(w http.ResponseWriter, body any) {
	b, err := json.Marshal(&body)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	_, err = w.Write(b)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func NewServer(addr string, api *API) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/create", api.Create)
	mux.HandleFunc("/{slug}", api.Redirect)

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
