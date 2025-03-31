# Distributed-Cache

A distributed cache system in Golang that leverages single-node concurrent caching and HTTP-based multi-node caching.

## Features
- LRU cache eviction strategy
- Golang's native locking mechanisms
- Consistent hashing mechanism
- Protocol Buffers (protobuf) for inter-node communication

## Directory Structure
```
Distributed-Cache/
    |--lru/
        |--lru.go  
        |--lru_test.go
    |--consistentHash/
        |--consistentHash.go
        |--consistentHash_test.go
    |--byteView.go 
    |--cache.go    
    |--geecache.go 
    |--http.go
    |--peers.go
    |--cacheFlex_test.go
    |--go.mod
    |--go.sum
    |--main.go
    |--README.md
    |--Dockerfile
    |--.gitignore
    |--LICENSE
```

## Getting Started

### Prerequisites
- Go 1.16 or higher
- Docker (optional)

### Running Tests
```bash
go test
```

### Building with Docker
```bash
docker build -t distributed-cache .
docker run -p 8080:8080 distributed-cache
```

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Created by Rudra