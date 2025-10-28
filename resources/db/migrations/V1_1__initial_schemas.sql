ALTER TABLE public.logs
ALTER COLUMN user_id TYPE varchar(255) USING user_id::varchar;

ALTER TABLE public.logs
ADD COLUMN tenant_cat jsonb NULL;