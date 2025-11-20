CREATE TABLE IF NOT EXISTS emails (
    id SERIAL PRIMARY KEY,
    message_id TEXT,
    in_reply_to TEXT,
    references TEXT,
    thread_id BIGINT,
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
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    subject TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
