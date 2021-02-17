CREATE TABLE IF NOT EXISTS coupons
(
    id                uuid DEFAULT gen_random_uuid(),
    parent_request_id uuid         NOT NULL,
    max_attempts      INT          NOT NULL,
    attempts          INT          NOT NULL,
    type              INT          NOT NUll,
    email             varchar(255) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS coupons_id_uindex
    ON coupons (id);

CREATE INDEX coupons_parent_request_id_index
    ON coupons (parent_request_id);

ALTER TABLE coupons
    ADD CONSTRAINT coupons_pk
        PRIMARY KEY (id);



