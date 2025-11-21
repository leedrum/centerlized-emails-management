CREATE SCHEMA "<tenant>";
SET search_path = "<tenant>";
-- rồi chạy toàn bộ schema bên dưới

-- THREADS TABLE
CREATE TABLE IF NOT EXISTS threads (
    id BIGSERIAL PRIMARY KEY,
    subject TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- EMAILS TABLE
CREATE TABLE IF NOT EXISTS emails (
    id BIGSERIAL PRIMARY KEY,

    message_id TEXT,
    in_reply_to TEXT,
    "references" TEXT,

    thread_id BIGINT REFERENCES threads(id),

    subject TEXT,
    sender_name TEXT,
    sender_email TEXT,
    receiver_name TEXT,
    receiver_email TEXT,

    raw_body TEXT,
    text_body TEXT,

    attachments JSONB,

    sent_at TIMESTAMPTZ,
    received_at TIMESTAMPTZ,

    s3_key TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW()
);
