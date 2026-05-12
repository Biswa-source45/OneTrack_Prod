# Audit Logging System Diagnosis

## Current Implementation Analysis

### ✅ What's Already Implemented

1. **Database Schema** (`migrations/008_create_audit_logs_system.sql`)
   - Table: `audit_logs` with all required fields
   - Indexes for performance
   - Materialized views for analytics

2. **Backend Models** (`internal/models/audit_log.go`)
   - Complete AuditLog struct with GORM tags
   - 60+ action constants defined
   - Proper table name mapping

3. **Audit Service** (`internal/services/audit_service.go`)
   - `LogAuthentication()` - For login/logout/password reset
   - `LogCRUD()` - For create/update/delete operations
   - `LogError()` - For error events
   - Helper functions to extract actor and request context

4. **Handler Integration** (`internal/handlers/auth.go`)
   - User login success/failure logging
   - Contact login success/failure logging
   - Logout logging (both user and contact)
   - Password reset logging

### 🔍 Potential Issues Identified

#### Issue 1: GORM Model Type Mismatch
**Problem**: The model uses `uint64` for ID but database uses `BIGSERIAL`
```go
// Model
ID uint64 `gorm:"primaryKey" json:"id"`

// Database
id BIGSERIAL PRIMARY KEY
```
**Impact**: This should work fine, but could cause issues with auto-increment

#### Issue 2: JSONB Field Handling
**Problem**: Model uses `string` for JSONB fields
```go
OldValues string `gorm:"type:jsonb" json:"old_values"`
NewValues string `gorm:"type:jsonb" json:"new_values"`
Metadata  string `gorm:"type:jsonb;default:'{}'" json:"metadata"`
```
**Impact**: GORM might not handle JSONB serialization properly

#### Issue 3: Missing Error Handling
**Problem**: Audit logging errors are silently ignored
```go
// In auth.go
auditService.LogAuthentication(...) // No error check
```
**Impact**: If logging fails, we won't know about it

#### Issue 4: Database Connection
**Problem**: No verification that audit_logs table exists
**Impact**: If migration wasn't run, all inserts will fail silently

#### Issue 5: Logout Handler Context
**Problem**: Logout handler expects user/contact in context but might not have it
```go
if userVal, exists := c.Get("user"); exists {
    // Log logout
}
```
**Impact**: If context is not set, logout won't be logged

## Recommended Fixes

### Fix 1: Add Error Logging and Debugging
Add error handling to see if audit logs are failing:

```go
// In auth.go - User Login Success
if err := auditService.LogAuthentication(...); err != nil {
    log.Printf("ERROR: Failed to log authentication: %v", err)
}
```

### Fix 2: Verify Table Exists
Add a check in main.go or during service initialization:

```go
func VerifyAuditLogTable(db *gorm.DB) error {
    var count int64
    if err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'audit_logs'").Scan(&count).Error; err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("audit_logs table does not exist - run migration 008")
    }
    return nil
}
```

### Fix 3: Use JSONB Type Properly
Change model to use proper JSONB handling:

```go
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
    return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }
    return json.Unmarshal(bytes, &j)
}
```

### Fix 4: Add Debug Logging
Add temporary debug logs to trace execution:

```go
// In audit_service.go - LogAuthentication
func (s *AuditService) LogAuthentication(...) error {
    log.Printf("DEBUG: Logging authentication - Actor: %s, Action: %s", actorName, action)
    
    log := models.AuditLog{...}
    
    if err := s.db.Create(&log).Error; err != nil {
        log.Printf("ERROR: Failed to create audit log: %v", err)
        return err
    }
    
    log.Printf("DEBUG: Successfully created audit log ID: %d", log.ID)
    return nil
}
```

### Fix 5: Ensure Logout Context
Make sure logout handler is called with proper authentication:

```go
// In routes.go
router.POST("/logout", middleware.AuthMiddleware(db), handlers.Logout(db))
```

## Testing Checklist

1. ✅ Verify `audit_logs` table exists in database
2. ✅ Test user login - check if log is created
3. ✅ Test user login failure - check if failure is logged
4. ✅ Test contact login - check if log is created
5. ✅ Test logout - check if log is created
6. ✅ Test password reset - check if log is created
7. ✅ Test CRUD operations - check if logs are created
8. ✅ Query audit_logs table directly to see if any records exist

## SQL Queries for Testing

```sql
-- Check if table exists
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_name = 'audit_logs'
);

-- Count total audit logs
SELECT COUNT(*) FROM audit_logs;

-- View recent audit logs
SELECT id, actor_name, action, entity_type, description, created_at 
FROM audit_logs 
ORDER BY created_at DESC 
LIMIT 10;

-- Check authentication logs
SELECT * FROM audit_logs 
WHERE action LIKE '%LOGIN%' OR action = 'LOGOUT'
ORDER BY created_at DESC;

-- Check for any errors
SELECT * FROM audit_logs 
WHERE status IN ('failure', 'error')
ORDER BY created_at DESC;
```

## Next Steps

1. Run SQL query to verify table exists
2. Add error logging to audit service
3. Test login/logout and check database
4. Add debug logging if needed
5. Fix any identified issues
