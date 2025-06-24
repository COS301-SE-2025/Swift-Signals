## 📦 User Service

---

## 🚀 Getting Started

### ✅ Prerequisites

* [Go](https://golang.org/doc/install)
* [grpcurl](https://github.com/fullstorydev/grpcurl) (for testing)

---

## 🐘 Install & Start MongoDB

### 🧑💻 macOS (via Homebrew)

```bash
brew tap mongodb/brew
brew install mongodb-community@6.0
brew services start mongodb-community@6.0
```

### 🪟 Windows

### 🐧 Linux (Debian/Ubuntu)

---

## 🔧 Create User and Database

Open a terminal and `mongosh` shell and run:

```bash
use UserService
```

Then insert a document:

```bash
db.Users.insertOne({
  name: "Chris",
  email: "chris@test.com",
  password: "abc123",
  created_at: new Date()
})```
---

## 📁 Setup `.env`

Create a `.env` file by copying the example:

```bash
cp .env.example .env  # (use `copy` on Windows)
```

Edit `.env` to match your local settings:

```env
# For MongoDB Local
MONGO_URI=mongodb://localhost:27017

# For MongoDB cluster
# MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net

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

#### Check it worked
```bash
db.Users.find().pretty()
```

#### To exit mongosh
```bash
exit
```

---
///////////////////////////
CURRENTLY NOT WORKING BELOW
///////////////////////////

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

