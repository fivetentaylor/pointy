-- Remove the trigger
DROP TRIGGER IF EXISTS update_document_versions_updated_at ON public.document_versions;

-- Remove the index
DROP INDEX IF EXISTS idx_document_versions_document_id;

-- Remove foreign key constraints
ALTER TABLE IF EXISTS public.document_versions
    DROP CONSTRAINT IF EXISTS document_versions_updated_by_fkey;

ALTER TABLE IF EXISTS public.document_versions
    DROP CONSTRAINT IF EXISTS document_versions_created_by_fkey;

ALTER TABLE IF EXISTS public.document_versions
    DROP CONSTRAINT IF EXISTS document_versions_document_id_fkey;

-- Remove primary key constraint
ALTER TABLE IF EXISTS public.document_versions
    DROP CONSTRAINT IF EXISTS document_versions_pkey;

-- Drop the document_versions table
DROP TABLE IF EXISTS public.document_versions;
