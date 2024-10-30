# ProtoText

## Kubernetes

- Apply kubectl
```
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```
- Port Foward
```
kubectport-forward service/app-service 8080:808080 -n default
```
- Check pods
```
kubectl get pods -n default
kubectl describe pod $NAME -n default
```

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
curl --http2 -X POST http://localhost:8080/v1/auth/protected      -H "Content-Type: application/json"      -H "Authorization: $TOKEN"      -d '{
           "text": "This is a protected request."
         }'
```

## Run Docker
```
docker build -t app8:latest .
docker run -d -p 50051:50051 -p 8080:8080 --name app8 app8:latest
```
