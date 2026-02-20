Create User module with:

Fields:
- id (uuid, primary key)
- email (unique)
- password (hashed)
- name
- role
- created_at

Endpoints:
- POST /users
- GET /users/:id
- GET /users (cursor pagination)

Include:
- migration
- sqlc queries
- repository
- service
- handler
- integration test
- OpenAPI update