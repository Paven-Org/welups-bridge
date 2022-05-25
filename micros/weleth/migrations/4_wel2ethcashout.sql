-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS wel_cashout_eth_trans (
  id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  eth_disperse_tx_hash varchar(100) DEFAULT '',
  wel_withdraw_tx_hash varchar(100) UNIQUE,
  eth_token_addr varchar(100),
  wel_token_addr varchar(100),
  network_id varchar(20),

  eth_wallet_addr varchar(100),
  wel_wallet_addr varchar(100),

  total varchar(40),
  amount varchar(40),
  commission_fee varchar(40),
  cashout_status varchar(20),
  disperse_status varchar(20),

  created_at timestamp DEFAULT NOW(),
  dispersed_at timestamp,

  CHECK (cashout_status IN ('unconfirmed', 'confirmed')),
  CHECK (disperse_status IN ('unconfirmed', 'confirmed', 'retry'))
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE wel_cashout_eth_trans CASCADE;
-- +goose StatementEnd
