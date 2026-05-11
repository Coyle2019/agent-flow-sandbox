# Agent Flow Sandbox

Backend service for user registration and authentication.

## API Endpoints

### POST /api/v1/register

Register a new user account.

**Request:**
```json
{
  "username": "string",
  "email": "string",
  "password": "string (min 8 characters)"
}
```

**Response (201):**
```json
{
  "user_id": "uuid",
  "username": "string",
  "email": "string",
  "access_token": "jwt-token"
}
```

**Error Responses:**
- `400`: Validation error (invalid email, password too short, missing fields)
- `409`: Email already registered

## Setup

```bash
pip install flask pytest
python src/app.py
```

## Testing

```bash
pytest tests/ -v
```
