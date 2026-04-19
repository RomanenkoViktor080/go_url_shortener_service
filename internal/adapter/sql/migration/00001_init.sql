-- +goose Up
CREATE TABLE url
(
    hash       VARCHAR(6) NOT NULL PRIMARY KEY,
    url        VARCHAR(255) NOT NULL,
    created_at TIMESTAMP  NOT NULL DEFAULT NOW()
);

CREATE TABLE hash
(
    hash VARCHAR(6) PRIMARY KEY
);

CREATE SEQUENCE unique_number_seq
    START WITH 1
    INCREMENT BY 1;

-- +goose Down
DROP TABLE url;
DROP TABLE hash;
DROP SEQUENCE unique_number_seq;
с