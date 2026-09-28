# KNOWLEDGE.md — terraform-provider-powerflex

<!-- yaml-metadata-start -->
scope_paths: ["./"]
capture_git_sha: "8bedd67e8fefe18fb46adf50d1cfb4adf1bea1bd"
status: "current"
auto_update: false
preview_before_apply: true
scaffold_version: "1.0"
# session_state: { is_complete: true }
<!-- yaml-metadata-end -->

<!-- quick-reference-start -->
## Agent Quick Reference

| Section | Heading | Summary | never_again_count |
|---------|---------|---------|-------------------|
| Component Overview | `## Component Overview` | Dell PowerFlex software-defined storage provider | — |
| Architectural Rationale | `## Architectural Rationale` | goscaleio SDK; Plugin Framework architecture | — |
| Failure Modes & Gotchas | `## Failure Modes & Gotchas` | SDK coupling, state corruption, auth edge cases | 3 |
| Implicit Contracts | `## Implicit Contracts` | Env var precedence, auth validation, TLS defaults | — |
<!-- quick-reference-end -->

## Five Questions Quick Reference

### What does it do?
Terraform provider for Dell PowerFlex (VxFlex OS) software-defined
storage. Exposes 28 resources and 22 data sources covering volumes,
SDCs, SDSs, protection domains, storage pools, fault sets, devices,
snapshots, snapshot policies, clusters, firmware repositories,
templates, compatibility management, NVMe hosts/targets,
replication, and more through HashiCorp's Terraform Plugin
Framework. Communicates with the system REST API via
`github.com/dell/goscaleio` v1.19.0.

### How do you modify it?
Create `*_resource.go` and `*_resource_schema.go` implementing
`resource.Resource`, add model structs in `powerflex/models/`,
register in `provider.go`, add unit tests with mockey mocks, add
acceptance tests, create example HCL, and run `make generate` for
docs. If SDK lacks needed methods, add custom client in `client/`.

### What breaks?
**SDK version mismatch is a blocking defect.** Acceptance tests
against live hardware create real resources — failed test runs may
leave orphaned resources. State files contain secrets — use
encrypted remote backends.

### What depends on it?
Terraform Core (gRPC go-plugin), `github.com/dell/goscaleio`
v1.19.0, `hashicorp/terraform-plugin-framework` v1.13.0.

### What's undocumented?
The `powerflexProvider` struct holds two SDK clients:
`goscaleio.Client` for MDM API operations and
`goscaleio.GatewayClient` for Gateway API operations. Custom HTTP
clients in `client/` handle template CRUD and other operations not
in the SDK. The three-file resource pattern
(`*_resource.go`, `*_resource_schema.go`, helpers) is a convention
not documented in the SDK.

---

## Component Overview

Terraform provider for Dell PowerFlex (VxFlex OS) software-defined
storage. 28 resources and 22 data sources covering volumes, SDCs,
SDSs, protection domains, storage pools, fault sets, devices,
snapshots, snapshot policies, clusters, firmware repositories,
templates, compatibility management, NVMe hosts/targets,
replication, and more. Resources use `*_resource.go` naming under
`powerflex/provider/`. The provider package is nested:
`powerflex/provider/`.

---

## Architectural Rationale

The provider follows the standard Terraform Plugin Framework
architecture — a standalone Go binary communicating with Terraform
Core over gRPC.

**SDK strategy (Public):** Uses a public, versioned Go module on
GitHub. Provider and SDK release independently. Update via
`go get github.com/dell/goscaleio@<version>; go mod tidy`.

**Dual client architecture:** Unlike single-SDK providers, PowerFlex
requires two SDK clients (`goscaleio.Client` for MDM,
`goscaleio.GatewayClient` for Gateway) plus custom HTTP clients in
`client/` for operations not yet in the SDK.

All providers in the Dell Terraform family share this architecture:
Terraform Plugin Framework interfaces, `resource.Resource` for CRUD
resources, `datasource.DataSource` for read-only queries, models
with `tfsdk` struct tags, and mockey-based unit testing.

### Evolution

Originally built on Terraform Plugin SDK v2, then migrated to
Terraform Plugin Framework. Major refactor patterns over time
include:

- Client abstraction cleanup
- Model-driven design
- Error handling standardization
- Async / polling improvements
- Testing maturity (mockey adoption)

---

## Failure Modes & Gotchas

### 1. SDK version coupling

Each provider release is tested against exactly one SDK version.
A mismatch between provider and SDK version is a blocking defect.
Never update `go.mod` SDK versions without verifying against the
corresponding provider release notes.

### 2. Sensitive attributes must be marked

All credential fields must have `Sensitive: true` in the schema.
Without this, passwords appear in `terraform plan` output and state
files. This is enforced by code convention, not by the framework.

### 3. State file contains secrets

Terraform state files contain full resource representations
including credentials. Always use encrypted remote backends
(S3+KMS, Terraform Cloud) in production.

### 4. Custom client divergence

The `client/` directory contains HTTP clients for operations not
in the `goscaleio` SDK. These clients can drift from the API if
the SDK is updated but the custom clients are not. Verify both
after SDK upgrades.

### 5. State corruption

State corruption can occur with large state files and many managed
resources. Always use remote backends with locking (S3+DynamoDB,
Terraform Cloud) to prevent concurrent state writes.

### 6. Authentication edge cases

Credential rotation during active Terraform runs, expired tokens,
and network timeouts during provider configuration can leave the
provider in an unrecoverable state requiring `terraform init`
re-run.

### 7. Resource cleanup failures

Failed acceptance test runs or interrupted `terraform destroy` can
leave orphaned resources on the PowerFlex system. These must be
cleaned up manually via the management UI or REST API.

