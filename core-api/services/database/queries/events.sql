-- name: CreateEvent :one
INSERT INTO "Events" ("event_id", "card_id", "client_email", "timestamp")
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetEventsByClientEmail :many
SELECT * FROM "Events" WHERE client_email = $1 ORDER BY timestamp DESC;