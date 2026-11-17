package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// stubDriver opens connections that need no database, so db.Ping succeeds.
type stubDriver struct{}

func (stubDriver) Open(string) (driver.Conn, error) { return stubConn{}, nil }

type stubConn struct{}

func (stubConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("not supported") }
func (stubConn) Close() error                        { return nil }
func (stubConn) Begin() (driver.Tx, error)           { return nil, errors.New("not supported") }

func TestHandleHealth(t *testing.T) {
	sql.Register("stub", stubDriver{})
	stub, err := sql.Open("stub", "")
	if err != nil {
		t.Fatalf("open stub db: %v", err)
	}
	defer stub.Close()
	db = stub

	rec := httptest.NewRecorder()
	handleHealth(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" || body["demo"] != "shared-preview" {
		t.Errorf("body = %v, want status ok and demo shared-preview", body)
	}
}
