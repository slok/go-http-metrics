# Changelog

## [Unreleased]

## [0.13.0] - 2024-09-05

### Added

- Added `IgnoredPaths` option to ignore measuring specific exact paths.

### Changed

- Dependency updates.
- Use Go 1.23.

## [0.12.0] - 2024-04-19

### Changed

- `Gin` middleware now records the path using the regex from `FullPath()` instead of the real URL, avoiding cardinality explosion.
- Dependency updates.
- Use Go 1.22.

## [0.11.0] - 2023-10-20

### Changed

- Update dependencies to latest versions.

### Added

- Support Iris library.

## [0.10.0] - 2021-11-16

### Added

- Support Fasthttp library.
- Support gorestful v3

### Changed

- Drop support gorestful v2

## [0.9.0] - 2020-10-03

### Added

- New helper method to get an `std` HTTP provider (`std.HandlerProvider`) (used by various frameworks like Gorilla).
- Support Chi library.
- Support Alice library.
- Support Gorilla library.

## [0.8.0] - 2020-06-04

### Added

- New middleware helper for the Goji framework.

## [0.7.0] - 2020-06-02

Breaking change: The library has been refactored to be more flexible when adding new framework/libraries.

### Added

- New middleware helper for the Echo framework.

### Changed

- Refactored internally how the Middleware works and gets the data to make it easier to extend and more reliable.
  - Added `Reporter` interface as the service responsible of getting the data to be measured.
  - All different framwork helpers now implement with the new Reporter way.
- Fixed Gin returning duplicated data (#31).
- (Breaking) Standard handler now is on `middleware/std` instead of `middleware`.

### Removed

- Middleware interface in favor of a struct.

## [0.6.1] - 2020-02-07

### Changed

- Internal response writer interceptor implements `http.Hijacker` and `http.Flusher` interface.

## [0.6.0] - 2019-12-11

### Breaking changes

- The Recorder methods now receive properties in a single argument, this will make less breaking changes and better API (there where too many arguments for a function).

### Added

- Added new `service` property to identify the service.

### Changed

- The Recorder methods now receive properties in a single argument, this will make less breaking changes and better API (there where too many arguments for a function)

## [0.5.0] - 2019-12-10

### Added

- Gin compatible middleware.

## [0.4.0] - 2019-03-27

### Breaking changes

- The Recorder methods now receive a context argument.

### Added

- OpenCensus recorder implementation.

## [0.3.0] - 2019-03-24

### Added

- Inflight requests metric per handler.

## [0.2.0] - 2019-03-22

### Added

- Metrics of HTTP response size in bytes.
- Make the label names of Prometheus recorder configurable.

## [0.1.0] - 2019-03-18

### Added

- Gorestful compatible middleware.
- Httprouter compatible middleware.
- Negroni compatible middleware.
- Option to group by status codes.
- Predefined handler label.
- URL infered handler label.
- Middleware.
- HTTP latency requests.
- Prometheus recorder.

[unreleased]: https://github.com/slok/go-http-metrics/compare/v0.13.0...HEAD
[0.13.0]: https://github.com/slok/go-http-metrics/compare/v0.12.0...v0.13.0
[0.12.0]: https://github.com/slok/go-http-metrics/compare/v0.11.0...v0.12.0
[0.11.0]: https://github.com/slok/go-http-metrics/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/slok/go-http-metrics/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/slok/go-http-metrics/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/slok/go-http-metrics/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/slok/go-http-metrics/compare/v0.6.1...v0.7.0
[0.6.1]: https://github.com/slok/go-http-metrics/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/slok/go-http-metrics/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/slok/go-http-metrics/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/slok/go-http-metrics/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/slok/go-http-metrics/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/slok/go-http-metrics/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/slok/go-http-metrics/releases/tag/v0.1.0
