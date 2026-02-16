# API Documentation

## Base URLs
```
Development:  http://localhost:8080/api/v1
Production:   https://netops-api.up.railway.app/api/v1
```

## Authentication

### Agent Endpoints
```
Header: Authorization: Bearer <api_key>
```

### User Endpoints  
```
Header: Authorization: Bearer <jwt_token>
```

## Core Endpoints

### Health Check
```http
GET /health

Response 200:
{
  "status": "healthy",
  "service": "netops-integration-api",
  "timestamp": "2026-02-16T10:30:00Z"
}
```

### Agent - Heartbeat
```http
POST /api/v1/agent/heartbeat

Headers:
  Authorization: Bearer <api_key>
  Content-Type: application/json

Body:
{
  "device_id": "router-bdg-01",
  "status": "online",
  "version": "7.14",
  "uptime": "1d 2h 30m",
  "cpu_load": "25%",
  "free_memory": "512MB"
}

Response 200:
{
  "success": true,
  "message": "Heartbeat recorded"
}
```

### Agent - Get Commands
```http
GET /api/v1/agent/commands/:deviceId

Response 200:
{
  "success": true,
  "commands": [
    {
      "id": "cmd_123",
      "type": "change_dns",
      "parameters": {
        "primary": "8.8.8.8",
        "secondary": "8.8.4.4"
      }
    }
  ]
}
```

### Devices - List
```http
GET /api/v1/devices

Response 200:
{
  "success": true,
  "data": [
    {
      "id": "dev_001",
      "name": "Router-BDG-01",
      "ip_address": "192.168.100.11",
      "location": "Bandung",
      "status": "online",
      "last_seen": "2026-02-16T10:25:00Z"
    }
  ],
  "total": 1
}
```

### Devices - Bulk Configure
```http
POST /api/v1/devices/bulk/configure

Body:
{
  "device_filter": {
    "all": true
  },
  "command": {
    "type": "change_dns",
    "parameters": {
      "primary": "8.8.8.8",
      "secondary": "8.8.4.4"
    }
  }
}

Response 200:
{
  "success": true,
  "execution_id": "exec_123",
  "total_devices": 80,
  "status": "queued"
}
```

### Executions - Track Progress
```http
GET /api/v1/executions/:id

Response 200:
{
  "success": true,
  "data": {
    "id": "exec_123",
    "total_devices": 80,
    "completed": 65,
    "failed": 2,
    "progress": 81,
    "status": "executing"
  }
}
```

## Error Responses
```json
{
  "success": false,
  "error": {
    "code": "DEVICE_NOT_FOUND",
    "message": "Device with ID 'dev_999' not found"
  }
}
```

## Rate Limits

- Agent endpoints: 100 requests/minute per device
- User endpoints: 1000 requests/hour per user

## Status Codes

- `200` OK
- `201` Created
- `400` Bad Request
- `401` Unauthorized
- `404` Not Found
- `429` Too Many Requests
- `500` Internal Server Error

---

**Interactive Docs:** `/swagger` (when deployed)  
**Postman Collection:** See `/postman/collection.json`
