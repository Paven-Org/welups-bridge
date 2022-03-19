-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS wel_sys_accounts (
  address varchar(256) NOT NULL,
  status varchar(10) NOT NULL DEFAULT 'locked',
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp NOT NULL DEFAULT NOW(),

  PRIMARY KEY (address),
  CHECK (status IN ('ok','locked'))
);

CREATE TABLE IF NOT EXISTS wel_sys_roles (
  role varchar(20) NOT NULL,
  PRIMARY KEY (role)
);

INSERT INTO wel_sys_roles(role) VALUES
  ('unauthorized'),
  ('super_admin'),
  ('treasury'),
  ('vault'),
  ('authenticator');

CREATE TABLE IF NOT EXISTS wel_sys_prikeys ( -- for those accounts we can afford to store keys
  address varchar(256) NOT NULL,
  prikey varchar(256) NOT NULL UNIQUE,

  PRIMARY KEY (address),
  FOREIGN KEY (address) REFERENCES wel_sys_accounts(address)
);

CREATE TABLE IF NOT EXISTS wel_sys_account_roles (
  address varchar(256) NOT NULL,
  role varchar(20) NOT NULL DEFAULT 'unauthorized',

  FOREIGN KEY (address) REFERENCES wel_sys_accounts(address),
  FOREIGN KEY (role) REFERENCES wel_sys_roles(role),
  PRIMARY KEY (address, role),
  CHECK (role IN ('unauthorized','super_admin','treasury','vault','authenticator'))
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE wel_sys_account_roles CASCADE;
DROP TABLE wel_sys_prikeys CASCADE;
DROP TABLE wel_sys_roles CASCADE;
DROP TABLE wel_sys_account CASCADE;
-- +goose StatementEnd
