import { gql } from "@/__generated__";

export const timelineEventFragment = gql(`
fragment TimelineEventFields on TimelineEvent {
  id
  replyTo
  createdAt
  authorId
  user {
    id
    name
    picture
  }
  event {
    __typename
    ... on TLJoinV1 {
      action
    }
    ... on TLMarkerV1 {
      title
    }
    ... on TLUpdateV1 {
      eventId
      title
      content
      startingContentAddress
      endingContentAddress
      flaggedVersionName
      flaggedVersionCreatedAt
      flaggedVersionID
      state
      flaggedByUser {
        name
      }
    }
    ... on TLEmpty {
      placeholder 
    }
    ... on TLAttributeChangeV1 {
      attribute
      oldValue
      newValue
    }
    ... on TLAccessChangeV1 {
      action
      userIdentifiers
    }
    ... on TLPasteV1 {
      contentAddressBefore
      contentAddressAfter
    }
    ... on TLMessageV1 {
      eventId 
      content
      contentAddress
      selectionStartId
      selectionEndId
      selectionMarkdown

      replies {
        id
        replyTo
        createdAt
        authorId
        user {
          id
          name
          picture
        }
        event {
          __typename
          ... on TLMessageV1 {
            eventId 
            content
            contentAddress
            selectionStartId
            selectionEndId
            selectionMarkdown
          }
          ... on TLMessageResolutionV1 {
            eventId
            resolutionSummary
            resolved
          }
        }
      }
    }
  }
}`);

export const GetDocumentTimeline = gql(`
  query GetDocumentTimeline($documentId: ID!, $filter: TimelineEventFilter) {
    getDocumentTimeline(documentId: $documentId, filter: $filter) {
      ...TimelineEventFields
    }
  }
`);

export const TimelineEventInserted = gql(`
  subscription TimelineEventInserted($documentId: ID!) {
    timelineEventInserted(documentId: $documentId) {
      ...TimelineEventFields
    }
  }
`);

export const TimelineEventUpdated = gql(`
  subscription TimelineEventUpdated($documentId: ID!) {
    timelineEventUpdated(documentId: $documentId) {
      ...TimelineEventFields
    }
  }
`);

export const TimelineEventDeleted = gql(`
  subscription TimelineEventDeleted($documentId: ID!) {
    timelineEventDeleted(documentId: $documentId) {
      id
      replyTo
      event {
        ... on TLMessageV1 {
          eventId
        }
      }
    }
  }
`);

export const CreateTimelineMessage = gql(`
  mutation CreateTimelineMessage($documentId: ID!, $input: TimelineMessageInput!) {
    createTimelineMessage(documentId: $documentId, input: $input) {
      id
    }
  }
`);

export const EditTimelineMessage = gql(`
  mutation EditTimelineMessage($documentId: ID!, $messageId: ID!, $input: EditTimelineMessageInput!) {
    editTimelineMessage(documentId: $documentId, messageId: $messageId, input: $input) {
      id
    }
  }
`);

export const UpdateMessageResolution = gql(`
  mutation UpdateMessageResolution($documentId: ID!, $messageId: ID!, $input: UpdateMessageResolutionInput!) {
    updateMessageResolution(documentId: $documentId, messageId: $messageId, input: $input) {
      id
    }
  }
`);

export const EditTimelineUpdateSummary = gql(`
  mutation EditTimelineUpdateSummary($documentId: ID!, $updateId: ID!, $summary: String!) {
    editTimelineUpdateSummary(documentId: $documentId, updateId: $updateId, summary: $summary) {
      id
      event {
        ... on TLUpdateV1 {
          eventId
          content
        }
      }
    }
  }
`);

export const EditMessageResolutionSummary = gql(`
  mutation EditMessageResolutionSummary($documentId: ID!, $messageId: ID!, $summary: String!) {
    editMessageResolutionSummary(documentId: $documentId, messageId: $messageId, summary: $summary) {
      id
      event {
      ... on TLMessageResolutionV1 {
        eventId
          resolutionSummary
        }
      }
    }
  }
`);

export const ForceTimelineUpdateSummary = gql(`
  mutation ForceTimelineUpdateSummary($documentId: ID!, $userId: String!) {
    forceTimelineUpdateSummary(documentId: $documentId, userId: $userId)      
  }
`);

export const DeleteTimelineMessage = gql(`
  mutation DeleteTimelineMessage($documentId: ID!, $messageId: ID!) {
    deleteTimelineMessage(documentId: $documentId, messageId: $messageId)
  }
`);
