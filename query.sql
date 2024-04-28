-- name: GetAuthor :one
SELECT *
FROM authors
WHERE id = ?
LIMIT 1;
-- name: ListAuthors :many
SELECT *
FROM authors
ORDER BY name;
-- name: CreateAuthor :one
INSERT INTO authors (name, bio)
VALUES (?, ?)
RETURNING *;
-- name: UpdateAuthor :one
UPDATE authors
set name = ?,
    bio = ?
WHERE id = ?
RETURNING *;
-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;
-- name: IncrementCounterAndReturn :one
INSERT INTO counters("name", "counter")
VALUES (?, ?) ON CONFLICT("name") DO
UPDATE
SET "name" = EXCLUDED."name",
    "counter" = EXCLUDED."counter"
RETURNING "counter";
-- name: GetCounter :one
SELECT CAST(COALESCE(MAX("counter"), 0) as INTEGER)
FROM counters
WHERE "name" = ?;