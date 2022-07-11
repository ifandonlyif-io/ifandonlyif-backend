ALTER TABLE users
ADD COLUMN "image_uri" varchar ;

ALTER TABLE users
ALTER COLUMN "image_uri" SET NOT NULL;