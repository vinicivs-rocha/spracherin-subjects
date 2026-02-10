-- +goose Up
CREATE TABLE subjects (
  id CHAR(36) PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description_text TEXT NOT NULL,
  description_links JSON NOT NULL,
  concepts JSON NOT NULL
);

-- +goose Down
DROP TABLE subjects;
