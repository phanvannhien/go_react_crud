# Transaction Pattern Skill

Purpose:
Ensure ACID compliance in service layer.

---

## Golden Rule

Transaction management belongs ONLY in service layer.

Repository must not:
- Start transaction
- Commit transaction
- Rollback transaction

---

## Standard Pattern

Service:

tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }

defer func() {
    if err != nil {
        tx.Rollback()
    } else {
        tx.Commit()
    }
}()

Repository:

- Accept *sql.Tx or interface
- Execute query using tx
- Return typed error

---

## Multi-step Operation Rule

If service:
- Creates entity
- Updates another table
- Writes audit log

Then MUST use transaction.

---

## No Nested Transactions

Never open nested transactions.
Use same tx across repository calls.

---

## Error Handling

Business logic error:
- Rollback

Validation error:
- Do not open transaction

System error:
- Rollback

---

## Forbidden

- Transaction inside handler
- Transaction inside repository
- Ignoring commit error
- Swallowing rollback error