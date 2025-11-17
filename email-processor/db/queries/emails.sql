-- name: CreateEmail :one
INSERT INTO emails (
    subject, sender_name, sender_email, receiver_name, receiver_email,
    raw_body, text_body, attachments, sent_at, received_at, s3_key
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10, $11
) RETURNING id;
