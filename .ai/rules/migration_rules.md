# Migration Rules

- Always create new migration file
- Never edit historical migration
- Must include Down migration
- Add index for foreign key
- Use snake_case naming
- Use IF NOT EXISTS carefully
- No destructive change without rollback note