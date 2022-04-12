-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO wel_sys_accounts(address, status) VALUES
('WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2','ok'); -- default all-purpose wel account as deployed on testnet. Should be changed appripriately in prod deployment i.e. no account should hold more than 1 role simultaneously

INSERT INTO wel_sys_prikeys(address, prikey) VALUES
('WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2','ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a');

INSERT INTO wel_sys_account_roles(address, role) VALUES
('WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2','activator');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM wel_sys_account_roles WHERE address = 'WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2';
DELETE FROM wel_sys_prikeys WHERE address = 'WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2';
DELETE FROM wel_sys_accounts WHERE address = 'WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2';
-- +goose StatementEnd
