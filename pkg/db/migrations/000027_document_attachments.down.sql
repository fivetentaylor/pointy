-- Drop indexes if they exist
DROP INDEX IF EXISTS idx_document_attachments_user_id;
DROP INDEX IF EXISTS idx_document_attachments_document_id;

-- Drop the document_attachments table
DROP TABLE IF EXISTS document_attachments CASCADE;
