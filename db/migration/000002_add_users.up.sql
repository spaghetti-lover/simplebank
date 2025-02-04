CREATE TABLE users (
  username varchar PRIMARY KEY,
  hashed_password varchar NOT NULL,
  fullname varchar NOT NULL,
  email varchar UNIQUE NOT NULL,
  password_changed_at timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00+00',
  created_at timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE accounts ADD FOREIGN KEY (owner) REFERENCES users (username);
--Lệnh này hoặc lệnh dưới đây cũng được
--CREATE UNIQUE INDEX ON accounts (owner, currency);
ALTER TABLE accounts ADD CONSTRAINT owner_currency_key UNIQUE (owner, currency);


