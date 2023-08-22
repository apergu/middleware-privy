create table if not exists top_up_datas (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL,
    transaction_id text NOT NULL,
    enterprise_id text NOT NULL,
    enterprise_name text NOT NULL,
    original_service_id text NOT NULL,
    service_id text NOT NULL,
    "service_name" text NOT NULL,
    quantity BIGINT NOT NULL,
    transaction_date BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    CONSTRAINT fk_merchant
        FOREIGN KEY(merchant_id)
        REFERENCES merchants(id)
        ON DELETE CASCADE
);