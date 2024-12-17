CREATE TABLE IF NOT EXISTS permissions (
  id bigserial PRIMARY KEY,
  code text NOT NULL
)

CREATE TABLE IF NOT EXISTS users_premissions (
  user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
  permission_id bigint NOT NULL REFERENCES premissions ON DELETE CASCADE,
  PRIMARY KEY(user_id, premissions_id)
)

INSERT INTO permissions (code)
VALUES ('movies:read'), ('movies:write');
