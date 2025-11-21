-- name: GetEmailByID :one
SELECT
  id,
  message_id,
  in_reply_to,
  "references",
  thread_id,
  subject,
  sender_name,
  sender_email,
  receiver_name,
  receiver_email,
  raw_body,
  text_body,
  attachments,
  sent_at,
  received_at,
  s3_key,
  created_at
FROM emails
WHERE id = $1 LIMIT 1;
