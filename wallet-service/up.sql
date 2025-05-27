CREATE TABLE IF NOT EXISTS payments(
    id UUID PRIMARY KEY,
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    note TEXT ,
    stat TEXT NOT NULL CHECK (stat IN ('PAINDING' , 'COMPLATE','ERROR','REVERSED')),
    imenpotency_key VARCHAR9(255) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);