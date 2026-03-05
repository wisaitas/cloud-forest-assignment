# Task 1: Database Schema Design

Data schema for the Cloud Servers RESTful API. Designed to work with **in-memory storage** (Task 2) and to be straightforward to map to a real database later.

---

## Entity Relationship Overview

```
┌─────────────┐       ┌─────────────┐       ┌──────────────────┐
│   Users     │──1:N──│   Servers   │       │  ActivityLogs    │
└─────────────┘       └─────────────┘       └──────────────────┘
       │                      │                        │
       └──────────────────────┴──────────1:N──────────┘
```

- **Users** own **Servers** (one user, many servers).
- **ActivityLogs** record actions performed via the API (linked to user and optionally to a server).

---

## 1. Users

Stores credentials for authentication.

| Field      | Type     | Constraints | Description                    |
|-----------|----------|-------------|--------------------------------|
| `id`      | string   | PK, unique  | Application-side user ID (e.g. UUID). |
| `email`   | string   | required, unique | Login email.              |
| `password`| string   | required    | Hashed password (e.g. bcrypt). |

**Note:** For the assignment seed data, use User ID `123123123` with email `john.smith@gmail.com` and password `not-so-secure-password`.

---

## 2. Servers

Stores server metadata and the **Infrastructure Resource ID** from the infrastructure microservice. The API uses its own **Server ID**; the Infrastructure Resource ID is only for calling the external service.

| Field                        | Type   | Constraints | Description |
|-----------------------------|--------|-------------|-------------|
| `id`                        | string | PK, unique  | **API Service Server ID** (e.g. UUID). Used in API paths like `POST /servers/:server-id/power`. |
| `user_id`                   | string | required, FK | Owner; references `Users.id`. |
| `infrastructure_resource_id`| string | required, unique | **Infrastructure Resource ID** returned by the microservice when provisioning. Used only when calling the infra API (e.g. power on/off). Must not be used as the primary Server ID in this API. |
| `sku`                       | string | required    | Server SKU (e.g. `C1-R1GB-D40GB`). Must be validated against `/v1/skus`. |
| `power_status`              | string | required    | Power state: `"on"` or `"off"`. Align with infra service (e.g. `running` → `on`, `stopped` → `off`). |
| `created_at`                | string (ISO 8601) or timestamp | optional | When the server was provisioned in this API. |

**Critical:**  
- **Server ID** = `id` (our primary key, used in responses and routes).  
- **Infrastructure Resource ID** = `infrastructure_resource_id` (stored for calling the microservice; not exposed as the main server identifier to clients if you want to keep internal vs external IDs clear).

---

## 3. ActivityLogs

Audit trail for actions performed via the API.

| Field          | Type   | Constraints | Description |
|----------------|--------|-------------|-------------|
| `id`           | string | PK, unique  | Log entry ID (e.g. UUID). |
| `user_id`      | string | required, FK | User who performed the action; references `Users.id`. |
| `action`       | string | required    | Action type (e.g. `login`, `provision_server`, `power_on`, `power_off`). |
| `resource_type`| string | optional    | Affected resource (e.g. `server`). |
| `resource_id`  | string | optional    | ID of the resource (e.g. server `id` or `infrastructure_resource_id` for traceability). |
| `details`      | string or JSON | optional | Extra context (e.g. request body, error message). |
| `created_at`   | string (ISO 8601) or timestamp | optional | When the action occurred. |

---

## In-Memory Implementation Notes (Task 2)

- Use **slices/maps keyed by ID** (e.g. `map[string]*Server`) for O(1) lookup.
- **Server ID:** Generate a new UUID (or unique string) when creating a server; use this as `id` in your store and in API paths.
- **Infrastructure Resource ID:** Persist the value returned from `POST /v1/resources` in `infrastructure_resource_id`; use it only when calling `GET/POST /v1/resources/{id}/power` and similar infra endpoints.
- **ActivityLogs:** Append a new log entry on login, server provision, power on/off, and other relevant operations.

---

## Example: Server Provisioning Flow

1. Client: `POST /servers` with `{ "sku": "C1-R1GB-D40GB" }`.
2. API: Validate SKU via infra `GET /v1/skus` (or equivalent).
3. API: Call infra `POST /v1/resources` with `{ "sku": "C1-R1GB-D40GB" }` → receive `id` (e.g. `"i-12345678"`).
4. API: Create server in store with:
   - `id` = new UUID (our Server ID),
   - `infrastructure_resource_id` = `"i-12345678"`,
   - `sku`, `user_id`, `power_status` (e.g. from infra response).
5. API: Append ActivityLog (e.g. `provision_server`, `resource_type=server`, `resource_id=our server id`).
6. API: Response `{ "success": true, "id": "<Infrastructure Resource ID>" }` (per task: respond with infra id in `id` field; keep our own server id in store for routing and power actions).

*(If the task expects the response `id` to be the Infrastructure Resource ID, use that in the response; internally still use your own Server ID for `POST /servers/:server-id/power` by mapping `server-id` from the path to your store—e.g. if you choose to expose our Server ID in the path, then `:server-id` is your `id`.)*

---

## File Location

This document satisfies **Task 1** and should be referenced from the project root `README.md` (e.g. “Task 1 schema: [docs/DATABASE_SCHEMA.md](docs/DATABASE_SCHEMA.md)”).
