CREATE TYPE product_type AS ENUM ('FOOD', 'DRUG');

CREATE TABLE products
(
    id              BIGSERIAL,
    chat_id         BIGINT,
    name            VARCHAR(256),
    type            product_type,
    expiration_date DATE
);

ALTER TABLE products ADD PRIMARY KEY (id);