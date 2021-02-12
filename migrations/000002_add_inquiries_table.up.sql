CREATE TABLE IF NOT EXISTS inquiries
(
    id         UUID      DEFAULT gen_random_uuid(),
    user_id    UUID                                NOT NULL,
    files      TEXT                                NOT NULL,
    type       INT                                 NOT NULL,
    prefix     VARCHAR(10)                         NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS inquiries_id_uindex ON inquiries (id);

ALTER TABLE IF EXISTS inquiries
    ADD CONSTRAINT inquiries_pk PRIMARY KEY (id);

