-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS tx_to_treasury (
  tx_id varchar(100) PRIMARY KEY,
  from_address varchar(100),
  treasury_address varchar(100),
  token_address varchar(100) DEFAULT '0x0000000000000000000000000000000000000000',
  amount varchar(40),
  tx_fee varchar(40),
  status varchar(20) DEFAULT 'unconfirmed',
  created_at timestamp DEFAULT NOW(),

  CHECK (status IN ('unconfirmed','isCashin', 'expired'))
);

CREATE TABLE IF NOT EXISTS eth_cashin_wel_trans (
  id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  eth_tx_hash varchar(100) REFERENCES tx_to_treasury(tx_id),
  wel_issue_tx_hash varchar(100) DEFAULT '',
  eth_token_addr varchar(100),
  wel_token_addr varchar(100),
  network_id varchar(20),

  eth_wallet_addr varchar(100),
  wel_wallet_addr varchar(100),

  total varchar(40),
  amount varchar(40) DEFAULT '',
  commission_fee varchar(40) DEFAULT '',
  status varchar(20),

  created_at timestamp DEFAULT NOW(),
  issued_at timestamp,

  CHECK (status IN ('unconfirmed', 'confirmed', 'failed'))
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE eth_cashin_wel_trans CASCADE;
DROP TABLE tx_to_treasury CASCADE;
-- +goose StatementEnd
