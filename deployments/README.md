# Getting Started with Docker Compose

Run the services:
```bash
cd deployments
docker compose up -d
```

Stop the services:
```bash
docker compose down
```

Stop the services and remove volumes:
```bash
docker compose down -v
```

Note since docker-compose.yml is for development to allow testing locally
```yml
ports:
      - "50051:50051"
```
and
```yml
ports:
      - "5432:5432"
```
are being exposed to localhost but will remove for production. 

Testing user-service:
```bash
grpcurl -plaintext -d '{                   
  "name": "exampleUser1",
  "email": "exampleUser1@test.com",
  "password": "abcd1234"
}' localhost:50051 swiftsignals.user.UserService/RegisterUser
```

Testing intersection-service:
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

```bash
grpcurl -plaintext -d '{}' localhost:50052 swiftsignals.intersection.IntersectionService/GetAllIntersections
```

