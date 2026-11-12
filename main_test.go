package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCountHandler(t *testing.T) {
	h := countHandler(func() (int, error) { return 3, nil })
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/items/count", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body struct {
		Service string `json:"service"`
		Count   int    `json:"count"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Service != "backend-1" || body.Count != 3 {
		t.Fatalf("body = %+v, want backend-1/3", body)
	}
}

func TestCountHandlerError(t *testing.T) {
	h := countHandler(func() (int, error) { return 0, errors.New("boom") })
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/items/count", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
}
