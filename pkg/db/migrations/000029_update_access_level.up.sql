ALTER TABLE public.document_access DROP CONSTRAINT document_access_access_level_check;

ALTER TABLE public.document_access
ADD CONSTRAINT document_access_access_level_check
CHECK (access_level = ANY (ARRAY['comment', 'write', 'owner', 'admin']));
