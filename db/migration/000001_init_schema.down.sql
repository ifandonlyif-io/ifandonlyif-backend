ALTER TABLE IF EXISTS "iff_nfts" DROP CONSTRAINT IF EXISTS "iff_nfts_user_wallet_address_fkey";
ALTER TABLE IF EXISTS "report_blocklists" DROP CONSTRAINT IF EXISTS "report_blocklists_user_wallet_address_fkey";
ALTER TABLE IF EXISTS "report_whitelists" DROP CONSTRAINT IF EXISTS "report_whitelists_user_wallet_address_fkey";
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS iff_nfts;
DROP TABLE IF EXISTS report_blocklists;
DROP TABLE IF EXISTS report_whitelists;
DROP TABLE IF EXISTS nft_projects;