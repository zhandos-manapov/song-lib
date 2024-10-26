CREATE TABLE IF NOT EXISTS "song" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" varchar NOT NULL,
  "group_id" uuid,
  "text" varchar NOT NULL,
  "link" varchar NOT NULL,
  "release_date" date NOT NULL
);

ALTER TABLE "song" ADD FOREIGN KEY ("group_id") REFERENCES "group" ("id");
