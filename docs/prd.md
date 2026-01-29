# MarketLens — Product Requirements Document (PRD)

## 1. Product Overview

**Product Name:** MarketLens  
**Positioning:** Hybrid (public-good + commercial extensibility)  
**Target Region (Pilot):** Nigeria  
**Primary Users:** Smallholder farmers (feature phones)  
**Secondary Users:** NGOs, cooperatives, analysts, planners  

MarketLens is a market price intelligence platform that provides
trusted, up-to-date crop prices through **USSD for farmers** and a
**web dashboard for analysts**, using a shared data pipeline.

---

## 2. Problem Statement

Smallholder farmers often sell produce without access to:
- current prices in nearby markets
- historical price trends
- confidence in the accuracy of price information

This information asymmetry reduces farmer income, increases
post-harvest losses, and strengthens dependency on intermediaries.

---

## 3. Product Goals (MVP)

### Primary Goals
1. Provide **reliable daily crop prices** to farmers via USSD
2. Increase **price transparency** across markets
3. Build **trust** through confidence scoring and auditability

### Secondary Goals
4. Enable NGOs and planners to analyse market trends
5. Demonstrate inclusive, emerging-market-ready system design
6. Serve as a deployable, portfolio-grade open-source product

---

## 4. Target Users & Use Cases

### 4.1 Farmer (USSD User)
- Owns a feature phone
- Limited internet or literacy
- Needs quick, actionable price information

**Key Use Case**
> “What is today’s price for my crop in a nearby market?”

---

### 4.2 Analyst / NGO / Aggregator (Web User)
- Works with farmer groups or markets
- Needs validated price data and trends
- Requires transparency into data quality

**Key Use Case**
> “Which markets show abnormal price volatility this week?”

---

## 5. MVP Scope (What We Will Build)

### 5.1 Supported Crops (Pilot)
- Maize
- Rice
- Cassava
- Tomato

### 5.2 Supported Markets (Pilot)
- Mile 12 (Lagos)
- Bodija (Ibadan)
- Wuse (Abuja)

*(Lists are configurable and expandable)*

---

### 5.3 USSD Features (Farmer-Facing)

USSD is **read-only** in MVP.

#### Required Features
- Select crop
- Select market
- View:
  - current price
  - price trend (↑ ↓ →)
  - confidence level
  - last update timestamp

#### Non-Functional Requirements
- Response time < 3 seconds
- Works on basic feature phones
- Stateless requests with session continuity

---

### 5.4 Web Dashboard Features (Admin/Analyst)

#### Required Features
- View aggregated prices by crop & market
- Daily and weekly price trends
- Confidence level indicators
- Outlier/anomaly flags
- Manual override with justification
- Audit log of all admin actions

---

## 6. Data & Intelligence Requirements

### 6.1 Data Ingestion
- Manual submissions from enumerators
- Partner/NGO uploads (future)
- API-ready ingestion interface

All raw data must be preserved.

---

### 6.2 Price Aggregation Logic

Prices are published only after aggregation.

**Rules**
- Median-based aggregation
- Configurable outlier threshold
- Time-window validity (e.g. last 24 hours)
- Minimum submission count for publication

Each published price includes:
- aggregated value
- trend indicator
- confidence score (low / medium / high)
- source count
- last updated timestamp

---

### 6.3 Confidence Scoring (MVP Logic)

| Condition                          | Confidence |
|-----------------------------------|------------|
| Few sources / high variance       | Low        |
| Moderate sources / stable prices  | Medium     |
| Many sources / low variance       | High       |

---

## 7. Non-Goals (Explicitly Out of Scope for MVP)

The following are **not** part of MVP:
- Payments or financial transactions
- Logistics booking
- Farmer profiles or personal data
- Machine learning models
- Messaging or chat features
- Multi-country support

These may be considered post-pilot.

---

## 8. Compliance & Ethics

- No personally identifiable farmer data stored
- Transparent data corrections with audit trails
- Explainable pricing logic (no black-box models)
- Designed for public-sector and NGO trust

---

## 9. Success Metrics (Pilot)

### Farmer Impact (Simulated or Real)
- % of queries returning valid prices
- Reduction in price uncertainty
- Adoption rate of USSD flow

### System Health
- Average USSD response time
- Data freshness (% prices < 24h old)
- Outlier detection rate

---

## 10. Risks & Mitigations

| Risk | Mitigation |
|----|-----------|
| Inaccurate data | Confidence scoring + admin review |
| Low adoption | Simple USSD flow |
| Abuse/spam | Read-only USSD in MVP |
| Scope creep | Frozen non-goals |

---

## 11. Acceptance Criteria (MVP Complete When…)

- Farmers can retrieve daily prices via USSD
- Prices include trend and confidence indicators
- Admin users can review and override anomalies
- All data changes are auditable
- System runs locally via Docker Compose
- Documentation allows a new contributor to onboard

---

## 12. Guiding Principle

> **If a farmer cannot use it on a basic phone, it is not MVP-ready.**
