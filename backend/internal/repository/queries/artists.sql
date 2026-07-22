-- name: CreateArtist :one
INSERT INTO artists (id, name, official_site_url)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetArtist :one
SELECT * FROM artists WHERE id = $1;

-- name: GetArtistByURL :one
SELECT * FROM artists WHERE official_site_url = $1;

-- name: ListArtists :many
SELECT * FROM artists
WHERE (sqlc.narg('query')::text IS NULL OR name ILIKE '%' || sqlc.narg('query') || '%')
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CountArtists :one
SELECT count(*) FROM artists
WHERE (sqlc.narg('query')::text IS NULL OR name ILIKE '%' || sqlc.narg('query') || '%');

-- name: ListAllArtists :many
SELECT * FROM artists ORDER BY name;

-- name: UpdateArtist :one
UPDATE artists
SET name = $2, official_site_url = $3
WHERE id = $1
RETURNING *;

-- name: DeleteArtist :execrows
DELETE FROM artists WHERE id = $1;
