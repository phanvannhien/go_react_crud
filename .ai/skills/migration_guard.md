# Migration Guard Skill

Purpose:
Protect database integrity in PostgreSQL + dbmate + sqlc environment.

---

## Core Principles

- Never modify historical migration files.
- Always create a new migration file.
- Every migration must include both Up and Down.
- Schema changes must be backward compatible unless explicitly stated.
- Every foreign key must be indexed.
- Every filtered column must be indexed.

---

## Validation Checklist

Before generating migration:

1. Is table name snake_case?
2. Is primary key UUID?
3. Is DEFAULT gen_random_uuid() used?
4. Are timestamps included?
5. Is NOT NULL applied correctly?
6. Is unique constraint explicit?
7. Is index created for filter columns?
8. Is Down migration safe?

---

## Forbidden Operations

- Editing existing migration
- Dropping column without replacement plan
- Dropping table without rollback safety
- Removing index without performance review

---

## Breaking Change Policy

If migration:
- Renames column
- Changes type
- Drops constraint

Then:
- Require data backfill plan
- Require backward compatibility note
- Require service adaptation

---

## sqlc Sync Requirement

After migration:
- Regenerate sqlc
- Ensure no compile errors
- Ensure repository still compiles
- Ensure integration tests pass

---

## High Risk Migration Types

- Auth table change
- Role change
- Foreign key modification
- Index removal

These require explicit explanation.