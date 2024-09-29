# url-shortener-backend

## Overview

Backend project for an URL shortener web app (zURL).

## Features

- User system
- Shareable url shortcodes
- Built with performance in mind

## Tech stack

- Golang
- Iris framework
- Redis for storage

## Prerequisites

**Golang** >= 2 and a **Redis cache**

## Usage

1. Define an **.env** file in the root directory (next to all `.go` files) with the following parameters:
   - Some are used for encryption: if you change them in any moment, data will be unavailable for access

```dotenv
# Redis storage connection string (address and port)
REDIS_URI=[REDIS_CACHE_ADDRESS]:[PORT]

# Used to pepper the password before storing
PEPPER=[PEPPER]

# Generate with your favorite encryption algorithm
JWT_SECRET=[JWT_SECRET]
```

2. Execute `go build` or `go run` on `main.go` entry file

3. Use the API endpoints at `localhost:8080`. They're listed in the `handler.go` file

## Project status

In development:

- [x] User system (authentication, register, login, edit user information)
- [ ] Session management
- [ ] URL module (endpoints, middleware, logic, etc.)
- [ ] Front end project
