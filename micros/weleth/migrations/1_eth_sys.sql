-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS wel_eth_sys (
  eth_last_scan_block bigint(20),
  wel_last_scan_block bigint(20)
);

INSERT INTO wel_eth_sys(eth_last_scan_block, wel_last_scan_block) VALUES (0,0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE wel_eth_sys CASCADE;
-- +goose StatementEnd
