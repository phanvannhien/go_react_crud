# Backend Coding Rules

- Use context.Context everywhere
- No panic in production code
- Wrap errors using fmt.Errorf("context: %w", err)
- Return typed errors for business logic
- Validate request payload
- No global DB access
- No raw SQL outside query files
- No SELECT *
- Prefer UUID primary key