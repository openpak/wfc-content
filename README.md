# wfc-content

Independent Wii/DS SAKE storage and Nintendo race/leaderboard HTTP service,
extracted from nn-wfc. AGPL-3.0, derived from WiiLink WFC.
Uses the existing WFC database and shared native login-ticket decoder. It creates
no replacement identity store and requires no data export/import.

Private HTTP listener WFC_HTTP_ADDR=:20122 serves existing SAKE/race paths,
guarded by WFC_INTERNAL_KEY. GET /healthz is public. NAS forwards the existing
console hostnames here; this service has no published host ports.

Required: WFC_DB_USER, WFC_DB_PASSWORD, WFC_DB_ADDRESS, WFC_API_SECRET,
WFC_TOKEN_KEY_FILE (same >=32-byte binary key as NAS/runtime),
WFC_INTERNAL_KEY (same >=32-byte internal credential). WFC_DB_NAME defaults to wwfc.

Build/test independently: `GOWORK=off go test -mod=vendor ./...`,
`GOWORK=off go vet -mod=vendor ./...`, `docker build -t wfc-content .`.
Shared code is vendored from sibling wfc-common. Image ghcr.io/openpak/wfc-content.
Read nn-wfc docs/service-split.md before deploying.
