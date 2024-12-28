## Overview

`shrunk` is a simple url shortener made with Go and postgreSQL with very simple setup and functionality

Functionalities:
- user is anonymous
- user send long url, `shrunk` shorten it to 7 characters randomly (can't create custom url)
- user access the shortener url, `shrunk` redirect it to original long url


## Database Setup

Create a PostgreSQL database named `shrunk` with the following schema:

```sql
CREATE TABLE urls (
    slug TEXT PRIMARY KEY,
    original TEXT NOT NULL,
    clicked INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Set up the PostgreSQL database and create the required table
4. Run the service:
   ```bash
   go run main.go
   ```


## API Endpoints

### Create Shortened URL
- **URL**: `/`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "url": "https://example.com/very/long/url"
  }
  ```
- **Response**: Text containing the generated slug
- **Status Codes**:
  - 200: Success
  - 400: Invalid request body
  - 500: Server error
- **Sample Curl**:
  ```bash
  curl -X POST \
    http://localhost:3002/ \
    -H 'Content-Type: application/json' \
    -d '{"url": "https://example.com/very/long/url"}'
  ```
- **Sample Response**:
  ```
  Mzk4NTY
  ```

### Access Shortened URL
- **URL**: `/{slug}`
- **Method**: `GET`
- **Response**: Redirects to the original URL
- **Status Codes**:
  - 302: Successful redirect
  - 404: Slug not found
  - 405: Method not allowed
  - 500: Server error
- **Sample Curl**:
  ```bash
  # Follow redirects with -L flag
  curl -L http://localhost:3002/Mzk4NTY

  # Just show headers without following redirect
  curl -I http://localhost:3002/Mzk4NTY
  ```
- **Sample Response Header**:
  ```
  HTTP/1.1 302 Found
  Location: https://example.com/very/long/url
  Date: Wed, 13 Nov 2024 10:00:00 GMT
  ```

## Motivation

The motivation is actually learn more on dockerize and load testing aim for high availability

## Todo
- [x] create dockerfile
- [x] create docker-compose
- [ ] setup k6
- [ ] current logging is bad, no log when request error or unhandled error
- [ ] Setup ci/cd


## How to Run (without docker-compose)
Application and database must be run in same network
1. create network by `docker network create shrunk` (`shrunk` is network name, you can change it)
2. run postgres container by
```
  docker run -d \
  -p 5433:5432 \
  --name shrunk-postgres \
  -e POSTGRES_PASSWORD=rahasia \
  -e POSTGRES_DB=shrunk \
  -v $(pwd)/init.sql:/docker-entrypoint-initdb.d/init.sql \
  --network shrunk \
  postgres
```
[Reference](https://hub.docker.com/_/postgres#:~:text=start%20a%20postgres%20instance)

`--network shrunk` means run the container in network `shrunk`

`-v $(pwd)/init.sql:/docker-entrypoint-initdb.d/init.sql \` mounts local `init.sql` file into the container's /docker-entrypoint-initdb.d/ directory to automatically initialize the database on first run. ps: only the first time the container is started

3. run application container by build the image first and run the container

build image
```
docker build -t shrunk:1.2 .
```

run container
```
docker run -d \
-p 3002:3002 \
--name shrunk \
--network shrunk \
shrunk:1.2
```

4. send request
```
curl -X POST http://localhost:3002/ \
-H "Content-Type: application/json" \
-d '{"url": "https://example.com"}'
```

## How to Run (with compose)
1. change directory into root folder
```
cd shrunk
```

2. run
```
docker-compose up -d
```

3. send request
```
curl -X POST http://localhost:3002/ \
-H "Content-Type: application/json" \
-d '{"url": "https://example.com"}'
```