create table if not exists customers (
    id BIGSERIAL PRIMARY KEY,
    customer_id text NOT NULL UNIQUE,
    customer_type text NOT NULL,
    customer_name text NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    phone_no text NOT NULL,
    "address" text NOT NULL,
    created_by BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);
create table if not exists customer_usages (
    id BIGSERIAL PRIMARY KEY,
    customer_id text NOT NULL,
    customer_name text NOT NULL,
    product_id text NOT NULL,
    product_name text NOT NULL,
    transaction_at BIGINT NOT NULL,
    balance BIGINT NOT NULL,
    balance_amount numeric(64,2) NOT NULL,
    usage BIGINT NOT NULL,
    usage_amount numeric(64,2) NOT NULL,
    created_by BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);
create table if not exists sales_order_headers (
    id BIGSERIAL PRIMARY KEY,
    "order_number" text NOT NULL,
    customer_id text NOT NULL,
    customer_name text NOT NULL,
    subtotal numeric(64,2) NOT NULL,
    tax numeric(64,2) NOT NULL,
    grandtotal numeric(64,2) NOT NULL,
    created_by BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);
create table if not exists sales_order_lines (
    id BIGSERIAL PRIMARY KEY,
    sales_order_header_id BIGINT NOT NULL,
    product_id text NOT NULL,
    product_name text NOT NULL,
    quantity BIGINT NOT NULL,
    rate_item numeric(64,2) NOT NULL,
    tax_rate numeric(64,2) NOT NULL,
    subtotal numeric(64,2) NOT NULL,
    grandtotal numeric(64,2) NOT NULL,
    created_by BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);