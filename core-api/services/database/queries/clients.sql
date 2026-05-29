-- name: CreateClient :one
INSERT INTO "Clients" (
    name, 
    email, 
    request_type,
    status, 
    priority, 
    amount)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING 
  id, 
  name, 
  email, 
  request_type, 
  status, 
  priority, 
  amount, 
  created_at, 
  updated_at;

-- name: GetClientByEmail :one
SELECT * FROM "Clients" WHERE email = $1;

-- name: UpdateClientStatus :one
UPDATE "Clients"
SET status = $2, updated_at = NOW()
WHERE email = $1
RETURNING *;