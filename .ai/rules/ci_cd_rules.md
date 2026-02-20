# CI/CD Rules

CI must:

- Run go vet
- Run go test
- Run sqlc generate
- Apply migration on test DB
- Build frontend
- Build Docker image

Fail pipeline if:
- SELECT * found
- Migration fails
- Test fails