-- Make document_id non-nullable
ALTER TABLE public.document_access
ALTER COLUMN document_id SET NOT NULL;

-- Make user_id non-nullable
ALTER TABLE public.document_access
ALTER COLUMN user_id SET NOT NULL;
