-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS total_rows_of (
  users bigint
);

INSERT INTO total_rows_of (users) SELECT COUNT(*) FROM users;

CREATE OR REPLACE FUNCTION new_user()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
AS $$
BEGIN
  UPDATE total_rows_of SET users=users+1;
  RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER new_user_trigger 
  AFTER INSERT ON users
  FOR EACH ROW EXECUTE PROCEDURE new_user();


CREATE OR REPLACE FUNCTION rm_user()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
AS $$
BEGIN
  UPDATE total_rows_of SET users=users-1;
  RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER rm_user_trigger AFTER DELETE ON users
  FOR EACH ROW EXECUTE PROCEDURE rm_user();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TRIGGER new_user_trigger ON users;
DROP FUNCTION new_user;
DROP TRIGGER rm_user_trigger ON users;
DROP FUNCTION rm_user;
DROP TABLE total_rows_of CASCADE;
-- +goose StatementEnd
