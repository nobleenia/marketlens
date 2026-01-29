# MarketLens — System Architecture

## 1. Architecture Goals

MarketLens is designed to:

- Operate reliably in low-connectivity environments
- Serve both low-tech (USSD) and high-tech (web) users
- Maintain data trust through validation and aggregation
- Remain simple enough for small teams and open-source contributors

The MVP prioritises **clarity, robustness, and auditability** over scale.

---

## 2. High-Level Architecture

MarketLens uses a **single Go backend (monolith)** in the MVP phase.

```
Market Enumerators / APIs
        ↓
Data Ingestion Layer
        ↓
Price Validation & Aggregation
        ↓
Central PostgreSQL Database
        ↓
┌──────────────┬──────────────┐
| USSD Handler | Web API      |
└──────────────┴──────────────┘
        ↓
Farmers (USSD) Analysts / NGOs (Web)
```


---

## 3. Core Components

### 3.1 Data Ingestion

**Sources**
- Market enumerators (manual entry)
- Partner organisations / NGOs
- Pilot APIs or scrapers (future)

**Responsibilities**
- Validate incoming price submissions
- Normalise units and currencies
- Store raw observations immutably

All ingested data is preserved for traceability.

---

### 3.2 Price Validation & Aggregation Engine

Raw price submissions are processed using:

- Median-based aggregation
- Outlier detection (configurable thresholds)
- Source count weighting
- Time-windowed validity (e.g. last 24h)

Each published price includes:
- Aggregated value
- Trend indicator
- Confidence level (low / medium / high)
- Last updated timestamp

This layer is the **core intelligence of MarketLens**.

---

### 3.3 Central Database (PostgreSQL)

Primary tables include:
- `price_observations`
- `aggregated_prices`
- `markets`
- `crops`
- `sources`
- `users`
- `audit_logs`

Design principles:
- Append-only raw data
- Explicit audit trails
- No silent overwrites

---

### 3.4 USSD Handler

The USSD interface is designed for:

- Stateless request handling
- Redis-backed session persistence
- Minimal text payloads
- Fast response times (< 3 seconds)

USSD is **read-only in MVP** to ensure reliability and prevent abuse.

---

### 3.5 Web API & Admin Dashboard

The web API exposes:
- Aggregated prices
- Market trends
- Anomaly flags
- Confidence metadata

The admin dashboard supports:
- Manual review of anomalies
- Data corrections with justification
- Monitoring of data freshness

All admin actions are logged.

---

## 4. Technology Choices

| Layer        | Technology     | Rationale |
|-------------|----------------|-----------|
| Backend     | Go              | Performance, simplicity, strong typing |
| Database    | PostgreSQL      | Relational integrity, analytics-ready |
| Cache       | Redis           | USSD session management |
| Frontend    | React           | Lightweight admin tooling |
| CI          | GitHub Actions  | Open-source friendly |
| Containers  | Docker          | Reproducible environments |

---

## 5. Architectural Decisions (ADR Summary)

- **Monolith over microservices:** reduces cognitive load and ops overhead
- **USSD-first:** ensures inclusion of non-smartphone users
- **Confidence scoring:** increases trust and transparency
- **Human-in-the-loop:** prevents blind automation in fragile data environments

Formal ADRs are stored in `docs/adr/`.

---

## 6. Security & Reliability Considerations

- Input validation on all ingestion endpoints
- Rate limiting on USSD endpoints
- Role-based access control for admin features
- No personal farmer data stored in MVP
- Daily database backups (deployment-dependent)

---

## 7. Future Evolution (Non-MVP)

Explicitly out of scope for MVP:
- Payments or financial transactions
- Logistics booking
- Machine learning models
- User-to-user communication

These may be introduced post-pilot.

---

## 8. Guiding Principle

> **Clarity beats cleverness. Trust beats scale.**
