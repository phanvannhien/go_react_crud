# CI Enforcement Contract

Every Pull Request MUST validate:

- Integration tests passing
- Migration includes rollback
- No hard stop violations
- Pagination enforced on list endpoints
- API contract unchanged without version bump

Deployment Gate:

- SLO validation complete
- Observability enabled
- Rollback plan verified

Failure = Block Merge