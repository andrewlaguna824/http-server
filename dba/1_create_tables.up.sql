/* SQL Commands to initialize all tables */

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS book (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    author VARCHAR(64) NOT NULL,
    title VARCHAR(64) NOT NULL
);
