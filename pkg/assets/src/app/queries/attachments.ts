import { gql } from "@/__generated__";

export const UploadAttachment = gql(`
mutation UploadAttachment($file: Upload!, $docId: ID!) {
  uploadAttachment(file: $file, docId: $docId) {
    id
    filename
    contentType
    createdAt
  }
}
`);

export const ListDocumentAttachments = gql(`
query ListDocumentAttachments($docId: ID!) {
  listDocumentAttachments(docId: $docId) {
    id
    filename
    contentType
    createdAt
  }
}
`);

export const ListUsersAttachments = gql(`
query ListUsersAttachments {
  listUsersAttachments {
    id
    filename
    contentType
    createdAt
  }
}
`);
