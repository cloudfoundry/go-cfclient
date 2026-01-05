# Changelog

## [Unreleased]

### Added

- Support for setting and reading application lifecycle (`buildpack`, `docker`, and `cnb`/Cloud Native Buildpacks) in V3 API.
- Custom marshaling/unmarshaling for `Lifecycle` struct to support multiple lifecycle types.
- Unit tests for all lifecycle types, including CNB.
- Support for Service Broker-provided metadata (labels and attributes) on Service Instances. This includes the `BrokerProvidedMetadata` field on `ServiceInstance`, `ServiceInstanceManagedCreate`, and `ServiceInstanceManagedUpdate` structs, along with fluent builder methods for managing labels and attributes.

### Changed

- All lifecycle-related test expectations updated to match new marshaling output (both `type` and `data` fields).

### Notes

- For CNB, only `buildpacks` and `stack` fields are currently supported. Additional fields may be added in the future as the API evolves.
- No breaking changes for existing users of buildpack or docker lifecycles.
