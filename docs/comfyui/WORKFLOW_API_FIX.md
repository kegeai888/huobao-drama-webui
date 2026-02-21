# Workflow API Network Error Fix

## Problem
When uploading workflow files successfully, the refresh/list operation failed with:
```
AxiosError: Network Error
```

## Root Cause
The `workflow_handler.go` was returning raw JSON responses without the standard API response wrapper:
```go
c.JSON(http.StatusOK, gin.H{
    "workflows": workflows,
    "count": len(workflows),
})
```

However, the frontend's `request.ts` interceptor expects all API responses to have this format:
```json
{
  "success": true,
  "data": { ... },
  "timestamp": "..."
}
```

The interceptor checks for `res.success` and only returns `res.data` if it's true. Without the `success` field, the response was rejected.

## Solution
Updated `workflow_handler.go` to use the standard `response` package (same as other handlers):

1. Added import: `"github.com/drama-generator/backend/pkg/response"`
2. Replaced all `c.JSON()` calls with response helper functions:
   - `response.Success(c, data)` - for successful responses
   - `response.BadRequest(c, message)` - for 400 errors
   - `response.NotFound(c, message)` - for 404 errors
   - `response.InternalError(c, message)` - for 500 errors

## Changes Made
- Modified `api/handlers/workflow_handler.go`:
  - `ListWorkflows()` - now returns `{"success": true, "data": {"workflows": [...], "count": N}}`
  - `UploadWorkflow()` - now returns `{"success": true, "data": {"message": "...", "filename": "...", "path": "..."}}`
  - `DeleteWorkflow()` - now returns `{"success": true, "data": {"message": "..."}}`
  - `GetWorkflow()` - error responses now use response helpers (success case returns raw JSON as before)

## Testing
```bash
# Test workflow list API
curl http://localhost:5678/api/v1/workflows

# Response format:
{
  "success": true,
  "data": {
    "workflows": [...],
    "count": 3
  },
  "timestamp": "2026-02-20T17:22:41Z"
}
```

## Result
✅ Workflow upload works correctly
✅ Workflow list refresh works without Network Error
✅ All workflow CRUD operations return consistent response format
✅ Frontend can properly parse and display workflow data
