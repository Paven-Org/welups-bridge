-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO wel_sys_accounts(address, status) VALUES
('WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2','ok'); -- default "root" wel account as deployed on testnet. Should be changed appripriately in prod deployment 

INSERT INTO wel_sys_account_roles(address, role) VALUES
('WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2','super_admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM wel_sys_account_roles WHERE address = 'WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2';
DELETE FROM wel_sys_accounts WHERE address = 'WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2';
-- +goose StatementEnd
