# Database Migration
To create a new database migration, follow these steps:
1. **Install the `migrate` CLI tool**:

   ```bash
   curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xz
   
   sudo mv migrate /usr/local/bin
   ```
    - Adjust the URL and commands for your operating system if needed.
3. **Create a new migration**:

   ```bash
   migrate create -ext sql -dir migrations -seq create_users_table
   ```
    - `-ext sql`: Specifies the file extension for the migration files.
    - `-dir migrations`: Specifies the directory for storing migration files.
    - `-seq`: Creates a sequential migration file.
    - `create_users_table`: Name of the migration (adjust as needed).