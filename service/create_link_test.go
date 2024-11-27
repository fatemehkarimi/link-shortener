package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateLink(t *testing.T) {
	db, mock, err := sqlmock.New()

	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO urls\\(original_url, short_code, create_date, expires_at\\)").
		WithArgs("www.google.com", "hello fatemeh", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	payload := []byte(`{"URL": "www.google.com"}`)
	req, err := http.NewRequest(http.MethodPost, "/v1/create-link", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	handler := &Handler{DB: db}
	rr := httptest.NewRecorder()
	handler.createLinkHandler(rr, req)

	assert.NoError(t, mock.ExpectationsWereMet())
}
