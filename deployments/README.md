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

Testing:
```bash
grpcurl -plaintext -d '{                   
  "name": "exampleUser1",
  "email": "exampleUser1@test.com",
  "password": "abcd1234"
}' localhost:50051 swiftsignals.user.UserService/RegisterUser
```


