CREATE TABLE comments (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL,
    thread_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    user_id UUID NOT NULL,
    body TEXT NOT NULL,
    FOREIGN KEY (document_id) REFERENCES documents(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_comments_document_id ON comments(document_id);
CREATE INDEX idx_comments_user_thread ON comments(user_id, thread_id);