CREATE TABLE IF NOT EXISTS payments(
    id UUID PRIMARY KEY,
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    note TEXT ,
    imenpotency_key VARCHAR(255) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);