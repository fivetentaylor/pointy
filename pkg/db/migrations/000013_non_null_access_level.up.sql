UPDATE public.document_access
SET access_level = 'comment'
WHERE access_level IS NULL;

ALTER TABLE public.document_access
ALTER COLUMN access_level SET NOT NULL;
