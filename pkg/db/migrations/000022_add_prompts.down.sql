DROP TRIGGER IF EXISTS update_prompts_updated_at ON prompts;
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP INDEX IF EXISTS idx_prompts_prompt_name;
DROP TABLE IF EXISTS prompts;
