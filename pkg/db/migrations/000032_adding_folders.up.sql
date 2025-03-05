ALTER TABLE documents
ADD COLUMN is_folder BOOLEAN DEFAULT FALSE,
ADD COLUMN folder_id UUID REFERENCES documents(id);

CREATE INDEX idx_documents_folder_id ON documents(folder_id);
CREATE INDEX idx_documents_is_folder ON documents(is_folder);
CREATE INDEX idx_documents_parent_folder ON documents(parent_id) WHERE is_folder = true;
