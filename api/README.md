# Swagger Documentation Update
If you have made changes to the Watchlist API and need to regenerate the Swagger documentation, follow these steps:
1. **Install Swag CLI** (if not already installed):
   
   ```bash
   go get -u github.com/swaggo/swag/cmd/swag
   ```
3. **Run the command** to regenerate documentation:
   
   ```bash
   swag init -g cmd/watchlist/main.go -o api
   ```
5. **Verify** that the `api` directory contains the updated Swagger files.
