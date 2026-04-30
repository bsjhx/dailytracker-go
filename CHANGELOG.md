# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-04-30

### Added
- Daily tracking for work and personal scores (0-5 scale)
- Weekly statistics and averages
- User authentication with sessions
- File-based database migrations using golang-migrate
- SQLite database with pure Go implementation (modernc.org/sqlite)
- Docker and Docker Compose support
- User creation via API endpoint
- Login/logout functionality
- View and edit last 30 entries
- Responsive web interface
- Project structure following Go standards

### Technical
- Migrations run automatically on app startup
- Idempotent migrations (safe to run multiple times)
- Session-based authentication
- Database auto-creation with parent directories

[Unreleased]: https://github.com/bsjhx/dailytracker-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/bsjhx/dailytracker-go/releases/tag/v0.1.0
