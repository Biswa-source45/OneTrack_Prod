# GlobX HD — Comprehensive System Documentation

> **Version:** 1.0  
> **Last Updated:** July 2025  
> **Application:** GlobX Help Desk (Ticketing System)

---

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Technology Stack](#2-technology-stack)
3. [Project Structure](#3-project-structure)
4. [Backend Architecture](#4-backend-architecture)
   - 4.1 [Entry Point & Configuration](#41-entry-point--configuration)
   - 4.2 [Database Schema & Models](#42-database-schema--models)
   - 4.3 [Authentication & Authorization](#43-authentication--authorization)
   - 4.4 [API Routes](#44-api-routes)
   - 4.5 [Handlers (Controllers)](#45-handlers-controllers)
   - 4.6 [Services (Business Logic)](#46-services-business-logic)
   - 4.7 [Middleware](#47-middleware)
   - 4.8 [Utilities](#48-utilities)
   - 4.9 [WebSocket System](#49-websocket-system)
5. [Frontend Architecture](#5-frontend-architecture)
   - 5.1 [Application Bootstrapping](#51-application-bootstrapping)
   - 5.2 [Router & Navigation Guards](#52-router--navigation-guards)
   - 5.3 [State Management (Pinia Stores)](#53-state-management-pinia-stores)
   - 5.4 [API Layer](#54-api-layer)
   - 5.5 [Component Hierarchy](#55-component-hierarchy)
   - 5.6 [Reusable UI Components](#56-reusable-ui-components)
   - 5.7 [Frontend Services](#57-frontend-services)
   - 5.8 [Utility Functions](#58-utility-functions)
6. [Database Schema](#6-database-schema)
   - 6.1 [Master Tables](#61-master-tables)
   - 6.2 [Core Entity Tables](#62-core-entity-tables)
   - 6.3 [Ticket Sub-Tables](#63-ticket-sub-tables)
   - 6.4 [Task Sub-Tables](#64-task-sub-tables)
   - 6.5 [Notification Tables](#65-notification-tables)
   - 6.6 [Audit Log Table](#66-audit-log-table)
   - 6.7 [Entity-Relationship Summary](#67-entity-relationship-summary)
7. [Data Flow & Key Workflows](#7-data-flow--key-workflows)
   - 7.1 [Authentication Flow](#71-authentication-flow)
   - 7.2 [Ticket Lifecycle](#72-ticket-lifecycle)
   - 7.3 [Task Lifecycle](#73-task-lifecycle)
   - 7.4 [Notification Flow](#74-notification-flow)
   - 7.5 [Approval Workflow](#75-approval-workflow)
   - 7.6 [Audit Logging Flow](#76-audit-logging-flow)
8. [n8n Integration & AI Pipeline](#8-n8n-integration--ai-pipeline)
   - 8.1 [Architecture Overview](#81-architecture-overview)
   - 8.2 [Webhook Endpoints](#82-webhook-endpoints)
   - 8.3 [Smart Resolver Engine](#83-smart-resolver-engine)
   - 8.4 [AI Agent Tool Endpoints](#84-ai-agent-tool-endpoints)
   - 8.5 [Gemini AI Service](#85-gemini-ai-service)
9. [Database Migrations](#9-database-migrations)
10. [Security](#10-security)
11. [Environment Configuration](#11-environment-configuration)

---

## 1. System Overview

**GlobX HD** is a full-stack help desk / ticketing system designed for customer support management. It supports multiple user roles, ticket lifecycle management, task assignment, real-time notifications, approval workflows, comprehensive audit logging, and AI-powered email-to-ticket automation via n8n workflows.

### Core Capabilities

| Capability | Description |
|---|---|
| **Multi-Role Access** | Admin, Manager, Engineer, Contact (customer) — each with distinct permissions and views |
| **Ticket Management** | Full CRUD, status tracking, priority, assignment, call logging, comments, approvals, attachments |
| **Task Management** | Independent tasks with assignment, status tracking, comments, and activity history |
| **Real-Time Notifications** | WebSocket-based live notifications with in-app bell, toast, and email alerts |
| **Audit Logging** | Comprehensive audit trail for all CRUD, authentication, and system events |
| **AI Email Integration** | n8n + Gemini AI pipeline that converts inbound emails into tickets automatically |
| **Master Data** | Configurable products, issues, roles, designations |
| **Dashboard** | Manager dashboard with ticket statistics and monthly filtering |

---

## 2. Technology Stack

### Backend
| Component | Technology | Version |
|---|---|---|
| Language | Go | 1.25.1 |
| Web Framework | Gin | v1.10.0 |
| ORM | GORM | v1.25.12 |
| Database Driver | gorm/driver/postgres | v1.5.9 |
| Database | PostgreSQL | — |
| Authentication | JWT (golang-jwt/jwt/v5) | v5.2.1 |
| WebSocket | gorilla/websocket | v1.5.3 |
| Password Hashing | golang.org/x/crypto (bcrypt) | — |
| UUID Generation | google/uuid | v1.6.0 |
| Env Loading | joho/godotenv | v1.5.1 |

### Frontend
| Component | Technology | Version |
|---|---|---|
| Framework | Vue.js 3 (Composition API) | ^3.5.18 |
| Build Tool | Vite | ^7.0.6 |
| Router | Vue Router | ^4.5.1 |
| State Management | Pinia | ^3.0.3 |
| HTTP Client | Axios | ^1.12.2 |
| CSS Framework | TailwindCSS | ^3.3.3 |
| Icons | Heroicons (Vue) | ^2.2.0 |
| UI Primitives | Headless UI (Vue) | ^1.7.23 |
| Language | TypeScript (type-check) / JavaScript (components) | ~5.8.0 |

### External Services
| Service | Purpose |
|---|---|
| n8n | Workflow automation (email → ticket pipeline) |
| Google Gemini 2.5 Flash | AI entity extraction from email content |
| SMTP | Email notifications |

---

## 3. Project Structure

```
globx_hd/
├── ticketing_backend/             # Go backend
│   ├── cmd/
│   │   └── main.go                # Application entry point
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go          # DB connection, env config
│   │   ├── handlers/              # HTTP request handlers (24 files)
│   │   │   ├── auth.go            # Auth middleware + login/logout/register/refresh
│   │   │   ├── tickets.go         # Ticket CRUD + manager operations
│   │   │   ├── engineer.go        # Engineer-specific handlers
│   │   │   ├── tasks.go           # Task CRUD
│   │   │   ├── n8n_webhook.go     # n8n webhook ticket creation
│   │   │   ├── n8n_smart_resolver.go  # AI-powered smart ticket resolver
│   │   │   ├── n8n_tools.go       # AI agent tool endpoints
│   │   │   ├── ticket_calls.go    # Call logging
│   │   │   ├── ticket_comments.go # Comments (internal + customer-visible)
│   │   │   ├── ticket_approvals.go# Approval workflow
│   │   │   ├── ticket_activities.go# Activity history
│   │   │   ├── attachments.go     # File upload/download
│   │   │   ├── notifications.go   # Notification endpoints
│   │   │   ├── audit_logs.go      # Audit log retrieval
│   │   │   ├── dashboard.go       # Manager dashboard stats
│   │   │   ├── accounts.go        # Account CRUD
│   │   │   ├── contacts.go        # Contact CRUD
│   │   │   ├── users.go           # User CRUD
│   │   │   ├── masters.go         # Master data CRUD
│   │   │   ├── product_issues.go  # Product issue CRUD
│   │   │   ├── dumped_query_handler.go # Unresolved email queries
│   │   │   ├── task_activities.go # Task activity history
│   │   │   └── task_comments.go   # Task comments
│   │   ├── middleware/
│   │   │   └── audit_middleware.go # Request context capture
│   │   ├── models/
│   │   │   ├── models.go          # All domain models
│   │   │   └── audit_log.go       # Audit log model + constants
│   │   ├── routes/
│   │   │   └── routes.go          # Route registration
│   │   ├── services/
│   │   │   ├── activity.go        # Ticket activity logging
│   │   │   ├── task_service.go    # Task activity logging
│   │   │   ├── notification_service.go  # Notification creation + WS broadcast
│   │   │   ├── audit_service.go   # Audit log creation + querying
│   │   │   ├── gemini_service.go  # Gemini AI API integration
│   │   │   └── email_notification_service.go  # SMTP email sending
│   │   ├── utils/
│   │   │   ├── utils.go           # Password hashing, ticket ID gen, customer codes
│   │   │   └── db_check.go        # Notification table diagnostics
│   │   └── websocket/
│   │       ├── hub.go             # WS hub (client registry, broadcast)
│   │       └── handler.go         # WS upgrade handler
│   ├── migrations/                # SQL migration files (13 files)
│   ├── go.mod / go.sum
│   └── .env.email                 # Email & n8n config
│
├── ticketing_frontend/            # Vue.js frontend
│   ├── src/
│   │   ├── api/                   # Axios API modules (11 files)
│   │   │   ├── api.js             # Centralized axios instance
│   │   │   ├── auth.js            # Auth API calls
│   │   │   ├── tickets.js         # Ticket API calls
│   │   │   ├── tasks.js           # Task API calls
│   │   │   ├── contacts.js        # Contact API calls
│   │   │   ├── engineer.js        # Engineer ticket API
│   │   │   ├── engineers.js       # Engineer listing API
│   │   │   ├── attachments.js     # File upload/download API
│   │   │   ├── notifications.js   # Notification API
│   │   │   ├── auditLogs.js       # Audit log API
│   │   │   └── dumpedQueries.js   # Dumped query API
│   │   ├── components/
│   │   │   ├── Layout.vue         # Main app layout (Header + Sidebar + Footer)
│   │   │   ├── Header.vue         # Top navigation bar
│   │   │   ├── Sidebar.vue        # Collapsible sidebar navigation
│   │   │   ├── Footer.vue         # Footer
│   │   │   ├── LoginUser.vue      # Staff login
│   │   │   ├── LoginContact.vue   # Customer login
│   │   │   ├── ResetPasswordUser.vue   # Staff password reset
│   │   │   ├── ResetPasswordContact.vue # Customer password reset
│   │   │   ├── Accounts.vue       # Account management
│   │   │   ├── Contacts.vue       # Contact management
│   │   │   ├── Users.vue          # User management
│   │   │   ├── Dashboard.vue      # Admin dashboard
│   │   │   ├── MasterData.vue     # Master data landing
│   │   │   ├── contacts/          # Customer-facing views (7 items)
│   │   │   ├── manager/           # Manager views (9 items)
│   │   │   ├── engineer/          # Engineer views (5 items)
│   │   │   ├── masterdata/        # Master data CRUD views (6 items)
│   │   │   ├── shared/            # Shared components (25 items)
│   │   │   │   ├── TicketDetailPage.vue  # Manager ticket detail
│   │   │   │   ├── TicketForm.vue        # Reusable ticket create/edit form
│   │   │   │   ├── StatusChangeModal.vue # Status change with remarks
│   │   │   │   ├── FileUpload.vue        # Drag-and-drop file upload
│   │   │   │   ├── NotificationBell.vue  # Header notification bell
│   │   │   │   ├── NotificationsPage.vue # Full notification page
│   │   │   │   ├── tabs/                 # Tab components (8 items)
│   │   │   │   │   ├── ConversationTab.vue
│   │   │   │   │   ├── CallsTab.vue
│   │   │   │   │   ├── ApprovalsTab.vue
│   │   │   │   │   ├── ManagerAttachmentsTab.vue
│   │   │   │   │   ├── EngineerAttachmentsTab.vue
│   │   │   │   │   ├── HistoryTab.vue
│   │   │   │   │   └── TaskCommentsTab.vue
│   │   │   │   └── ...
│   │   │   └── ui/                # Reusable UI primitives (11 items)
│   │   │       ├── Button.vue
│   │   │       ├── Modal.vue
│   │   │       ├── DataTable.vue
│   │   │       ├── FormField.vue
│   │   │       ├── FormLayout.vue
│   │   │       ├── FuzzySearchDropdown.vue
│   │   │       ├── Pagination.vue
│   │   │       ├── ConfirmDialog.vue
│   │   │       └── ...
│   │   ├── router/
│   │   │   └── index.ts           # Route definitions + navigation guards
│   │   ├── stores/
│   │   │   ├── auth.js            # Authentication state (Pinia)
│   │   │   └── notifications.js   # Notification state + WS integration
│   │   ├── services/
│   │   │   └── websocket.js       # WebSocket client service
│   │   └── utils/
│   │       ├── jwt.js             # JWT decode helper
│   │       ├── user.js            # Name formatting utilities
│   │       └── date.js            # Date formatting
│   ├── package.json
│   ├── tailwind.config.js
│   └── vite.config.ts
│
└── SYSTEM_DOCUMENTATION.md        # This file
```

---

## 4. Backend Architecture

### 4.1 Entry Point & Configuration

**File:** `cmd/main.go`

The application boots in this sequence:

1. **Load environment variables** from `.env.email` via `godotenv`
2. **Connect to PostgreSQL** via `config.ConnectDatabase()` (GORM auto-migrate all models)
3. **Initialize WebSocket Hub** — `websocket.NewHub()` runs in a background goroutine
4. **Setup Gin router** — `routes.SetupRouter(db, hub)` registers all routes
5. **Start HTTP server** on `:8080` (configurable via `SERVER_ADDRESS` env var)

**File:** `internal/config/config.go`

- Reads `DATABASE_URL` (fallback: `host=localhost user=postgres password=postgres dbname=globx_hd port=5432 sslmode=disable TimeZone=Asia/Kolkata`)
- Reads `SERVER_ADDRESS` (fallback: `:8080`)
- Exposes `config.DB` as a package-level database handle
- Auto-migrates all models on startup

---

### 4.2 Database Schema & Models

**File:** `internal/models/models.go` — All domain models  
**File:** `internal/models/audit_log.go` — Audit log model + constants

#### Master Tables
| Model | Table | Purpose |
|---|---|---|
| `MasterProduct` | `master_products` | Products the company supports |
| `MasterUserDesignation` | `master_user_designations` | Staff job titles |
| `MasterContactDesignation` | `master_contact_designations` | Customer contact titles |
| `MasterRole` | `master_roles` | System roles (Admin=1, Manager=2, Engineer=3) |
| `MasterProductIssue` | `master_product_issues` | Issue categories per product |

#### Core Entities
| Model | Table | Key Fields |
|---|---|---|
| `Account` | `accounts` | account_name, customer_code (unique 3-digit), account_type (Govt./Private), account_owner |
| `Contact` | `contacts` | first_name, last_name, email, mobile, account_id (nullable for Individual), customer_code, contact_type |
| `User` | `users` | first_name, last_name, email, employee_id, role_id, designation_id, password, first_login |
| `Ticket` | `tickets` | ticket_id (formatted), account_id, contact_id, product_id, subject, ticket_details, ticket_status, priority, assigned_engineer, channel |
| `Task` | `tasks` | subject, description, status, priority, assigned_to, created_by, ticket_id (optional link) |

#### Ticket Sub-Entities
| Model | Table | Purpose |
|---|---|---|
| `TicketAttachment` | `ticket_attachments` | File uploads per ticket |
| `TicketComment` | `ticket_comments` | Internal + customer-visible comments |
| `TicketCall` | `ticket_calls` | Call logs (inbound/outbound, scheduled/completed) |
| `TicketActivity` | `ticket_activities` | Activity timeline (status changes, assignments, etc.) |
| `TicketApproval` | `ticket_approvals` | Approval requests (PENDING/APPROVED/REJECTED) |

#### Other
| Model | Table | Purpose |
|---|---|---|
| `TaskComment` | `task_comments` | Comments on tasks |
| `TaskActivity` | `task_activities` | Task activity timeline |
| `Notification` | `notifications` | In-app notifications |
| `NotificationTemplate` | `notification_templates` | Templates for notification generation |
| `DumpedQuery` | `dumped_queries` | Unresolved email queries from n8n |
| `AuditLog` | `audit_logs` | System-wide audit trail |

---

### 4.3 Authentication & Authorization

**File:** `internal/handlers/auth.go`

#### JWT Token System
- **Signing method:** HMAC (HS256) with `JWT_SECRET` env var (fallback: hardcoded secret)
- **Access token validity:** 24 hours
- **Refresh token validity:** 7 days
- **Token claims:** `sub` (user/contact ID), `type` ("user" or "contact"), `role` (role name), `exp` (expiry)

#### Auth Endpoints
| Endpoint | Handler | Purpose |
|---|---|---|
| `POST /auth/login` | `UserLogin` | Staff login (returns access_token, refresh_token, first_login, user) |
| `POST /auth/contact/login` | `ContactLogin` | Customer login |
| `POST /auth/refresh` | `RefreshToken` | Refresh access token |
| `POST /auth/logout` | `Logout` | Stateless logout (audit log only) |
| `PUT /auth/reset-password` | `ResetUserPassword` | Staff password reset (sets first_login=false) |
| `PUT /auth/contact/reset-password` | `ResetContactPassword` | Customer password reset |

#### Auth Middleware — `AuthMiddleware(db)`
Applied to all protected routes. Extracts Bearer token, validates JWT, loads user/contact from DB, and sets context:
- **Staff users:** `c.Set("user", user)` — full `models.User` struct
- **Customer contacts:** `c.Set("contact", contact)` and `c.Set("contact_id", id)`

#### Role Check Helpers
- `IsManager(c)` — checks if user's role_id == 2
- `IsEngineer(c)` — checks if user's role_id == 3
- Used within handlers for fine-grained access control

#### First-Login Password Reset Enforcement
- New users are created with `first_login = true`
- Login response includes `first_login` field
- Frontend router guard forces redirect to `/reset-password/user` if `firstLogin == true`
- After successful reset, `first_login` is set to `false` and auth is cleared

---

### 4.4 API Routes

**File:** `internal/routes/routes.go`

All routes use the **Gin** framework. The setup applies CORS middleware and the global `AuditMiddleware`.

#### Public Routes (No Auth)
```
POST   /auth/login
POST   /auth/contact/login
POST   /auth/register             (staff registration)
POST   /auth/contact/register     (customer registration)
```

#### Protected Routes (AuthMiddleware)

**Authentication:**
```
POST   /auth/refresh
POST   /auth/logout
PUT    /auth/reset-password
PUT    /auth/contact/reset-password
```

**Admin/Manager — CRUD (both admin & manager access):**
```
GET/POST         /accounts
GET/PUT/DELETE   /accounts/:id
GET/POST         /contacts
GET/PUT/DELETE   /contacts/:id
GET/POST         /users
GET/PUT/DELETE   /users/:id
GET/POST         /products, /products/:id, /products/:id/issues
GET/POST/PUT/DEL /designations, /roles, /contact-designations, /product-issues
```

**Manager-Specific:**
```
GET    /manager/dashboard/stats
POST   /manager/tickets                    (create ticket)
GET    /manager/tickets                    (list with filters)
GET    /manager/tickets/:id                (detail with counts)
PUT    /manager/tickets/:id                (update)
PUT    /manager/tickets/:id/status         (change status)
PUT    /manager/tickets/:id/assign         (assign engineer)
GET    /manager/audit-logs                 (audit log list)
GET    /manager/audit-logs/:id             (audit log detail)
GET    /manager/audit-logs/stats           (audit statistics)
POST   /manager/tasks                      (create task)
GET/PUT/DEL  /manager/tasks/:id
PUT    /manager/tasks/:id/status
GET    /manager/engineers                  (list engineers)
```

**Engineer-Specific:**
```
GET    /engineer/tickets                   (assigned tickets)
GET    /engineer/tickets/:id
PUT    /engineer/tickets/:id/status
POST   /engineer/tickets                   (create ticket)
GET    /engineer/tasks
GET    /engineer/tasks/:id
PUT    /engineer/tasks/:id/status
```

**Ticket Sub-Resources (shared):**
```
GET/POST       /tickets/:id/comments
GET/POST       /tickets/:id/calls
PUT            /tickets/:id/calls/:callId
PUT            /tickets/:id/calls/:callId/complete
PUT            /tickets/:id/calls/:callId/cancel
GET            /tickets/:id/activities
POST           /tickets/:id/attachments
GET            /tickets/:id/attachments/:attachmentId/download
DELETE         /tickets/:id/attachments/:attachmentId
POST           /tickets/:id/approvals
GET            /tickets/:id/approvals
PATCH          /tickets/:id/approvals/:approvalId/approve
PATCH          /tickets/:id/approvals/:approvalId/reject
```

**Task Sub-Resources:**
```
GET/POST       /tasks/:id/comments
GET            /tasks/:id/activities
```

**Notifications:**
```
GET    /notifications
GET    /notifications/unread-count
PUT    /notifications/:id/read
PUT    /notifications/read-all
DELETE /notifications/:id
```

**WebSocket:**
```
GET    /ws/notifications                   (WS upgrade)
```

**Contact (Customer) Routes:**
```
POST   /customer/tickets                   (raise ticket)
GET    /customer/my-tickets
GET    /customer/my-tickets/:id
POST   /customer/my-tickets/:id/comments
```

**n8n Webhook Routes (API Key Auth):**
```
GET    /n8n/health
POST   /n8n/ticket                         (basic ticket creation)
POST   /n8n/smart-ticket                   (AI-powered smart resolver)
GET    /n8n/lookup/accounts
GET    /n8n/lookup/contacts
GET    /n8n/lookup/products
```

**AI Agent Tool Endpoints (API Key Auth):**
```
GET    /n8n/tools/search-contact-by-email
GET    /n8n/tools/search-contact-by-phone
GET    /n8n/tools/search-contact-by-name
GET    /n8n/tools/search-account-by-name
GET    /n8n/tools/search-product-by-name
GET    /n8n/tools/list-products
POST   /n8n/tools/create-ticket
POST   /n8n/tools/extract-email-with-gemini
POST   /n8n/tools/dump-unresolved-query
```

---

### 4.5 Handlers (Controllers)

Each handler file follows the pattern: parse request → validate → database operation → activity/audit logging → notification → response.

| Handler File | Key Functions | Purpose |
|---|---|---|
| `auth.go` | `AuthMiddleware`, `UserLogin`, `ContactLogin`, `RefreshToken`, `Logout`, `ResetUserPassword` | Authentication lifecycle |
| `tickets.go` | `ManagerCreateTicketHandler`, `GetTicketsHandler`, `GetTicketDetailHandler`, `UpdateTicketHandler`, `ManagerChangeStatusHandler`, `AssignEngineerHandler` | Full ticket CRUD for managers |
| `engineer.go` | `EngineerChangeStatusHandler`, `EngineerChangeTaskStatusHandler`, `EngineerCreateTicketHandler`, `GetEngineerTicketsHandler`, `GetEngineerTasksHandler` | Engineer operations |
| `ticket_comments.go` | `CreateTicketCommentHandler`, `GetTicketCommentsHandler` | Internal/customer-visible comments |
| `ticket_calls.go` | `CreateTicketCallHandler`, `UpdateTicketCallHandler`, `CompleteTicketCallHandler`, `CancelTicketCallHandler` | Call log management |
| `ticket_approvals.go` | `CreateApprovalHandler`, `ListApprovalsHandler`, `ApproveHandler`, `RejectHandler` | Approval workflow |
| `ticket_activities.go` | `GetTicketActivitiesHandler` | Activity history retrieval |
| `attachments.go` | `UploadTicketAttachments`, `DownloadAttachment`, `DeleteAttachment` | File management (dual auth: user + contact) |
| `tasks.go` | `CreateTaskHandler`, `GetTasksHandler`, `UpdateTaskHandler`, `DeleteTaskHandler` | Task CRUD |
| `task_comments.go` | `CreateTaskCommentHandler`, `GetTaskCommentsHandler` | Task comment management |
| `notifications.go` | `GetNotificationsHandler`, `MarkAsReadHandler`, `MarkAllAsReadHandler`, `DeleteNotificationHandler`, `GetUnreadCountHandler` | Notification management |
| `audit_logs.go` | `GetAuditLogsHandler`, `GetAuditLogStatsHandler`, `GetAuditLogDetailHandler` | Audit log retrieval with filtering |
| `dashboard.go` | `ManagerDashboardStatsHandler` | Ticket statistics by status, with monthly filtering |
| `accounts.go` | `CreateAccountHandler`, `GetAccountsHandler`, `UpdateAccountHandler`, `DeleteAccountHandler` | Account CRUD |
| `contacts.go` | `CreateContactHandler`, `GetContactsHandler`, `UpdateContactHandler`, `DeleteContactHandler` | Contact CRUD with customer code generation |
| `users.go` | `CreateUserHandler`, `GetUsersHandler`, `UpdateUserHandler`, `DeleteUserHandler` | User CRUD (password hashing on create) |
| `masters.go` | CRUD for products, designations, roles, contact designations | Master data management |
| `product_issues.go` | `CreateProductIssue`, `UpdateProductIssue`, `DeleteProductIssue` | Product issue CRUD with audit logging |
| `dumped_query_handler.go` | `ListDumpedQueries`, `GetDumpedQuery`, `DeleteDumpedQuery`, `UpdateDumpedQueryStatus` | Unresolved email query management |
| `n8n_webhook.go` | `N8nCreateTicketHandler`, `N8nLookupAccountsHandler`, `N8nLookupContactsHandler`, `N8nLookupProductsHandler`, `N8nHealthCheckHandler` | n8n webhook integration |
| `n8n_smart_resolver.go` | `SmartTicketHandler` → `SmartResolver.ResolveAndCreateTicket` | AI-powered email → ticket with waterfall entity resolution |
| `n8n_tools.go` | `ToolSearchContactByEmail`, `ToolSearchContactByPhone`, `ToolSearchContactByName`, `ToolSearchAccountByName`, `ToolSearchProductByName`, `ToolCreateTicket`, `ToolExtractEmailWithGemini`, `ToolDumpUnresolvedQuery` | AI agent callable endpoints |

---

### 4.6 Services (Business Logic)

| Service | File | Purpose |
|---|---|---|
| `ActivityService` | `activity.go` | Logs ticket activities (creation, status change, assignment, comments, calls, approvals). Supports remarks for status changes. |
| `TaskService` | `task_service.go` | Logs task activities (creation, status change, assignee change, comments). |
| `NotificationService` | `notification_service.go` | Creates notifications using templates, replaces variables (`{{ticket_id}}`, `{{user_name}}`), broadcasts via WebSocket, sends email for new tickets. |
| `AuditService` | `audit_service.go` | Logs detailed audit entries (actor, action, entity, old/new values, HTTP context). Provides filtered retrieval with pagination. |
| `GeminiService` | `gemini_service.go` | Calls Google Gemini 2.5 Flash API to extract structured data (phone numbers, person names, org names, product hints, priority hints) from email content. |
| `EmailNotificationService` | `email_notification_service.go` | Sends HTML email notifications via SMTP for ticket creation events. |

#### NotificationService Key Design

The notification service uses a **template-based approach**:
1. Templates stored in `notification_templates` table with `notification_type`, `title_template`, `message_template`, `default_priority`
2. Variables like `{{ticket_id}}`, `{{user_name}}`, `{{status}}` are replaced at runtime
3. After DB insert, notifications are **broadcast via WebSocket** in real-time to the target user
4. Email notifications are sent asynchronously for ticket creation events

---

### 4.7 Middleware

| Middleware | File | Scope | Purpose |
|---|---|---|---|
| `AuditMiddleware` | `audit_middleware.go` | Global (all routes) | Generates unique request_id (UUID), captures client IP and User-Agent, stores in Gin context |
| `AuthMiddleware` | `auth.go` | All protected routes | JWT validation, user/contact loading into context |
| `N8nAPIKeyMiddleware` | `n8n_webhook.go` | `/n8n/*` routes | Validates API key via `X-N8N-API-Key` header or `api_key` query param |
| CORS | `routes.go` | Global | Allows origins `http://localhost:5173` and `http://127.0.0.1:5173` |

---

### 4.8 Utilities

**File:** `internal/utils/utils.go`

| Function | Purpose |
|---|---|
| `HashPassword(password)` | bcrypt hash generation |
| `CheckPasswordHash(password, hash)` | bcrypt comparison |
| `GenerateUnique3Digit(db)` | Generates unique 3-digit customer code for accounts |
| `GenerateUniqueCustomerCode(db)` | Same but checks uniqueness across BOTH accounts and contacts tables |
| `GetNextTicketSequence(db, accountID, date)` | Annual sequential numbering per account (handles NULL account_id for individuals) |
| `FormatDateForTicketID(date)` | Formats date as DDMMYY |
| `FormatTicketID(customerCode, dateStr, seq)` | Produces ticket ID: `{code}-{DDMMYY}-{0001}` |
| `GetUserName(db, userID)` | Returns full name by user ID |

**File:** `internal/utils/db_check.go`

| Function | Purpose |
|---|---|
| `CheckNotificationTables(db)` | Diagnostic tool to verify notification system tables and templates exist |

---

### 4.9 WebSocket System

**Files:** `internal/websocket/hub.go`, `internal/websocket/handler.go`

#### Architecture
- **Hub** — Singleton that maintains a map of connected clients keyed by `{userType}:{userID}`
- **Client** — Represents a single WebSocket connection with send/receive channels
- **Registration** — Happens on WS upgrade at `/ws/notifications` (requires auth middleware)

#### Connection Flow
1. Client sends HTTP upgrade request with JWT auth
2. `WebSocketHandler` validates auth, creates `Client`, registers with Hub
3. Sends welcome message + current unread notification count
4. Starts `ReadPump` (reads client messages, handles ping/pong) and `WritePump` (sends hub messages, manages keep-alive pings)

#### Message Types
| Type | Direction | Purpose |
|---|---|---|
| `connected` | Server → Client | Connection confirmation with client_id |
| `count.update` | Server → Client | Updated unread notification count |
| `notification.new` | Server → Client | New notification data |
| `notification.read` | Server → Client | Notification marked as read |
| `notification.all_read` | Server → Client | All notifications marked as read |
| `notification.deleted` | Server → Client | Notification deleted |
| `ping` / `pong` | Bidirectional | Keep-alive |

#### Broadcasting
`Hub.BroadcastToUser(userID, userType, messageType, data)` sends to all connections of a specific user. Used by `NotificationService` after creating notifications.

---

## 5. Frontend Architecture

### 5.1 Application Bootstrapping

The Vue 3 app is bootstrapped with:
- **Vite** as build tool
- **Pinia** for state management
- **Vue Router** for navigation
- **TailwindCSS** for styling
- **Heroicons** for iconography
- **Headless UI** for accessible UI primitives

### 5.2 Router & Navigation Guards

**File:** `src/router/index.ts`

#### Route Structure

| Route Prefix | Role Required | View |
|---|---|---|
| `/login/user` | Public | Staff login |
| `/login/contact` | Public | Customer login |
| `/reset-password/*` | Public | Password reset |
| `/dashboard` | admin | Admin dashboard |
| `/accounts`, `/contacts`, `/users`, `/master-data/*` | admin | Admin CRUD views |
| `/contacts/raise-ticket` | contact | Customer ticket creation |
| `/contacts/my-tickets` | contact | Customer ticket list |
| `/contacts/my-tickets/:id` | contact | Customer ticket detail |
| `/manager/dashboard` | manager | Manager dashboard |
| `/manager/tickets`, `/manager/tickets/:id` | manager | Ticket management |
| `/manager/tasks`, `/manager/tasks/:id` | manager | Task management |
| `/manager/engineers` | manager | Engineer listing |
| `/manager/dumped-queries` | manager | Unresolved queries |
| `/manager/audit-logs` | manager | Audit log viewer |
| `/manager/accounts`, `/contacts`, `/users`, `/master-data/*` | manager | Admin features for managers |
| `/engineer/tickets`, `/engineer/tickets/:id` | engineer | Assigned tickets |
| `/engineer/tasks`, `/engineer/tasks/:id` | engineer | Assigned tasks |
| `/notifications` | any authenticated | Notification page |

#### Navigation Guard Logic
```
1. If public page → allow (redirect authenticated users to their landing page)
2. If not authenticated → redirect to /login/user
3. If firstLogin=true AND not on reset-password → force redirect to /reset-password/user
4. If route requires a role AND user doesn't have it → redirect to correct landing
5. Allow navigation
```

#### Landing Pages by Role
| Role | Landing |
|---|---|
| admin | `/dashboard` |
| manager | `/manager/dashboard` |
| engineer | `/engineer/tickets` |
| contact | `/contacts/my-tickets` |

---

### 5.3 State Management (Pinia Stores)

#### Auth Store (`stores/auth.js`)

| State | Purpose |
|---|---|
| `token` | JWT access token |
| `userType` | Role (decoded from JWT) |
| `firstLogin` | First-login flag for password reset enforcement |
| `user` | Full user object |

| Action | Purpose |
|---|---|
| `setAuth(token, userType, firstLogin, user)` | Save auth state + localStorage |
| `clearAuth()` | Clear all auth state + localStorage |
| `loadAuth()` | Restore from localStorage (decodes JWT for role) |

#### Notification Store (`stores/notifications.js`)

| State | Purpose |
|---|---|
| `notifications` | Array of notification objects |
| `unreadCount` | Badge count |
| `filters` | Category, priority, isRead filtering |
| `pagination` | currentPage, limit, hasMore |
| `wsConnected` | WebSocket connection status |

| Action | Purpose |
|---|---|
| `fetchNotifications(reset?)` | Load notifications from API with filtering |
| `fetchUnreadCount()` | Get unread badge count |
| `markAsRead(id)` | Mark single notification read |
| `markAllAsRead()` | Mark all as read |
| `deleteNotification(id)` | Remove notification |
| `initializeWebSocket()` | Connect WS, subscribe to events |
| `cleanupWebSocket()` | Disconnect WS, unsubscribe |
| `initialize()` | Fetch data + start WS |

**WebSocket Integration:** The store subscribes to `notification.new`, `count.update`, `notification.read`, `notification.all_read`, and `notification.deleted` events, updating local state in real-time.

---

### 5.4 API Layer

**File:** `src/api/api.js` — Centralized axios instance

```javascript
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || 'http://localhost:8080',
  headers: { 'Content-Type': 'application/json' },
  timeout: 10000,
})
// Request interceptor adds Authorization: Bearer <token>
```

All API modules import this centralized instance for consistent headers, auth, and base URL.

| API Module | Key Functions |
|---|---|
| `auth.js` | `loginUser`, `loginContact`, `registerUser`, `refreshToken`, `resetPassword`, `logout` |
| `tickets.js` | `getTickets`, `getTicket`, `createTicket`, `updateTicket`, `changeTicketStatus`, `assignEngineer`, `getTicketComments`, `createComment`, plus manager and customer-specific variants |
| `tasks.js` | `getTasks`, `getTask`, `createTask`, `updateTask`, `deleteTask`, `changeTaskStatus`, `getTaskComments`, `createTaskComment`, `getTaskActivities` |
| `engineer.js` | `getEngineerTickets`, `getEngineerTicket`, `changeEngineerTicketStatus`, `getEngineerTasks`, `changeEngineerTaskStatus` |
| `engineers.js` | `getEngineers` |
| `contacts.js` | `getContacts` |
| `attachments.js` | `uploadAttachments`, `downloadAttachment`, `deleteAttachment` |
| `notifications.js` | `getNotifications`, `getUnreadCount`, `markAsRead`, `markAllAsRead`, `deleteNotification` |
| `auditLogs.js` | `getAuditLogs`, `getAuditLog`, `getStats`, `getRecent`, `getCritical`, `getFailed`, `getByEntity`, `getByActor` |
| `dumpedQueries.js` | `getDumpedQueries`, `getDumpedQuery`, `deleteDumpedQuery`, `updateDumpedQueryStatus` |

---

### 5.5 Component Hierarchy

#### Layout Structure
```
App.vue
├── LoginUser.vue / LoginContact.vue          (public)
├── ResetPasswordUser.vue / ResetPasswordContact.vue  (public)
└── Layout.vue                                (authenticated)
    ├── Header.vue                            (top bar with NotificationBell)
    ├── Sidebar.vue                           (collapsible, role-based menu)
    ├── <router-view />                       (content area, flex-1)
    └── Footer.vue
```

#### Manager Views
| Component | Description |
|---|---|
| `ManagerDashboard.vue` | Stats cards + monthly filter |
| `Tickets.vue` | Ticket list with search, filters, pagination |
| `TicketDetailPage.vue` (shared) | Tabbed detail: Conversation, Calls, Approvals, Attachments, History |
| `RaiseTicket.vue` → `TicketForm.vue` (shared) | Ticket creation form with fuzzy search |
| `EditTicketPage.vue` (shared) | Ticket property editing |
| `Tasks.vue` | Task list with filters |
| `CreateTask.vue` | Task creation form |
| `TaskDetailPage.vue` (shared) | Task detail with comments + history tabs |
| `Engineers.vue` | Engineer listing |
| `DumpedQueries.vue` | Unresolved query list |
| `DumpedQueryDetail.vue` | Query detail with resolve/ignore actions |
| `AuditLogs.vue` | Audit log table with filters, stats, detail modal |

#### Engineer Views
| Component | Description |
|---|---|
| `EngineerTickets.vue` | List of assigned tickets |
| `EngineerTicketDetailPage.vue` | Ticket detail (same tabs as manager, engineer-specific actions) |
| `EngineerTasks.vue` | Assigned task list |
| `EngineerTaskDetailPage.vue` | Task detail |
| `RaiseTicket.vue` | Ticket creation (engineer) |

#### Contact (Customer) Views
| Component | Description |
|---|---|
| `MyTickets.vue` | Customer's ticket list |
| `CustomerTicketDetailPage.vue` | Customer ticket detail (conversation + attachments) |
| `RaiseTicket.vue` | Customer ticket creation |

#### Shared Tab Components
| Tab | Used In | Purpose |
|---|---|---|
| `ConversationTab.vue` | Ticket detail | Comments thread (internal + customer-visible) |
| `CallsTab.vue` | Ticket detail | Call logs with add/edit/complete/cancel |
| `ApprovalsTab.vue` | Ticket detail | Approval request + approve/reject workflow |
| `ManagerAttachmentsTab.vue` | Manager ticket detail | File upload/download/preview |
| `EngineerAttachmentsTab.vue` | Engineer ticket detail | Same as manager attachments |
| `HistoryTab.vue` | Ticket detail | Activity timeline with old/new values + remarks |
| `TaskCommentsTab.vue` | Task detail | Task comment thread + activity history |

---

### 5.6 Reusable UI Components

| Component | Props | Purpose |
|---|---|---|
| `Button.vue` | variant, size, loading | Themed button with loading state |
| `Modal.vue` | open, title | Generic modal wrapper |
| `DataTable.vue` | columns, data | Generic sortable table |
| `FormField.vue` | label | Form field wrapper with label |
| `FormLayout.vue` | — | Form grid layout |
| `FuzzySearchDropdown.vue` | items, searchKeys, displayKey, allowCreate | Fuzzy search with keyboard nav + "Create New" |
| `Pagination.vue` | currentPage, totalPages | Page navigation controls |
| `ConfirmDialog.vue` | open, message | Confirmation dialog |
| `PageHeader.vue` | title, subtitle | Page header |
| `DropdownMenu.vue` | — | Dropdown menu wrapper |
| `NotificationToast.vue` | — | Toast notification display |

#### Specialized Shared Components
| Component | Purpose |
|---|---|
| `StatusChangeModal.vue` | Status change with required remarks, reused for close/reopen + approvals |
| `FileUpload.vue` | Drag-and-drop upload with progress, file type validation, 3MB limit |
| `NotificationBell.vue` | Header bell icon with unread count badge and dropdown |
| `NotificationsPage.vue` | Full notification management page with filters |
| `ContactCreateModal.vue` | Inline contact creation from ticket form |
| `ProductCreateModal.vue` | Inline product creation from ticket form |
| `AccountCreateModal.vue` | Inline account creation |
| `InlineEditDropdown.vue` | Inline editable dropdown field |
| `PropertySection.vue` | Ticket property display section |

---

### 5.7 Frontend Services

#### WebSocket Service (`services/websocket.js`)

A singleton class that manages the WebSocket connection lifecycle:

| Feature | Implementation |
|---|---|
| **Auto-reconnect** | Exponential backoff: 1s → 2s → 4s → ... → 30s max, 10 attempts max |
| **Keep-alive** | Ping every 30 seconds |
| **Message queue** | Queues messages when disconnected, sends on reconnect |
| **Event system** | `on(type, handler)` returns unsubscribe function, supports wildcard `*` |
| **Connection states** | `connecting`, `connected`, `disconnected`, `reconnecting`, `error`, `failed` |

---

### 5.8 Utility Functions

| File | Functions | Purpose |
|---|---|---|
| `utils/jwt.js` | `decodeJWT(token)` | Base64 decode JWT payload (no verification) |
| `utils/user.js` | `formatUserName(user)`, `getUserInitials(user)` | Smart name formatting to handle database inconsistencies (e.g., full name in first_name field) |
| `utils/date.js` | Date formatting helpers | Display-friendly date formatting |
| `utils/debug-jwt.js` | Debug JWT decoding | Development helper |

---

## 6. Database Schema

### 6.1 Master Tables

#### `master_products`
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK, auto-increment |
| product_name | varchar(255) | NOT NULL |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

#### `master_roles`
| Column | Type | Notes |
|---|---|---|
| id | uint | PK (1=Admin, 2=Manager, 3=Engineer) |
| role_name | varchar(50) | NOT NULL |

#### `master_user_designations`
| Column | Type |
|---|---|
| id | uint | PK |
| designation_name | varchar(100) |

#### `master_contact_designations`
| Column | Type |
|---|---|
| id | uint | PK |
| designation_name | varchar(100) |

#### `master_product_issues`
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK |
| product_id | uint | FK → master_products |
| issue_name | varchar(255) | NOT NULL |

---

### 6.2 Core Entity Tables

#### `accounts`
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK |
| account_name | varchar(255) | NOT NULL |
| customer_code | varchar(10) | UNIQUE, auto-generated 3-digit |
| account_type | varchar(50) | Govt. / Private |
| account_owner | varchar(255) | |
| address, city, state, pincode | varchar | Location fields |
| created_at, updated_at | timestamp | auto |

#### `contacts`
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK |
| first_name, last_name | varchar(255) | NOT NULL |
| email | varchar(255) | UNIQUE, NOT NULL |
| mobile | varchar(20) | |
| account_id | uint | FK → accounts (NULL for Individual) |
| designation_id | uint | FK → master_contact_designations |
| customer_code | varchar(10) | Auto-generated (unique across accounts + contacts) |
| contact_type | varchar(50) | Govt. / Private / Individual |
| password | varchar(255) | bcrypt hash |
| first_login | boolean | Default true |

#### `users`
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK |
| first_name, last_name | varchar(255) | NOT NULL |
| email | varchar(255) | UNIQUE, NOT NULL |
| employee_id | varchar(50) | UNIQUE |
| mobile | varchar(20) | |
| role_id | uint | FK → master_roles |
| designation_id | uint | FK → master_user_designations |
| password | varchar(255) | bcrypt hash |
| first_login | boolean | Default true |

#### `tickets`
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK |
| ticket_id | varchar(50) | UNIQUE (format: `{code}-{DDMMYY}-{0001}`) |
| account_id | uint | FK → accounts (NULL for Individual) |
| contact_id | uint | FK → contacts, NOT NULL |
| product_id | uint | FK → master_products |
| subject | text | NOT NULL |
| ticket_details | text | NOT NULL |
| ticket_status | varchar(20) | OPEN, IN PROGRESS, RESOLVED, CLOSED |
| priority | varchar(10) | High, Medium, Low |
| assigned_engineer | uint | FK → users (nullable) |
| channel | varchar(50) | Web, Mail |
| created_at, updated_at | timestamp | auto |

#### `tasks`
| Column | Type | Constraints |
|---|---|---|
| id | uint | PK |
| subject | varchar(255) | NOT NULL |
| description | text | |
| status | varchar(20) | Open, In Progress, Completed, Closed |
| priority | varchar(10) | High, Medium, Low |
| assigned_to | uint | FK → users |
| created_by | uint | FK → users |
| ticket_id | uint | FK → tickets (nullable — task can be independent) |
| due_date | timestamp | |

---

### 6.3 Ticket Sub-Tables

#### `ticket_attachments`
| Column | Type |
|---|---|
| id | uint |
| ticket_id | varchar(50) | FK → tickets.ticket_id |
| original_filename | varchar(255) |
| stored_filename | varchar(255) |
| file_path | text |
| file_size | int |
| mime_type | varchar(100) |
| uploaded_by | uint |

#### `ticket_comments`
| Column | Type |
|---|---|
| id | uint |
| ticket_id | uint | FK → tickets |
| user_id | uint | FK → users (nullable for contacts) |
| contact_id | uint | FK → contacts (nullable for users) |
| comment_type | varchar(20) | internal, customer |
| content | text |

#### `ticket_calls`
| Column | Type |
|---|---|
| id | uint |
| ticket_id | uint | FK → tickets |
| scheduled_by | uint | FK → users |
| subject | varchar(255) |
| direction | varchar(10) | Inbound, Outbound |
| status | varchar(20) | Completed, Scheduled |
| start_time | timestamp |
| description | text |
| call_type | varchar(50) |

#### `ticket_activities`
| Column | Type |
|---|---|
| id | uint |
| ticket_id | uint | FK → tickets |
| user_id | uint | FK → users (nullable) |
| activity_type | varchar(50) | Enum constants |
| description | text |
| old_value | varchar(255) |
| new_value | varchar(255) |
| remarks | text | For status change explanations |
| created_at | timestamp |

**Activity Types:** TICKET_CREATED, STATUS_CHANGED, ASSIGNED, UNASSIGNED, PRIORITY_CHANGED, COMMENT_ADDED, RESOLUTION_ADDED, CALL_SCHEDULED, CALL_COMPLETED, CALL_CANCELLED, PRODUCT_CHANGED, SUBJECT_CHANGED, APPROVAL_REQUESTED, APPROVAL_APPROVED, APPROVAL_REJECTED, TASK_CREATED, TASK_STATUS_CHANGED, TASK_ASSIGNEE_CHANGED, TASK_PRIORITY_CHANGED, TASK_COMMENT_ADDED

#### `ticket_approvals`
| Column | Type |
|---|---|
| id | uint |
| ticket_id | uint | FK → tickets |
| requester_id | uint | FK → users |
| approver_id | uint | FK → users |
| subject | varchar(255) |
| status | varchar(20) | PENDING, APPROVED, REJECTED |
| remarks | text |

---

### 6.4 Task Sub-Tables

#### `task_comments`
| Column | Type |
|---|---|
| id | uint |
| task_id | uint | FK → tasks |
| user_id | uint | FK → users |
| content | text |

#### `task_activities`
| Column | Type |
|---|---|
| id | uint |
| task_id | uint | FK → tasks |
| user_id | uint | FK → users |
| activity_type | varchar(50) |
| description | text |
| old_value, new_value | varchar(255) |

---

### 6.5 Notification Tables

#### `notifications`
| Column | Type |
|---|---|
| id | uint |
| recipient_id | uint |
| recipient_type | varchar(10) | user / contact |
| title | varchar(255) |
| message | text |
| notification_type | varchar(100) | e.g., ticket.assigned_to_you |
| priority | varchar(10) | low, normal, high, urgent |
| category | varchar(20) | ticket, task, communication, system |
| related_id | uint | ID of related entity |
| related_type | varchar(50) | ticket, task |
| is_read | boolean |
| read_at | timestamp |

#### `notification_templates`
| Column | Type |
|---|---|
| id | uint |
| notification_type | varchar(100) | UNIQUE |
| title_template | text | e.g., "Ticket {{ticket_id}} assigned to you" |
| message_template | text |
| default_priority | varchar(10) |
| category | varchar(20) |
| is_active | boolean |

---

### 6.6 Audit Log Table

#### `audit_logs`
| Column | Type | Purpose |
|---|---|---|
| id | uint64 | PK |
| actor_id | uint | User/Contact ID |
| actor_type | varchar(20) | user, contact, system, n8n |
| actor_name | varchar(255) | Display name |
| actor_email | varchar(255) | |
| actor_ip_address | varchar(45) | Client IP |
| action | varchar(100) | e.g., ticket.created, user.login |
| entity_type | varchar(50) | ticket, user, account, etc. |
| entity_id | uint | |
| entity_name | varchar(255) | |
| description | text | Human-readable description |
| old_values | jsonb | Pre-change state |
| new_values | jsonb | Post-change state |
| changes_summary | text | |
| http_method | varchar(10) | |
| endpoint | varchar(255) | |
| user_agent | text | |
| request_id | varchar(100) | UUID from AuditMiddleware |
| severity | varchar(20) | info, warning, critical |
| status | varchar(20) | success, failure |
| error_message | text | |
| metadata | jsonb | |
| created_at | timestamp | |

---

### 6.7 Entity-Relationship Summary

```
master_roles ──1:N──> users
master_user_designations ──1:N──> users
master_contact_designations ──1:N──> contacts
master_products ──1:N──> master_product_issues
master_products ──1:N──> tickets

accounts ──1:N──> contacts
accounts ──1:N──> tickets

contacts ──1:N──> tickets (as creator)
users ──1:N──> tickets (as assigned_engineer)

tickets ──1:N──> ticket_comments
tickets ──1:N──> ticket_calls
tickets ──1:N──> ticket_attachments
tickets ──1:N──> ticket_activities
tickets ──1:N──> ticket_approvals
tickets ──1:N──> tasks (optional link)

users ──1:N──> tasks (as assigned_to)
users ──1:N──> tasks (as created_by)
tasks ──1:N──> task_comments
tasks ──1:N──> task_activities

users/contacts ──1:N──> notifications
```

---

## 7. Data Flow & Key Workflows

### 7.1 Authentication Flow

```
Frontend                          Backend
────────                          ───────
LoginUser.vue
  └─ POST /auth/login ──────────> auth.go: UserLogin
     {email, password}              ├─ Find user by email
                                    ├─ Verify bcrypt hash
                                    ├─ Generate JWT (access + refresh)
                                    ├─ Audit log: login event
  <────────────────────────────── └─ Return {access_token, refresh_token, first_login, user}
  │
  ├─ auth.setAuth(token, ...)
  ├─ If first_login → redirect /reset-password
  └─ Else → redirect to landing page

Router Guard (every navigation):
  ├─ Load auth from localStorage
  ├─ Decode JWT, check expiry
  ├─ Enforce first_login redirect
  └─ Role-based route protection
```

### 7.2 Ticket Lifecycle

```
1. CREATION
   ├─ Manager: POST /manager/tickets (TicketForm.vue with fuzzy search)
   ├─ Engineer: POST /engineer/tickets
   ├─ Contact: POST /customer/tickets
   └─ n8n: POST /n8n/ticket or /n8n/smart-ticket
   
   Backend:
   ├─ Resolve contact, account, product
   ├─ Generate ticket_id: {customer_code}-{DDMMYY}-{annual_seq}
   ├─ Create ticket record (status=OPEN)
   ├─ Process attachments (if any)
   ├─ Log activity: TICKET_CREATED
   ├─ Create notification (template-based)
   ├─ Broadcast notification via WebSocket
   ├─ Send email notification
   └─ Log audit entry

2. ASSIGNMENT
   PUT /manager/tickets/:id/assign
   ├─ Update assigned_engineer
   ├─ Log activity: ASSIGNED / UNASSIGNED
   ├─ Notify engineer + managers
   └─ Audit log

3. STATUS CHANGES
   PUT /manager/tickets/:id/status
   PUT /engineer/tickets/:id/status
   ├─ Validate status transition
   ├─ StatusChangeModal collects remarks
   ├─ Log activity with remarks: STATUS_CHANGED
   ├─ Notify relevant parties
   └─ Audit log

4. COMMENTS
   POST /tickets/:id/comments
   ├─ Internal comments (staff only) or customer-visible
   ├─ Log activity: COMMENT_ADDED
   └─ Notify contact or engineer

5. CALLS
   POST /tickets/:id/calls
   PUT  /tickets/:id/calls/:callId/complete
   PUT  /tickets/:id/calls/:callId/cancel
   ├─ Log call with subject, direction, description
   ├─ Log activity: CALL_SCHEDULED / COMPLETED / CANCELLED
   └─ Notify relevant parties

6. APPROVALS
   POST  /tickets/:id/approvals (request)
   PATCH /tickets/:id/approvals/:id/approve
   PATCH /tickets/:id/approvals/:id/reject
   ├─ Only managers can approve/reject
   ├─ Log activity with remarks
   └─ Notify requester/approver only

7. CLOSURE / REOPENING
   ├─ Close: StatusChangeModal pre-selects CLOSED
   ├─ Reopen: StatusChangeModal pre-selects OPEN
   ├─ Both require mandatory remarks
   └─ Dynamic button text + color
```

### 7.3 Task Lifecycle

```
1. CREATE: POST /manager/tasks
   ├─ Subject, description, priority, assignee, optional ticket link
   ├─ Log task activity: TASK_CREATED
   └─ Notify assigned engineer

2. UPDATE: PUT /manager/tasks/:id
   ├─ Update any field
   ├─ Log field change activities
   └─ Notify relevant parties

3. STATUS: PUT /manager/tasks/:id/status or /engineer/tasks/:id/status
   ├─ Open → In Progress → Completed → Closed
   └─ Log TASK_STATUS_CHANGED

4. COMMENTS: POST /tasks/:id/comments
   └─ Log TASK_COMMENT_ADDED
```

### 7.4 Notification Flow

```
1. TRIGGER: Any handler calls NotificationService
   ├─ notificationService.CreateNotification(type, recipientID, recipientType, variables, relatedID, relatedType)

2. TEMPLATE LOOKUP
   ├─ Find active template by notification_type
   ├─ Replace variables: {{ticket_id}}, {{user_name}}, {{status}}, etc.

3. DATABASE INSERT
   ├─ Create notification record

4. REAL-TIME BROADCAST
   ├─ hub.BroadcastToUser(recipientID, recipientType, "notification.new", notificationData)
   ├─ hub.BroadcastToUser(recipientID, recipientType, "count.update", {unread_count})

5. FRONTEND RECEPTION
   ├─ WebSocket service receives message
   ├─ Notification store updates local state
   ├─ NotificationBell updates badge count
   ├─ NotificationToast shows popup
   └─ NotificationsPage updates list
```

### 7.5 Approval Workflow

```
Requester (Manager/Engineer)           Approver (Manager)
─────────────────────────             ──────────────────
ApprovalsTab: "Add Approval"
  ├─ Select manager as approver
  ├─ Enter subject/reason
  └─ POST /tickets/:id/approvals
     ├─ Create approval (PENDING)        
     ├─ Notify approver ────────────────> Notification received
     └─ Log activity                      │
                                          ├─ ApprovalsTab: See pending request
                                          ├─ Click Approve/Reject
                                          ├─ StatusChangeModal: Enter remarks
                                          └─ PATCH /tickets/:id/approvals/:id/approve
                                             ├─ Update status
Notification received <──────────────────── ├─ Notify requester
                                             └─ Log activity with remarks
```

### 7.6 Audit Logging Flow

```
1. AuditMiddleware (global)
   ├─ Generate request_id (UUID)
   ├─ Capture client_ip, user_agent
   └─ Store in Gin context

2. Handler calls AuditService
   ├─ auditService.LogCRUD(c, action, entityType, entityID, entityName, description, oldValues, newValues)
   ├─ auditService.LogAuthentication(actorType, actorID, actorName, actorEmail, action, success, ip, ua, error)
   └─ Extracts request_id, client_ip, user_agent from context

3. Audit log stored with full context:
   ├─ WHO: actor_id, actor_type, actor_name, actor_email
   ├─ WHAT: action, entity_type, entity_id, entity_name
   ├─ CHANGES: old_values (jsonb), new_values (jsonb), description
   ├─ HOW: http_method, endpoint, user_agent, request_id
   └─ META: severity, status, error_message, metadata

4. Retrieval: GET /manager/audit-logs
   ├─ Filters: action, entity_type, actor_type, severity, status, search, date range
   ├─ Pagination: page, limit
   └─ Stats: total_logs, critical_logs, failed_logs, auth_events
```

---

## 8. n8n Integration & AI Pipeline

### 8.1 Architecture Overview

```
┌─────────────┐     ┌──────────┐     ┌─────────────────┐     ┌──────────────┐
│  Email       │────>│  n8n     │────>│  Gemini AI      │────>│  Backend     │
│  (IMAP)      │     │ Workflow │     │  (2.5 Flash)    │     │  API         │
└─────────────┘     └──────────┘     └─────────────────┘     └──────────────┘
                         │                    │                      │
                         │  1. Fetch email    │ 2. Extract entities  │ 3. Smart resolve
                         │  2. Parse content  │    - phone numbers   │    - Contact (email/phone/name)
                         │  3. Call Gemini    │    - person names    │    - Account (from contact/org)
                         │  4. Call backend   │    - org names       │    - Product (from hints)
                         │                    │    - product hints   │ 4. Create ticket
                         │                    │    - priority hints  │ 5. Send notification
                         │                    │                      │
                    OR: AI Agent Mode                               │
                         │                                          │
                    n8n AI Agent node                               │
                    with function calling ──────────────────────────┘
                    uses /n8n/tools/* endpoints
```

### 8.2 Webhook Endpoints

**Authentication:** All `/n8n/*` routes require `X-N8N-API-Key` header or `api_key` query param.

| Endpoint | Method | Purpose |
|---|---|---|
| `/n8n/health` | GET | Health check |
| `/n8n/ticket` | POST | Basic ticket creation (direct IDs or name-based lookup) |
| `/n8n/smart-ticket` | POST | AI-powered smart ticket creation with waterfall entity resolution |
| `/n8n/lookup/accounts` | GET | Fuzzy account search with match scoring |
| `/n8n/lookup/contacts` | GET | Fuzzy contact search by name/email |
| `/n8n/lookup/products` | GET | Product search/listing |

### 8.3 Smart Resolver Engine

**File:** `n8n_smart_resolver.go`

The `SmartResolver` implements a **waterfall resolution strategy**:

#### Contact Resolution (highest priority)
1. **Exact email match** (confidence: 95-100)
2. **Phone number match** from AI-extracted hints
3. **Name match** from email signature/body
4. **Domain-based account match** (email domain → account → first contact)

#### Account Resolution
1. **From resolved contact** (contact.account_id)
2. **Organization name match** from AI-extracted hints
3. **Email domain match** against account records

#### Product Resolution
1. **Product hint match** from AI-extracted keywords
2. **Default product** fallback (first product in system)

#### Priority Determination
- Analyzes AI-extracted priority hints (urgent, critical, ASAP, etc.)
- Normalizes to High/Medium/Low

#### Ticket Creation
- Generates ticket_id with proper format
- Sets channel to "Mail"
- Processes base64-encoded attachments
- Logs activity and sends notifications
- Returns resolution confidence scores and warnings

### 8.4 AI Agent Tool Endpoints

**File:** `n8n_tools.go`

These endpoints are designed for **n8n's AI Agent node** with function calling:

| Endpoint | Purpose | Return |
|---|---|---|
| `search-contact-by-email` | Exact email match | contact_id, account info |
| `search-contact-by-phone` | Phone number match | contact_id, account info |
| `search-contact-by-name` | Fuzzy name search | Ranked contact list |
| `search-account-by-name` | Fuzzy account search | Ranked account list |
| `search-product-by-name` | Product search | Product list |
| `list-products` | All products | Full product list |
| `create-ticket` | Final ticket creation | ticket_id, details |
| `extract-email-with-gemini` | AI entity extraction | Structured hints |
| `dump-unresolved-query` | Save failed resolution | dumped_query_id |

### 8.5 Gemini AI Service

**File:** `services/gemini_service.go`

| Method | Purpose |
|---|---|
| `CallGemini(prompt)` | Send prompt to Gemini API, return text response |
| `ExtractEmailData(prompt)` | Extract structured data (phones, names, orgs, products, priority) from email |
| `ParseEmailFromPrompt(prompt)` | Regex-based extraction of From, Subject, Body from email text |
| `parseExtractedData(response)` | Parse JSON from Gemini response (handles markdown code blocks) |

**Configuration:** Uses `GEMINI_API_KEY` environment variable. Temperature set to 0.1 for consistent extraction.

---

## 9. Database Migrations

| Migration | Purpose |
|---|---|
| `001_add_ticket_enhancements.sql` | Initial ticket schema enhancements |
| `002_enhance_ticket_calls.sql` | Call logging schema improvements |
| `003_recreate_ticket_calls.sql` | Call table restructure |
| `003_update_call_status_values.sql` | Call status value updates |
| `004_add_notifications_system.sql` | Notification + template tables |
| `004_add_remarks_to_ticket_activities.sql` | Remarks column for activities |
| `005_add_ticket_attachments.sql` | Attachment table |
| `005_create_ticket_approvals.sql` | Approval workflow table |
| `006_add_approval_notification_templates.sql` | Approval notification templates |
| `006_enhance_call_logging.sql` | Enhanced call fields |
| `007_add_channel_to_tickets.sql` | Channel field (Web/Mail) |
| `007_add_close_remarks_to_calls.sql` | Close remarks for calls |
| `008_create_audit_logs_system.sql` | Audit logs table + indexes |

> **Note:** GORM auto-migration handles most schema changes. SQL migrations are for complex alterations, constraints, and seed data.

---

## 10. Security

| Feature | Implementation |
|---|---|
| **Password Storage** | bcrypt hashing (cost factor: default) |
| **Authentication** | JWT Bearer tokens (HS256) |
| **Token Expiry** | Access: 24h, Refresh: 7d |
| **First-Login Enforcement** | Mandatory password reset for new users (router guard + backend flag) |
| **Role-Based Access** | Route-level (meta.requiresRole) + handler-level (IsManager/IsEngineer checks) |
| **CORS** | Restricted to localhost dev origins |
| **n8n API Auth** | Separate API key mechanism (header or query param) |
| **WebSocket Auth** | JWT validated during upgrade, user loaded from DB |
| **Input Validation** | Gin's ShouldBindJSON with struct tags |
| **SQL Injection** | Prevented by GORM's parameterized queries |
| **File Upload** | 3MB size limit, file type validation |
| **Audit Trail** | All CRUD, auth, and system events logged with full context |

---

## 11. Environment Configuration

### `.env.email` (Backend)
```env
# Database
DATABASE_URL=host=localhost user=postgres password=postgres dbname=globx_hd port=5432 sslmode=disable TimeZone=Asia/Kolkata

# Server
SERVER_ADDRESS=:8080

# JWT
JWT_SECRET=your-secret-key

# n8n Integration
N8N_API_KEY=your-n8n-api-key

# Gemini AI
GEMINI_API_KEY=your-gemini-api-key

# Email Notifications
EMAIL_SMTP_SERVER=smtp.example.com
EMAIL_SMTP_PORT=587
EMAIL_USERNAME=notifications@example.com
EMAIL_PASSWORD=your-email-password
EMAIL_NOTIFICATION_ADDRESS=support@example.com

# File Uploads
UPLOAD_DIR=./uploads
```

### Frontend Environment
```env
# .env or .env.local
VITE_API_BASE=http://localhost:8080
VITE_API_URL=localhost:8080    # For WebSocket
```

---

*End of Documentation*
