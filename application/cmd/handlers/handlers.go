// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"url-shortener/internal/db"
	"url-shortener/internal/hash"

	"github.com/go-chi/chi"
)

// Health endpoint returns 200 status code
func (s *Server) health(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	resp := map[string]string{"Status": "OK"}
	byteData, _ := json.Marshal(resp)
	if _, err := w.Write(byteData); err != nil {
		log.Println("error on response writing", err)
	}
}

// Shorten endpoint creates a database record with short url
func (s *Server) shorten(w http.ResponseWriter, r *http.Request) {

	var requestBody URLShortenInput
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, errors.New("error: "+err.Error()).Error(), http.StatusBadRequest)
		return
	}

	if requestBody.URL == "" {
		http.Error(w, errors.New("error: url value is empty").Error(), http.StatusBadRequest)
		return
	}

	var expiry int64
	var defaultTTL bool
	if requestBody.TTL == "" {
		expiry = time.Now().UTC().Add(time.Hour * 24 * 90).Unix()
		defaultTTL = true
	} else {
		i, err := strconv.ParseInt(requestBody.TTL, 10, 64)
		if err != nil {
			http.Error(w, errors.New("error: ttl value is not valid").Error(), http.StatusBadRequest)
			return
		}
		expiry = i
		defaultTTL = false
	}

	data, err := s.DB.DDBPutItem(db.Item{
		LongURL:    requestBody.URL,
		TTL:        expiry,
		DefaultTTL: defaultTTL,
	})
	if err != nil {
		http.Error(w, errors.New("error: "+err.Error()).Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(data)
	w.Write(jsonData)
}

// redirect endpoint read shorturl corresponding value from cache
// if it returns empty result.
// reads shorturl key from database, puts it to cache and extends its TTL 90 days more
// Redirect!
// if it returns the longsurl value
// Redirect!
func (s *Server) redirect(w http.ResponseWriter, r *http.Request) {

	var shortURL string
	var longURL string

	// checking cache
	shortURL = strconv.FormatUint(hash.Decode(chi.URLParam(r, "shortenurl")), 10)
	val, _ := s.Cache.Get(shortURL)
	if val != "" {
		// found in cache
		longURL = val
	} else {
		// checking database
		item, err := s.DB.GetItembyPK(shortURL)
		if err != nil {
			http.Error(w, errors.New("error: url not found").Error(), http.StatusNotFound)
			return
		}
		if item.LongURL == "" {
			http.Error(w, errors.New("error: url not found").Error(), http.StatusNotFound)
			return
		}
		// found in database
		// putting it to cache
		longURL = item.LongURL
		err = s.Cache.Set(shortURL, longURL)
		if err != nil {
			http.Error(w, errors.New("error: cache cant be updated, reason: "+err.Error()).Error(), http.StatusBadGateway)
			return
		}

		// extending TTL in database 90 days more
		// TODO: make it async
		err = s.DB.ExtendTTL(shortURL)
		if err != nil {
			http.Error(w, errors.New("error: ttl can't be extended, reason: "+err.Error()).Error(), http.StatusBadGateway)
			return
		}
	}

	http.Redirect(w, r, string(longURL), http.StatusMovedPermanently)
}
