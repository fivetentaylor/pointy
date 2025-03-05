ALTER TABLE public.documents
ADD COLUMN is_public boolean NOT NULL DEFAULT true;
