CREATE TABLE account (
    id          SERIAL PRIMARY KEY,
    email       TEXT NOT NULL UNIQUE,
    passphrase  TEXT NOT NULL,
    fullname    TEXT,
    phone       VARCHAR(20)
);

CREATE TABLE payment_method (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE payment_gateway (
    id     SERIAL PRIMARY KEY,
    name    VARCHAR(50) NOT NULL
);

CREATE TABLE rule (
    id          SERIAL PRIMARY KEY,
    gateway_id  INT REFERENCES payment_gateway(id),
    regex       TEXT NOT NULL,
    status      BOOLEAN DEFAULT true NOT NULL
);

CREATE TYPE tx_state AS ENUM ('pending', 'timeout', 'confirmed');

CREATE TABLE receipt (
    id              SERIAL PRIMARY KEY,
    raw_receipt     TEXT NOT NULL,
    phonenumber     TEXT NOT NULL,
    create_at       TIMESTAMP NOT NULL,
    parsed_amount   FLOAT,
    parsed_account  VARCHAR(26),
    parsed_code     VARCHAR(8),
    parsed_balance  FLOAT,
    sys_create_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE transaction (
    id          SERIAL PRIMARY KEY,
    user_id     INT REFERENCES account(id) NOT NULL,
    receipt_id  INT REFERENCES receipt(id),
    method_id   INT REFERENCES payment_method(id),
    gateway_id  INT REFERENCES payment_gateway(id),
    status    tx_state DEFAULT 'pending' NOT NULL,
    code     VARCHAR(8) NOT NULL,

    create_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_payment_method (
    id          SERIAL PRIMARY KEY,
    user_id     INT REFERENCES account(id) NOT NULL,
    method_id   INT REFERENCES payment_method(id) NOT NULL
);


