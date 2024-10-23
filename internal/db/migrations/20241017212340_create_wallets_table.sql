-- migrate:up

CREATE TABLE wallets (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,  -- Primary key for wallet (UUID)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),  -- When the wallet was created
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()  -- When the wallet was last updated
);
-- migrate:down
DROP TABLE IF EXISTS wallets;