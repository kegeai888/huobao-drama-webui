# ComfyUI URL Double Slash Fix

## Problem
ComfyUI image generation was failing with error:
```
[ComfyUI-Image] Sending request to: http://192.168.2.76:8080//prompt
ERROR: API error (status 405): 405: Method Not Allowed
```

Notice the double slash `//` in the URL path.

## Root Cause
The URL construction in both ComfyUI clients had a bug:

```go
url := c.BaseURL + c.Endpoint
```

When `BaseURL` ends with `/` (e.g., `http://192.168.2.76:8080/`) and `Endpoint` starts with `/` (e.g., `/prompt`), the result is a double slash: `http://192.168.2.76:8080//prompt`.

Many web servers (including ComfyUI) treat `//prompt` differently from `/prompt` and may return 405 Method Not Allowed or 404 errors.

## Solution
Normalized URL construction in both ComfyUI clients:

1. **Trim trailing slashes from BaseURL**:
   ```go
   baseURL = strings.TrimRight(baseURL, "/")
   ```

2. **Ensure endpoints start with a slash**:
   ```go
   if !strings.HasPrefix(endpoint, "/") {
       endpoint = "/" + endpoint
   }
   ```

This ensures URLs are always constructed correctly:
- `http://192.168.2.76:8080` + `/prompt` = `http://192.168.2.76:8080/prompt` ✅
- `http://192.168.2.76:8080/` + `/prompt` = `http://192.168.2.76:8080/prompt` ✅

## Files Modified
1. `pkg/image/comfyui_image_client.go` - Fixed `NewComfyUIImageClient()`
2. `pkg/video/comfyui_client.go` - Fixed `NewComfyUIClient()`

## Testing
After the fix, ComfyUI requests should work correctly regardless of whether the BaseURL in the AI config has a trailing slash or not:

```
✅ http://192.168.2.76:8080  → http://192.168.2.76:8080/prompt
✅ http://192.168.2.76:8080/ → http://192.168.2.76:8080/prompt
```

## Configuration Example
In AI Config settings, both formats now work:
```yaml
# Without trailing slash (recommended)
base_url: http://192.168.2.76:8080

# With trailing slash (also works now)
base_url: http://192.168.2.76:8080/
```
