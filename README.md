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
   go mod init urlshortener
   go get github.com/gorilla/mux
   go get github.com/lib/pq
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
- [ ] create dockerfile
- [ ] create docker-compose
- [ ] setup k6
- [ ] current logging is bad, no log when request error or unhandled error
- [ ] Setup ci/cd
