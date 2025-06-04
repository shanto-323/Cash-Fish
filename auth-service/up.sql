CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(250) NOT NULL,
    email VARCHAR(50) NOT NULL,
    r_token VARCHAR(250) NOT NULL
);

CREATE TABLE IF NOT EXISTS cards(
    id  SERIAL PRIMARY KEY, 
    user_id UUID NOT NULL,
    card_number VARCHAR(20) NOT NULL,
    brand VARCHAR(20) NOT NULL, 
    expiry_month INT NOT NULL, 
    expiry_year INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);