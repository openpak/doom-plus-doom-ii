# Ports

| Port | Use |
| --- | --- |
| 23630 | plain HTTP API (behind whatever terminates TLS for the clients' https hosts) |
| 23631 | `/healthz` |
| 23632–23639 | reserved (Bethesda.net notification websocket listener, Party/QoS, future split) |
| 443/tcp | fixed by the client (`<titleId>.playfabapi.com`, `api.bethesda.net`, `ugcmods.bethesda.net`, `nso-entitlements.osiris2.quakeapi.com` all https; notifications are wss); published on the host and forwarded to 23630 |
