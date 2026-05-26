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