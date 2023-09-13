ALTER TABLE customer_usages
ADD COLUMN IF NOT EXISTS enterprise_id text NOT null default '',
ADD COLUMN IF NOT EXISTS enterprise_name text NOT null default '',
ADD COLUMN IF NOT EXISTS channel_name text NOT null default '',
ADD COLUMN IF NOT EXISTS trx_id text NOT null default '',
ADD COLUMN IF NOT EXISTS service_id text NOT null default '',
ADD COLUMN IF NOT EXISTS unit_price text NOT null default '';