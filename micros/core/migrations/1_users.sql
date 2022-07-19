-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS roles (
  role varchar(20),
  PRIMARY KEY (role)
);

CREATE TABLE IF NOT EXISTS users (
  id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
  username varchar(100) UNIQUE NOT NULL,
  password varchar(100) NOT NULL,
  email varchar(100) UNIQUE NOT NULL,
  status varchar(10) NOT NULL DEFAULT 'ok',

  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp NOT NULL DEFAULT NOW(),
  CHECK (status IN ('ok','locked','banned','permabanned'))
);

CREATE TABLE user_roles (
    user_id integer,
    role varchar(20),

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role) REFERENCES roles(role),
    PRIMARY KEY (user_id, role)
);

INSERT INTO roles (role) VALUES 
  ('root'),
  --('service'),
  ('admin');

-- default password: root
-- must be changed immediately upon operation of course
INSERT INTO users (username, password, email) VALUES
  ('root', '$2a$10$ubmld8cryzM0bULIwFHmwOHkRzylFwzhI4q9sGGtAhRDYBzwrMESC', 'welbridgeroot@gmail.com');
--  ('weleth_bridge', '$2a$10$TP.3Z1/JJyGJwtDaX0xSs.FrO76FEz17DNwujGO0FKpBxr9gKCZpm', 'welbridgeroot@gmail.com');

INSERT INTO user_roles(user_id, role)
  SELECT id, 'root' FROM users WHERE username = 'root';
  --UNION
  --SELECT id, 'service' FROM users WHERE username = 'weleth_bridge';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users CASCADE;
DROP TABLE roles CASCADE;
DROP TABLE user_roles CASCADE;
-- +goose StatementEnd
