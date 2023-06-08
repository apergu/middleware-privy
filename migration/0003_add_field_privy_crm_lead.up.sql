ALTER TABLE customers
ADD COLUMN "address_1" text not null default '',
ADD COLUMN "npwp" text not null default '',
ADD COLUMN "state" text not null default '',
ADD COLUMN "city" text not null default '',
ADD COLUMN "zip_code" text not null default '';