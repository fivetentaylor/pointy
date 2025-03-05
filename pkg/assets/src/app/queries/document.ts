import { gql } from "@/__generated__";

export const DocumentFragment = gql(`
  fragment DocumentFields on Document {
    id
    title
    isPublic
    isFolder
    folderID
    updatedAt
    access
    ownedBy {
      id
    }
    editors {
      id
      name
      displayName
      email
      picture
    }
    preferences {
      enableFirstOpenNotifications
      enableMentionNotifications
      enableDMNotifications
      enableAllCommentNotifications
    }
  }
`);

export const GetDocument = gql(`
  query GetDocument($id: ID!) {
    document(id: $id) {
      ...DocumentFields
    }
  }
`);

export const UpdateDocumentTitle = gql(`
  mutation UpdateDocumentTitle($id: ID!, $input: DocumentInput!) {
    updateDocument(id: $id, input: $input) {
      id
      title
    }
  }
`);

export const GetDocuments = gql(`
  query GetDocuments($offset: Int, $limit: Int) {
    documents(offset: $offset, limit: $limit) {
      totalCount
      edges {
        node {
          id
          title
          isFolder  
          folderID
          updatedAt
          ownedBy {
            id
          }
        }
      }
      pageInfo {
        hasNextPage
      }
    }
  }
`);

export const GetBaseDocuments = gql(`
  query GetBaseDocuments($offset: Int, $limit: Int) {
    baseDocuments(offset: $offset, limit: $limit) {
      totalCount
      edges {
        node {
          id
          title
          isFolder  
          folderID
          updatedAt
          ownedBy {
            id
          }
        }
      }
      pageInfo {
        hasNextPage
      }
    }
  }
`);

export const GetSharedDocuments = gql(`
  query GetSharedDocuments($offset: Int, $limit: Int) {
    sharedDocuments(offset: $offset, limit: $limit) {
      totalCount
      edges {
        node {
          id
          title
          isFolder  
          folderID
          updatedAt
          ownedBy {
            id
          }
        }
      }
      pageInfo {
        hasNextPage
      }
    }
  }
`);

export const GetFolderDocuments = gql(`
  query GetFolderDocuments($folderId: ID!, $offset: Int, $limit: Int) {
    folderDocuments(folderID: $folderId, offset: $offset, limit: $limit) {
      totalCount
      edges {
        node {
          id
          title
          isFolder
          folderID
          updatedAt
          ownedBy {
            id
          }
        }
      }
      pageInfo {
        hasNextPage
      }
    }
  }
`);

export const SearchDocuments = gql(`
  query SearchDocuments($query: String!, $offset: Int, $limit: Int) {
    searchDocuments(query: $query, offset: $offset, limit: $limit) {
      totalCount
      edges {
        node {
          id
          title
          isFolder  
          folderID
          updatedAt
          ownedBy {
            id
          }
        }
      }
      pageInfo {
        hasNextPage
      }
    }
  }
`);

export const CreateDocument = gql(`
  mutation CreateDocument {
    createDocument {
      id
    }
  }
`);

export const CreateFolder = gql(`
  mutation CreateFolder {
    createFolder {
      ...DocumentFields
    }
  }
`);

export const DeleteDocument = gql(`
  mutation DeleteDocument($id: ID!, $deleteChildren: Boolean) {
    deleteDocument(id: $id, deleteChildren: $deleteChildren)
  }
`);

export const ShareDocument = gql(`
  mutation ShareDocument($id: ID!, $emails: [String!]!, $message: String) {
    shareDocument(documentID: $id, emails: $emails, message: $message) {
      inviteLink
    }
  }
`);

export const UnshareDocument = gql(`
  mutation UnshareDocument($docId: ID!, $editorId: ID!) {
    unshareDocument(documentID: $docId, editorID: $editorId) {
       id
       editors {
        id
        name
        email
        picture
      }
    }
  }
`);

export const UpdateSharedLink = gql(`
  mutation UpdateSharedLink($inviteLink: String!, $isActive: Boolean!) {
    updateShareLink(inviteLink: $inviteLink, isActive: $isActive) {
      inviteLink
    }
  }
`);

export const UpdateDocumentVisibility = gql(`
  mutation UpdateDocumentVisibility($id: ID!, $input: DocumentInput!) {
    updateDocument(id: $id, input: $input) {
      id
      isPublic
    }
  }
`);

export const SharedDocumentLinks = gql(`
  query SharedDocumentLinks($id: ID!) {
    sharedLinks(documentID: $id) {
      inviteLink
      inviteeEmail
      invitedBy {
        name
      }
      isActive
    }
  }
`);

export const UpdateDocumentPreferences = gql(`
mutation UpdateDocumentPreferences($documentId: ID!, $input: DocumentPreferenceInput!) {
  updateDocumentPreference(id: $documentId, input: $input) {
      enableFirstOpenNotifications
      enableMentionNotifications
      enableDMNotifications
      enableAllCommentNotifications
  }
}
`);

export const CreateFlaggedVersion = gql(`
  mutation CreateFlaggedVersion($documentId: ID!, $input: FlaggedVersionInput!) {
    createFlaggedVersion(documentId: $documentId, input: $input)
  }
`);

export const EditFlaggedVersion = gql(`
  mutation EditFlaggedVersion($flaggedVersionId: ID!, $input: FlaggedVersionInput!) {
    editFlaggedVersion(flaggedVersionId: $flaggedVersionId, input: $input)
  }
`);

export const DeleteFlaggedVersion = gql(`
  mutation DeleteFlaggedVersion($flaggedVersionId: ID!, $timelineEventId: ID!) {
    deleteFlaggedVersion(flaggedVersionId: $flaggedVersionId, timelineEventId: $timelineEventId)
  }
`);

export const MoveDocument = gql(`
  mutation MoveDocument($documentId: ID!, $folderId: ID) {
    moveDocument(id: $documentId, folderID: $folderId) {
      ...DocumentFields
    }
  }
`);

export const DocumentUpdatedSubscription = gql(`
  subscription DocumentUpdated($documentId: ID!) {
    documentUpdated(documentId: $documentId) {
      ...DocumentFields
    }
  }
`);

export const DocumentInsertedSubscription = gql(`
  subscription DocumentInserted($userId: ID!) {
    documentInserted(userId: $userId) {
      ...DocumentFields
    }
  }
`);
