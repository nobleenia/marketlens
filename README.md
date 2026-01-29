# MarketLens 🌾📍

**MarketLens** is an open-source market price intelligence platform designed for
smallholder farmers in emerging markets.

It provides **real-time crop price visibility via USSD** for feature phones,
and a **web dashboard** for analysts, NGOs, and aggregators—using the same
trusted data pipeline.

---

## 🎯 Problem

Smallholder farmers often sell crops without knowing:
- current prices in nearby markets
- price trends over time
- whether transport to another market is worth it

This information gap reduces income, increases post-harvest losses, and
reinforces dependence on middlemen.

---

## 💡 Solution

MarketLens delivers **the same price intelligence through two channels**:

- **USSD** – for farmers using basic phones  
- **Web dashboard** – for analysts, NGOs, cooperatives, and planners

Data is validated, aggregated, and scored for confidence before publication.

---

## 🧱 System Overview
```
Market Enumerators / APIs  
        ↓  
Data Ingestion API  
        ↓  
Price Validation & Aggregation  
        ↓  
Central Price Database  
        ↓  
┌──────────────┬──────────────┐
| USSD Gateway | Web Dashboard|  
└──────────────┴──────────────┘
```


---

## 📱 USSD Features (Farmer-Facing)

- Check crop prices by market
- View daily price trend (↑ ↓ →)
- Confidence indicator (low / medium / high)
- Last update timestamp

**Example**
```
347800#
1. Check Crop Price

2. Nearby Markets

3. Weather Alert

→ Tomato (Mile 12): ₦18,500
→ Trend: ↑ Rising
→ Confidence: High
```


---

## 📊 Web Dashboard Features

- Market price tables and trends
- Outlier and anomaly detection
- Manual admin review & overrides
- Confidence score visibility
- Exportable datasets

---

## 🛠 Tech Stack

**Backend**
- Go (REST API)
- PostgreSQL (data storage)
- Redis (USSD session handling)

**Frontend**
- React (admin dashboard)

**Infrastructure**
- Docker & Docker Compose
- GitHub Actions (CI)

---

## 🗂 Repository Structure
```
marketlens/
    apps/       # Web applications
    services/   # Backend services (Go)
    docs/       # Architecture, PRD, ADRs
    infra/      # Docker, deployment
    brand/      # Logo & colours
```

---

## 🚀 Getting Started (Local)

```bash
git clone https://github.com/marketlens/marketlens.git
cd marketlens
docker-compose up
```
API will be available at http://localhost:8080

---

## 📄 Documentation


- Product requirements: docs/prd.md

- Architecture: docs/architecture.md

- USSD flow: docs/ussd-flow.md

- API contract: docs/api/openapi.yaml

---

## 🧩 Contributing

We welcome contributors—especially from agriculture, logistics,
and emerging-market tech backgrounds.

See CONTRIBUTING.md
 to get started.

 ---

## 📜 License

MIT License