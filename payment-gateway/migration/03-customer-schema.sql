CREATE TABLE customer (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES account(id) NOT NULL,
    cus_stripe_id   TEXT NOT NULL UNIQUE
);