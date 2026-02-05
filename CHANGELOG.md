# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-02-06

### ⚠️ BREAKING CHANGES

- Complete API rewrite with new endpoint structure
- All JSON responses now use `snake_case` notation (was `CamelCase`)
- Service now requires root privileges for `wg-quick up/down` operations
- Removed legacy OpenAPI-generator models

### Added

- **Clean Architecture**: Domain-driven design with separate layers (domain/usecase/infrastructure/interface)
- **Fiber v2**: Migrated from Echo v4 to Fiber v2 web framework
- **Device Discovery**: `GET /devices/` now returns both running interfaces AND config-only devices
- **Running Status**: New `running` field on Device entity indicates if interface is up
- **Multi-Directory Config Search**: Searches `/etc/wireguard`, `/usr/local/etc/wireguard`, `/opt/homebrew/etc/wireguard`
- **Config Write-Back**: Saves configs back to the directory where they were found
- **macOS Support**: Automatic `utun*` interface name resolution via `/var/run/wireguard/*.name`
- **Swagger UI**: Integrated documentation at `/swagger/`
- **Config Dump Service**: Periodic background service saves running state to disk
- **Platform Defaults**: Platform-aware default configuration directories (Linux/macOS/FreeBSD)
- **GitHub Actions CI**: Multi-platform builds, automated packaging, and releases

### Changed

- **Authentication**: Bearer token middleware with configurable token
- **JSON Format**: All responses use `snake_case` (e.g., `public_key`, `allowed_ips`)
- **OpenAPI**: Clean model names (`Device`, `Peer`, `Error`) instead of `entity.*` prefixes
- **Packaging**: Service runs as root (required for wg-quick operations)
- **wg-quick Commands**: 30-second timeout prevents hanging on sudo password prompts

### Removed

- Legacy Echo v4 handlers (`handlers/`)
- OpenAPI-generator models (`models/`)
- Legacy storage layer (`storage/`)
- Legacy utility functions (`utils/`)
- Drone CI configuration (`.drone.yml`)

### Fixed

- Device listing now shows configs even when interface is not running
- macOS WireGuard integration with userspace `utun` interfaces
- Non-interactive wg-quick execution (no sudo password prompts)

## [1.x.x] - Previous

See git history for changes prior to the v2.0.0 rewrite.
