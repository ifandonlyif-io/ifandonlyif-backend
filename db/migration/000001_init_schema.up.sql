CREATE TABLE "users" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "full_name" varchar,
  "wallet" varchar UNIQUE ,
  "created_at" timestamptz DEFAULT now(),
  "country_code" varchar,
  "email_address" varchar,
  "kyc_date" timestamptz,
  "twitter_name" varchar,
  "blockpass_id" bigint
);

CREATE TABLE "iff_nfts" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "project_id" bigint NOT NULL,
  "user_wallet_address" varchar UNIQUE NOT NULL,
  "nft_projects_contract_address" varchar UNIQUE NOT NULL,
  "nft_projects_collection_name" varchar UNIQUE NOT NULL,
  "mint_date" timestamptz NOT NULL,
  "mint_transaction" varchar NOT NULL
);

CREATE TABLE "report_blocklists" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "http_address" varchar NOT NULL,
  "verified_at" timestamptz,
  "user_wallet_address" varchar,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "report_whitelists" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "http_address" varchar NOT NULL,
  "verified_at" timestamptz,
  "user_wallet_address" varchar,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "nft_projects" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "name" varchar NOT NULL,
  "contract_address" varchar UNIQUE NOT NULL,
  "collection_name" varchar UNIQUE NOT NULL
);

CREATE INDEX ON "users" ("wallet");

CREATE INDEX ON "iff_nfts" ("user_wallet_address");

CREATE INDEX ON "report_blocklists" ("user_wallet_address");

CREATE INDEX ON "report_whitelists" ("user_wallet_address");

ALTER TABLE "iff_nfts" ADD FOREIGN KEY ("user_wallet_address") REFERENCES "users" ("wallet");

ALTER TABLE "iff_nfts" ADD FOREIGN KEY ("nft_projects_contract_address") REFERENCES "nft_projects" ("contract_address");

ALTER TABLE "iff_nfts" ADD FOREIGN KEY ("nft_projects_collection_name") REFERENCES "nft_projects" ("collection_name");

ALTER TABLE "report_blocklists" ADD FOREIGN KEY ("user_wallet_address") REFERENCES "users" ("wallet");

ALTER TABLE "report_whitelists" ADD FOREIGN KEY ("user_wallet_address") REFERENCES "users" ("wallet");
