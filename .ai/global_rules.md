# Global Rules

## Backend Rules
- No SELECT * in queries.
- No raw SQL outside sqlc query files (adapt PHP: use PDO prepared statements; Node: use parameterized queries in pg/sequelize).
- No DB call in handler/controller.
- No business logic in handler/controller.
- Prefer UUID primary key.
- Use context.Context everywhere (adapt PHP: use request context in Laravel; Node: use req.context in Express).
- Wrap errors using fmt.Errorf (adapt PHP: throw custom exceptions; Node: use Error wrapping).
- Validate all input.
- No panic in production code.
- Prefer typed errors over string comparison.
- No global DB access.

## Frontend Rules
- Feature-based folder structure.
- No direct API call outside service layer (adapt React with Axios/RTK Query).
- Use typed DTO matching backend.
- Centralized API client.
- No business logic inside component.
- Pagination must support cursor.

## Performance Rules
- Add index for filtered columns.
- Use EXPLAIN ANALYZE for heavy queries.
- Avoid N+1 queries.
- Limit max page size.
- Use cursor pagination.
- Avoid loading full table.

## Security Rules
- Hash password using argon2id (adapt PHP: password_hash; Node: argon2).
- Never store plaintext password.
- JWT expiration required.
- Validate all input.
- Sanitize query filters.
- Never log secrets.
- Use RBAC middleware for protected routes.
- CSRF protection for cookie-based auth.

## Migration Rules
- Always create new migration file.
- Never edit historical migration.
- Must include Down migration.
- Add index for foreign key.
- Use snake_case naming.
- Use IF NOT EXISTS carefully.
- No destructive change without rollback note.

## Testing Rules
- Integration test required for every endpoint.
- Validation failure test required.
- Auth failure test required.
- Pagination test required (if applicable).
- Use table-driven test.
- Coverage target: >= 80% for service layer.

## CI/CD Rules
- Run go vet (adapt PHP: phpstan; Node: eslint).
- Run tests.
- Run sqlc generate (adapt: run migrations).
- Apply migration on test DB.
- Build frontend.
- Build Docker image.
- Fail pipeline if: SELECT * found, migration fails, test fails.