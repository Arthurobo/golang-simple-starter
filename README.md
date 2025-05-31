# How To Add New Columns to Tables with Raw SQL Commands

## Steps:

1. Enter the PostgreSQL database container shell:
   ```bash
   docker compose exec -it <container_id_or_name> bash
   ```

2. Connect to the PostgreSQL client:
   ```bash
   psql -U postgres -d postgres
   ```

3. Add a new column to the desired table:
   ```sql
   ALTER TABLE users ADD COLUMN username VARCHAR(255) NOT NULL DEFAULT '';
   ```

4. To confirmed it works, use this command below:
    ```
    \d users
    ```

## Notes:
- Replace `<container_id_or_name>` with your actual container ID or name
- The example adds a `username` column to the `users` table
- The column is defined as VARCHAR(255) with a NOT NULL constraint and empty string as default value


# To Drop a table entire:
`DROP TABLE IF EXISTS users CASCADE;`