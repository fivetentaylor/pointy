-- Make document_id nullable
ALTER TABLE public.document_access
ALTER COLUMN document_id DROP NOT NULL;

-- Make user_id nullable
ALTER TABLE public.document_access
ALTER COLUMN user_id DROP NOT NULL;
