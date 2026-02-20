# Skill: Security Audit

Backend Checks:
- SQL injection (must use sqlc only)
- Missing auth middleware
- Missing input validation
- Missing rate limit
- Hardcoded secrets
- CORS misconfiguration

Auth:
- JWT expiration
- Token validation
- Role-based access check

Frontend:
- No token in localStorage
- No secret in client bundle
- No direct eval
- Escape dynamic HTML

Database:
- No superuser connection
- Principle of least privilege

Output:
- Vulnerability list
- Severity (critical/high/medium/low)
- Fix recommendation

Additional Security Checks:

- Sort parameter must be validated (prevent SQL injection)
- Filter fields must be validated
- LIMIT must have max cap (e.g. max 100)
- Input must be sanitized
- No blind trust of client-side validation
- Password must be hashed using bcrypt in service layer