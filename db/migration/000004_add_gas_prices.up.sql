CREATE TABLE "gas_prices" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "average" integer,
  "created_at" timestamptz DEFAULT now()
);
