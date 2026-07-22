-- name: CreateConcert :one
INSERT INTO concerts (
    id, artist_id, title, venue_name, venue_location, date,
    co_performers, is_festival, source_url, raw_text
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (artist_id, date, venue_name) DO NOTHING
RETURNING *;

-- name: GetConcert :one
SELECT c.*, a.name AS artist_name
FROM concerts c
JOIN artists a ON a.id = c.artist_id
WHERE c.id = $1;

-- name: ListConcerts :many
SELECT c.*, a.name AS artist_name
FROM concerts c
JOIN artists a ON a.id = c.artist_id
WHERE (sqlc.narg('artist_id')::uuid IS NULL OR c.artist_id = sqlc.narg('artist_id'))
  AND (sqlc.narg('from_date')::date IS NULL OR c.date >= sqlc.narg('from_date'))
  AND (sqlc.narg('to_date')::date IS NULL OR c.date <= sqlc.narg('to_date'))
  AND (sqlc.narg('is_festival')::boolean IS NULL OR c.is_festival = sqlc.narg('is_festival'))
ORDER BY c.date ASC
LIMIT $1 OFFSET $2;

-- name: CountConcerts :one
SELECT count(*)
FROM concerts c
WHERE (sqlc.narg('artist_id')::uuid IS NULL OR c.artist_id = sqlc.narg('artist_id'))
  AND (sqlc.narg('from_date')::date IS NULL OR c.date >= sqlc.narg('from_date'))
  AND (sqlc.narg('to_date')::date IS NULL OR c.date <= sqlc.narg('to_date'))
  AND (sqlc.narg('is_festival')::boolean IS NULL OR c.is_festival = sqlc.narg('is_festival'));

-- name: ListDashboardConcerts :many
SELECT c.*, a.name AS artist_name
FROM concerts c
JOIN artists a ON a.id = c.artist_id
JOIN user_artists ua ON ua.artist_id = c.artist_id
WHERE ua.user_id = $1 AND c.date >= CURRENT_DATE
ORDER BY c.date ASC
LIMIT $2 OFFSET $3;

-- name: CountDashboardConcerts :one
SELECT count(*)
FROM concerts c
JOIN user_artists ua ON ua.artist_id = c.artist_id
WHERE ua.user_id = $1 AND c.date >= CURRENT_DATE;
