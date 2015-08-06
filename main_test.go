package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockDatabase struct {
	d map[string]string
}

func (r MockDatabase) Store(k string, v string) {
	r.d[k] = v
}

func (r MockDatabase) Read(k string) (v string) {
	return r.d[k]
}

func TestShortenStringReturnsOk(t *testing.T) {
	h := shortenHandler(MockDatabase{d: make(map[string]string)})
	req, _ := http.NewRequest("POST", "/www.google.se", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Shorten did not return %v", http.StatusOK)
	}
}

func TestRedirectStringReturnsRedirect(t *testing.T) {
	db := MockDatabase{d: make(map[string]string)}
	db.Store("test", "qwert")
	h := redirectHandler(db)
	req, _ := http.NewRequest("GET", "/qwert", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Errorf("Redirect did not return status %v", http.StatusFound)
	}
}
