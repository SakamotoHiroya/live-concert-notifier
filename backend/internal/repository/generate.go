package repository

// Regenerate internal/repository/sqlcgen from migrations/*.sql and
// internal/repository/queries/*.sql (sqlc.yaml lives at the backend module root).
//go:generate sh -c "cd ../.. && go tool sqlc generate"
