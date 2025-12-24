# Chirpy

Chirpy is social media backend server written in **Go.** It features a RESTful API for managing users and "chirps"(posts), complete with a authentication system using **JWT Access Token** and **Refresh Tokens**.

## Features
- **User Managment:** Secure user registration and login with bcrypt password hashing.
- **Chirps:** Create, retrieve and filter short messages.
- **Hybrid Authentication**
  - **Access Tokens:** Stateless JWTs for secure, short session access (1 hour expiry).
  - **Refresh Tokens:** Database-backend tokens for long sessions (60 day expiry) with revocation support.
- **Filtering:** Support for query parameters to filter chirps by author.
- **Clean Architecture:** Organized into packages for auth, database and internal handlers.

## Tech Stack 
- **Language:** Golang
- **Database:** PostgresSQL
- **SQL Tooling** SQLC
- **Authentication:** JWT and Custom Refresh Tokens.

## API Endpoints
**Authentication**
> `POST`  `/api/login` =  Login and receive access + Refresh Tokens
> `POST` `/api/refresh` = User a Refresh Token to get a new Access Token
> `POST` `/api/revoke` = Revoke a Refresh Token (Logout)

**Users**
> `POST` `/api/users` = Register a new user
> `PUT` `/api/users` = Update user details (requires auth permissions)

**Chirps**
> `GET` `/api/chirps` = List all chirps (Optional `?author_id={author_id}` and `?sort=desc filter)
> `POST` `/api/chips` = Create a new chirp (requires auth)
> `GET` `/api/chirps/{chirpID}` = Get chirp by id

## Setup

inside your terminal

**1. Clone the repository:**


`git clone https://github.com/jaxhemopo/Chirpy.git
cd Chirpy`

**2. Set up environment variables: Create a .env file in the root directory:**

`DB_URL=postgres://user:password@localhost:5432/chirpy
JWT_SECRET=your_super_secret_key
PLATFORM=development`

**3. Run database migrations:**

`
cd sql/schema
goose postgres <your_db_url> up
`

**4. Build and Run:**

`go build -o chirpy && ./chirp`

