CREATE TABLE author_ids (
    id SERIAL PRIMARY KEY,
    author_id INT NOT NULL,
    document_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(author_id, document_id)
);

CREATE OR REPLACE FUNCTION increment_author_id(doc_id UUID)
RETURNS INT AS $$
DECLARE
    new_author_id INT;
BEGIN
    SELECT COALESCE(MAX(author_id), 0) + 1 INTO new_author_id
    FROM author_ids WHERE document_id = doc_id;

    RETURN new_author_id;
END;
$$ LANGUAGE plpgsql;

CREATE INDEX idx_author_document ON author_ids (author_id, document_id);
