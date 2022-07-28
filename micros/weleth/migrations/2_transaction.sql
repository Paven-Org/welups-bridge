-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS wel_cashin_eth_trans (
  id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  request_id varchar(100) DEFAULT '',
  deposit_tx_hash varchar(100) UNIQUE,
  claim_tx_hash varchar(100) DEFAULT '', 
  wel_token_addr varchar(100),
  eth_token_addr varchar(100),
  wel_wallet_addr varchar(100),
  eth_wallet_addr varchar(100),
  network_id varchar(20),
  amount varchar(40),
  fee varchar(40),
  deposit_status varchar(20),
  claim_status varchar(20) DEFAULT 'unconfirmed',
  deposit_at timestamp,
  claim_at timestamp,
  CHECK (claim_status IN ('confirmed','unconfirmed','pending'))
);

CREATE TABLE IF NOT EXISTS wel_cashin_eth_req (
  request_id varchar(100) PRIMARY KEY,
  tx_id integer REFERENCES wel_cashin_eth_trans(id),
  status varchar(20) DEFAULT 'pending',
  expired_at timestamp with time zone,
  CHECK (status IN ('success','pending','expired','doubleclaimed'))
);

CREATE TABLE IF NOT EXISTS eth_cashout_wel_trans (
  id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  request_id varchar(100) DEFAULT '',
  deposit_tx_hash varchar(100) UNIQUE,
  claim_tx_hash varchar(100) DEFAULT '', 
  wel_token_addr varchar(100),
  eth_token_addr varchar(100),
  wel_wallet_addr varchar(100),
  eth_wallet_addr varchar(100),
  network_id varchar(20),
  amount varchar(40),
  fee varchar(40),
  deposit_status varchar(20),
  claim_status varchar(20) DEFAULT 'unconfirmed',
  deposit_at timestamp,
  claim_at timestamp,
  CHECK (claim_status IN ('confirmed','unconfirmed','pending'))
);

CREATE TABLE IF NOT EXISTS eth_cashout_wel_req (
  request_id varchar(100) PRIMARY KEY,
  tx_id integer REFERENCES eth_cashout_wel_trans(id),
  status varchar(20) DEFAULT 'pending',
  expired_at timestamp with time zone,
  CHECK (status IN ('success','pending','expired','doubleclaimed'))
);

CREATE UNIQUE INDEX wel_cashin_eth_deposit_tx_index ON wel_cashin_eth_trans(deposit_tx_hash); 
CREATE UNIQUE INDEX eth_cashout_wel_deposit_tx_index ON eth_cashout_wel_trans(deposit_tx_hash); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX wel_cashin_eth_deposit_tx_index;
DROP INDEX eth_cashout_wel_deposit_tx_index;
DROP TABLE wel_cashin_eth_req CASCADE;
DROP TABLE eth_cashout_wel_req CASCADE;
DROP TABLE wel_cashin_eth_trans CASCADE;
DROP TABLE eth_cashout_wel_trans CASCADE;
-- +goose StatementEnd
