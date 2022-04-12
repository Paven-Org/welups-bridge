-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO eth_sys_accounts(address, status) VALUES
('0x25e8370E0e2cf3943Ad75e768335c892434bD090','ok'); -- default "root" eth account as deployed on rinkeby testnet. Should be changed appripriately in prod deployment 

INSERT INTO eth_sys_account_roles(address, role) VALUES
('0x25e8370E0e2cf3943Ad75e768335c892434bD090','operator'),
('0x25e8370E0e2cf3943Ad75e768335c892434bD090','AUTHENTICATOR');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM eth_sys_account_roles WHERE address = '0x25e8370E0e2cf3943Ad75e768335c892434bD090';
DELETE FROM eth_sys_accounts WHERE address = '0x25e8370E0e2cf3943Ad75e768335c892434bD090';
-- +goose StatementEnd
