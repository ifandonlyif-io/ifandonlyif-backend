CREATE TABLE "gas_prices" (
  "id" bigserial PRIMARY KEY,
  "average" integer,
  "created_at" timestamptz DEFAULT now()
);
