# Global Architecture Rules

Backend and frontend must be decoupled.

Backend:
- RESTful
- Stateless
- No session storage
- JWT for authentication

Frontend:
- Token stored in memory or httpOnly cookie
- No localStorage for sensitive data

Cross Rules:
- All timestamps must be RFC3339
- All IDs must be UUID
- Pagination standard:
  {
    "data": [],
    "meta": {
      "page": 1,
      "limit": 20,
      "total": 100
    }
  }