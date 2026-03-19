# Godis
Slim, learning project inspired by redis powered by Golang

## What is Godis
It's Redis clone written in Golang. Learning project for me to learn Resp protocol. Obviously not intended to be production ready.

### Build 
```bash
go build -o godis-server
```

- Connect to server using redis-cli
```
redis-cli -p 6379
```
## Supported Commands

- [x] GET
- [x] SET
- [x] DEL
- [ ] EXISTS
- [ ] EXPIRE
- [ ] TTL
- [x] HSET
- [x] HGET
- [x] HGETALL
- [ ] HDEL
- [ ] HLEN
- [ ] HKEYS
- [ ] HVALS

## Data Persistence
Currently Godis uses Append Only File to persist state.