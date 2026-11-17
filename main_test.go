package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleItemsCount(t *testing.T) {
	h := handleItemsCount(func() (int, error) { return 3, nil })
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/items/count", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var body struct {
		Service string `json:"service"`
		Count   int    `json:"count"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Service != "backend-1" || body.Count != 3 {
		t.Errorf("body = %+v, want service backend-1 and count 3", body)
	}
}

func TestHandleItemsCountError(t *testing.T) {
	h := handleItemsCount(func() (int, error) { return 0, errors.New("boom") })
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/items/count", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
