ALTER TABLE merchants
DROP CONSTRAINT merchants_enterprise_id_key;

ALTER TABLE merchants
DROP CONSTRAINT fk_customer;

ALTER TABLE channels
DROP CONSTRAINT fk_merchant;

ALTER TABLE divissions
DROP CONSTRAINT fk_channel;

ALTER TABLE channels
DROP COLUMN merchant_id;

ALTER TABLE channels
ADD merchant_id text NOT NULL;
