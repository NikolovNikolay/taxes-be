CREATE TABLE IF NOT EXISTS inquiries
(
    id                    uuid      DEFAULT gen_random_uuid(),
    user_id               uuid                                NOT NULL,
    files                 TEXT                                NOT NULL,
    type                  INT                                 NOT NULL,
    year                  INT                                 NOT NULL,
    prefix                VARCHAR(10)                         NOT NULL,
    paid                  BOOLEAN                             NOT NULL,
    email                 VARCHAR(255)                        NOT NULL,
    full_name             VARCHAR(255)                        NOT NULL,
    generated_with_coupon BOOLEAN                             NOT NULL,
    created_at            TIMESTAMP DEFAULT current_timestamp NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS inquiries_id_uindex ON inquiries (id);

ALTER TABLE IF EXISTS inquiries
    ADD CONSTRAINT inquiries_pk PRIMARY KEY (id);

