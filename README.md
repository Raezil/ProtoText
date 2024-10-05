# protobuf-go

## Sample curl requests
Register:
```
curl --http2 -X POST http://localhost:8080/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
           "email": "newuser@example.com",
           "password": "password123",
           "name": "John",
           "surname": "Doe",
           "age": 30
         }'
```
Log in
```
curl --http2 -X POST http://localhost:8080/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{
           "email": "user@example.com",
           "password": "password123"
         }'

```

Protected request
```
curl --http2 -X POST http://localhost:8080/v1/auth/protected      -H "Content-Type: application/json"      -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im5ld3VzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3MjgyMTc0ODZ9.1HkqukYROgGsxa6HHdjAXADzKa6cyKKir5gk7TwK0Og"      -d '{
           "text": "This is a protected request."
         }'
```

