## 📦 User Service

---

## 🚀 Getting Started

### ✅ Prerequisites

* [Go](https://golang.org/doc/install)
* [grpcurl](https://github.com/fullstorydev/grpcurl) (for testing)

---

## 🐘 Install & Start PostgreSQL

### 🧑💻 macOS (via Homebrew)

```bash
brew install postgresql
brew services start postgresql
```

### 🪟 Windows

* Download installer: [https://www.postgresql.org/download/windows/](https://www.postgresql.org/download/windows/)
* Run the installer and remember the username, password, and port (default: 5432).
* Use the pgAdmin GUI or `psql` in the command line to create the database and user (see below).

### 🐧 Linux (Debian/Ubuntu)

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo service postgresql start
```

---

## 🔧 Create User and Database

Open a terminal and `psql` shell (See psqlREADME.md) and run:

```sql
CREATE USER user_service WITH PASSWORD 'password';
CREATE DATABASE user_service_db OWNER user_service;
```

Then apply the schema:

```bash
psql -U user_service -d user_service_db -f db/schema.sql
```

If you get a connection error, try:

```bash
psql -h localhost -p 5432 -U user_service -d user_service_db
```

---

## 📁 Setup `.env`

Create a `.env` file by copying the example:

```bash
cp .env.example .env  # (use `copy` on Windows)
```

Edit `.env` to match your local settings:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=user_service
DB_PASSWORD=password
DB_NAME=user_service_db
APP_PORT=50051
```

---

## 🏁 Run the Service

From the `user-service` directory:

```bash
go run ./cmd
```

You should see:

```
gRPC server running on :50051
```

---

## 💬 Test with `grpcurl`

> All requests use plaintext and connect to `localhost:50051`

---

### 👤 Register a User

```bash
grpcurl -plaintext -d '{
  "name": "exampleUser",
  "email": "exampleUser@test.com",
  "password": "abc123"
}' localhost:50051 user.UserService/RegisterUser
```

---

### 🔐 Login

```bash
grpcurl -plaintext -d '{
  "email": "exampleUser@test.com",
  "password": "abc123"
}' localhost:50051 user.UserService/LoginUser
```

Returns a token like:

```json
{
  "token": "mock-token-for-<user_id>"
}
```

---

### 🚪 Logout

```bash
grpcurl -plaintext -d '{
  "user_id": "<user_id>"
}' localhost:50051 user.UserService/LogoutUser
```

---

## 📁 Project Structure

```
user-service/
├── .env.example
├── cmd/                  # Entry point
├── db/                   # Database repo & schema
├── internal/             # gRPC handlers & business logic
├── models/               # Domain models
├── proto/                # .proto files
```

---
mongosh

use UserService

db.Users.insertOne({
  name: "Chris",
  email: "chris@test.com",
  password: "abc123",
  created_at: new Date()
})

db.Users.find().pretty()

exit
