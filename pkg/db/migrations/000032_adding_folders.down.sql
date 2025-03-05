DROP INDEX idx_documents_parent_folder;
DROP INDEX idx_documents_is_folder;
DROP INDEX idx_documents_folder_id;
ALTER TABLE documents
DROP COLUMN folder_id,
DROP COLUMN is_folder;
