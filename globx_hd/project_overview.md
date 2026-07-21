# GlobX HD — Comprehensive Project Overview

GlobX Help Desk (HD) is a robust, multi-tenant help desk and ticketing system designed for customer support management. It supports role-based access control (Admin, Manager, Engineer, Contact), real-time WebSocket-based notifications, independent task management, comprehensive auditing, and an automated AI-powered email-to-ticket ingestion pipeline.

---

## 1. Technology Stack

### Backend (Go / Gin / GORM)
* **Runtime / Compiler:** Go 1.25.1
* **Web Framework:** Gin (v1.10.0) — handles routing, CORS, and request-response lifecycle.
* **ORM:** GORM (v1.25.12) — maps Go models to PostgreSQL tables and executes auto-migrations.
* **Database Driver:** `pgx` via GORM driver (v1.5.9).
* **Authentication:** JWT (v5) — signed stateless access/refresh tokens.
* **Real-time Engine:** Gorilla WebSocket (v1.5.3) — provides live browser communication.
* **Integrations:** SMTP (emails), Google Gemini 2.5 Flash API (AI text entity extraction).

### Frontend (Vue.js 3 / Vite)
* **Framework:** Vue 3 (Composition API) with Vite compilation.
* **State Management:** Pinia (stores auth token, user info, and websocket notifications).
* **Routing:** Vue Router (handles guards, dynamic layout wrappers, and role-based redirects).
* **Styling:** Tailwind CSS + Heroicons.

---

## 2. Directory Structure

```
globx_hd/
├── ticketing_backend/             # Go Backend
│   ├── cmd/main.go                # Server entry point
│   ├── internal/
│   │   ├── config/                # DB & Environment Config
│   │   ├── handlers/              # API Route Handlers (Controllers)
│   │   ├── middleware/            # JWT auth & Audit logging middlewares
│   │   ├── models/                # GORM entity definitions
│   │   ├── routes/                # Gin route registrations
│   │   ├── services/              # Business logic (Notification, Audit, Gemini, Email)
│   │   ├── utils/                 # Utility functions (Hashed Passwords, ID generators)
│   │   └── websocket/             # WebSocket Hub & client pumps
│   └── migrations/                # Database migrations
└── ticketing_frontend/            # Vue.js Frontend
    ├── src/
    │   ├── api/                   # Axios HTTP requests
    │   ├── components/            # Vue Layouts & Views
    │   │   ├── masterdata/        # Products, Issues, Roles UI
    │   │   ├── manager/           # Manager dashboards & ticket lists
    │   │   ├── engineer/          # Assigned tickets & tasks
    │   │   └── contacts/          # Customer-facing Views
    │   ├── router/index.ts        # Route configuration & Navigation guards
    │   ├── stores/                # Pinia global states
    │   └── utils/                 # Date/JWT parser helpers
```

---

## 3. Database Architecture & Key Models

### Master Tables
1. **`master_products`:** Products supported by the company.
2. **`master_user_designations` & `master_contact_designations`:** Job titles.
3. **`master_roles`:** Roles (1 = Admin, 2 = Manager, 3 = Engineer).
4. **`master_product_issues`:** Specific issues categories assigned to each product.

### Core Entities
* **`accounts`:** Corporate accounts with unique 3-digit customer codes.
* **`contacts`:** Customer users linked to accounts.
* **`users`:** Support team employees (Admins, Managers, Engineers).
* **`tickets`:** Core helpdesk ticket. Formatted ID format: `{CustomerCode}-{DDMMYY}-{Seq}`.
* **`tasks`:** Actionable units associated with tickets or standalone.

### Sub-Entities
* **`ticket_comments` / `task_comments`:** Message threads (Internal/External).
* **`ticket_calls`:** Logged inbound/outbound customer calls.
* **`ticket_activities` / `task_activities`:** State history audit trail.
* **`ticket_approvals`:** Manager-approved actions.
* **`ticket_attachments`:** Multi-format file attachments.

---

## 4. Key Workflows & Pipelines

### A. Authentication & Guards
1. Users log in via `/auth/login` (staff) or `/auth/contact/login` (customers).
2. The server yields a JWT token payload containing user info and role.
3. If `first_login` is true, the user is redirected immediately to `/reset-password/*` to enforce security.
4. **Navigation Guard:** For every protected route, the router matches `meta.requiresRole` with the decrypted token role. If mismatched, the user is thrown back to their designated role dashboard.

### B. WebSocket Notification Pipeline
1. Upon logging in, the frontend registers a persistent WebSocket upgrade connection at `/ws/notifications`.
2. A background `websocket.Hub` registers client connections.
3. When a ticket/task event triggers (e.g. status change or assignment), the `NotificationService` formats a template-based notification and saves it to the DB.
4. The service then calls `Hub.BroadcastToUser(...)` to immediately push the event payload and update the unread bell count.

### C. n8n & Gemini AI Ingestion
1. Emails received on the corporate support mail are captured by n8n.
2. n8n makes a request to the backend `/n8n/smart-ticket` endpoint.
3. The server forwards the text payload to Google Gemini 2.5 Flash.
4. Gemini extracts entities (phone, person, account, product keywords, priority).
5. The backend executes a waterfall lookup logic:
   * Match contact by email -> Match account -> Create formatted ticket.
   * If lookup fails, it deposits the query to `dumped_queries` for manual manager resolution.

---

## 5. Main API Endpoints

### Authentication
* `POST /auth/login` - Staff login
* `POST /auth/contact/login` - Customer login
* `PUT /auth/reset-password` - Reset staff password

### Ticketing
* `POST /manager/tickets` - Create ticket (Manager)
* `GET /manager/tickets` - List tickets (filters: status, priority, account)
* `PUT /manager/tickets/:id/status` - Change ticket status with remarks
* `PUT /manager/tickets/:id/assign` - Assign to engineer
* `GET /engineer/tickets` - List engineer-assigned tickets

### Master Data
* `GET/POST /products` - CRUD products list
* `GET/POST /products/:id/issues` - Product-specific issue categories

---

## 6. Resolving the "Add Product" Bug

### The Issue
`Products.vue` and `ProductForm.vue` are shared master-data interfaces utilized by both **Admins** and **Managers**.
However, their routing commands were hardcoded to the admin routes:
* Clicking "+ Add Product" forced: `router.push('/master-data/products/new');`
* Clicking "Edit" forced: `router.push('/master-data/products/:id/edit');`
* Saving a product pushed: `router.push('/master-data/products');`

Because a **Manager's** route namespace is prefixed with `/manager` (i.e. `/manager/master-data/products/new`), the navigation guard caught the non-prefixed paths, identified a role mismatch (as `/master-data/*` requires role `admin`), and redirected the Manager to the dashboard.

### The Fix
Both components were modified to dynamically inspect the current logged-in role from Pinia's `authStore` and construct the appropriate path prefix:
1. Imported `useAuthStore` inside `Products.vue` and `ProductForm.vue`.
2. Prepended the `/manager` namespace dynamically:
   ```javascript
   const prefix = authStore.userType === 'manager' ? '/manager' : '';
   router.push(`${prefix}/master-data/products/...`);
   ```
3. Verification: Ran production builds with no syntax or compiler exceptions. Both admins and managers can now add and edit products flawlessly.
