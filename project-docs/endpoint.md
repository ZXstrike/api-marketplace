Here’s the updated and complete list of your REST API endpoints, categorized by folder (domain), now including **API Key management** as requested:

---

### 📁 `auth/` – Authentication & Authorization

| Method | Endpoint             | Description                |
| ------ | -------------------- | -------------------------- |
| POST   | `/auth/register` | Register new user (client) |
| POST   | `/auth/login`    | Login and get JWT token    |
| POST   | `/auth/refresh`  | Refresh access token       |

---

### 📁 `user/` – User Account Management

| Method | Endpoint             | Description                              |
| ------ | -------------------- | ---------------------------------------- |
| GET    | `/user/profile`  | Get current user profile (client/seller) |
| PUT    | `/user/profile`  | Update user profile                      |
| PUT    | `/user/password` | Change password                          |
| POST   | `/user/avatar`   | Upload or update profile picture         |

---

### 📁 `store/` – Store Management (Seller)

| Method | Endpoint               | Description                         |
| ------ | ---------------------- | ----------------------------------- |
| POST   | `/store`           | Create a store                      |
| GET    | `/store`           | Get your own store info             |
| PUT    | `/store`           | Update store profile                |
| GET    | `/store/:username` | Publicly get seller's store profile |

---

### 📁 `api/` – API Product Management (Seller)

| Method | Endpoint        | Description                     |
| ------ | --------------- | ------------------------------- |
| POST   | `/apis`     | Create new API product          |
| GET    | `/apis`     | List all APIs (public)          |
| GET    | `/apis/:id` | Get API product details         |
| PUT    | `/apis/:id` | Update own API product          |
| DELETE | `/apis/:id` | Delete own API product          |
| GET    | `/my-apis`  | List APIs of the current seller |

---

### 📁 `subscription/` – API Subscription Management

| Method | Endpoint                   | Description                          |
| ------ | -------------------------- | ------------------------------------ |
| POST   | `/subscriptions`       | Subscribe to an API                  |
| DELETE | `/subscriptions/:id`   | Unsubscribe from an API              |
| GET    | `/subscriptions`       | View current user's subscriptions    |
| GET    | `/subscribers/:api_id` | Seller: View subscribers for own API |

---

### 📁 `billing/` – Payments, Balance & Top-up

| Method | Endpoint               | Description                       |
| ------ | ---------------------- | --------------------------------- |
| POST   | `/billing/topup`   | Top-up balance                    |
| GET    | `/billing/history` | View top-up & transaction history |
| GET    | `/billing/balance` | Check current balance             |

---

### 📁 `analytics/` – Usage & Revenue Stats

| Method | Endpoint                 | Description                         |
| ------ | ------------------------ | ----------------------------------- |
| GET    | `/analytics/usage`   | View API usage logs (seller/client) |
| GET    | `/analytics/revenue` | View total revenue (seller)         |
| GET    | `/analytics/api/:id` | View detailed usage of specific API |

---

### 📁 `apikey/` – API Key Management (Client)

| Method | Endpoint                 | Description                       |
| ------ | ------------------------ | --------------------------------- |
| GET    | `/apikey`            | View current API key              |
| POST   | `/apikey/regenerate` | Regenerate API key                |
| DELETE | `/apikey`            | Revoke current API key (optional) |

---

