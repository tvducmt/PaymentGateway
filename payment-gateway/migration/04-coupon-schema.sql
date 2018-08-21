
ALTER TABLE receipt ADD COLUMN coupon_value FLOAT;

CREATE TYPE coupon_state AS ENUM ('spend', 'unspend');

CREATE TABLE coupon (
    id SERIAL PRIMARY KEY,
    tx_id INT REFERENCES transaction(id),
    status coupon_state DEFAULT 'unspend' NOt NULl,
    code VARCHAR(10) NOT NULL  UNIQUE,
    create_at TIMESTAMP NOT NULl DEFAULT NOW(),
    value FLOAT
);

ALTER TABLE transaction  ALTER COLUMN code DROP NOT NULL;
ALTER TABLE receipt  ALTER COLUMN raw_receipt DROP NOT NULL;
ALTER TABLE receipt  ALTER COLUMN phonenumber DROP NOT NULL;
ALTER TABLE receipt  ALTER COLUMN create_at SET DEFAULT NOW();

ALTER TABLE coupon ADD CONSTRAINT tx_id_unique UNIQUE (tx_id);

ALTER TABLE coupon  ALTER COLUMN code SET DATA TYPE VARCHAR(12);
ALTER TABLE coupon ADD COLUMN currency VARCHAR(10);
