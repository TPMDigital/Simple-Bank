
-- DROP INDEX ON "accounts" ("owner", "currency");
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";

ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

DROP INDEX IF EXISTS "accounts_owner_idx";

DROP TABLE IF EXISTS "users";