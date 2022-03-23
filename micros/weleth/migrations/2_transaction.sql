-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS wel_eth_trans (
  id varchar(70),
  wel_eth boolean,
  deposit_tx_hash varchar(70) UNIQUE,
  claim_tx_hash varchar(70) UNIQUE, 
  wel_token_addr varchar(70),
  eth_token_addr varchar(70),
  wel_wallet_addr varchar(70),
  eth_wallet_addr varchar(70),
  network_id varchar(20),
  amount varchar(40),
  fee varchar(40),
  deposit_status varchar(20),
  claim_status varchar(20),
  deposit_at timestamp,
  claim_at timestamp,

  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX deposit_tx_index ON wel_eth_trans(deposit_tx_hash); 

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE wel_eth_trans DROP INDEX deposit_tx_index;
DROP TABLE wel_eth_trans CASCADE;
-- +goose StatementEnd