### Never Again

#### NA-001: State corruption from concurrent applies
- **Impact:** State file corruption when multiple engineers ran
  `terraform apply` simultaneously without state locking.
- **Constraint:** Must use remote backend with locking enabled.
- **Applies to:** All Dell Terraform providers.

#### NA-002: Orphaned resources from test failures
- **Impact:** Acceptance test resources left on system after test
  failure, consuming capacity.
- **Constraint:** Manual cleanup required; `TF_ACC=1` gating.
- **Applies to:** All Dell Terraform providers.

#### NA-003: Custom client API drift
- **Impact:** Custom HTTP clients in `client/` returned incorrect
  data after SDK upgrade changed API response format.
- **Constraint:** Custom clients must be verified after every SDK
  upgrade.
- **Applies to:** terraform-provider-powerflex.

### Evolution

Failure modes evolved with the SDK v2 → Plugin Framework migration.
Error handling was standardized during the model-driven design
refactor. The dual client architecture (SDK + custom HTTP clients)
introduced a new failure surface around client version mismatches.

---

## Performance Characteristics

**Large state files:** Performance degrades with many managed
resources in a single state file. Recommend splitting into multiple
Terraform workspaces or state files when managing >100 resources.

**API rate limiting:** PowerFlex systems may enforce API rate
limits. Bulk operations may hit these limits, causing transient
errors. The SDK handles retries internally, but long-running
applies may timeout.

**Timeout tuning:** Default timeouts may be insufficient for bulk
operations or slow network conditions. Increase for large
deployments.

### Evolution

Timeout was made configurable via environment variable after
production deployments hit the original hardcoded limit.

---

## Implicit Contracts

**Environment variable precedence:** env vars (`POWERFLEX_*`)
override HCL provider block values when both are set. This is
implemented in `Configure()` and is not documented as an explicit
contract.

**Authentication validation:** `Configure()` validates credentials
before any resource operations proceed. If this call fails, all
resource operations are blocked.

**TLS verification default:** `insecure` defaults to `false` —
TLS verification is on by default. Setting `insecure = true` is
a lab-only setting and must never be used in production.

**Acceptance test gating:** tests guarded by `TF_ACC=1` — never
run without live hardware credentials. Tests create real resources
that must be cleaned up manually if the test run fails.

### Evolution

Environment variable precedence was established during the SDK v2
era and carried forward into Plugin Framework. The authentication
validation call was added after production incidents with invalid
credentials causing cascading resource failures.

---

## Threading & Synchronization

Terraform Plugin Framework handles concurrency at the provider
level. Individual resource operations are not concurrent by default,
but Terraform Core may invoke multiple resource operations in
parallel during `terraform apply` (controlled by `-parallelism`
flag, default 10).

**Concurrent API access:** Multiple resources hitting the same
PowerFlex API endpoint simultaneously can cause contention. The
SDK clients are shared across all resource operations within a
single provider instance.

**Dual client concurrency:** Both `goscaleio.Client` and
`goscaleio.GatewayClient` plus custom HTTP clients are initialized
in `Configure()` and shared. No mutex protects concurrent access —
the SDK clients are expected to be thread-safe, but edge cases
exist under high parallelism.

### Evolution

Migration from SDK v2 to Plugin Framework changed the concurrency
model. SDK v2 serialized all operations; Plugin Framework allows
parallel resource operations. The dual-client architecture
introduced additional concurrency surface area.

---

## Build System & Configuration

Standard Makefile targets shared across all Dell Terraform
providers:

| Target | Purpose | Hardware Required |
|--------|---------|-------------------|
| `make build` | Compile provider binary | No |
| `make install` | Install to `~/.terraform.d/plugins/` | No |
| `make test` | Run unit tests | No |
| `make testacc` | Run acceptance tests | **Yes** |
| `make check` | Format, lint, vet | No |
| `make gosec` | Security scan | No |
| `make cover` | Generate coverage report | No |
| `make generate` | Generate documentation | No |

GoReleaser configuration: CGO_ENABLED=0, platforms (freebsd,
windows, linux, darwin), architectures (amd64, 386, arm, arm64).

### Evolution

Build system evolved from basic `go build` to Makefile with
linting, security scanning (gosec), and GoReleaser for
cross-platform releases. Testing maturity improved from minimal
acceptance tests to comprehensive mockey-based unit tests.

---

## Operational Knowledge

**Unit tests:** `bytedance/mockey` for runtime function patching.
No hardware required. Run with `make test`.

**Acceptance tests:** `terraform-plugin-testing` against live
hardware. Creates real resources. Run with `TF_ACC=1 make testacc`.
Clean up manually if tests fail mid-run.

### Evolution

Operational patterns matured with the mockey adoption for unit
tests, reducing dependence on live hardware for development
feedback loops.

---

## General Context

### Open Issues

No TODO/FIXME/HACK markers found in non-test source files.

### Glossary

| Term | Definition |
|------|------------|
| Plugin Framework | HashiCorp's Terraform Plugin Framework (`terraform-plugin-framework`) |
| mockey | `bytedance/mockey` — runtime function patching for unit tests |
| POWERFLEX | Environment variable prefix for this provider |
| goscaleio | `github.com/dell/goscaleio` — PowerFlex REST API SDK |
| MDM | Metadata Manager — core PowerFlex management component |
| Gateway | PowerFlex Manager — web-based management interface |

---

## References

- [Terraform Plugin Framework Docs](https://developer.hashicorp.com/terraform/plugin/framework)
- [Dell Terraform Registry](https://registry.terraform.io/namespaces/dell)

---

## Governance Spec Discrepancies

No discrepancies detected between code/SME knowledge and loaded
governance specs.
