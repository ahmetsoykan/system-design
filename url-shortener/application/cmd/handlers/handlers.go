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

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
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

	ctx := r.Context()
	tracer := otel.Tracer("shorten")
	ctx, span := tracer.Start(ctx, "/shorten")

	defer span.End()

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
	var extendTTL bool
	if requestBody.TTL == "" {
		expiry = time.Now().UTC().Add(time.Hour * 24 * 90).Unix()
		extendTTL = true
	} else {
		i, err := strconv.ParseInt(requestBody.TTL, 10, 64)
		if err != nil {
			http.Error(w, errors.New("error: ttl value is not valid").Error(), http.StatusBadRequest)
			return
		}
		expiry = i
		extendTTL = false
	}

	ctx, dbSpan := tracer.Start(ctx, "/shorten - dynamodb put item")
	data, err := s.DB.DDBPutItem(db.Item{
		LongURL:   requestBody.URL,
		TTL:       expiry,
		ExtendTTL: extendTTL,
	})
	if err != nil {
		http.Error(w, errors.New("error: "+err.Error()).Error(), http.StatusBadRequest)
		return
	}
	dbSpan.End()

	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(data)
	w.Write(jsonData)
}

// redirect endpoint read short path corresponding value from cache
// if it returns empty result.
// reads short path key from database, puts it to cache and extends its TTL 90 days more
// Redirect!
// if it returns the longsurl value
// Redirect!
func (s *Server) redirect(w http.ResponseWriter, r *http.Request) {

	var short string
	var longURL string
	ctx := r.Context()

	tracer := otel.Tracer("redirect")

	// checking cache
	params := mux.Vars(r)
	short = strconv.FormatUint(hash.Decode(params["short"]), 10)

	_, span := tracer.Start(ctx, params["short"])
	defer span.End()

	_, redisSpan := tracer.Start(ctx, "redisCheck")
	val, _ := s.Cache.Get(short)
	if val != "" {
		// found in cache
		longURL = val
		redisSpan.End()
	} else {
		// checking database
		_, databaseSpan := tracer.Start(ctx, "databaseCheck")
		
		item, err := s.DB.GetItembyPK(short)
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
		err = s.Cache.Set(short, longURL)
		if err != nil {
			http.Error(w, errors.New("error: cache cant be updated, reason: "+err.Error()).Error(), http.StatusBadGateway)
			return
		}

		// extending TTL in database 90 days more
		err = s.DB.ExtendTTL(short)
		if err != nil {
			http.Error(w, errors.New("error: ttl can't be extended, reason: "+err.Error()).Error(), http.StatusBadGateway)
			return
		}

		databaseSpan.End()
	}

	defer redisSpan.End()

	http.Redirect(w, r, string(longURL), http.StatusMovedPermanently)
}
