# CI/CD Configuration for Watchlist API
This document provides instructions on setting up the CI/CD workflows using GitHub Actions.

## 🛠️ Installation and Setup
### Initializing the project 🚀
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
### Setting Up SSH Keys 🔐
1. Generate SSH keys on your computer:
```sh
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
```
2. View the public and private keys:
```sh
cat ~/.ssh/id_rsa.pub

cat ~/.ssh/id_rsa # SSH_PRIVATE_KEY
```
3. Access the remote server:
```sh
ssh `user@host`
```
4. (Optional) Remove old keys from the server:
```sh
ssh-keygen -R `host`
```
5. Add the public key to the server:
```sh
nano ~/.ssh/authorized_keys
```
- Paste the contents of the public key into the file.
- Save the file with `Ctrl + O -> Enter`, then exit with `Ctrl + X`.

### Installing Docker, Docker Compose and Loki Docker Driver 🐳
1. Ensure you are connected to the server. If not:

   ```bash
   ssh user@host
   ```
3. Install Docker:

   ```bash
   sudo apt update

   curl -fsSL https://get.docker.com | sudo sh

   sudo docker --version
   ```
5. Install Docker Compose:

   ```bash
   sudo curl -L "https://github.com/docker/compose/releases/download/v2.22.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

   sudo chmod +x /usr/local/bin/docker-compose

   docker-compose --version
   ```
7. Install Loki Docker Driver

   ```bash
   docker plugin install grafana/loki-docker-driver:2.9.4 --alias loki --grant-all-permissions
   ```

## ⚙️ Configuration
To prepare your repository for CI/CD, you need to set up GitHub Secrets. Follow these steps:
1. Go to `Your repository -> Settings -> Security -> Secrets and variables -> Actions`.
2. Add the following secrets:

   ```txt
   SSH_PRIVATE_KEY: The private key for SSH access.
   
   SSH_HOST: The hostname of your server.
   
   SSH_USER: The SSH username for accessing the server.
   
   FULL_REBUILD: (Optional) Flag to trigger a full rebuild.
   
   GRAFANA_PASSWORD: Password for Grafana.
   
   POSTGRES_DB: PostgreSQL database name.
   
   POSTGRES_USER: PostgreSQL username.
   
   POSTGRES_PASSWORD: PostgreSQL password.
   
   POSTGRES_PORT: PostgreSQL port number.
   
   APP_DSN: Database connection string.
   
   APP_MIGRATIONS: Path to database migration files.

   APP_PORT: Port number for the API server.
   
   APP_SECRET: Secret password for creating JWT tokens.
   
   APP_ENV: Environment setting (local, dev, prod).
   ```

## 🔄 GitHub Actions Workflow
By default, the CI/CD workflows are configured to run on pushes to the `main` branch.

This is set up in the GitHub Actions [configuration file](deploy.yml) with the following setting:
```yaml
on:
  push:
    branches:
      - main
```
You can customize this setting as needed to fit your workflow.

## 🌐 Available Services After Deployment
- Watchlist API: `http://SSH_HOST:APP_PORT`
- Swagger API Documentation: `http://SSH_HOST:APP_PORT/swagger/index.html`
- Grafana (Log Aggregator): `http://SSH_HOST:3000`
