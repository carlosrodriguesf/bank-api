CREATE TABLE accounts
(
    id          VARCHAR(36)              NOT NULL PRIMARY KEY DEFAULT uuid(),
    name        TEXT                     NOT NULL,
    document    VARCHAR(11)              NOT NULL UNIQUE,
    balance     BIGINT                   NOT NULL             DEFAULT 0,
    secret      TEXT                     NOT NULL,
    secret_salt TEXT                     NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL             DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transfers
(
    id                VARCHAR(36)              NOT NULL PRIMARY KEY DEFAULT uuid(),
    origin_account_id VARCHAR(36)              NOT NULL REFERENCES accounts (id),
    target_account_id VARCHAR(36)              NOT NULL REFERENCES accounts (id),
    amount            BIGINT                   NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL             DEFAULT CURRENT_TIMESTAMP,

    CHECK ( amount > 0 )
);