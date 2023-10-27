# american-state-backend


## Project setup
```
go build
```

## Dockerized MongoDB set up (Here we set the container name as state-mongo and run on port 27018)

```
docker run -d -p 27018:27017 --name state-mongo mongo:latest
```

## Project execute
```
go run server.go
```


