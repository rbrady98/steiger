-- name: GetJoke :one
SELECT * FROM jokes
WHERE id = ? LIMIT 1;

-- name: CreateJoke :exec
INSERT INTO jokes (
    joke, nsfw, created_at
) VALUES (
    ?, ?, datetime('now')
);

-- name: ListJokes :many
SELECT * FROM jokes
LIMIT 50;
