-- name: RegisterStatusPage :one
INSERT INTO status_pages (
    hostname,
    page_data_url
) VALUES ($1, $2) RETURNING id;

-- name: GetStatusPage :one
SELECT page_data_url FROM status_pages WHERE hostname = $1 LIMIT 1;

