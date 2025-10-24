CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL CHECK (price >= 0),
    start_date DATE NOT NULL,
    end_date DATE
);
