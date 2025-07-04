# Getting Started with Docker Compose

Run the services:
```bash
cd deployments
docker compose up -d --build
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

```bash
grpcurl -plaintext -d '{}' localhost:50052 swiftsignals.intersection.IntersectionService/GetAllIntersections
```

Testing api-gateway:
```bash
curl localhost:9090/register -d '{"email": "ham@ferrari.com", "password": "1234abcd", "username":"lh44"}'
```

```bash
curl localhost:9090/login -d '{"email": "ham@ferrari.com", "password": "1234abcd"}'                                           
```

```bash
curl localhost:9090/intersections -H "Authorization: Bearer fjaskdf;a" -d '{                                           
  "default_parameters": {
    "green": 10,
    "intersection_type": "t-junction",
    "red": 6,
    "seed": 3247128304,
    "speed": 60,
    "yellow": 2
  },
  "details": {
    "address": "Corner of Foo and Bar",
    "city": "Pretoria",
    "province": "Gauteng"
  },
  "name": "My Intersection",
  "traffic_density": "high"
}'
```

```bash
curl localhost:9090/intersections -H "Authorization: Bearer fjaskdf;a" 
```
