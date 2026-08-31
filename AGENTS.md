# AGENTS.md - Dell Terraform Provider for PowerFlex

## Project Overview

This is the Terraform provider for Dell PowerFlex (VxFlex OS) software-defined storage. It implements resources and data sources using HashiCorp's Terraform Plugin Framework, enabling infrastructure-as-code management of PowerFlex systems.

- **Language:** Go 1.24
- **Module path:** `terraform-provider-powerflex`
- **Terraform Plugin Framework:** v1.13.0
- **SDK:** `github.com/dell/goscaleio` v1.19.0
- **Registry address:** `registry.terraform.io/dell/powerflex`
- **License:** Mozilla Public License 2.0

## Architecture

The provider follows the standard Terraform Plugin Framework architecture. It runs as a gRPC server that Terraform Core communicates with to manage PowerFlex resources.

### Provider Configuration

The provider authenticates to a PowerFlex Gateway using endpoint, username, and password. Configuration can be supplied via HCL provider block or environment variables (`POWERFLEX_ENDPOINT`, `POWERFLEX_USERNAME`, `POWERFLEX_PASSWORD`, `POWERFLEX_INSECURE`, `POWERFLEX_TIMEOUT`).

The provider maintains two client instances:
- **`goscaleio.Client`** — For PowerFlex MDM API operations.
- **`goscaleio.GatewayClient`** — For PowerFlex Manager (Gateway) API operations.

### SDK Strategy

Uses `goscaleio` — a public, versioned Go module on GitHub (`github.com/dell/goscaleio`). The provider and SDK release independently. For operations not covered by the SDK (e.g., template CRUD), a custom HTTP client (`client/`) makes direct REST API calls.

### Resources and Data Sources

The provider exposes approximately 34 resources and 19 data sources covering PowerFlex entities such as volumes, SDSs, SDCs, storage pools, protection domains, snapshots, snapshot policies, fault sets, devices, clusters, firmware repositories, templates, compatibility management, NVMe hosts/targets, replication, and more.

## Directory Structure

```
main.go                           Entry point (providerserver.Serve)
powerflex/
  provider/
    provider.go                   Provider configuration, resource/datasource registration
    *_resource.go                 Resource implementations (Create, Read, Update, Delete)
    *_resource_schema.go          Resource schema definitions
    *_datasource.go               Data source implementations
    *_datasource_schema.go        Data source schema definitions
    *_test.go                     Unit and acceptance tests
  helper/                         Shared helper functions (API-to-Terraform mapping)
  models/                         Terraform state model structs
  constants/                      Shared constants
  resource-test/                  Additional test utilities
client/                           Custom HTTP clients for operations not in goscaleio SDK
examples/                         Example HCL configurations
  resources/                      Resource examples
  data-sources/                   Data source examples
  provider/                       Provider configuration example
docs/                             Generated documentation
templates/                        Documentation templates
about/                            Provider metadata
```

## Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Compile the provider binary |
| `make install` | Build and install to `~/.terraform.d/plugins/` |
| `make test` | Run formatting, linting, vetting, and unit tests |
| `make testacc` | Run acceptance tests (`TF_ACC=1`, requires live hardware) |
| `make check` | Run `gofmt`, `golangci-lint`, `go vet` |
| `make gosec` | Run security scan with `gosec` |
| `make cover` | Generate HTML coverage report |
| `make generate` | Run `go generate` (docs generation) |

## Testing

### Unit Tests (mockey)

- Test files follow `*_test.go` convention alongside source files in `powerflex/provider/`.
- Frameworks: `github.com/stretchr/testify` (assertions), `github.com/bytedance/mockey` (function-level mocking).
- Unit tests are gated by `os.Getenv("TF_ACC") == "1"` skip — they run when `TF_ACC` is NOT set.
- Run with `make test`.
- No hardware required.

### Acceptance Tests (terraform-plugin-testing)

- **Requires live PowerFlex hardware** with credentials set via environment variables.
- Acceptance tests are gated by `os.Getenv("TF_ACC") != "1"` skip — they run only when `TF_ACC=1`.
- Creates real resources — clean up after failures.
- Run with `make testacc`.
- Tests use `resource.TestCase` with `ProtoV6ProviderFactories`.
- Environment-specific values are read from `powerflex.env` file or environment variables.

### Running Tests

```bash
# Unit tests (no hardware)
make test

# Acceptance tests (requires live hardware)
export POWERFLEX_ENDPOINT="https://gateway-ip"
export POWERFLEX_USERNAME="admin"
export POWERFLEX_PASSWORD="secret"
export POWERFLEX_INSECURE="true"
make testacc
```

## Code Style and Conventions

### Formatting and Linting

- All Go code must pass `gofmt`, `go vet`, and `golangci-lint`.

### Code Organization Patterns

- **Resource pattern:** Each resource has up to three files:
  - `<name>_resource.go` — CRUD lifecycle implementation.
  - `<name>_resource_schema.go` — Schema definition with attributes, validators, and defaults.
  - Helper functions in `powerflex/helper/` for mapping between `goscaleio` types and Terraform models.
- **Data source pattern:** Similar three-file pattern with `_datasource` suffix.
- **Models:** Terraform state structs are in `powerflex/models/` using `tfsdk` struct tags.
- **Custom clients:** `client/` contains HTTP clients for operations not available in the `goscaleio` SDK (e.g., `TemplateClient` for template CRUD).
- **Provider struct:** `powerflexProvider` holds `goscaleio.Client`, `goscaleio.GatewayClient`, and any custom clients. Resources receive this via `Configure()`.

### File Header

All source files must include the Dell copyright and MPL 2.0 license header:

```go
/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
...
*/
```

## Common Development Tasks

### Adding a New Resource

1. Create `powerflex/provider/<name>_resource.go` implementing `resource.Resource`.
2. Create `powerflex/provider/<name>_resource_schema.go` with the schema definition.
3. Create the model struct in `powerflex/models/`.
4. Add helper functions in `powerflex/helper/` for API-to-Terraform mapping.
5. If the `goscaleio` SDK lacks needed methods, add a custom client in `client/`.
6. Register the resource in `powerflex/provider/provider.go` `Resources()` method.
7. Add unit tests using mockey mocks and acceptance tests.
8. Create example HCL in `examples/resources/powerflex_<name>/`.
9. Run `make generate` to produce documentation.

### Adding a New Data Source

1. Create `powerflex/provider/<name>_datasource.go` implementing `datasource.DataSource`.
2. Create the schema and model files.
3. Register in `powerflex/provider/provider.go` `DataSources()` method.
4. Add tests and examples.

### Updating the SDK

```bash
go get github.com/dell/goscaleio@<version>
go mod tidy
```

## CI/CD

GitHub Actions workflows in `.github/workflows/`. GoReleaser configuration in `.goreleaser.yaml` builds cross-platform binaries (linux, darwin, windows, freebsd; amd64, arm64, 386, arm).

## Code Ownership

All files are owned by the maintainers defined in `.github/CODEOWNERS`.
