-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS current_eth_treasury (
  singleton bool PRIMARY KEY DEFAULT TRUE,
  address varchar(256) NOT NULL,
  role varchar(20) NOT NULL,
  FOREIGN KEY (address, role) REFERENCES eth_sys_account_roles(address, role),
  CONSTRAINT onerow CHECK (singleton),
  CONSTRAINT treasury CHECK (role IN ('treasury')) 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE current_eth_treasury;
-- +goose StatementEnd
