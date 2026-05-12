# Audit Logging System - Testing & Verification Guide

## 🔍 What Was Done to Refine the System

### 1. **Added Comprehensive Debug Logging**
- Added `[AUDIT]` prefix logs to track every audit operation
- Added `[AUDIT ERROR]` logs to identify failures
- Added `[AUDIT SUCCESS]` logs to confirm successful saves
- Logs show: Actor name, action type, entity details, and database operation results

### 2. **Added Error Handling**
- All audit logging calls now check for errors
- Errors are logged to console with WARNING prefix
- Authentication handlers won't fail if audit logging fails (non-blocking)

### 3. **Created Test Endpoint**
- New endpoint: `GET /test-audit-log`
- Tests 4 things:
  1. Can create audit log entry
  2. Does audit_logs table exist
  3. Count of audit logs in database
  4. Retrieve recent logs

### 4. **Enhanced Audit Service**
- Better error messages
- Detailed logging at each step
- Proper variable naming (changed `log` to `auditLog` to avoid confusion)

---

## 🧪 Testing Procedures

### **Step 1: Verify Database Table Exists**

Run this SQL query in pgAdmin:

```sql
-- Check if audit_logs table exists
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_name = 'audit_logs'
);

-- If false, run the migration:
-- \i 'C:/Users/KIIT0001/globx_hd/globx_hd/ticketing_backend/migrations/008_create_audit_logs_system.sql'
```

### **Step 2: Restart Backend Server**

```bash
cd ticketing_backend
# Stop the server (Ctrl+C)
go run cmd/main.go
```

Watch the console output for `[AUDIT]` logs.

### **Step 3: Test Using the Test Endpoint**

```bash
# Call the test endpoint
curl http://localhost:8080/test-audit-log
```

**Expected Response:**
```json
{
  "test_1_create_log": {
    "success": true,
    "error": "<nil>"
  },
  "test_2_table_exists": {
    "exists": true,
    "error": "<nil>"
  },
  "test_3_count_logs": {
    "count": 5,
    "error": "<nil>"
  },
  "test_4_recent_logs": {
    "count": 5,
    "logs": [...],
    "error": "<nil>"
  }
}
```

**If test fails:**
- Check console logs for `[AUDIT ERROR]` messages
- Verify database connection
- Ensure migration was run

### **Step 4: Test Login Logging**

#### **Test User Login:**
```bash
curl -X POST http://localhost:8080/login/user \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }'
```

**Watch console for:**
```
[AUDIT] Logging authentication - Actor: John Doe, Action: USER_LOGIN_SUCCESS, Success: true
[AUDIT SUCCESS] Created audit log ID: 123
```

#### **Test Contact Login:**
```bash
curl -X POST http://localhost:8080/login/contact \
  -H "Content-Type: application/json" \
  -d '{
    "username": "contact@example.com",
    "password": "password"
  }'
```

**Watch console for:**
```
[AUDIT] Logging authentication - Actor: Jane Smith, Action: CONTACT_LOGIN_SUCCESS, Success: true
[AUDIT SUCCESS] Created audit log ID: 124
```

#### **Test Failed Login:**
```bash
curl -X POST http://localhost:8080/login/user \
  -H "Content-Type: application/json" \
  -d '{
    "username": "wrong_user",
    "password": "wrong_pass"
  }'
```

**Watch console for:**
```
[AUDIT] Logging authentication - Actor: wrong_user, Action: USER_LOGIN_FAILURE, Success: false
[AUDIT SUCCESS] Created audit log ID: 125
```

### **Step 5: Verify Logs in Database**

```sql
-- Count total audit logs
SELECT COUNT(*) FROM audit_logs;

-- View recent authentication logs
SELECT 
    id,
    actor_name,
    actor_email,
    action,
    status,
    actor_ip_address,
    created_at
FROM audit_logs 
WHERE entity_type = 'authentication'
ORDER BY created_at DESC 
LIMIT 10;

-- View all logs from today
SELECT 
    id,
    actor_name,
    action,
    entity_type,
    entity_name,
    description,
    status,
    severity,
    created_at
FROM audit_logs 
WHERE DATE(created_at) = CURRENT_DATE
ORDER BY created_at DESC;

-- Check for any failed operations
SELECT * FROM audit_logs 
WHERE status IN ('failure', 'error')
ORDER BY created_at DESC;
```

### **Step 6: Test CRUD Operations Logging**

#### **Create a User:**
```bash
curl -X POST http://localhost:8080/users \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "first_name": "Test",
    "last_name": "User",
    "password": "password123",
    "role_id": 3
  }'
```

**Watch console for:**
```
[AUDIT] LogWithContext - Action: USER_CREATED, Entity: user, Name: Test User
[AUDIT] Actor extracted - Type: user, Name: Admin User
[AUDIT SUCCESS] Created audit log ID: 126
```

**Verify in database:**
```sql
SELECT * FROM audit_logs 
WHERE action = 'USER_CREATED' 
ORDER BY created_at DESC 
LIMIT 1;
```

#### **Update a User:**
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Updated",
    "last_name": "Name"
  }'
```

**Check for old_values and new_values in database:**
```sql
SELECT 
    id,
    action,
    entity_name,
    old_values,
    new_values,
    created_at
