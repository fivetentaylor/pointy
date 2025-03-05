import { gql } from "@/__generated__";

export const GetAIThreads = gql(`
  query GetAIThreads($documentId: ID!) {
    getAskAiThreads(documentId: $documentId) {
      id
      title
      updatedAt

      user {
        id
        name
      }
    }
  }
`);

export const messageFragment = gql(`
fragment MessageFields on Message {
  id
  containerId
  content
  createdAt
  lifecycleStage
  lifecycleReason
  authorId
  hidden
  user {
    id
    name
    picture
  }
  aiContent {
    concludingMessage
    feedback
  }
  metadata {
    allowDraftEdits
    contentAddressBefore
    contentAddress
    contentAddressAfter
    contentAddressAfterTimestamp
    revisionStatus
  }
  attachments {
    __typename
    ... on Selection {
      start
      end
      content
    }
    ... on Revision {
      start
      end
      updated
      beforeAddress
      afterAddress
      appliedOps
    }
    ... on Suggestion {
      content
    }
    ... on AttachmentContent {
      text
    }
    ... on AttachmentError {
      title
      text
    }
    ... on AttachmentFile {
      id
      filename
      contentType
    }
    ... on AttachedRevisoDocument {
      id
      title
    }
  }
}`);

export const GetAIThreadMessages = gql(`
  query GetAIThreadMessages($documentId: ID!, $threadId: ID!) {
    getAskAiThreadMessages(documentId: $documentId, threadId: $threadId) {
      ...MessageFields
    }
  }
`);

export const MessageUpserted = gql(`
subscription MessageUpserted($documentId: ID!, $channelId: ID!) {
  messageUpserted(documentId: $documentId, channelId: $channelId) {
    ...MessageFields
  }
}
`);

export const threadFragment = gql(`
fragment ThreadFields on Thread {
  id
  title
  updatedAt
}`);

export const UpdateMessageRevisionStatus = gql(`
  mutation UpdateMessageRevisionStatus($containerId: ID!, $messageId: ID!, $status: MessageRevisionStatus!, $contentAddress: String!) {
    updateMessageRevisionStatus(containerId: $containerId, messageId: $messageId, status: $status, contentAddress: $contentAddress) {
      ...MessageFields
    }
  }
`);

export const CreateAIThread = gql(`
  mutation CreateAIThread($documentId: ID!) {
    createAskAiThread(documentId: $documentId) {
      id
    }
  }
`);

export const CreateAIThreadMessage = gql(`
  mutation CreateAIThreadMessage($documentId: ID!, $threadId: ID!, $input: MessageInput!) {
    createAskAiThreadMessage(documentId: $documentId, threadId: $threadId, input: $input) {
      ...MessageFields
    }
  }
`);

export const ThreadUpserted = gql(`
subscription ThreadUpserted($documentId: ID!) {
  threadUpserted(documentId: $documentId) {
    ...ThreadFields
  }
}
`);
