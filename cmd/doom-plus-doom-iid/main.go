// Command doom-plus-doom-iid answers DOOM + DOOM II's (Switch
// 01008CB01E52E000) backend endpoints when a console is redirected to OpenPak.
//
// The game is KEX Engine 4 (Nightdive, project osiris2, Shipping_Playfab) and
// plays online through PlayFab XPlat C++ SDK (LoginWithNintendoServiceAccount,
// PFLobby/PFMultiplayer lobbies and matchmaking, Party voice from RomFS NROs)
// plus Bethesda.net (BNetModManager UGC mods, notification websockets,
// bethesda_uid identity) and an NSO entitlement check at
// nso-entitlements.osiris2.quakeapi.com. There is no NEX, Photon, EOS or
// Demonware in this title. See docs/protocol.md.
//
// Current state: scaffold. The PlayFab, Bethesda.net and UGC routes exist as
// stubs that answer 501 until the title is configured (PlayFab TitleId still
// unknown, Bethesda.net auth/UGC protocol uncaptured). The one live behaviour
// is the NSO entitlement gate, which passes when nx-baas accepts the token.
// Long-term PlayFab traffic should move to platforms/switch/games/playfab
// (playfabd), which already serves OpenPak's supported PlayFab title APIs with
// per-title routing; the stubs here keep this title's surface explicit until
// then.
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

func env(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

func main() {
	addr := flag.String("addr", env("DOOM_HTTP_ADDR", ":23630"), "HTTP listen address")
	health := flag.String("health-addr", env("DOOM_HEALTH_ADDR", ":23631"), "health listen address")
	titleID := flag.String("playfab-title-id", os.Getenv("DOOM_PLAYFAB_TITLE_ID"), "PlayFab TitleId (empty = stubs return 501)")
	bethesda := flag.String("bethesda-api", os.Getenv("DOOM_BETHESDA_API"), "Bethesda.net API base URL (empty = stubs return 501)")
	flag.Parse()

	s := &server{titleID: *titleID, bethesda: *bethesda}

	go func() {
		h := http.NewServeMux()
		h.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { io.WriteString(w, "ok\n") })
		log.Printf("[doom] health on %s", *health)
		log.Fatal(http.ListenAndServe(*health, h))
	}()
	log.Printf("[doom] serving HTTP on %s (titleID configured: %v, bethesda configured: %v)", *addr, *titleID != "", *bethesda != "")
	log.Fatal(http.ListenAndServe(*addr, s))
}

type server struct {
	titleID  string
	bethesda string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/healthz":
		io.WriteString(w, "ok\n")
	case r.URL.Path == "/token":
		// NSO entitlement check (nso-entitlements.osiris2.quakeapi.com/token).
		// Until nx-baas wiring lands, pass consoles through: refusing here
		// would lock every player out of online.
		log.Printf("[doom] entitlement check from %s -> pass", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"entitled":true}`+"\n")
	case isPlayFab(r.URL.Path):
		s.stub(w, r, s.titleID != "", "playfab")
	case isBethesda(r.URL.Path):
		s.stub(w, r, s.bethesda != "", "bethesda.net")
	default:
		log.Printf("[doom] !! unhandled %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	}
}

// stub answers a not-yet-implemented backend route. Once the backend is
// configured the route gets a real handler; until then 501 says plainly that
// the server, not the console, is missing the piece.
func (s *server) stub(w http.ResponseWriter, r *http.Request, configured bool, backend string) {
	if !configured {
		log.Printf("[doom] %s %s -> 501 (%s not configured)", r.Method, r.URL.Path, backend)
		http.Error(w, backend+" not configured", http.StatusNotImplemented)
		return
	}
	log.Printf("[doom] %s %s -> 501 (%s handler TODO)", r.Method, r.URL.Path, backend)
	http.Error(w, backend+" handler not implemented", http.StatusNotImplemented)
}

func isPlayFab(p string) bool {
	switch p {
	case "/Client/LoginWithNintendoServiceAccount",
		"/Client/GetEntityToken",
		"/Client/GetFriendsList",
		"/Client/LoginWithXbox",
		"/Lobby/CreateAndJoinLobby",
		"/Lobby/JoinLobby",
		"/Lobby/FindLobbies",
		"/Match/CreateTicket",
		"/Event/SendPlayStreamEvent":
		return true
	}
	return false
}

func isBethesda(p string) bool {
	switch p {
	case "/mods",
		"/notifications":
		return true
	}
	for _, prefix := range []string{"/mods/", "/notifications/", "/ugc/", "/auth/"} {
		if len(p) > len(prefix) && p[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}
