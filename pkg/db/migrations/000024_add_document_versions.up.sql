-- Create the document_versions table
CREATE TABLE public.document_versions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    document_id uuid NOT NULL,
    name text NOT NULL,
    content_address text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL
);

-- Add primary key constraint
ALTER TABLE ONLY public.document_versions
    ADD CONSTRAINT document_versions_pkey PRIMARY KEY (id);

-- Add foreign key constraints
ALTER TABLE ONLY public.document_versions
    ADD CONSTRAINT document_versions_document_id_fkey FOREIGN KEY (document_id) REFERENCES public.documents(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.document_versions
    ADD CONSTRAINT document_versions_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);

ALTER TABLE ONLY public.document_versions
    ADD CONSTRAINT document_versions_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.users(id);

-- Create an index on document_id for faster lookups
CREATE INDEX idx_document_versions_document_id ON public.document_versions(document_id);

-- Add a trigger to update the updated_at timestamp
CREATE TRIGGER update_document_versions_updated_at
BEFORE UPDATE ON public.document_versions
FOR EACH ROW
EXECUTE FUNCTION public.update_updated_at_column();
