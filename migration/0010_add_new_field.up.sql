ALTER TABLE top_up_datas
ADD COLUMN IF NOT EXISTS customer_internalid BIGINT NOT null default 0,
ADD COLUMN IF NOT EXISTS merchant_internalid BIGINT NOT null default 0,
ADD COLUMN IF NOT EXISTS channel_internalid BIGINT NOT null default 0,
ADD COLUMN IF NOT EXISTS merchant_code text NOT null default '',
ADD COLUMN IF NOT EXISTS channel_id text NOT null default '',
ADD COLUMN IF NOT EXISTS channel_code text NOT null default '',
ADD COLUMN IF NOT EXISTS topup_id text NOT null default '';
ALTER TABLE top_up_datas
DROP COLUMN if exists merchant_id;
ALTER TABLE top_up_datas
ADD merchant_id text NOT null default '';
ALTER TABLE customers
ADD COLUMN IF NOT EXISTS customer_internalid BIGINT NOT null default 0;
ALTER TABLE merchants
ADD COLUMN IF NOT EXISTS customer_internalid BIGINT NOT null default 0,
ADD COLUMN IF NOT EXISTS merchant_internalid BIGINT NOT null default 0;
ALTER TABLE channels
ADD COLUMN IF NOT EXISTS customer_internalid BIGINT NOT null default 0,
ADD COLUMN IF NOT EXISTS merchant_internalid BIGINT NOT null default 0,
ADD COLUMN IF NOT EXISTS channel_internalid BIGINT NOT null default 0;
