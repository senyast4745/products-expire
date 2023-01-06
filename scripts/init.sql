CREATE TYPE product_type AS ENUM ('FOOD', 'DRUG');

CREATE TABLE products
(
    id              BIGSERIAL,
    chat_id         BIGINT,
    name            VARCHAR(256),
    type            product_type,
    expiration_date DATE,
    is_expired      BOOLEAN DEFAULT FALSE
);

ALTER TABLE products
    ADD PRIMARY KEY (id);