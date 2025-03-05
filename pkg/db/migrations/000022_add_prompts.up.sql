CREATE TABLE prompts (
    id SERIAL PRIMARY KEY,
    prompt_name VARCHAR(255) NOT NULL,
    prompt_text TEXT NOT NULL,
    version TEXT NOT NULL DEFAULT '',
    provider VARCHAR(255) NOT NULL,
    model_name VARCHAR(255) NOT NULL, 
    temperature FLOAT DEFAULT 1.0,
    max_tokens INT DEFAULT 0,
    top_p FLOAT DEFAULT 1.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_prompts_prompt_name ON prompts (prompt_name);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to automatically update the updated_at column on updates
CREATE TRIGGER update_prompts_updated_at
BEFORE UPDATE ON prompts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
