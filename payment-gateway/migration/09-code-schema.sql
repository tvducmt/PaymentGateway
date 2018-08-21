CREATE TABLE codetable (
    id     SERIAL PRIMARY KEY,
    code    VARCHAR NOT NULL UNIQUE,
    active boolean NOT NULL DEFAULT false
);