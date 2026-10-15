# drill-app

WARNING: This is a test repository containing a purposefully vulnerable app. DO NOT USE, ever!

The live-drill twin of [hardened-app](https://github.com/mm-weber/hardened-app)
for the [dhc-pipeline](https://github.com/mm-weber/dhc-pipeline) hardened-image
catalogue (task group 16, Req 7.11): the same tiny Go HTTP service, built the
same way (`CGO_ENABLED=0`, git tag + commit checksum, `dhi.io/build`, nonroot
UID 65532), with two dependencies pinned at versions that carry one published
vulnerability each, on purpose, so that every triage mechanism of the
catalogue can be exercised live and recorded.

- `GET /` → `200`, version banner (`-ldflags "-X main.version=..."`)
- `GET /healthz` → `200 ok` (liveness/readiness probe target)
- `GET /whoami` → `200 hello <sub>` for a bearer token signed with
  `DRILL_JWT_KEY`, `401` otherwise
- Every response passes the CORS middleware (allowed origin
  `https://drill.example`)
- Listens on `:8080` (override with `PORT`)

The pins are the drill, not a mistake: `go.mod` names the versions and the
catalogue's `triage/LOG.md` records what each finding was used to exercise.
Nothing here is meant to run anywhere but the drill.

```bash
go test ./...
go run .
```

Releases are git tags (`v<semver>`); the catalogue pins one and tracks new
tags via Renovate. Release commits and tags are signed and verify on GitHub,
as hardened-app's are.
