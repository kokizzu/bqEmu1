
# bqEmu1

example how to use bigquery emulator for local project.

[bigquery-emulator](//github.com/goccy/bigquery-emulator) is masterpiece by goccy, that one dev that also creates [one of the fastest](//github.com/kokizzu/kokizzu-benchmark/blob/master/ser-deser/README.md) json serialization-library 
[go-json](//github.com/goccy/go-json).

it's quite useful because you might hate sqlite interface, and you prefer bigquery's interface.

So if your project became super large, you can move it to bigquery so it can be mass-run by a lot of people (as long as you can pay the bigquery bill).

## Requirement

- download `bigquery-emulator-linux-amd64` from [bigquery-emulator](//github.com/goccy/bigquery-emulator/releases)
- `go run main.go`

## Maintenance Checklist

- [x] Update the Go runtime directive to Go 1.26.5.
- [x] Refresh BigQuery, protobuf, gRPC, and Google API dependencies.
- [x] Add Makefile targets for vulnerability checks and arbitrary commands.
- [x] Run `make test`.
- [x] Run `make verify-dependency-security`.
- [x] Run `make vulncheck`; no reachable vulnerabilities were found.
