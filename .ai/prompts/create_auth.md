Create Auth module:

Endpoints:
- POST /auth/register
- POST /auth/login
- POST /auth/refresh

Requirements:
- argon2id password hashing
- JWT access token
- refresh token support
- RBAC middleware
- integration tests