# URL Shortener API 

A simple, secure, and fast URL shortener backend built with **Go**, **Echo**, **MongoDB**, and **JWT authentication**.

##  Features

-  JWT-based user authentication
-  Shorten long URLs
- Click tracking (optional)
- MongoDB for persistent storage
- Lightweight & fast with Echo framework

---

## 🛠 Tech Stack

- **Language**: Go
- **Web Framework**: [Echo](https://echo.labstack.com/)
- **Database**: MongoDB
- **Auth**: JWT (JSON Web Tokens)

---

##  Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Thompson6626/Echo-Url-Shortener.git
   cd Echo-Url-Shortener
   ```

2. **Set up your environment variables**
   Create a `.env` file , follow the `.env.example`.

3. **Run the app**
   ```bash
   go mod tidy
   go run main.go
   ```

---

##  API Endpoints

###  Auth

- **POST** `api/v1/auth/register`  
  Register a new user  
  **Body**:
  ```json
  {
    "username": "exampleusername",
    "email": "user@example.com",
    "password": "yourpassword"
  }
  ```

- **POST** `api/v1/auth/login`  
  Login and receive JWT token  
  **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "yourpassword"
  }
  ```

---

### 🔗 URL Shortener

- **GET** `/api/v1/urls/:shortcode`  
  Redirect to the original URL

> All endpoints below require a Bearer token in the `Authorization` header.

- **POST** `/api/v1/urls/shorten`  
  Create a short URL  
  **Body**:
  ```json
  {
    "original_url": "https://example.com/very/long/url"
  }
  ```

- **GET** `/api/v1/urls/`  
  Get all URLs created by the authenticated user

- **DELETE** `/api/v1/urls/:shortCode`  
  Delete a shortened URL by ID



---

## 📁 Project Structure

```
.
├── cmd
│   └── api
│       ├── api.go
│       ├── auth.go
│       ├── errors.go
│       ├── health.go
│       ├── json.go
│       ├── main.go
│       ├── middleware.go
│       ├── urls.go
│       ├── users.go
│       └── utils.go
├── internal
│   ├── auth
│   │   ├── auth.go
│   │   └── jwt.go
│   ├── base62
│   │   └── base62.go
│   ├── database
│   │   ├── database.go
│   │   └── database_test.go
│   ├── env
│   │   └── env.go
│   ├── ratelimiter
│   │   ├── fixed-window.go
│   │   └── ratelimiter.go
│   └── store
│       ├── cache
│       │   ├── redis.go
│       │   ├── storage.go
│       │   └── users.go
│       ├── storage.go
│       ├── urls.go
│       └── users.go
```
