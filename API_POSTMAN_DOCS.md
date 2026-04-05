# Finance Data Processing API Reference (Postman Ready)

## Base Setup

- Base URL: `http://localhost:8080`
- API Prefix: `/api`
- Content Type: `application/json`

## Authentication

- Auth Type: Bearer Token (JWT)
- Header:

```http
Authorization: Bearer <access_token>
```

- Login endpoint returns token from `data.access_token`.

## Standard Response Shapes

### Success

```json
{
  "data": {}
}
```

### Error

```json
{
  "code": "invalid_input",
  "message": "invalid request payload",
  "details": null
}
```

## Role Access Matrix

- `viewer`: read records + read dashboard
- `analyst`: viewer permissions + create/update/delete records
- `admin`: full access including user management

---

## 1) Auth APIs

### 1.1 Login

- Method: `POST`
- URL: `/api/auth/login`
- Auth: No

#### Request Body

```json
{
  "email": "admin@example.com",
  "password": "ChangeMeNow123!"
}
```

#### Success Response (`200`)

```json
{
  "data": {
    "access_token": "<jwt_token>",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

#### Error Responses

- `400` invalid payload
- `401` invalid credentials / inactive user

```json
{
  "code": "unauthorized",
  "message": "invalid credentials",
  "details": null
}
```

---

## 2) User APIs (Admin Only)

> All user endpoints require Bearer token with admin role.

### 2.1 Create User

- Method: `POST`
- URL: `/api/users`
- Auth: Yes (`admin`)

#### Request Body

```json
{
  "name": "Analyst One",
  "email": "analyst@example.com",
  "password": "Pass@12345",
  "role": "analyst"
}
```

Allowed role values: `viewer`, `analyst`, `admin`

#### Success Response (`201`)

```json
{
  "data": {
    "id": "1bc1f9ca-8938-4af5-9a7f-f15f0c94ec48"
  }
}
```

#### Error Responses

- `400` invalid role / missing fields / invalid body
- `401` missing or invalid token
- `403` insufficient role

---

### 2.2 List Users

- Method: `GET`
- URL: `/api/users`
- Auth: Yes (`admin`)

#### Query Params

- `limit` (optional, default `25`)
- `offset` (optional, default `0`)

#### Example

`/api/users?limit=10&offset=0`

#### Success Response (`200`)

```json
{
  "data": [
    {
      "id": "1bc1f9ca-8938-4af5-9a7f-f15f0c94ec48",
      "name": "Analyst One",
      "email": "analyst@example.com",
      "role": "analyst",
      "is_active": true,
      "created_at": "2026-04-05T12:00:00Z",
      "updated_at": "2026-04-05T12:00:00Z"
    }
  ]
}
```

---

### 2.3 Update User

- Method: `PATCH`
- URL: `/api/users/:id`
- Auth: Yes (`admin`)

#### Path Params

- `id` (UUID user id)

#### Request Body

```json
{
  "name": "Analyst Renamed",
  "role": "viewer"
}
```

#### Success Response (`200`)

```json
{
  "data": {
    "updated": true
  }
}
```

#### Error Responses

- `400` invalid input / user not found

---

### 2.4 Activate/Deactivate User

- Method: `PATCH`
- URL: `/api/users/:id/active`
- Auth: Yes (`admin`)

#### Path Params

- `id` (UUID user id)

#### Request Body

```json
{
  "active": false
}
```

#### Success Response (`200`)

```json
{
  "data": {
    "updated": true
  }
}
```

---

## 3) Financial Record APIs

## Record Object

```json
{
  "id": "31f4a738-7f4e-4f65-b5cd-6dd72d4c346f",
  "user_id": "1bc1f9ca-8938-4af5-9a7f-f15f0c94ec48",
  "amount": "1200.50",
  "type": "income",
  "category": "salary",
  "date": "2026-04-01T08:30:00Z",
  "notes": "April salary",
  "created_at": "2026-04-05T12:00:00Z",
  "updated_at": "2026-04-05T12:00:00Z"
}
```

### 3.1 Create Record

- Method: `POST`
- URL: `/api/records`
- Auth: Yes (`analyst`, `admin`)

#### Request Body

```json
{
  "amount": "1200.50",
  "type": "income",
  "category": "salary",
  "date": "2026-04-01T08:30:00Z",
  "notes": "April salary"
}
```

Notes:
- `type` must be `income` or `expense`
- `date` must be RFC3339; if omitted/empty, server uses current UTC time

#### Success Response (`201`)

```json
{
  "data": {
    "id": "31f4a738-7f4e-4f65-b5cd-6dd72d4c346f"
  }
}
```

---

### 3.2 List Records

- Method: `GET`
- URL: `/api/records`
- Auth: Yes (`viewer`, `analyst`, `admin`)

#### Query Params

- `from` (optional, RFC3339)
- `to` (optional, RFC3339)
- `category` (optional)
- `type` (optional: `income` or `expense`)
- `limit` (optional, default `25`)
- `offset` (optional, default `0`)

#### Example

`/api/records?from=2026-04-01T00:00:00Z&to=2026-04-30T23:59:59Z&type=expense&category=food&limit=20&offset=0`

#### Success Response (`200`)

```json
{
  "data": [
    {
      "id": "31f4a738-7f4e-4f65-b5cd-6dd72d4c346f",
      "user_id": "1bc1f9ca-8938-4af5-9a7f-f15f0c94ec48",
      "amount": "1200.50",
      "type": "income",
      "category": "salary",
      "date": "2026-04-01T08:30:00Z",
      "notes": "April salary",
      "created_at": "2026-04-05T12:00:00Z",
      "updated_at": "2026-04-05T12:00:00Z"
    }
  ]
}
```

Access scope:
- `admin`: can view all records
- `viewer`/`analyst`: only their own records

---

### 3.3 Get Record By ID

- Method: `GET`
- URL: `/api/records/:id`
- Auth: Yes (`viewer`, `analyst`, `admin`)

#### Path Params

- `id` (UUID record id)

#### Success Response (`200`)

```json
{
  "data": {
    "id": "31f4a738-7f4e-4f65-b5cd-6dd72d4c346f",
    "user_id": "1bc1f9ca-8938-4af5-9a7f-f15f0c94ec48",
    "amount": "1200.50",
    "type": "income",
    "category": "salary",
    "date": "2026-04-01T08:30:00Z",
    "notes": "April salary",
    "created_at": "2026-04-05T12:00:00Z",
    "updated_at": "2026-04-05T12:00:00Z"
  }
}
```

#### Error Responses

- `404` not found or not accessible by requester

---

### 3.4 Update Record

- Method: `PATCH`
- URL: `/api/records/:id`
- Auth: Yes (`analyst`, `admin`)

#### Path Params

- `id` (UUID record id)

#### Request Body

```json
{
  "amount": "1100.00",
  "type": "income",
  "category": "salary",
  "date": "2026-04-01T08:30:00Z",
  "notes": "Corrected amount"
}
```

#### Success Response (`200`)

```json
{
  "data": {
    "updated": true
  }
}
```

Access scope:
- `admin`: can update any record
- `analyst`: only own records

---

### 3.5 Delete Record (Soft Delete)

- Method: `DELETE`
- URL: `/api/records/:id`
- Auth: Yes (`analyst`, `admin`)

#### Path Params

- `id` (UUID record id)

#### Success Response (`204`)

No response body.

#### Error Responses

- `404` record not found

Access scope:
- `admin`: can delete any record
- `analyst`: only own records

---

## 4) Dashboard APIs

> All dashboard endpoints require auth. `viewer`, `analyst`, and `admin` are allowed.

Shared Query Params:
- `from` (optional, RFC3339)
- `to` (optional, RFC3339)

Scope:
- `admin`: all users data
- `viewer`/`analyst`: own data only

### 4.1 Dashboard Summary

- Method: `GET`
- URL: `/api/dashboard/summary`

#### Example

`/api/dashboard/summary?from=2026-04-01T00:00:00Z&to=2026-04-30T23:59:59Z`

#### Success Response (`200`)

```json
{
  "data": {
    "total_income": "5000.00",
    "total_expenses": "1200.00",
    "net_balance": "3800.00"
  }
}
```

---

### 4.2 Category Totals

- Method: `GET`
- URL: `/api/dashboard/categories`

#### Success Response (`200`)

```json
{
  "data": [
    {
      "category": "salary",
      "total": "5000.00"
    },
    {
      "category": "food",
      "total": "300.00"
    }
  ]
}
```

---

### 4.3 Trends (Monthly)

- Method: `GET`
- URL: `/api/dashboard/trends`

#### Success Response (`200`)

```json
{
  "data": [
    {
      "period": "2026-03",
      "income": "3000.00",
      "expense": "900.00"
    },
    {
      "period": "2026-04",
      "income": "5000.00",
      "expense": "1200.00"
    }
  ]
}
```

---

### 4.4 Recent Transactions

- Method: `GET`
- URL: `/api/dashboard/recent`

#### Query Params

- `limit` (optional, default `10`)
- `from` (optional, RFC3339)
- `to` (optional, RFC3339)

#### Success Response (`200`)

```json
{
  "data": [
    {
      "id": "31f4a738-7f4e-4f65-b5cd-6dd72d4c346f",
      "user_id": "1bc1f9ca-8938-4af5-9a7f-f15f0c94ec48",
      "amount": "1200.50",
      "type": "income",
      "category": "salary",
      "date": "2026-04-01T08:30:00Z",
      "notes": "April salary",
      "created_at": "2026-04-05T12:00:00Z",
      "updated_at": "2026-04-05T12:00:00Z"
    }
  ]
}
```

---

## 5) Common HTTP Statuses

- `200` success
- `201` created
- `204` deleted (no content)
- `400` bad request / validation failure
- `401` unauthorized (missing/invalid token)
- `403` forbidden (role restriction)
- `404` resource not found
- `500` internal error

---

## 6) Suggested Postman Environment Variables

- `baseUrl` = `http://localhost:8080`
- `accessToken` = (set after login)
- `adminUserId` = (optional)
- `recordId` = (optional)

Use Authorization tab or header:

```http
Authorization: Bearer {{accessToken}}
```

---

## 7) Quick Postman Collection Flow

1. Login (save `data.access_token` to `accessToken`)
2. Create User (viewer/analyst)
3. Create Record
4. List Records (with filters)
5. Update Record
6. Dashboard Summary / Categories / Trends / Recent
7. Delete Record
