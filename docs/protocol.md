# DOOM + DOOM II — wire protocol (recon facts)

Wire facts only, from the workspace recon at
`scratch/decompile-doom-plus-doom-ii/RECON.md` (RECON) and
`scratch/decompile-doom-plus-doom-ii/findings.json` (findings). Keep this file
in step with `cmd/doom-plus-doom-iid`.

## Title

- Title ID `01008CB01E52E000`, NACP `DOOM + DOOM II`, Bethesda Softworks LLC,
  display version 1.0.0. SDK 17.5.3.0 (RECON §Dump).
- Engine: KEX Engine 4 (Nightdive), project `osiris2`, build config
  `Shipping_Playfab` (`osiris2_NX64.nss`/`.nrr`). Network/service code lives in
  `kexEngine-PlayFabSDK/KexAPI/...`; UI is CoHTML/Gameface (RECON §Engine).

## PlayFab (XPlat C++ SDK + Multiplayer/Lobby + Party)

- SDK `XPlatCppSdk-3.60.211025` and `-3.91.220908` (`XPlatCppSdk-Private-Switch`,
  vendored jsoncpp); `X-PlayFabSDK` header; endpoints
  `<titleId>.playfabapi.com` + `/pubsub` (SignalR) (RECON §Network).
- Auth: `LoginWithNintendoServiceAccount` (Nintendo BaaS access-token via
  `nn::account::baas` + `acc:u0`) in `socialSignInPlayfab.cpp` /
  `PlayFabAuthenticationDriver(NX)`; `LoginWithXbox` also linked;
  `GetEntityToken`/`PFMultiplayerSetEntityToken`; friends via
  `/Client/GetFriendsList`, `FriendPlayFabId` (RECON §Network).
- Multiplayer: full `PFMultiplayer*` lobby API (`CreateAndJoin/Join/JoinArranged`,
  `FindLobbies`, matchmaking tickets + server backfill, server-claim lobbies,
  invites, search properties) via KEX `playfabLobbyDriver`/`playfabServiceDriver`;
  PlayStream events + statistics (RECON §Network).
- Party/voice: RomFS NROs loaded at runtime —
  `rom:/nro/libPlayFabMultiplayer.nro` + `Party.nro`; transport is Xbox Secure
  Sockets + DTLS + Azure relay with NAT-traversal tokens; voice is Opus with
  text chat and Azure speech transcription (`.tts.speech.microsoft.com`)
  (RECON §Network).
- PlayFab TitleId `A0C71` (recovered 2026-10-27: rodata literal VA
  `0x1ea4e48`, assigned to `PlayFabSettings::titleId` at `0x19a8c8c`;
  overridable at runtime only via `-playfab_serviceid` argv, which
  Switch retail never passes). Emulated shape:
  `https://A0C71.playfabapi.com/Client/*`.

## Bethesda.net

- `BNetModManagerGeneric/NX.cpp`: UGC mods + notifications, identity
  `bethesda_uid`, full dev/int/staging/test/cert/prod env ladder:
  `https://api.bethesda.net`, `https://cert-api.bethesda.net`,
  `ugcmods.bethesda.net`, `wss://prod-notification-websocket.bethesda.net`,
  `*.bnet.run` equivalents (RECON §Network).
- Auth/token flow and UGC mod protocol are uncaptured; see
  [bethesda-net](bethesda-net.md) for the capture plan (RECON §Gaps).

## NSO entitlements

- `https://nso-entitlements.osiris2.quakeapi.com/token` (RECON §Network).
  Served live by `doom-plus-doom-iid` as a pass-through until nx-baas wiring
  lands.

## Transport

- Nintendo `libcurl/7.64.1` (`nn_libcurl` + nghttp2) + `NintendoSDK_WebSocket`
  + `ssl` service; no bundled certs; no hardcoded game-server IPs (loopback
  only) (RECON §Network).
- NPDM network services: `bsd:u`, `bsdcfg`, `nsd:u`, `nifm:u`, `sfdnsres`,
  `ssl`, `acc:u0`, `friend:u`, `bcat:u`, `pctl`, `ldn:u`, `hwopus`,
  `mnpp:app`, `caps:*`, `prepo:u`, `mii:u` (RECON §Engine). `ldn:u` looks
  manifest-only: no `nn::ldn`/`nn::pia` usage strings (RECON §Gaps).

## What is not here

- NEX (`nn::nex`/PRUDP/`srv.nintendo.net`: 0 hits), Photon, EOS, Demonware,
  Steamworks (0 hits) (RECON §Network).
