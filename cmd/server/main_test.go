package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func request(t *testing.T, h http.Handler, method, path, body string, status int) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(method, "http://short.example"+path, strings.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != status {
		t.Fatalf("%s %s: status %d, want %d: %s", method, path, w.Code, status, w.Body.String())
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing CORS")
	}
	return w
}
func decodeObject(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	if w.Header().Get("Content-Type") != "application/json" {
		t.Fatal("incorrect JSON content type")
	}
	var value map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &value); err != nil {
		t.Fatal(err)
	}
	return value
}
func TestURLLifecycle(t *testing.T) {
	h := newHandler("../../web/static")
	empty := request(t, h, "GET", "/api/list", "", 200)
	if strings.TrimSpace(empty.Body.String()) != "[]" {
		t.Fatalf("empty list: %s", empty.Body.String())
	}
	created := decodeObject(t, request(t, h, "POST", "/api/shorten", `{"url":"  https://example.org/path?q=value  "}`, 201))
	code, ok := created["short_code"].(string)
	if !ok || len(code) != 6 {
		t.Fatalf("invalid code: %v", created)
	}
	if len(created) != 3 || created["original"] != "https://example.org/path?q=value" || created["short_url"] != "https://short.example/"+code {
		t.Fatalf("unexpected response: %v", created)
	}
	before := decodeObject(t, request(t, h, "GET", "/api/stats/"+code, "", 200))
	if before["clicks"] != float64(0) || before["short_code"] != code || before["original_url"] != created["original"] || len(before) != 4 {
		t.Fatalf("stats: %v", before)
	}
	if _, err := time.Parse(time.RFC3339Nano, before["created_at"].(string)); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		redirected := request(t, h, "GET", "/"+code, "", 302)
		if redirected.Header().Get("Location") != created["original"] {
			t.Fatal("wrong redirect")
		}
	}
	after := decodeObject(t, request(t, h, "GET", "/api/stats/"+code, "", 200))
	if after["clicks"] != float64(2) || after["created_at"] != before["created_at"] {
		t.Fatalf("stats: %v", after)
	}
	var list []map[string]any
	if err := json.Unmarshal(request(t, h, "GET", "/api/list", "", 200).Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0]["short_code"] != code || list[0]["clicks"] != float64(2) {
		t.Fatalf("list: %v", list)
	}
	deleted := decodeObject(t, request(t, h, "DELETE", "/api/delete/"+code, "", 200))
	if deleted["message"] != "удалено" {
		t.Fatalf("delete: %v", deleted)
	}
	request(t, h, "GET", "/api/stats/"+code, "", 404)
	request(t, h, "DELETE", "/api/delete/"+code, "", 404)
	if got := request(t, h, "GET", "/api/list", "", 200).Body.String(); strings.TrimSpace(got) != "[]" {
		t.Fatal(got)
	}
	request(t, h, "GET", "/"+code, "", 200)
}
func TestHTTPCompatibility(t *testing.T) {
	h := newHandler("../../web/static")
	cases := []struct {
		method, path, body string
		status             int
		message            string
	}{
		{"POST", "/api/shorten", "{", 400, "неверный формат JSON"},
		{"POST", "/api/shorten", `{"url":42}`, 400, "неверный формат JSON"},
		{"POST", "/api/shorten", `{"url":"  "}`, 400, "поле url обязательно"},
		{"POST", "/api/shorten", `{"url":"ftp://example.com"}`, 400, "URL должен начинаться с http:// или https://"},
		{"GET", "/api/shorten", "", 405, "метод не поддерживается"},
		{"GET", "/api/delete/abcdef", "", 405, "метод не поддерживается"},
		{"POST", "/abcdef", "", 405, "метод не поддерживается"},
		{"POST", "/", "", 405, "метод не поддерживается"},
		{"HEAD", "/", "", 405, "метод не поддерживается"},
		{"GET", "/api/stats/", "", 400, "укажите короткий код"},
		{"GET", "/api/stats/missing", "", 404, "ссылка не найдена"},
		{"DELETE", "/api/delete/", "", 404, "ссылка не найдена"},
		{"POST", "/api/list", "", 200, ""},
		{"POST", "/api/stats/missing", "", 404, "ссылка не найдена"},
	}
	for _, tc := range cases {
		t.Run(tc.method+tc.path+tc.body, func(t *testing.T) {
			w := request(t, h, tc.method, tc.path, tc.body, tc.status)
			if tc.message != "" {
				value := decodeObject(t, w)
				if len(value) != 1 || value["error"] != tc.message {
					t.Fatalf("error: %v", value)
				}
			}
		})
	}
	for _, path := range []string{"/", "/missing", "/favicon.svg", "/static/css/style.css", "/static/js/app.js"} {
		w := request(t, h, "GET", path, "", 200)
		if w.Body.Len() == 0 {
			t.Fatalf("empty asset: %s", path)
		}
	}
	// OPTIONS short-circuits routing, including otherwise unknown routes.
	w := request(t, h, "OPTIONS", "/api/shorten", "", 200)
	if w.Body.Len() != 0 || w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, DELETE, OPTIONS" || w.Header().Get("Access-Control-Allow-Headers") != "Content-Type" {
		t.Fatal("preflight changed")
	}
	// Legacy validation checks only the scheme prefix.
	request(t, h, "POST", "/api/shorten", `{"url":"http://"}`, 201)
}
