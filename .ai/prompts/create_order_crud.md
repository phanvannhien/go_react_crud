# Prompt: Create Order CRUD Feature

Entry Point: create_feature
Feature Name: Order CRUD
Environment: development

## Objective

Implement full CRUD operations for Order module including:

- Create Order
- Get Order by ID
- Update Order
- Delete Order
- List Orders with pagination

## Constraints

- Must follow workflows/feature_workflow.md
- Must apply runtime risk_model.md
- Must pass guard_policy.md validation
- Must satisfy ci_contract.md

## Schema Requirements

Tables:

orders
- id
- user_id (FK → users.id)
- status (pending, paid, cancelled, shipped)
- total_amount
- created_at
- updated_at

order_items
- id
- order_id (FK → orders.id)
- product_id
- quantity
- price

Migration Requirements:
- Include rollback
- No destructive DROP without backup
- Index foreign keys

## Architecture Requirements

Must include:

- Domain Model
- Repository Layer
- Service Layer (business logic)
- Controller Layer
- Input Validation
- Standardized API Response
- Transaction boundary for write operations

## Business Rules

- Cannot update order if status = shipped
- Cannot delete order if status = paid
- total_amount must be recalculated from order_items
- Only owner can access their order

## API Requirements

List endpoint:
- Must include pagination (page, limit)
- Must enforce maximum limit
- Must return standardized error format

## Testing Requirements

Minimum tests:

- Create order success
- Create order invalid payload fail
- Update shipped order fail
- Delete paid order fail
- Pagination works correctly

## Output Type Declaration

Expected output types:
- migration_script
- api_endpoint
- service_logic
- test_suite

All outputs must pass output-type guard validation.

## Risk Declaration

Feature includes:
- migration_change

No:
- auth_modification
- production_environment

Estimated Base Risk: Medium