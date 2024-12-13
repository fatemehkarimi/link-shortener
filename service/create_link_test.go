package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct {
	formattedTime string
}

func (a AnyTime) Match(v driver.Value) bool {
	time, ok := v.(time.Time)
	return ok && time.Format("2006-01-02") == a.formattedTime
}

func TestCreateLink(t *testing.T) {
	now = func() time.Time { return time.Date(2024, time.December, 6, 11, 7, 12, 12, time.UTC) }
	db, mock, err := sqlmock.New()

	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO urls\\(original_url, short_code, create_date, expires_at\\)").
		WithArgs("www.google.com", "480f2ef8fc12", AnyTime{"2024-12-06"}, AnyTime{"2024-12-13"}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	payload := []byte(`{"URL": "www.google.com"}`)
	req, err := http.NewRequest(http.MethodPost, "/v1/create-link", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	handler := &Handler{DB: db}
	rr := httptest.NewRecorder()
	handler.createLinkHandler(rr, req)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateDifferentLinkHashForSameLink(t *testing.T) {
	now = func() time.Time { return time.Date(2024, time.December, 6, 11, 8, 12, 12, time.UTC) }
	db, mock, err := sqlmock.New()

	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO urls\\(original_url, short_code, create_date, expires_at\\)").
		WithArgs("www.google.com", "e80fec96bcb4", AnyTime{"2024-12-06"}, AnyTime{"2024-12-13"}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	payload := []byte(`{"URL": "www.google.com"}`)
	req, err := http.NewRequest(http.MethodPost, "/v1/create-link", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	handler := &Handler{DB: db}
	rr := httptest.NewRecorder()
	handler.createLinkHandler(rr, req)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLinkByHashFound(t *testing.T) {
	db, mock, err := sqlmock.New()

	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT original_url FROM urls WHERE short_code=\$1`).
		WithArgs("e577cd4e510c").
		WillReturnRows(sqlmock.NewRows([]string{"original_url"}).AddRow("https://www.google.com"))

	payload := []byte(`{"hash": "e577cd4e510c"}`)
	req, err := http.NewRequest(http.MethodGet, "/v1/get-link-by-hash", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	handler := &Handler{DB: db}
	rr := httptest.NewRecorder()

	handler.getURLByHash(rr, req)
	assert.Equal(t, rr.Result().StatusCode, http.StatusFound)
	assert.Equal(t, rr.Header().Get("Location"), "https://www.google.com")
}

func TestGetLinkByHashNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	mock.ExpectQuery(`SELECT original_url FROM urls WHERE short_code=\$1`).
		WithArgs("e577cd4e510c").
		WillReturnError(sql.ErrNoRows)

	payload := []byte(`{"hash": "e577cd4e510c"}`)
	req, err := http.NewRequest(http.MethodGet, "/v1/get-link-by-hash", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	handler := &Handler{DB: db}
	rr := httptest.NewRecorder()

	handler.getURLByHash(rr, req)
	assert.Equal(t, rr.Result().StatusCode, http.StatusNotFound)
}
