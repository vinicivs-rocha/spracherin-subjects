# Database Setup

**Environment Variables**
1. `DB_HOST`
2. `DB_PORT` (optional, default `3306`)
3. `DB_USER`
4. `DB_PASSWORD`
5. `DB_NAME`
6. `DB_PARAMS` (optional, query string like `charset=utf8mb4&collation=utf8mb4_general_ci`)

**Goose Migrations**
1. `goose -dir migrations mysql "$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true" up`
2. `goose -dir migrations mysql "$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true" down`
