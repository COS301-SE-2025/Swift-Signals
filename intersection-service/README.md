# The Intersection Microservice
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
use IntersectionService
```

Then insert a document:

```bash
db.Users.insertOne({
})
```
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

APP_PORT=50052

```

---

## 🏁 Run the Service

From the `intersection-service` directory:

```bash
go run ./cmd
```

You should see:

```
gRPC server running on :50052
```

---

## 💬 Test with `grpcurl`

> All requests use plaintext and connect to `localhost:50052`

---

### 👤 Create an Intersection

```bash
grpcurl -plaintext -d '{
  "name": "Main Street & 1st Avenue",
  "details": {
    "address": "123 Main Street",
    "city": "Johannesburg",
    "province": "Gauteng"
  },
  "trafficDensity": "TRAFFIC_DENSITY_LOW",
  "defaultParameters": {
    "optimisationType": "OPTIMISATION_TYPE_GRIDSEARCH",
    "parameters": {
      "intersectionType": "INTERSECTION_TYPE_TRAFFICLIGHT",
      "green": 30,
      "yellow": 3,
      "red": 25,
      "speed": 50,
      "seed": 12345
    }
  }
}' localhost:50052 swiftsignals.intersection.IntersectionService/CreateIntersection


```

#### Check it worked
```bash
db.Intersections.find().pretty()
```

#### To exit mongosh
```bash
exit
```

---

