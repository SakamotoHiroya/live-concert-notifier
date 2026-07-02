-- name: CreateScrapeJob :one
INSERT INTO scrape_jobs (id, artist_id, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetScrapeJob :one
SELECT sj.*, a.name AS artist_name
FROM scrape_jobs sj
JOIN artists a ON a.id = sj.artist_id
WHERE sj.id = $1;

-- name: ListScrapeJobs :many
SELECT sj.*, a.name AS artist_name
FROM scrape_jobs sj
JOIN artists a ON a.id = sj.artist_id
WHERE (sqlc.narg('artist_id')::uuid IS NULL OR sj.artist_id = sqlc.narg('artist_id'))
  AND (sqlc.narg('status')::text IS NULL OR sj.status = sqlc.narg('status'))
ORDER BY sj.id
LIMIT $1 OFFSET $2;

-- name: CountScrapeJobs :one
SELECT count(*)
FROM scrape_jobs sj
WHERE (sqlc.narg('artist_id')::uuid IS NULL OR sj.artist_id = sqlc.narg('artist_id'))
  AND (sqlc.narg('status')::text IS NULL OR sj.status = sqlc.narg('status'));

-- name: UpdateScrapeJobStatus :one
UPDATE scrape_jobs
SET status = $2, started_at = $3, finished_at = $4, error_message = $5
WHERE id = $1
RETURNING *;
