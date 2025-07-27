# Brand Activations: Make Sponsorships Smarter

A platform that provides brands with meaningful insights about their sponsorships while maintaining a non-intrusive user experience.

## Project Overview

This project implements a system for tracking and analyzing brand sponsorship effectiveness through:
- Event attendance tracking
- Product reception measurement
- Post-event engagement tracking
- Visual content collection
- Purchase attribution

## Repository Structure

- `/mobile-client` - React Native mobile application (TypeScript)
- `/backend` - Go-based API and services
- `/database` - SQLite database schema and migrations
- `/design-specs` - Project requirements and design documents

## Development Setup

### Prerequisites
- Node.js 16+
- Go 1.20+
- SQLite 3

### Mobile Client Setup
```bash
cd mobile-client
npm install / yarn install
npm start / yarn start
```

### Backend Setup
```bash
cd backend
go mod download
go run cmd/api/main.go
```

### Database Setup
```bash
cd backend/data
sqlite3 brand_activations.db < ../../schema.sql
```

### note
to login into brand use email: brand@example.com
password: password

you should be able to create user account by signup

## Git Branching Strategy

We follow a structured branching strategy:

- `main` - Production-ready code
- `develop` - Integration branch for features

### Branch Naming Conventions
- Feature branches: `feature/short-description`
- Bug fixes: `bugfix/issue-number-description`
- Hotfixes: `hotfix/issue-number-description`
- Release branches: `release/version-number`
- Documentation: `docs/what-changed`
- Refactoring: `refactor/what-changed`

### Commit Message Standards
- Use the imperative mood: "Add feature" not "Added feature"
- Start with a capital letter
- Keep the subject line under 50 characters
- Reference issue numbers when applicable: "Fix login bug (fixes #123)"
- Structure: `<type>(<scope>): <subject>`
  - Types: feat, fix, docs, style, refactor, test, chore

### Merge Workflow
- Create pull requests for merging into develop/main
- Require at least one code review before merging
- Squash commits when merging feature branches
- Delete branches after merging