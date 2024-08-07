ALTER TABLE merchants
DROP CONSTRAINT if exists merchants_enterprise_id_key;
ALTER TABLE merchants
DROP CONSTRAINT if exists fk_customer;
-- ALTER TABLE merchants
-- DROP COLUMN if exists enterprise_id;
ALTER TABLE channels
DROP CONSTRAINT if exists fk_merchant;
ALTER TABLE divissions
DROP CONSTRAINT if exists fk_channel;
ALTER TABLE channels
DROP COLUMN if exists merchant_id;
ALTER TABLE channels
ADD merchant_id text NOT null default '';