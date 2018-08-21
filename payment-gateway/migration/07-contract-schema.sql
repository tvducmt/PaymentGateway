CREATE TABLE transaction_ethereum (
    id SERIAL PRIMARY KEY,
    sender VARCHAR NOT NULL ,
    txhash VARCHAR NOT NULL UNIQUE,
    amount FLOAT NOT NULL,
    currency VARCHAR NOT NULL DEFAULT 'eth' 
);
ALTER TABLE transaction  ADD COLUMN contract_id integer UNIQUE;
INSERT INTO payment_method (name) VALUES ('Ether Transfer');
-- INSERT INTO rule (gateway_id, regex) VALUES (1, '^.*CODE(?P<transaction_code>\w{6,6}).*$');