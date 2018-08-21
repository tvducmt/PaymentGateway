ALTER TABLE customer ADD COLUMN fingerprint TEXT NOT NULL UNIQUE;

ALTER TABLE transaction  ADD COLUMN charge_id TEXT UNIQUE;

ALTER type tx_state ADD VALUE 'failed';
ALTER type tx_state ADD VALUE 'successed';