CREATE TABLE IF NOT EXISTS "phone_book" (
    "id" UUID NOT NULL,
    "fullname" VARCHAR(255) NOT NULL,
    "phone_number" VARCHAR(255) NOT NULL,
    "address" VARCHAR(255) NOT NULL,
    "created_date_utc" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" VARCHAR(255) NOT NULL,
    "updated_date_utc" TIMESTAMPTZ,
    "updated_by" VARCHAR(255),
    "deleted_date_utc" TIMESTAMPTZ,
    "deletedby" VARCHAR(255),
    CONSTRAINT "ak_phone_book_id" UNIQUE("id"),
    CONSTRAINT "pk_phone_book" PRIMARY KEY("ud")
);
CREATE INDEX IF NOT EXISTS "ix_phone_book_id" ON "phone_book" USING btree("id");