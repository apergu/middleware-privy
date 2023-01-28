ALTER TABLE customers
ADD COLUMN "crm_lead_id" text not null default '',
ADD COLUMN "enterprise_privy_id" text not null default '';