# REST API Endpoints

Here’s the updated and complete list of your REST API endpoints, categorized by domain.

---

### 📁 `auth/` – Authentication & Authorization

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/auth/register` | Register a new user account. |
| `POST` | `/auth/login` | Log in to an existing account and receive a JWT token. |
| `POST` | `/auth/refresh` | Refresh an expired access token using a refresh token. |

---

### 📁 `user/` – User Account Management

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/user/{user_id}` | Retrieve a user's profile by ID. |
| `POST` | `/user/update` | Update the profile of the currently authenticated user. |
| `PUT` | `/user/change-password` | Change the password for the currently authenticated user. |
| `POST` | `/user/update-profile-picture` | Upload or update the user's profile picture. |
| `GET` | `/user/profile-picture/{image_id}` | Retrieve a user's profile picture by image ID. |

---

### 📁 `store/` – Store Management

*These endpoints are for seller accounts.*

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/store/create` | Create a new store. |
| `PUT` | `/store/update` | Update your store's profile information. |
| `GET` | `/store/username/{username}` | Retrieve a seller's store profile by username. |
| `GET` | `/store/all` | Get a list of all available stores. |

---

### 📁 `api/` – API Product Management

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/create` | Create a new API product (for sellers). |
| `GET` | `/api/all` | Get a public list of all available APIs. |
| `GET` | `/api/{id}` | Retrieve the details of a specific API product. |
| `PUT` | `/api/update/{id}` | Update an API product that you own. |
| `DELETE` | `/api/delete/{id}` | Delete an API product that you own. |
| `POST` | `/api/create-endpoint` | Define new endpoints for an API version you own. |
| `GET` | `/api/api-endpoints/{api_version_id}` | Get defined endpoints for a specific API version. |
| `PUT` | `/api/update-endpoint` | Update existing endpoints for an API version you own. |
| `DELETE` | `/api/delete-endpoint/{id}` | Delete a specific endpoint from an API version you own. |

---

### 📁 `subscription/` – API Subscription Management

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/subscriptions/subscribe` | Subscribe to an API plan. |
| `POST` | `/subscriptions/unsubscribe` | Unsubscribe from an API plan (requires `subscription_id`). |
| `GET` | `/subscriptions/get-by-user` | View all of the current user's active subscriptions. |
| `GET` | `/subscriptions/get` | Retrieve details of a specific subscription (using `subscriptionID` query parameter). |

---

### 📁 `api-keys/` – API Key Management

*These endpoints are for client (buyer) accounts to manage their access.*

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api-keys/create` | Create a new API key for a specific subscription. |
| `DELETE` | `/api-keys/delete` | Revoke an API key. |
