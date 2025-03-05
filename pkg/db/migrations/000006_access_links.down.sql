-- Drop the shared_document_links table
DROP TABLE shared_document_links;

-- Drop the wrapper function for default access link length
DROP FUNCTION generate_default_invite_link();

-- Drop the function to generate a custom access link
DROP FUNCTION generate_access_link(INT);
