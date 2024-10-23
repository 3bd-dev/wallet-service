-- migrate:up
CREATE TYPE transaction_type AS ENUM ('deposit', 'withdrawal');
CREATE TYPE transaction_status AS ENUM ('created', 'pending', 'completed', 'failed');
CREATE TYPE payment_gateway AS ENUM ('gateway_a', 'gateway_b');

CREATE TABLE transactions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,  -- Primary key for transaction (UUID)
    wallet_id uuid NOT NULL,  -- Foreign key to wallet table (UUID)
    amount DECIMAL(18, 2) NOT NULL,  -- Transaction amount with 2 decimal precision
    type transaction_type NOT NULL,  -- Transaction type (deposit/withdrawal)
    status transaction_status NOT NULL,  -- Transaction status (pending/completed/failed)
    payment_gateway payment_gateway NOT NULL,  -- Payment gateway used for transaction
    payment_method_details JSONB NOT NULL,  -- Payment method details (e.g., card details)
    payment_method VARCHAR(255) NOT NULL,  -- Payment method used for transaction (e.g., card, bank transfer)
    reference_id VARCHAR(255) UNIQUE,  -- Unique reference ID for the transaction (optional)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),  -- When the transaction was created
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),  -- When the transaction was last updated
    CONSTRAINT fk_wallet FOREIGN KEY (wallet_id) REFERENCES wallets(id)  -- Foreign key constraint
);
-- migrate:down
DROP TABLE IF EXISTS transactions;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS payment_gateway;


