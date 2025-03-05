ALTER TABLE documents DROP CONSTRAINT fk_root_parent_id;

ALTER TABLE documents DROP COLUMN root_parent_id;
ALTER TABLE documents DROP COLUMN parent_address;

-- Migration Down: Re-add CASCADE DELETE
ALTER TABLE documents DROP CONSTRAINT fk_parent_id;
ALTER TABLE documents ADD CONSTRAINT fk_parent_id
FOREIGN KEY (parent_id) REFERENCES documents(id) ON DELETE CASCADE;

DROP TRIGGER IF EXISTS set_root_parent_id_trigger ON documents;
DROP FUNCTION IF EXISTS set_default_root_parent_id();
