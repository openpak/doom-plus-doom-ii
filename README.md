# DOOM + DOOM II

OpenPak's DOOM + DOOM II server (scaffold).

| | |
| --- | --- |
| Game | DOOM + DOOM II (Bethesda Softworks) |
| Network | Nintendo Switch |
| Title ID | `01008CB01E52E000` |
| Licence | AGPL-3.0-only |

## What it serves

| Host | What | Served by |
| --- | --- | --- |
| `<titleId>.playfabapi.com` | Auth, lobbies, matchmaking, PlayStream (PlayFab XPlat C++ SDK) | this repo (`doom-plus-doom-iid`, stubs: 501 until configured) |
| `api.bethesda.net`, `ugcmods.bethesda.net` | Bethesda.net account + UGC mods | this repo (stubs: 501 until configured) |
| `prod-notification-websocket.bethesda.net` | Notification websocket | not yet served (see TODOs) |
| `nso-entitlements.osiris2.quakeapi.com` | NSO entitlement check | this repo (passes through) |

No NEX, Photon, EOS or Demonware in this title. Long-term, PlayFab traffic
should move to [`playfab`](https://github.com/openpak/playfab) (`playfabd`),
which already serves OpenPak's supported PlayFab title APIs.

## Run

```sh
go run ./cmd/doom-plus-doom-iid
```

| Variable | Meaning |
| --- | --- |
| `DOOM_HTTP_ADDR` | HTTP listen address (default `:23630`) |
| `DOOM_HEALTH_ADDR` | health listen address (default `:23631`) |
| `DOOM_PLAYFAB_TITLE_ID` | PlayFab TitleId (unknown; see TODOs) |
| `DOOM_BETHESDA_API` | Bethesda.net API base URL override |

## Docs

- [Protocol](docs/protocol.md)
- [Bethesda.net](docs/bethesda-net.md)
- [Ports](docs/ports.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
