-- Add the parent_id column
ALTER TABLE public.documents
ADD COLUMN parent_id uuid;

-- Add the foreign key constraint
ALTER TABLE public.documents
ADD CONSTRAINT fk_parent_id
FOREIGN KEY (parent_id)
REFERENCES public.documents(id)
ON DELETE CASCADE;
