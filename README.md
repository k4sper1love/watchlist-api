# Watchlist API 🎬

**Watchlist API** is a REST API that allows you to save and manage a list of films you want to watch. It includes features for registration, managing film collections, adding comments, and leaving feedback.

## 🔴 Live Server
- **API**: [https://k4sper1love.kz/api](https://k4sper1love.kz/api)
- **Swagger Documentation**: [https://k4sper1love.kz/swagger/index.html](https://k4sper1love.kz/swagger/index.html)
- **Grafana Dashboards**: 
  - [https://k4sper1love.kz/logs](https://k4sper1love.kz/logs)
  - [https://k4sper1love.kz/monitoring](https://k4sper1love.kz/monitoring)

## 🔎 Navigation
- [Main Features](#-main-features)
  - [API Functionality](#api-functionality)
  - [Authorization](#authorization)
  - [Additional features](#additional-features)
- [Technology Stack](#-technology-stack)
- [Project Requirements](#-project-requirements)
- [GitHub Actions (CI/CD)](#-github-actions-cicd)
- [Installation and Setup](#%EF%B8%8F-installation-and-setup)
- [Configuration](#%EF%B8%8F-configuration)
- [Run the application](#-run-the-application)
  - [Using Terminal](#using-terminal)
  - [Using Docker Compose](#using-docker-compose)
- [API Documentation](#-api-documentation)
- [Account Authentication](#-account-authentication)
  - [Using Credentials](#using-credentials)
  - [Via Telegram Bot](#via-telegram-bot)
- [Testing with Postman](#-testing-with-postman)
  - [Postman tests](#postman-tests)
  - [Running tests](#running-tests)
- [Monitoring System (Grafana)](#-monitoring-system-grafana)
- [Watchlist REST API Endpoints](#-watchlist-rest-api-endpoints)
- [Database Structure](#-database-structure)
- [License](#-license)
- [Contact](#-contact)

## ⭐ Main Features
### API Functionality
- **Registration and Authorization**: Secure user registration and login with JWT.
- **Collections Management**: Create and manage film collections.
- **Comments**: Add and manage comments on films.
- **Viewing Status**: Mark films as viewed.
- **Ratings and Reviews**: Rate films and write reviews.

### Authorization
Authorization in Watchlist API is handled using **JWT (JSON Web Token)** in the `Authorization` header. This ensures secure access to the API by verifying the token sent with each request.
```bash
Authorization: Bearer <JWT_TOKEN>
```
### Additional Features
- **Permissions**: Flexible permission system to control access to different IP endpoints based on permissions.
- **Validator**: Automatic request validation to ensure incoming data is properly formatted and meets required conditions before processing.
- **Filters**: Filtering options for API requests to allow users to filter films, collections, and other resources based on specific criteria.

## 🚀 Technology Stack
- **Programming Language**: Go
- **Authorization**: JWT
- **Database**: PostgreSQL
- **API Documentation**: Swagger
- **Log Aggregation**: Loki
- **Metrics**: Prometheus + Node Exporter
- **Visualization**: Grafana (for logs and metrics)
- **Testing**: Postman
- **Deployment**: Docker, Docker Compose
- **CI/CD**: GitHub Actions

## 📝 Project Requirements
- **Go**: 1.18+
- **PostgreSQL**: 13+
### For Deployment
- **Docker**: 20.10+
- **Docker Compose**: 1.29+
- **Loki Docker Driver**
### For Testing
- **Postman**

## 🔄 GitHub Actions (CI/CD)
The project uses **GitHub Actions** to automate testing, building, and deployment processes on remote server.

[Instructions on setting up the workflows](.github/workflows/README.md)

❗**Note:** If you are not using GitHub Actions, you can remove the [.github/workflows/deploy.yml](.github/workflows/deploy.yml).


## 🛠️ Installation and Setup
1. Clone the repository:
```bash
git clone https://github.com/k4sper1love/watchlist-api.git
```
2. Go to the project directory:
```bash
cd watchlist-api
```
3. Install dependencies:
```bash
go mod tidy
```
4. Install Loki Docker Driver (required for Depoyment):
```bash
docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions
```

## ⚙️ Configuration
Create an `.env` file in the root directory and configure the environment variables. Use [.env.example](.env.example) as a reference.
```txt
(Optional) VERSION=unknown 

(Optional) GRAFANA_PASSWORD=password

APP_PORT=8001

APP_ENV=local

APP_MIGRATIONS=file://migrations

APP_SECRET=TOKENPASSWORD 

(Optional) APP_TELEGRAM=TOKENPASSWORD

POSTGRES_DB=watchlist

POSTGRES_PORT=5432

POSTGRES_USER=postgres

POSTGRES_PASSWORD=password

POSTGRES_HOST=localhost
```
❗**Note:** For Docker Compose, replace **localhost** with **db** in the `POSTGRES_HOST` value.

## ⚡ Run the application
### Using Terminal
Run the application directly with Go:
```bash
go run ./cmd/watchlist
```
#### Available Flags
- `-p`, `--port`: Port number for the API server (default: `8001`).
- `-e`, `--env`: Environment setting (`local`, `dev`, `prod`) (default: `local`).
- `-m`, `--migrations`: Path to migration files (e.g., `file://migrations`).
- `-s`, `--secret`: Secret password for creating JWT tokens (default: `secretPass`).
- `-t`, `--telegram`: Secret password for checking verification token (default: `secretPass`).

### Using Docker Compose
Start the project with Docker Compose:
```bash
docker-compose up --build
```
❗**Note:** Use the `--env-file` flag if your `.env` file is not in the root directory.

## 📋 API Documentation
You can find the Swagger documents and test the API functionality at:
```bash
http://localhost:8001/swagger/index.html
```
[Swagger Documentation Update](api/README.md)

❗**Note:** Use the port on which your application is running.

## 🛡️ Account Authentication
There are two methods to register or log in to your account:
### Using Credentials
- Endpoints: `/auth/register`, `/auth/login`
- Use this method to register or log in with your username and password.
### Via Telegram Bot
- Endpoints: `/auth/register/telegram`, `/auth/login/telegram`
- The Telegram bot generates a token by signing it with the `APP_TELEGRAM` secret. 
- The token contains the `Telegram ID` as an integer in the claims.
- This token is sent in the header with the key Verification.
- The API reads the token, extracts the `Telegram ID`, and generates a random username for the user.

## 👨🏻‍💻 Testing with Postman
Watchlist API uses Postman for automated API testing.

### Postman Tests
To run the tests, you need the Postman Collection and Environment files.

These files are located in the [tests/postman](tests/postman) directory:
- **Postman Collection**: [postman_collection.json](tests/postman/postman_collection.json) - Contains a set of pre-defined API requests for testing various endpoints.
- **Postman Environment**: [postman_environment.json](tests/postman/postman_environment.json) - Provides environment-specific variables such as base URLs and authentication tokens.

### Running Tests
1. Import Collection and Environment into Postman.
2. Select the Environment.
3. Run the Collection.

## 📁 Monitoring System (Grafana)
**Grafana** is used to collect and monitor logs using Loki and metrics using Prometheus.
### Setting Up Monitoring 🔧
1. Navigate to Grafana at:
```bash
http://localhost:3000
```
4. Log in using the credentials (by default: `admin`, `admin`).
5. Create a new dashboard and select the preset data sources:
- For logs, use **Loki**.
- For metrics, use **Prometheus**.
7. Set up a query, such as `{compose_service="app"}`, and save the dashboard.

## 🌐 Watchlist REST API Endpoints
```bash
# Image Section
POST /upload
GET /images/:filename

# Auth section
POST /api/v1/auth/register
POST /api/v1/auth/register/telegram
POST /api/v1/auth/login
POST /api/v1/auth/login/telegram
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET /api/v1/auth/check-token

# User section
GET /api/v1/user
PUT /api/v1/user
DELETE /api/v1/user

# Films section
GET /api/v1/films
POST /api/v1/films
GET /api/v1/films/:film_id
PUT /api/v1/films/:film_id
DELETE /api/v1/films/:film_id

# Collections section
GET /api/v1/collections
POST /api/v1/collections
GET /api/v1/collections/:collection_id
PUT /api/v1/collections/:collection_id
DELETE /api/v1/collections/:collection_id

# Collection_films section
GET /api/v1/collections/:collection_id/films
POST /api/v1/collections/:collection_id/films
POST /api/v1/collections/:collection_id/films/:film_id
GET /api/v1/collections/:collection_id/films/:film_id
PUT /api/v1/collections/:collection_id/films/:film_id
DELETE /api/v1/collections/:collection_id/films/:film_id
```

## 📊 Database Structure
The database schema is detailed in the [schema.dbml](docs/database/schema.dbml) file, which you can view in the [docs/database](docs/database) folder.

![Database Schema](docs/database/schema.png)

## 📜 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE.txt) file for details.

## 📫 Contact
For any questions or feedback, please contact:
- **Email**: s_yelkin@proton.me
- **Telegram**: [k4sper1love](https://t.me/k4sper1love)
- **GitHub**: [k4sper1love](https://github.com/k4sper1love)
