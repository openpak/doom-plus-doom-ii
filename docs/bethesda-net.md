# Bethesda.net — auth/token + UGC protocol plan

Recon source: `scratch/decompile-doom-plus-doom-ii/RECON.md` (§Network,
§Gaps) and `scratch/decompile-doom-plus-doom-ii/findings.json`
(`middleware.bethesda_net`).

## What recon says

- Client code: `BNetModManagerGeneric/NX.cpp` (UGC mods) + notification
  websockets; identity key `bethesda_uid`.
- Hosts: `https://api.bethesda.net`, `https://cert-api.bethesda.net`,
  `ugcmods.bethesda.net` (+ cert twin), `wss://prod-notification-websocket.bethesda.net`,
  `*.bnet.run` dev/int/staging/test equivalents.
- Transport is the shared Nintendo `libcurl`/WebSocket stack (see
  [protocol](protocol.md)); no bundled certs, no hardcoded IPs.

## Recovered values (2026-10-27, static)

- Paths: `/login`, `/login/token-login`, `/v1/connect`,
  `/ugcmods/v2/curated`, `/ugcmods/v2/content`,
  `/ugcmods/v2/content/subscribe`, `/ugcmods/v2/stats/`,
  `/notification/v1/system`. NSO entitlements:
  `https://nso-entitlements.osiris2.quakeapi.com/token`.
- Static `X-BNET-Key` header blob (128 chars, rodata VA `0x1e2a962`,
  sole consumer the header setter at `0x1839c10`, sent verbatim
  beside `X-BNET-Agent`):

    CzdAaFVGDGpcSnUwQ0M3Lk5ZZWBzTEsqKCRCSHczZEwpUXchPmlNQERKK2g3OW5eVVlRSVswV3kpOCoxMDJddlJqLiVDc1w0dXhsOQozdzJAIU5rJEtVWTlnOEYiWAxJ

- No OAuth client pair by design (key + session headers instead).
  The runtime `X-BNET-Key` variant (singleton `0x24f6350+0x2e0`,
  session-derived) and the auth/token handshake still need a live
  capture.

## Auth/token plan (uncaptured)

1. Capture the token grant against `api.bethesda.net`: which credential the
   Switch client presents (Nintendo BaaS token? PlayFab entity? device key?)
   and the token shape it gets back.
2. Capture how `bethesda_uid` links to the PlayFab entity from
   `LoginWithNintendoServiceAccount` — one login driving the other, or two
   independent sessions.
3. Decide the server side: reimplement the grant against nx-baas-issued
   identities, or stub a fixed test identity until captures exist.

## UGC mod protocol plan (uncaptured)

1. Capture mod list/search/download calls against `ugcmods.bethesda.net`
   (paths, query params, auth headers, response schema, paging).
2. Capture upload/publish calls if the Switch client offers them.
3. Capture notification websocket frames (`prod-notification-websocket`):
   connect/auth, subscribe topics, message schema, keepalive.
4. Server side: serve a static mod catalogue first (list + download of
   repo-local WADs/mods), uploads and live notifications later.

## Current server state

`doom-plus-doom-iid` answers `/mods`, `/mods/*`, `/ugc/*`, `/notifications`,
`/notifications/*` and `/auth/*` with 501 until `DOOM_BETHESDA_API` is set,
and 501-with-config until the handlers above land. The wss endpoint is not
served yet; reserve a port in 23632–23639 when it is.
