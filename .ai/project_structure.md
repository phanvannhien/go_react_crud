# Project Structure Contract

Backend (Go Clean Architecture)

Dependency Direction:
handler → service → sqlc
domain independent

Rules:
- No DB call in handler
- No raw SQL outside /db/query
- sqlc generated code must stay in /db/sqlc
- No circular imports
- No business logic in handler
- No Echo dependency in service

Frontend (React SPA)

Rules:
- Feature-based structure
- No API call in components
- All network logic in api.ts
- No shared model between backend and frontend

Violation = reject generation