-- name: FollowArtist :exec
INSERT INTO user_artists (user_id, artist_id)
VALUES ($1, $2);

-- name: UnfollowArtist :execrows
DELETE FROM user_artists WHERE user_id = $1 AND artist_id = $2;

-- name: IsFollowing :one
SELECT EXISTS (
    SELECT 1 FROM user_artists WHERE user_id = $1 AND artist_id = $2
);

-- name: ListFollowedArtists :many
SELECT a.*
FROM artists a
JOIN user_artists ua ON ua.artist_id = a.id
WHERE ua.user_id = $1
ORDER BY a.name;
