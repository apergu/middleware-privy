create table if not exists merchants (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    enterprise_id text NOT NULL UNIQUE,
    merchant_id text NOT NULL,
    merchant_name text NOT NULL,
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
    CONSTRAINT fk_customer
        FOREIGN KEY(customer_id)
        REFERENCES customers(id)
        ON DELETE CASCADE
);
create table if not exists channels (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL,
    channel_id text NOT NULL,
    channel_name text NOT NULL,
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
    CONSTRAINT fk_merchant
        FOREIGN KEY(merchant_id)
        REFERENCES merchants(id)
        ON DELETE CASCADE
);
create table if not exists divissions (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL,
    divission_id text NOT NULL,
    divission_name text NOT NULL,
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
    CONSTRAINT fk_channel
        FOREIGN KEY(channel_id)
        REFERENCES channels(id)
        ON DELETE CASCADE
);