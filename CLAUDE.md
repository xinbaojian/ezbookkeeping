# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ezBookkeeping is a full-stack personal finance application with a unique dual-UI architecture:

- **Backend**: Go 1.25 with Gin framework, supporting SQLite/MySQL/PostgreSQL
- **Frontend**: Vue 3 + TypeScript with dual UI modes:
    - Desktop: Vuetify (Material Design) - served via `desktop.html`
    - Mobile: Framework7 (iOS/Android native feel) - served via `mobile.html`

## Common Development Commands

### Frontend Development

```bash
npm run serve          # Start dev server (port 8081, proxies to backend on 8080)
npm run build          # Build for production
npm run lint           # Run ESLint + Vue TypeScript checks
npm run test           # Run Jest tests
```

### Backend Development

```bash
go run ezbookkeeping.go server run    # Run backend directly
go test ./... -v                        # Run all backend tests
go test ./... -v -skip "TestName"      # Skip specific tests
```

### Full Build

```bash
./build.sh backend     # Build Go binary only
./build.sh frontend    # Build frontend assets only
./build.sh package     # Build complete package (both)
./build.sh docker      # Build Docker image
```

## Architecture Overview

### Backend Structure (`pkg/`)

- `api/`: REST API endpoints organized by domain (accounts, transactions, etc.)
- `auth/`: Authentication logic including OAuth2 providers
- `datastore/`: Database abstraction layer with implementations for SQLite, MySQL, PostgreSQL
- `services/`: Business logic orchestration
- `models/`: Data models
- `converters/`: Data import/export (CSV, OFX, QIF, etc.)
- `mcp/`: Model Context Protocol for AI integration

### Frontend Structure (`src/`)

- `components/desktop/`: Vuetify components for desktop UI
- `components/mobile/`: Framework7 components for mobile UI
- `stores/`: Pinia stores organized by domain
- `lib/`: Core utilities and helpers
- `locales/`: Internationalization files

### Key Architectural Decisions

1. **Dual UI from Single Codebase**: The application detects device type and serves appropriate UI framework
    - Desktop uses Vuetify 3 for Material Design
    - Mobile uses Framework7 for native feel
    - Shares common components and logic where possible

2. **Proxy Development Setup**: Vite dev server (8081) proxies all API calls to Go backend (8080)
    - `/api/*` - Main API endpoints
    - `/oauth2/*` - OAuth2 flow
    - `/mcp/*` - Model Context Protocol

3. **Build Configuration**: Multi-entry Vite build creates separate bundles for desktop/mobile UIs

4. **Database Migration System**: All schema changes handled via Go migrations in `pkg/datastore/dbupgrade/`

### Testing Approach

- **Frontend**: Jest with TypeScript support, tests in `**/__tests__/` directories
- **Backend**: Standard Go testing with `*_test.go` files
- **CI/CD**: GitHub Actions workflows for multiple platforms (see `.github/workflows/`)

### Development Constraints

- Go 1.25+ required for backend development
- Node.js 24+ for frontend builds
- GCC required for CGO (MySQL/PostgreSQL drivers)
- Frontend linting includes both ESLint and Vue TypeScript checks
- Build scripts automatically run linting and tests before compilation

### Common Development Tasks

1. **Adding a New API Endpoint**:
    - Add route in appropriate `pkg/api/*_routes.go` file
    - Implement handler in corresponding `pkg/api/*_handlers.go`
    - Add any new models to `pkg/models/`
    - Add frontend service methods in `src/services/`

2. **Modifying UI**:
    - Desktop: Edit components in `src/components/desktop/`
    - Mobile: Edit components in `src/components/mobile/`
    - Shared logic: Edit `src/lib/` utilities

3. **Database Schema Changes**:
    - Create new migration in `pkg/datastore/dbupgrade/`
    - Follow existing naming pattern (e.g., `upgrade_0009.go`)
    - Implement both `upgrade()` and `downgrade()` functions

4. **Adding Translations**:
    - Add new keys to all `src/locales/*.json` files
    - Follow existing structure for consistency
