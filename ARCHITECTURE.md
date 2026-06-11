# Architecture: terraform-provider-powerflex

## Metadata

<!-- yaml-metadata-start -->
scope_paths: ["./"]
capture_git_sha: "8bedd67e8fefe18fb46adf50d1cfb4adf1bea1bd"
status: "current"
auto_update: false
preview_before_apply: true
scaffold_version: "1.0"
<!-- yaml-metadata-end -->

---

## Purpose and Structure

Terraform provider for Dell PowerFlex (VxFlex OS) software-defined
storage. Implements 28 managed resources and 22 data sources using
HashiCorp's Terraform Plugin Framework, enabling
infrastructure-as-code management of PowerFlex systems via their
REST API.

The provider is a standalone Go binary that communicates with
Terraform Core over gRPC (go-plugin protocol). It uses the public
`goscaleio` SDK for core operations and a custom HTTP client
(`client/`) for operations not yet covered by the SDK (e.g.,
template CRUD).

---

## Components

| Component | Path | Responsibility |
|-----------|------|---------------|
| Entry point | `main.go` | `providerserver.Serve` — starts gRPC server |
| Provider | `powerflex/provider/provider.go` | Schema, Configure, resource/datasource registration |
| Resources | `powerflex/provider/*_resource.go` | CRUD lifecycle for 28 managed resources |
| Resource schemas | `powerflex/provider/*_resource_schema.go` | Schema definitions with validators and defaults |
| Data sources | `powerflex/provider/*_datasource.go` | Read-only queries for 22 data sources |
| Data source schemas | `powerflex/provider/*_datasource_schema.go` | Data source schema definitions |
| Helpers | `powerflex/helper/` | API-to-Terraform type mapping functions |
| Models | `powerflex/models/` | Terraform state model structs (`tfsdk` tags) |
| Constants | `powerflex/constants/` | Shared constants |
| Custom clients | `client/` | HTTP clients for SDK-uncovered operations |
| Examples | `examples/` | HCL configurations for resources and data sources |
| Docs | `docs/` | Generated provider documentation |

---

## Key Behaviors

### Authentication

**GIVEN** a user configures the provider with endpoint, username,
and password (via HCL block or environment variables)
**WHEN** `Configure()` runs
**THEN** (1) env vars `POWERFLEX_ENDPOINT`, `POWERFLEX_USERNAME`,
`POWERFLEX_PASSWORD`, `POWERFLEX_INSECURE`, `POWERFLEX_TIMEOUT`
override HCL values, (2) two SDK clients are initialized:
`goscaleio.Client` for MDM API and `goscaleio.GatewayClient` for
Gateway API, (3) authentication is validated before any resource
operations proceed

### Resource CRUD Lifecycle

**GIVEN** a resource definition in HCL
**WHEN** `terraform apply` runs
**THEN** the resource's `Create()` reads the plan into a model
struct, calls the SDK to create the resource, maps the API response
back to Terraform state, and sets `resp.State`

### Drift Detection

**GIVEN** a resource exists in Terraform state
**WHEN** `terraform plan` or `terraform refresh` runs
**THEN** `Read()` calls the SDK to fetch current state from the
system, compares it with stored state, and updates the state if
drifted

### Import

**GIVEN** a resource exists on the system but not in Terraform state
**WHEN** `terraform import powerflex_<resource>.<name> <id>` runs
**THEN** `ImportState()` fetches the resource by ID and populates
state

---

## Interfaces

### Provider Configuration Schema

| Attribute | Type | Env Var | Default | Description |
|-----------|------|---------|---------|-------------|
| `endpoint` | string | `POWERFLEX_ENDPOINT` | — | Gateway server URL (inclusive of port) |
| `username` | string | `POWERFLEX_USERNAME` | — | API username |
| `password` | string (sensitive) | `POWERFLEX_PASSWORD` | — | API password |
| `insecure` | bool | `POWERFLEX_INSECURE` | `false` | Skip TLS verification |
| `timeout` | int64 | `POWERFLEX_TIMEOUT` | — | HTTPS timeout (seconds) |

### SDK Client Layer

The `powerflexProvider` struct holds two SDK clients:
- `goscaleio.Client` — for PowerFlex MDM API operations
- `goscaleio.GatewayClient` — for PowerFlex Manager (Gateway) API
Custom HTTP clients in `client/` handle operations not in the SDK
(e.g., `TemplateClient` for template CRUD).

---

## Dependencies

| Depends On | For |
|------------|-----|
| `github.com/dell/goscaleio` v1.19.0 | PowerFlex REST API SDK |
| `client/` (local) | Custom HTTP client for SDK-uncovered ops |
| `hashicorp/terraform-plugin-framework` v1.13.0 | Core provider interfaces |
| `hashicorp/terraform-plugin-framework-validators` v0.15.0 | Attribute validation |
| `hashicorp/terraform-plugin-log` | Structured logging |
| `hashicorp/terraform-plugin-testing` | Acceptance test harness |
| `bytedance/mockey` v1.2.13 | Unit test function-level mocking |
| `stretchr/testify` v1.10.0 | Test assertions |

---

## Known Constraints

1. **Terraform Plugin Framework only** — no SDK v2 code.
2. **CGO_ENABLED=0** — static binaries for all platforms.
3. **Sensitive attributes marked** — credentials never in plan output.
4. **ImportState required** — all resources support `terraform import`.
5. **Environment variable fallback** — all credentials support env vars.
6. **Acceptance tests gated** — never run without `TF_ACC=1`.
7. **Dual client strategy** — `goscaleio.Client` for MDM API,
   `goscaleio.GatewayClient` for Gateway API, plus custom `client/`
   for operations not in the SDK.
8. **Three-file resource pattern** — each resource may have
   `*_resource.go`, `*_resource_schema.go`, and helpers.

---

## Change History

| Date | Feature | What Changed | Author |
|------|---------|-------------|--------|
| 2026-06-11 | SME-validated architecture | Provider-specific architecture with verified counts and SME input | architecture-agent |