FROM audit_logs 
WHERE action = 'USER_UPDATED' 
ORDER BY created_at DESC 
LIMIT 1;
```

### **Step 7: Test Logout Logging**

```bash
curl -X POST http://localhost:8080/logout \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Watch console for:**
```
[AUDIT] Logging authentication - Actor: John Doe, Action: LOGOUT, Success: true
[AUDIT SUCCESS] Created audit log ID: 127
```

---

## 🐛 Troubleshooting

### **Issue: No logs appearing in database**

**Check 1: Table exists?**
```sql
SELECT tablename FROM pg_tables WHERE tablename = 'audit_logs';
```
If empty, run migration 008.

**Check 2: Database connection?**
Look for GORM connection errors in console.

**Check 3: Console logs?**
Look for `[AUDIT ERROR]` messages showing the actual error.

**Check 4: Permissions?**
```sql
SELECT has_table_privilege('audit_logs', 'INSERT');
```
Should return `true`.

### **Issue: [AUDIT ERROR] appears in console**

**Common errors:**

1. **"relation 'audit_logs' does not exist"**
   - Solution: Run migration 008

2. **"column 'xyz' does not exist"**
   - Solution: Model and table schema mismatch, check migration

3. **"violates check constraint"**
   - Solution: Invalid value for actor_type, severity, or status

4. **"null value in column 'description' violates not-null constraint"**
   - Solution: Description is required, check audit service calls

### **Issue: Logs created but not showing in frontend**

**Check 1: API endpoint working?**
```bash
curl http://localhost:8080/manager/audit-logs \
  -H "Authorization: Bearer YOUR_MANAGER_TOKEN"
```

**Check 2: Frontend API service configured?**
Check `src/api/auditLogs.js` exists and has correct base URL.

**Check 3: Route registered?**
Check `src/router/index.ts` has `/manager/audit-logs` route.

---

## ✅ Success Criteria

After testing, you should see:

1. ✅ **Console Logs**: `[AUDIT]` and `[AUDIT SUCCESS]` messages for every operation
2. ✅ **Database Records**: Audit logs in `audit_logs` table
3. ✅ **Login Tracking**: Both successful and failed login attempts logged
4. ✅ **Logout Tracking**: Logout events logged with user details
5. ✅ **CRUD Tracking**: Create/Update/Delete operations logged with old/new values
6. ✅ **Actor Information**: Correct user/contact name, email, IP address
7. ✅ **Timestamps**: Accurate created_at timestamps
8. ✅ **No Errors**: No `[AUDIT ERROR]` messages in console

---

## 📊 Sample Audit Log Entry

```json
{
  "id": 123,
  "actor_id": 12,
  "actor_type": "user",
  "actor_name": "John Doe",
  "actor_email": "john@example.com",
  "actor_ip_address": "192.168.1.100",
  "action": "USER_LOGIN_SUCCESS",
  "entity_type": "authentication",
  "entity_id": null,
  "entity_name": "",
  "description": "USER_LOGIN_SUCCESS: john@example.com",
  "old_values": null,
  "new_values": null,
  "changes_summary": null,
  "http_method": "POST",
  "endpoint": "/login/user",
  "user_agent": "Mozilla/5.0...",
  "request_id": "uuid-here",
  "severity": "info",
  "status": "success",
  "error_message": null,
  "metadata": "{}",
  "created_at": "2026-01-14T12:30:45.123Z"
}
```

---

## 🔧 Maintenance

### **Remove Test Endpoint (Production)**

Before deploying to production, remove the test endpoint:

In `internal/routes/routes.go`, delete:
```go
// Test endpoint for audit logging (remove in production)
r.GET("/test-audit-log", handlers.TestAuditLog(db))
```

### **Disable Debug Logging (Production)**

In `internal/services/audit_service.go`, comment out debug logs:
```go
// log.Printf("[AUDIT] Logging authentication - Actor: %s...", ...)
// log.Printf("[AUDIT SUCCESS] Created audit log ID: %d", auditLog.ID)
```

Keep only error logs:
```go
log.Printf("[AUDIT ERROR] Failed to create audit log: %v", err)
```

---

## 📈 Monitoring Queries

### **Daily Activity Summary**
```sql
SELECT 
    DATE(created_at) as date,
    COUNT(*) as total_logs,
    COUNT(*) FILTER (WHERE status = 'success') as successful,
    COUNT(*) FILTER (WHERE status = 'failure') as failed,
    COUNT(*) FILTER (WHERE action LIKE '%LOGIN%') as logins
FROM audit_logs
WHERE created_at >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY DATE(created_at)
ORDER BY date DESC;
```

### **Most Active Users**
```sql
SELECT 
    actor_name,
    actor_email,
    COUNT(*) as action_count
FROM audit_logs
WHERE actor_type = 'user'
    AND created_at >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY actor_name, actor_email
ORDER BY action_count DESC
LIMIT 10;
```

### **Failed Operations**
```sql
SELECT 
    actor_name,
    action,
    entity_type,
    error_message,
    created_at
FROM audit_logs
WHERE status IN ('failure', 'error')
    AND created_at >= CURRENT_DATE - INTERVAL '7 days'
ORDER BY created_at DESC;
```

---

## 🎯 Next Steps

1. **Run Step 1-3** to verify basic functionality
2. **Run Step 4** to test login/logout logging
3. **Run Step 5** to verify database entries
4. **Run Step 6** to test CRUD logging
5. **Check console** for any `[AUDIT ERROR]` messages
6. **Query database** to see actual log entries
7. **Test frontend** to view logs in UI

If any step fails, check the Troubleshooting section and console logs for specific error messages.
