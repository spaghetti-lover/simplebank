CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
<<<<<<< HEAD
  "is_blocked" boolean NOT NULL DEFAULT FALSE,
=======
  "is_blocked" boolean NOT NULL DEFAULT false,
>>>>>>> d4d0e58 (refactor)
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
