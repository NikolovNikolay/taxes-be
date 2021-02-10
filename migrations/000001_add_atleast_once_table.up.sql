CREATE TABLE IF NOT EXISTS at_least_once_tasks
(
    key  varchar(255),
    id   varchar(36),
    done boolean     NOT NULL default false,
    time TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (key, id)
);
