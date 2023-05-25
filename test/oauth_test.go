package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateConfig(t *testing.T) {
	router := NewRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/oauth/config", nil)
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.JSONEq(t, `{"enabled":false,"passwordLoginEnabled":true}`,
		w.Body.String())
}
