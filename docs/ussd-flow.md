# USSD Flow – Marketlens (MVP)

This document defines the **USSD menu tree**, **session state machine**, and **response formats** for the Market Price Tracker. It is written to be directly usable by backend engineers implementing Africa’s Talking–compatible USSD services.

---

## 1. Design Principles

- **Minimal depth**: Reach price in ≤ 4 screens
- **Numbered menus only** (feature-phone friendly)
- **Fail fast, recover gracefully** on invalid input
- **Stateless gateway, stateful backend**
- **Explicit END on value delivery** (price screen)

---

## 2. Menu Tree (MVP)

```
START
 └── MAIN_MENU
      ├── 1. Latest price
      │     └── SELECT_CROP
      │           └── SELECT_MARKET
      │                 └── SHOW_PRICE (END)
      ├── 2. Help
      │     └── SHOW_HELP (END)
      └── 3. Exit
            └── END
```

---

## 3. Session State Machine

### 3.1 States

- **START**: New session initiated by gateway
- **MAIN_MENU**: User chooses high-level action
- **SELECT_CROP**: User selects crop
- **SELECT_MARKET**: User selects market
- **SHOW_PRICE**: System displays price and terminates
- **ERROR**: Invalid input handler
- **END**: Terminal state

---

### 3.2 State Transitions

```
START -> MAIN_MENU

MAIN_MENU -> SELECT_CROP       (input = 1)
MAIN_MENU -> SHOW_HELP         (input = 2)
MAIN_MENU -> END               (input = 3)

SELECT_CROP -> SELECT_MARKET   (valid crop choice)

SELECT_MARKET -> SHOW_PRICE    (valid market choice)

SHOW_PRICE -> END
SHOW_HELP  -> END

Any invalid input -> ERROR
ERROR -> same state OR MAIN_MENU (after N retries)
```

---

### 3.3 Error Handling Rules

- Each state allows **max 2 invalid attempts**
- On 3rd invalid attempt:
  - Reset to MAIN_MENU
- ERROR state does **not** persist user intent

---

## 4. Session Storage

### 4.1 Redis Key

```
ussd:session:{sessionId}
```

Fallback (if no sessionId):
```
ussd:session:{phoneNumber}
```

---

### 4.2 Stored Fields

| Field | Description |
|-----|------------|
| state | Current session state |
| last_input | Last user input |
| chosen_crop_id | Selected crop |
| chosen_market_id | Selected market |
| tries | Invalid input counter |

---

### 4.3 TTL

- **2–5 minutes** (gateway-dependent)
- Session auto-expires on inactivity

---

## 5. Response Format (Africa’s Talking)

### 5.1 Rules

- Responses must be **plain text**
- Must begin with:
  - `CON ` → continue session
  - `END ` → terminate session
- No trailing spaces

---

## 6. Screen Templates

### 6.1 MAIN_MENU

```
CON Welcome to MarketLens
1. Latest price
2. Help
3. Exit
```

---

### 6.2 SELECT_CROP

```
CON Select crop:
1. Maize
2. Rice
3. Tomato
```

(Pagination rule – future):
```
9. Next
```

---

### 6.3 SELECT_MARKET

```
CON Select market:
1. Mile 12 (Lagos)
2. Bodija (Ibadan)
3. Wuse (Abuja)
```

---

### 6.4 SHOW_PRICE

```
END {crop} @ {market}
Price: ₦{min_price} – ₦{max_price} per {unit}
Trend: {trend}
Confidence: {confidence}
Updated: {timestamp}
```

Example:
```
END Maize @ Mile 12 (Lagos)
₦115,000 – ₦125,000 per bag
Trend: ↑ Rising
Confidence: High
Updated: Today 10:30am
```

---

### 6.5 SHOW_HELP

```
END MarketLens helps you check daily crop prices.
Prices are indicative and updated daily.
Dial this code anytime to check again.
```

---

### 6.6 ERROR (Generic)

```
CON Invalid choice. Please try again.
```

(Then re-render previous state menu)

---

## 7. Notes for Backend Implementation

- USSD gateway sends **entire input string** (e.g. `1*2*3`)
- Backend must parse **last token only**
- State is authoritative in Redis, not derived from input depth
- Always prefer `END` once value is delivered

---

## 8. MVP Guardrails

- No authentication
- No language selection
- No free-text input
- Numeric choices only

---

This flow is considered **final for MVP**. Any additions (language, location auto-detect, farm-gate reporting) belong to post-MVP phases.
