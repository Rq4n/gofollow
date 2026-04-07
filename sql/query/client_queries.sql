-- name: CreateNewClient :one
INSERT INTO clients (user_id, name, email, invoice_link)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetClientByID :one
SELECT * FROM clients 
WHERE id = $1;

