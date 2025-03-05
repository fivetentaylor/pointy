-- Drop the foreign key constraint
ALTER TABLE public.documents
DROP CONSTRAINT fk_parent_id;

-- Drop the parent_id column
ALTER TABLE public.documents
DROP COLUMN parent_id;
