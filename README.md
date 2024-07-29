## Usage:

```sh
docker-compose up -d
```
Create mercure schema.
Use ```cmd/mercure/create_updates.sql``` to create the table.

```
cd cmd/mercure/
go build
./mercure
```

### In browser:

[http://localhost:3001/.well-known/mercure?topic=https://example.com/foo](http://localhost:3001/.well-known/mercure?topic=https://example.com/foo)

### In terminal:

```sh
curl -d "topic=https://example.com/foo&data=Hello" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZXJjdXJlIjp7InB1Ymxpc2giOlsiKiJdfX0.iHLdpAEjX4BqCsHJEegxRmO-Y6sMxXwNATrQyRNt3GY" http://localhost:3001/.well-known/mercure
```