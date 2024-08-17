-- +goose Up
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email TEXT NOT NULL,
  username TEXT,
  profileIndex TEXT
)

-- +goose Down
DROP TABLE users