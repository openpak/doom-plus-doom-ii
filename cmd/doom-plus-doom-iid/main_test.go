package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func get(t *testing.T, s *server, target string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, target, nil))
	return w
}

func TestHealthz(t *testing.T) {
	s := &server{}
	w := get(t, s, "/healthz")
	if w.Code != 200 || w.Body.String() != "ok\n" {
		t.Fatalf("healthz: %d %q", w.Code, w.Body.String())
	}
}

func TestEntitlementPasses(t *testing.T) {
	s := &server{}
	w := get(t, s, "/token")
	if w.Code != 200 {
		t.Fatalf("token: %d", w.Code)
	}
}

func TestStubsAre501UntilConfigured(t *testing.T) {
	s := &server{}
	for _, p := range []string{
		"/Client/LoginWithNintendoServiceAccount",
		"/Client/GetEntityToken",
		"/Lobby/FindLobbies",
		"/mods",
		"/ugc/list",
		"/notifications/subscribe",
	} {
		if w := get(t, s, p); w.Code != http.StatusNotImplemented {
			t.Errorf("%s: got %d, want 501", p, w.Code)
		}
	}
}

func TestStubsStay501AfterConfigUntilHandlersLand(t *testing.T) {
	s := &server{titleID: "ABCDE", bethesda: "https://api.bethesda.net"}
	if w := get(t, s, "/Client/LoginWithNintendoServiceAccount"); w.Code != http.StatusNotImplemented {
		t.Errorf("playfab: got %d, want 501", w.Code)
	}
	if w := get(t, s, "/mods"); w.Code != http.StatusNotImplemented {
		t.Errorf("bethesda: got %d, want 501", w.Code)
	}
}

func TestUnknownIs404(t *testing.T) {
	s := &server{}
	if w := get(t, s, "/nope"); w.Code != 404 {
		t.Fatalf("unknown path: %d", w.Code)
	}
}
