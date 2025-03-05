-- query that adds root_parent_id string and parent_address string to the branch_copies table, both can be null
ALTER TABLE documents ADD COLUMN root_parent_id uuid NULL;
ALTER TABLE documents ADD COLUMN parent_address TEXT NULL;

UPDATE documents SET root_parent_id = id;
UPDATE documents SET parent_id = NULL;

ALTER TABLE documents ALTER COLUMN root_parent_id SET NOT NULL;
ALTER TABLE documents ADD CONSTRAINT fk_root_parent_id
FOREIGN KEY (root_parent_id) REFERENCES documents(id);

-- Migration Up: Remove CASCADE DELETE
ALTER TABLE documents DROP CONSTRAINT fk_parent_id;
ALTER TABLE documents ADD CONSTRAINT fk_parent_id
FOREIGN KEY (parent_id) REFERENCES documents(id);

-- Create a function to set root_parent_id
CREATE OR REPLACE FUNCTION set_default_root_parent_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.root_parent_id IS NULL THEN
       NEW.root_parent_id := NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to call the function before insert
CREATE TRIGGER set_root_parent_id_trigger
BEFORE INSERT ON documents
FOR EACH ROW
EXECUTE FUNCTION set_default_root_parent_id();
