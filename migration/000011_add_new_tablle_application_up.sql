
-- ALTER TABLE applications
-- ADD COLUMN "address" TEXT NOT NULL,
-- ADD COLUMN email TEXT NOT NULL,
-- ADD COLUMN phone_no TEXT NOT NULL,
-- ADD COLUMN "state" TEXT NOT NULL,
-- ADD COLUMN city TEXT NOT NULL,
-- ADD COLUMN zip_code TEXT NOT NULL,
-- ADD COLUMN application_code TEXT NOT NULL DEFAULT '',
-- ADD COLUMN customer_internalid BIGINT NOT NULL DEFAULT 0,
-- ADD COLUMN application_internalid BIGINT NOT NULL DEFAULT 0;


create table if not exists applications (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    enterprise_id BIGINT NOT NULL,
    application_id text NULL,
    application_name text NULL,
    "address" text NOT NULL,
    email text NOT NULL,
    phone_no text NOT NULL,
    "state" text NOT NULL,
    city text NOT NULL,
    zip_code text NOT NULL,
    created_by BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    application_code text NOT NULL DEFAULT '',
    customer_internalid BIGINT NOT null default 0,
    application_internalid BIGINT NOT null default 0,
    CONSTRAINT fk_customer
        FOREIGN KEY(customer_id)
        REFERENCES customers(id)
        ON DELETE CASCADE
);