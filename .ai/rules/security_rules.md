# Security Rules

- Hash password using argon2id
- Never store plaintext password
- JWT expiration required
- Validate all input
- Sanitize query filters
- Never log secrets
- Use RBAC middleware for protected routes
- CSRF protection for cookie-based auth