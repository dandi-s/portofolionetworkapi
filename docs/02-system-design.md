# System Design Document

## Architecture Overview

### High-Level Architecture
```
┌─────────────────────────────────────┐
│     External Systems Layer          │
│  ┌────────┐ ┌─────────┐ ┌────────┐ │
│  │ Zabbix │ │MikroTik │ │HubSpot │ │
│  └───┬────┘ └────┬────┘ └───┬────┘ │
└──────┼───────────┼──────────┼───────┘
       │           │          │
       └───────────┼──────────┘
                   │
       ┌───────────▼──────────┐
       │   API Gateway        │
       │  - Auth              │
       │  - Rate Limit        │
       │  - Logging           │
       └───────────┬──────────┘
                   │
       ┌───────────▼──────────┐
       │  Business Logic      │
       │  - Device Service    │
       │  - Command Service   │
       │  - Incident Service  │
       └───────────┬──────────┘
                   │
       ┌───────────▼──────────┐
       │    Data Layer        │
       │  ┌────────┐┌──────┐ │
       │  │Postgres││Redis │ │
       │  └────────┘└──────┘ │
       └──────────────────────┘
```

## Technology Stack Rationale

### Golang
**Why chosen:**
- Compiled language: Fast execution
- Goroutines: Perfect for concurrent device operations
- Single binary: Easy deployment
- Strong typing: Fewer runtime errors

**Use cases in project:**
- Parallel device status checks
- Concurrent command execution
- Webhook processing

### PostgreSQL
**Why chosen:**
- ACID compliance for critical data
- JSONB for flexible metadata
- Full-text search for incidents
- Battle-tested reliability

**Schema design:**
- Normalized for consistency
- Proper indexing for performance
- Foreign keys for integrity

### Redis
**Why chosen:**
- In-memory: Microsecond latency
- Data structures: Lists, sets, sorted sets
- Pub/sub: Real-time events
- Persistence: AOF for durability

**Use cases:**
- Device status cache (30s TTL)
- Command queue management
- Rate limiting counters
- Session storage

## Data Models

### Core Entities
```go
type Device struct {
    ID          uuid.UUID
    Name        string
    IPAddress   string
    Location    string
    Status      string    // online, offline
    LastSeen    time.Time
    Metadata    JSONB
}

type Command struct {
    ID          uuid.UUID
    Type        string    // firmware_update, config_change
    Parameters  JSONB
    DeviceIDs   []uuid.UUID
    Status      string    // pending, executing, completed
    Results     map[string]string
}

type Customer struct {
    ID            uuid.UUID
    Name          string
    Tier          string  // gold, silver, bronze
    SLAPercentage float64
    Devices       []Device
}
```

## API Design Principles

1. **RESTful:** Resource-based URLs
2. **Versioned:** /api/v1/ for compatibility
3. **Consistent:** Standard response format
4. **Documented:** OpenAPI/Swagger specs

### Endpoint Structure
```
/api/v1/
  /agent/           # Device agent communication
  /devices/         # Device management
  /customers/       # Customer operations
  /commands/        # Command execution
  /incidents/       # Incident tracking
  /dashboard/       # Analytics data
```

### Standard Response
```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "timestamp": "2026-02-16T10:30:00Z",
    "request_id": "req_abc123"
  }
}
```

## Security Architecture

### Authentication Layers
- **Agent Auth:** API Key (stored in device config)
- **User Auth:** JWT tokens with refresh mechanism
- **RBAC:** Role-based access control

### Security Measures
- HTTPS only (TLS 1.3)
- Input validation
- SQL injection prevention
- Rate limiting (100 req/min per IP)
- Audit logging (all write operations)

## Scalability Strategy

### Horizontal Scaling
- Stateless API servers
- Load balancer (nginx/Railway)
- Shared cache (Redis)
- Database connection pooling

### Performance Targets
- API latency: p95 <200ms
- Database queries: p95 <50ms
- Cache hit rate: >80%
- Concurrent users: 100+

## Deployment Architecture

### Containerization
- Docker for consistency
- Multi-stage builds (optimization)
- Health checks built-in
- Non-root user for security

### Environments
- **Development:** Local Docker Compose
- **Staging:** Railway (testing)
- **Production:** Railway (HA setup)

---

**Version:** 1.0  
**Last Updated:** February 16, 2026  
**Author:** Dandi Sugiarto
