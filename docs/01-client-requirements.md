# Client Requirements Document

## Executive Summary

**Client:** PT NetCom Solutions (Fictional ISP Case Study)  
**Project:** Network Operations Integration API  
**Timeline:** 4-6 weeks MVP  
**Problem:** Manual network device management causing operational delays

## Business Problem

### Current Pain Points

#### 1. Manual Device Configuration (High Impact)
- **Time Cost:** 5-10 minutes per device × 80 devices = 6-13 hours
- **Risk:** Human error, inconsistent configurations
- **Impact:** Cannot perform updates during business hours

#### 2. No Business Context in Monitoring
```
Current workflow:
Zabbix Alert: "Router-BDG-01 DOWN"
↓
Engineer manually checks:
- Which customers affected? (Excel lookup)
- SLA tier? (Another spreadsheet)
- When will SLA breach? (Manual calculation)
- Who to notify? (Contact list)

Time: 15-20 minutes before customer notification
```

#### 3. Disconnected Systems
- Monitoring: Zabbix
- Ticketing: HubSpot  
- Customer DB: Google Sheets
- Communication: Manual calls/WhatsApp

**Result:** Data silos, duplicate entry, slow response

## Solution Requirements

### Functional Requirements

**FR-1: Centralized Device Management**
- Bulk configuration across multiple devices
- Firmware update orchestration
- No VPN required (push model)
- Real-time execution tracking

**FR-2: Monitoring Integration** 
- Consume Zabbix webhooks
- Enrich with business context
- Auto-create tickets with customer impact

**FR-3: Business Logic Layer**
- Device-to-customer mapping
- SLA compliance tracking
- Automatic severity calculation
- Impact analysis

**FR-4: Workflow Automation**
- Auto-ticket creation
- Engineer assignment (location-based)
- Customer notifications (email, SMS)
- Incident lifecycle tracking

### Non-Functional Requirements

**Performance:**
- API response: <200ms (p95)
- Support 100 concurrent requests
- Handle 1000+ devices

**Scalability:**
- Horizontal scaling capability
- Queue-based bulk operations
- Redis caching

**Security:**
- API key authentication (agents)
- JWT authentication (users)
- Rate limiting
- Audit logging

**Reliability:**
- 99.9% uptime
- Automatic failover
- Daily backups
- Graceful degradation

## Success Metrics

### Quantitative
- Incident response: <1 minute (vs 15-20 minutes)
- Bulk config deployment: <5 minutes for 80 devices
- Cost savings: >Rp 50M/year operational efficiency

### Qualitative
- Improved NOC team workflow satisfaction
- Reduced customer complaints
- Higher SLA achievement rate

## Technology Stack

**Backend:** Golang (performance, concurrency)  
**Database:** PostgreSQL (ACID compliance)  
**Cache:** Redis (sub-millisecond latency)  
**Deployment:** Docker (consistency, isolation)

## Out of Scope (Phase 2)

- Mobile app for field engineers
- ML-based predictive maintenance
- Customer self-service portal
- Multi-vendor device support

---

**Status:** Approved  
**Last Updated:** February 16, 2026  
**Document Owner:** Dandi Sugiarto
