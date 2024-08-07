ALTER TABLE top_up_datas
ADD COLUMN IF NOT EXISTS transaction_type SMALLINT NOT null default 1;