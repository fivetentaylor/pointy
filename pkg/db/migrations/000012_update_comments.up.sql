ALTER TABLE comments
ADD COLUMN selection text,
ADD COLUMN selection_start text,
ADD COLUMN selection_finish text,
ADD COLUMN notes text;

INSERT INTO public.users (id, name, email, provider, password_hash)
VALUES ('00000000-0000-0000-0000-000000000000', 'Reviso', 'reviso@revi.so', 'manual', NULL)
