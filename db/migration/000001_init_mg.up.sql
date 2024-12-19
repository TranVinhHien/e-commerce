CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  -- "email" varchar UNIQUE NOT NULL,
  "is_active" bool NOT NULL DEFAULT true,
  "create_at" TIMESTAMPTZ NOT NULL
);
-- CREATE TABLE "sessions" (
--   "id" uuid PRIMARY KEY,
--   "username" varchar NOT NULL,
--   "refresh_token" varchar NOT NULL,
--   "user_agent" varchar NOT NULL,
--   "client_ip" varchar NOT NULL,
--   "is_blocked" boolean NOT NULL DEFAULT false,
--   "expires_at" timestamptz NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now())
-- );