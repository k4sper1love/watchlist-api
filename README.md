# Watchlist API 🎬

**Watchlist API** is a REST API that allows you to save and manage a list of films you want to watch. It includes features for registration, managing film collections, adding comments, and leaving feedback.

## 🔴 Live Server
- **Healthcheck Endpoint**: [http://91.243.71.175:8001/api/v1/healthcheck](http://91.243.71.175:8001/api/v1/healthcheck)
- **Swagger Documentation**: [http://91.243.71.175:8001/swagger/index.html#/monitoring/get_healthcheck](http://91.243.71.175:8001/swagger/index.html)
- **Grafana Dashboards**: [http://91.243.71.175:3000/dashboards](http://91.243.71.175:3000)

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
- [Testing with Postman](#-testing-with-postman)
  - [Postman tests](#postman-tests)
  - [Running tests](#running-tests)
- [Logging System](#-logging-system)
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
- **Authentication**: JWT
- **Database**: PostgreSQL
- **API Documentation**: Swagger
- **Logging**: Grafana + Loki
- **Testing**: Postman
- **Deployment**: Docker, Docker Compose
- **CI/CD**: GitHub Actions

## 📝 Project Requirements
- **Go**: 1.18+
- **PostgreSQL**
### For Deployment
- **Docker**
- **Docker Compose**
- **Loki Docker Driver**
### For Testing
- **Postman**

## 🔄 GitHub Actions (CI/CD)
The project uses **GitHub Actions** to automate testing, building, and deployment processes on remote server.

[Instructions on setting up the workflows](.github/workflows/README.md)

❗**Note:** If you are not using GitHub Actions for deployment, you can safely remove the [.github/workflows/deploy.yml](.github/workflows/deploy.yml) file.


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
docker plugin install grafana/loki-docker-driver:2.9.4 --alias loki --grant-all-permissions
```

## ⚙️ Configuration
Create an `.env` file in the root directory and configure the environment variables. Use [.env.example](.env.example) as a reference.
```text
VERSION=unknown

GRAFANA_PASSWORD=password

APP_SECRET=JWTPASSWORD

APP_PORT=8001

APP_ENV=local

APP_MIGRATIONS=file://migrations

APP_DSN=postgres://postgres:postgres@localhost:5432/example?sslmode=disable

POSTGRES_DB=watchlist

POSTGRES_PORT=5432

POSTGRES_USER=postgres
```
❗**Note:** For Docker Compose, replace **localhost** with **db** in the `APP_DSN` value.

## ⚡ Run the application
### Using Terminal
Run the application directly with Go:
```bash
go run ./cmd/watchlist
```
#### Available Flags
- `-p`, `--port`: Port number for the API server (default: `8001`).
- `-e`, `--env`: Environment setting (`local`, `dev`, `prod`) (default: `local`).
- `-d`, `--dsn`: PostgreSQL DSN for database connection (e.g., `postgres://user:pass@host:port/dbname?sslmode=disable`).
- `-m`, `--migrations`: Path to migration files (e.g., `file://migrations`).
- `-s`, `--secret`: Secret password for creating JWT tokens (default: `secretPass`).

### Using Docker Compose
Start the project with Docker Compose:
```bash
docker-compose up --build
```
❗**Note:** Use the `--env-file` flag if your `.env` file is not in the root directory.

## 📚 API Documentation
You can find the Swagger documents and test the API functionality at
```bash
http://localhost:8001/swagger/index.html
```
❗**Note:** Use the port on which your application is running.

## 🧪 Testing with Postman
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

## 📁 Logging System
Grafana + Loki is used to collect and monitor logs.
### Setting Up Logs 🔧
1. Navigate to Grafana at:
```bash
http://localhost:3000
```
4. Log in using the credentials (by default: `admin`, `admin`).
5. Create a new dashboard and add a `Loki` data source.
6. Set up a query, such as `{compose_service="app"}`, and save the dashboard.

## 🌐 Watchlist REST API Endpoints
```bash
# Auth section
POST /auth/register
POST /auth/login
POST /auth/refresh
POST /auth/logout

# User section
GET /user
PUT /user
DELETE /user

# Films section
GET /films
POST /films
GET /films/:film_id
PUT /films/:film_id
DELETE /films/:film_id

# Collections section
GET /collections
POST /collections
GET /collections/:collection_id
PUT /collections/:collection_id
DELETE /collections/:collection_id

# Collection_films section
GET /collections/:collection_id/films
POST /collections/:collection_id/films/:film_id
GET /collections/:collection_id/films/:film_id
PUT /collections/:collection_id/films/:film_id
DELETE /collections/:collection_id/films/:film_id
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
