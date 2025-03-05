/* eslint-disable */
import * as types from "./graphql";
import type { TypedDocumentNode as DocumentNode } from "@graphql-typed-document-node/core";

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
  "\nmutation UploadAttachment($file: Upload!, $docId: ID!) {\n  uploadAttachment(file: $file, docId: $docId) {\n    id\n    filename\n    contentType\n    createdAt\n  }\n}\n":
    types.UploadAttachmentDocument,
  "\nquery ListDocumentAttachments($docId: ID!) {\n  listDocumentAttachments(docId: $docId) {\n    id\n    filename\n    contentType\n    createdAt\n  }\n}\n":
    types.ListDocumentAttachmentsDocument,
  "\nquery ListUsersAttachments {\n  listUsersAttachments {\n    id\n    filename\n    contentType\n    createdAt\n  }\n}\n":
    types.ListUsersAttachmentsDocument,
  "\n  fragment DocumentFields on Document {\n    id\n    title\n    isPublic\n    isFolder\n    folderID\n    updatedAt\n    access\n    ownedBy {\n      id\n    }\n    editors {\n      id\n      name\n      displayName\n      email\n      picture\n    }\n    preferences {\n      enableFirstOpenNotifications\n      enableMentionNotifications\n      enableDMNotifications\n      enableAllCommentNotifications\n    }\n  }\n":
    types.DocumentFieldsFragmentDoc,
  "\n  query GetDocument($id: ID!) {\n    document(id: $id) {\n      ...DocumentFields\n    }\n  }\n":
    types.GetDocumentDocument,
  "\n  mutation UpdateDocumentTitle($id: ID!, $input: DocumentInput!) {\n    updateDocument(id: $id, input: $input) {\n      id\n      title\n    }\n  }\n":
    types.UpdateDocumentTitleDocument,
  "\n  query GetDocuments($offset: Int, $limit: Int) {\n    documents(offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n":
    types.GetDocumentsDocument,
  "\n  query GetBaseDocuments($offset: Int, $limit: Int) {\n    baseDocuments(offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n":
    types.GetBaseDocumentsDocument,
  "\n  query GetSharedDocuments($offset: Int, $limit: Int) {\n    sharedDocuments(offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n":
    types.GetSharedDocumentsDocument,
  "\n  query GetFolderDocuments($folderId: ID!, $offset: Int, $limit: Int) {\n    folderDocuments(folderID: $folderId, offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder\n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n":
    types.GetFolderDocumentsDocument,
  "\n  query SearchDocuments($query: String!, $offset: Int, $limit: Int) {\n    searchDocuments(query: $query, offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n":
    types.SearchDocumentsDocument,
  "\n  mutation CreateDocument {\n    createDocument {\n      id\n    }\n  }\n":
    types.CreateDocumentDocument,
  "\n  mutation CreateFolder {\n    createFolder {\n      ...DocumentFields\n    }\n  }\n":
    types.CreateFolderDocument,
  "\n  mutation DeleteDocument($id: ID!, $deleteChildren: Boolean) {\n    deleteDocument(id: $id, deleteChildren: $deleteChildren)\n  }\n":
    types.DeleteDocumentDocument,
  "\n  mutation ShareDocument($id: ID!, $emails: [String!]!, $message: String) {\n    shareDocument(documentID: $id, emails: $emails, message: $message) {\n      inviteLink\n    }\n  }\n":
    types.ShareDocumentDocument,
  "\n  mutation UnshareDocument($docId: ID!, $editorId: ID!) {\n    unshareDocument(documentID: $docId, editorID: $editorId) {\n       id\n       editors {\n        id\n        name\n        email\n        picture\n      }\n    }\n  }\n":
    types.UnshareDocumentDocument,
  "\n  mutation UpdateSharedLink($inviteLink: String!, $isActive: Boolean!) {\n    updateShareLink(inviteLink: $inviteLink, isActive: $isActive) {\n      inviteLink\n    }\n  }\n":
    types.UpdateSharedLinkDocument,
  "\n  mutation UpdateDocumentVisibility($id: ID!, $input: DocumentInput!) {\n    updateDocument(id: $id, input: $input) {\n      id\n      isPublic\n    }\n  }\n":
    types.UpdateDocumentVisibilityDocument,
  "\n  query SharedDocumentLinks($id: ID!) {\n    sharedLinks(documentID: $id) {\n      inviteLink\n      inviteeEmail\n      invitedBy {\n        name\n      }\n      isActive\n    }\n  }\n":
    types.SharedDocumentLinksDocument,
  "\nmutation UpdateDocumentPreferences($documentId: ID!, $input: DocumentPreferenceInput!) {\n  updateDocumentPreference(id: $documentId, input: $input) {\n      enableFirstOpenNotifications\n      enableMentionNotifications\n      enableDMNotifications\n      enableAllCommentNotifications\n  }\n}\n":
    types.UpdateDocumentPreferencesDocument,
  "\n  mutation CreateFlaggedVersion($documentId: ID!, $input: FlaggedVersionInput!) {\n    createFlaggedVersion(documentId: $documentId, input: $input)\n  }\n":
    types.CreateFlaggedVersionDocument,
  "\n  mutation EditFlaggedVersion($flaggedVersionId: ID!, $input: FlaggedVersionInput!) {\n    editFlaggedVersion(flaggedVersionId: $flaggedVersionId, input: $input)\n  }\n":
    types.EditFlaggedVersionDocument,
  "\n  mutation DeleteFlaggedVersion($flaggedVersionId: ID!, $timelineEventId: ID!) {\n    deleteFlaggedVersion(flaggedVersionId: $flaggedVersionId, timelineEventId: $timelineEventId)\n  }\n":
    types.DeleteFlaggedVersionDocument,
  "\n  mutation MoveDocument($documentId: ID!, $folderId: ID) {\n    moveDocument(id: $documentId, folderID: $folderId) {\n      ...DocumentFields\n    }\n  }\n":
    types.MoveDocumentDocument,
  "\n  subscription DocumentUpdated($documentId: ID!) {\n    documentUpdated(documentId: $documentId) {\n      ...DocumentFields\n    }\n  }\n":
    types.DocumentUpdatedDocument,
  "\n  subscription DocumentInserted($userId: ID!) {\n    documentInserted(userId: $userId) {\n      ...DocumentFields\n    }\n  }\n":
    types.DocumentInsertedDocument,
  "\n  mutation UploadImage($file: Upload!, $docId: ID!) {\n    uploadImage(file: $file, docId: $docId) {\n      id\n      docId\n      url\n      createdAt\n      mimeType\n      status\n      error\n    }\n  }\n":
    types.UploadImageDocument,
  "\n  query ListDocumentImages($docId: ID!) {\n    listDocumentImages(docId: $docId) {\n      id\n      docId\n      url\n      createdAt\n      mimeType\n      status\n      error\n    }\n  }\n":
    types.ListDocumentImagesDocument,
  "\n  query GetImageSignedUrl($docId: ID!, $imageId: ID!) {\n    getImageSignedUrl(docId: $docId, imageId: $imageId) {\n      url\n      expiresAt\n    }\n  }\n":
    types.GetImageSignedUrlDocument,
  "\n  query GetImage($docId: ID!, $imageId: ID!) {\n    getImage(docId: $docId, imageId: $imageId) {\n      id\n      docId\n      url\n      createdAt\n      mimeType\n      status\n      error\n    }\n  }\n":
    types.GetImageDocument,
  "\n  query GetAIThreads($documentId: ID!) {\n    getAskAiThreads(documentId: $documentId) {\n      id\n      title\n      updatedAt\n\n      user {\n        id\n        name\n      }\n    }\n  }\n":
    types.GetAiThreadsDocument,
  "\nfragment MessageFields on Message {\n  id\n  containerId\n  content\n  createdAt\n  lifecycleStage\n  lifecycleReason\n  authorId\n  hidden\n  user {\n    id\n    name\n    picture\n  }\n  aiContent {\n    concludingMessage\n    feedback\n  }\n  metadata {\n    allowDraftEdits\n    contentAddressBefore\n    contentAddress\n    contentAddressAfter\n    contentAddressAfterTimestamp\n    revisionStatus\n  }\n  attachments {\n    __typename\n    ... on Selection {\n      start\n      end\n      content\n    }\n    ... on Revision {\n      start\n      end\n      updated\n      beforeAddress\n      afterAddress\n      appliedOps\n    }\n    ... on Suggestion {\n      content\n    }\n    ... on AttachmentContent {\n      text\n    }\n    ... on AttachmentError {\n      title\n      text\n    }\n    ... on AttachmentFile {\n      id\n      filename\n      contentType\n    }\n    ... on AttachedRevisoDocument {\n      id\n      title\n    }\n  }\n}":
    types.MessageFieldsFragmentDoc,
  "\n  query GetAIThreadMessages($documentId: ID!, $threadId: ID!) {\n    getAskAiThreadMessages(documentId: $documentId, threadId: $threadId) {\n      ...MessageFields\n    }\n  }\n":
    types.GetAiThreadMessagesDocument,
  "\nsubscription MessageUpserted($documentId: ID!, $channelId: ID!) {\n  messageUpserted(documentId: $documentId, channelId: $channelId) {\n    ...MessageFields\n  }\n}\n":
    types.MessageUpsertedDocument,
  "\nfragment ThreadFields on Thread {\n  id\n  title\n  updatedAt\n}":
    types.ThreadFieldsFragmentDoc,
  "\n  mutation UpdateMessageRevisionStatus($containerId: ID!, $messageId: ID!, $status: MessageRevisionStatus!, $contentAddress: String!) {\n    updateMessageRevisionStatus(containerId: $containerId, messageId: $messageId, status: $status, contentAddress: $contentAddress) {\n      ...MessageFields\n    }\n  }\n":
    types.UpdateMessageRevisionStatusDocument,
  "\n  mutation CreateAIThread($documentId: ID!) {\n    createAskAiThread(documentId: $documentId) {\n      id\n    }\n  }\n":
    types.CreateAiThreadDocument,
  "\n  mutation CreateAIThreadMessage($documentId: ID!, $threadId: ID!, $input: MessageInput!) {\n    createAskAiThreadMessage(documentId: $documentId, threadId: $threadId, input: $input) {\n      ...MessageFields\n    }\n  }\n":
    types.CreateAiThreadMessageDocument,
  "\nsubscription ThreadUpserted($documentId: ID!) {\n  threadUpserted(documentId: $documentId) {\n    ...ThreadFields\n  }\n}\n":
    types.ThreadUpsertedDocument,
  "\n  query SubscriptionPlans {\n    subscriptionPlans {\n      id\n      name\n      priceCents\n      currency\n      interval\n    }\n  }\n":
    types.SubscriptionPlansDocument,
  "\n  mutation Checkout($id: ID!) {\n    checkoutSubscriptionPlan(id: $id) {\n      url\n    }\n  }\n":
    types.CheckoutDocument,
  "\n  mutation BillingPortalSession {\n    billingPortalSession {\n      url\n    }\n  }\n":
    types.BillingPortalSessionDocument,
  "\nfragment TimelineEventFields on TimelineEvent {\n  id\n  replyTo\n  createdAt\n  authorId\n  user {\n    id\n    name\n    picture\n  }\n  event {\n    __typename\n    ... on TLJoinV1 {\n      action\n    }\n    ... on TLMarkerV1 {\n      title\n    }\n    ... on TLUpdateV1 {\n      eventId\n      title\n      content\n      startingContentAddress\n      endingContentAddress\n      flaggedVersionName\n      flaggedVersionCreatedAt\n      flaggedVersionID\n      state\n      flaggedByUser {\n        name\n      }\n    }\n    ... on TLEmpty {\n      placeholder \n    }\n    ... on TLAttributeChangeV1 {\n      attribute\n      oldValue\n      newValue\n    }\n    ... on TLAccessChangeV1 {\n      action\n      userIdentifiers\n    }\n    ... on TLPasteV1 {\n      contentAddressBefore\n      contentAddressAfter\n    }\n    ... on TLMessageV1 {\n      eventId \n      content\n      contentAddress\n      selectionStartId\n      selectionEndId\n      selectionMarkdown\n\n      replies {\n        id\n        replyTo\n        createdAt\n        authorId\n        user {\n          id\n          name\n          picture\n        }\n        event {\n          __typename\n          ... on TLMessageV1 {\n            eventId \n            content\n            contentAddress\n            selectionStartId\n            selectionEndId\n            selectionMarkdown\n          }\n          ... on TLMessageResolutionV1 {\n            eventId\n            resolutionSummary\n            resolved\n          }\n        }\n      }\n    }\n  }\n}":
    types.TimelineEventFieldsFragmentDoc,
  "\n  query GetDocumentTimeline($documentId: ID!, $filter: TimelineEventFilter) {\n    getDocumentTimeline(documentId: $documentId, filter: $filter) {\n      ...TimelineEventFields\n    }\n  }\n":
    types.GetDocumentTimelineDocument,
  "\n  subscription TimelineEventInserted($documentId: ID!) {\n    timelineEventInserted(documentId: $documentId) {\n      ...TimelineEventFields\n    }\n  }\n":
    types.TimelineEventInsertedDocument,
  "\n  subscription TimelineEventUpdated($documentId: ID!) {\n    timelineEventUpdated(documentId: $documentId) {\n      ...TimelineEventFields\n    }\n  }\n":
    types.TimelineEventUpdatedDocument,
  "\n  subscription TimelineEventDeleted($documentId: ID!) {\n    timelineEventDeleted(documentId: $documentId) {\n      id\n      replyTo\n      event {\n        ... on TLMessageV1 {\n          eventId\n        }\n      }\n    }\n  }\n":
    types.TimelineEventDeletedDocument,
  "\n  mutation CreateTimelineMessage($documentId: ID!, $input: TimelineMessageInput!) {\n    createTimelineMessage(documentId: $documentId, input: $input) {\n      id\n    }\n  }\n":
    types.CreateTimelineMessageDocument,
  "\n  mutation EditTimelineMessage($documentId: ID!, $messageId: ID!, $input: EditTimelineMessageInput!) {\n    editTimelineMessage(documentId: $documentId, messageId: $messageId, input: $input) {\n      id\n    }\n  }\n":
    types.EditTimelineMessageDocument,
  "\n  mutation UpdateMessageResolution($documentId: ID!, $messageId: ID!, $input: UpdateMessageResolutionInput!) {\n    updateMessageResolution(documentId: $documentId, messageId: $messageId, input: $input) {\n      id\n    }\n  }\n":
    types.UpdateMessageResolutionDocument,
  "\n  mutation EditTimelineUpdateSummary($documentId: ID!, $updateId: ID!, $summary: String!) {\n    editTimelineUpdateSummary(documentId: $documentId, updateId: $updateId, summary: $summary) {\n      id\n      event {\n        ... on TLUpdateV1 {\n          eventId\n          content\n        }\n      }\n    }\n  }\n":
    types.EditTimelineUpdateSummaryDocument,
  "\n  mutation EditMessageResolutionSummary($documentId: ID!, $messageId: ID!, $summary: String!) {\n    editMessageResolutionSummary(documentId: $documentId, messageId: $messageId, summary: $summary) {\n      id\n      event {\n      ... on TLMessageResolutionV1 {\n        eventId\n          resolutionSummary\n        }\n      }\n    }\n  }\n":
    types.EditMessageResolutionSummaryDocument,
  "\n  mutation ForceTimelineUpdateSummary($documentId: ID!, $userId: String!) {\n    forceTimelineUpdateSummary(documentId: $documentId, userId: $userId)      \n  }\n":
    types.ForceTimelineUpdateSummaryDocument,
  "\n  mutation DeleteTimelineMessage($documentId: ID!, $messageId: ID!) {\n    deleteTimelineMessage(documentId: $documentId, messageId: $messageId)\n  }\n":
    types.DeleteTimelineMessageDocument,
  "\n  query GetMe {\n    me {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n      subscriptionStatus\n    }\n  }\n":
    types.GetMeDocument,
  "\n  query GetUser($id: ID!) {\n    user(id: $id) {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n    }\n  }\n":
    types.GetUserDocument,
  "\n  query GetUsers($ids: [ID!]!) {\n    users(ids: $ids) {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n    }\n  }\n":
    types.GetUsersDocument,
  "\n  query GetMyPreference {\n    myPreference {\n      enableActivityNotifications\n    }\n  }\n":
    types.GetMyPreferenceDocument,
  "\n  mutation UpdateMyPreference($input: UpdateUserPreferenceInput!) {\n    updateMyPreference(input: $input) {\n      enableActivityNotifications\n    }\n  }\n":
    types.UpdateMyPreferenceDocument,
  "\n  mutation UpdateMe($input: UpdateUserInput!) {\n    updateMe(input: $input) {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n    }\n  }\n":
    types.UpdateMeDocument,
  "\n  query GetUsersInMyDomain($includeSelf: Boolean = false) {\n    usersInMyDomain(includeSelf: $includeSelf) {\n      id\n      name\n      displayName\n      email\n      picture\n      isAdmin\n    }\n  }\n":
    types.GetUsersInMyDomainDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\nmutation UploadAttachment($file: Upload!, $docId: ID!) {\n  uploadAttachment(file: $file, docId: $docId) {\n    id\n    filename\n    contentType\n    createdAt\n  }\n}\n",
): (typeof documents)["\nmutation UploadAttachment($file: Upload!, $docId: ID!) {\n  uploadAttachment(file: $file, docId: $docId) {\n    id\n    filename\n    contentType\n    createdAt\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\nquery ListDocumentAttachments($docId: ID!) {\n  listDocumentAttachments(docId: $docId) {\n    id\n    filename\n    contentType\n    createdAt\n  }\n}\n",
): (typeof documents)["\nquery ListDocumentAttachments($docId: ID!) {\n  listDocumentAttachments(docId: $docId) {\n    id\n    filename\n    contentType\n    createdAt\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\nquery ListUsersAttachments {\n  listUsersAttachments {\n    id\n    filename\n    contentType\n    createdAt\n  }\n}\n",
): (typeof documents)["\nquery ListUsersAttachments {\n  listUsersAttachments {\n    id\n    filename\n    contentType\n    createdAt\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  fragment DocumentFields on Document {\n    id\n    title\n    isPublic\n    isFolder\n    folderID\n    updatedAt\n    access\n    ownedBy {\n      id\n    }\n    editors {\n      id\n      name\n      displayName\n      email\n      picture\n    }\n    preferences {\n      enableFirstOpenNotifications\n      enableMentionNotifications\n      enableDMNotifications\n      enableAllCommentNotifications\n    }\n  }\n",
): (typeof documents)["\n  fragment DocumentFields on Document {\n    id\n    title\n    isPublic\n    isFolder\n    folderID\n    updatedAt\n    access\n    ownedBy {\n      id\n    }\n    editors {\n      id\n      name\n      displayName\n      email\n      picture\n    }\n    preferences {\n      enableFirstOpenNotifications\n      enableMentionNotifications\n      enableDMNotifications\n      enableAllCommentNotifications\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetDocument($id: ID!) {\n    document(id: $id) {\n      ...DocumentFields\n    }\n  }\n",
): (typeof documents)["\n  query GetDocument($id: ID!) {\n    document(id: $id) {\n      ...DocumentFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation UpdateDocumentTitle($id: ID!, $input: DocumentInput!) {\n    updateDocument(id: $id, input: $input) {\n      id\n      title\n    }\n  }\n",
): (typeof documents)["\n  mutation UpdateDocumentTitle($id: ID!, $input: DocumentInput!) {\n    updateDocument(id: $id, input: $input) {\n      id\n      title\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetDocuments($offset: Int, $limit: Int) {\n    documents(offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n",
): (typeof documents)["\n  query GetDocuments($offset: Int, $limit: Int) {\n    documents(offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetBaseDocuments($offset: Int, $limit: Int) {\n    baseDocuments(offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n",
): (typeof documents)["\n  query GetBaseDocuments($offset: Int, $limit: Int) {\n    baseDocuments(offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetSharedDocuments($offset: Int, $limit: Int) {\n    sharedDocuments(offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n",
): (typeof documents)["\n  query GetSharedDocuments($offset: Int, $limit: Int) {\n    sharedDocuments(offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetFolderDocuments($folderId: ID!, $offset: Int, $limit: Int) {\n    folderDocuments(folderID: $folderId, offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder\n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n",
): (typeof documents)["\n  query GetFolderDocuments($folderId: ID!, $offset: Int, $limit: Int) {\n    folderDocuments(folderID: $folderId, offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder\n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query SearchDocuments($query: String!, $offset: Int, $limit: Int) {\n    searchDocuments(query: $query, offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n",
): (typeof documents)["\n  query SearchDocuments($query: String!, $offset: Int, $limit: Int) {\n    searchDocuments(query: $query, offset: $offset, limit: $limit) {\n      totalCount\n      edges {\n        node {\n          id\n          title\n          isFolder  \n          folderID\n          updatedAt\n          ownedBy {\n            id\n          }\n        }\n      }\n      pageInfo {\n        hasNextPage\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation CreateDocument {\n    createDocument {\n      id\n    }\n  }\n",
): (typeof documents)["\n  mutation CreateDocument {\n    createDocument {\n      id\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation CreateFolder {\n    createFolder {\n      ...DocumentFields\n    }\n  }\n",
): (typeof documents)["\n  mutation CreateFolder {\n    createFolder {\n      ...DocumentFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation DeleteDocument($id: ID!, $deleteChildren: Boolean) {\n    deleteDocument(id: $id, deleteChildren: $deleteChildren)\n  }\n",
): (typeof documents)["\n  mutation DeleteDocument($id: ID!, $deleteChildren: Boolean) {\n    deleteDocument(id: $id, deleteChildren: $deleteChildren)\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation ShareDocument($id: ID!, $emails: [String!]!, $message: String) {\n    shareDocument(documentID: $id, emails: $emails, message: $message) {\n      inviteLink\n    }\n  }\n",
): (typeof documents)["\n  mutation ShareDocument($id: ID!, $emails: [String!]!, $message: String) {\n    shareDocument(documentID: $id, emails: $emails, message: $message) {\n      inviteLink\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation UnshareDocument($docId: ID!, $editorId: ID!) {\n    unshareDocument(documentID: $docId, editorID: $editorId) {\n       id\n       editors {\n        id\n        name\n        email\n        picture\n      }\n    }\n  }\n",
): (typeof documents)["\n  mutation UnshareDocument($docId: ID!, $editorId: ID!) {\n    unshareDocument(documentID: $docId, editorID: $editorId) {\n       id\n       editors {\n        id\n        name\n        email\n        picture\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation UpdateSharedLink($inviteLink: String!, $isActive: Boolean!) {\n    updateShareLink(inviteLink: $inviteLink, isActive: $isActive) {\n      inviteLink\n    }\n  }\n",
): (typeof documents)["\n  mutation UpdateSharedLink($inviteLink: String!, $isActive: Boolean!) {\n    updateShareLink(inviteLink: $inviteLink, isActive: $isActive) {\n      inviteLink\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation UpdateDocumentVisibility($id: ID!, $input: DocumentInput!) {\n    updateDocument(id: $id, input: $input) {\n      id\n      isPublic\n    }\n  }\n",
): (typeof documents)["\n  mutation UpdateDocumentVisibility($id: ID!, $input: DocumentInput!) {\n    updateDocument(id: $id, input: $input) {\n      id\n      isPublic\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query SharedDocumentLinks($id: ID!) {\n    sharedLinks(documentID: $id) {\n      inviteLink\n      inviteeEmail\n      invitedBy {\n        name\n      }\n      isActive\n    }\n  }\n",
): (typeof documents)["\n  query SharedDocumentLinks($id: ID!) {\n    sharedLinks(documentID: $id) {\n      inviteLink\n      inviteeEmail\n      invitedBy {\n        name\n      }\n      isActive\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\nmutation UpdateDocumentPreferences($documentId: ID!, $input: DocumentPreferenceInput!) {\n  updateDocumentPreference(id: $documentId, input: $input) {\n      enableFirstOpenNotifications\n      enableMentionNotifications\n      enableDMNotifications\n      enableAllCommentNotifications\n  }\n}\n",
): (typeof documents)["\nmutation UpdateDocumentPreferences($documentId: ID!, $input: DocumentPreferenceInput!) {\n  updateDocumentPreference(id: $documentId, input: $input) {\n      enableFirstOpenNotifications\n      enableMentionNotifications\n      enableDMNotifications\n      enableAllCommentNotifications\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation CreateFlaggedVersion($documentId: ID!, $input: FlaggedVersionInput!) {\n    createFlaggedVersion(documentId: $documentId, input: $input)\n  }\n",
): (typeof documents)["\n  mutation CreateFlaggedVersion($documentId: ID!, $input: FlaggedVersionInput!) {\n    createFlaggedVersion(documentId: $documentId, input: $input)\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation EditFlaggedVersion($flaggedVersionId: ID!, $input: FlaggedVersionInput!) {\n    editFlaggedVersion(flaggedVersionId: $flaggedVersionId, input: $input)\n  }\n",
): (typeof documents)["\n  mutation EditFlaggedVersion($flaggedVersionId: ID!, $input: FlaggedVersionInput!) {\n    editFlaggedVersion(flaggedVersionId: $flaggedVersionId, input: $input)\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation DeleteFlaggedVersion($flaggedVersionId: ID!, $timelineEventId: ID!) {\n    deleteFlaggedVersion(flaggedVersionId: $flaggedVersionId, timelineEventId: $timelineEventId)\n  }\n",
): (typeof documents)["\n  mutation DeleteFlaggedVersion($flaggedVersionId: ID!, $timelineEventId: ID!) {\n    deleteFlaggedVersion(flaggedVersionId: $flaggedVersionId, timelineEventId: $timelineEventId)\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation MoveDocument($documentId: ID!, $folderId: ID) {\n    moveDocument(id: $documentId, folderID: $folderId) {\n      ...DocumentFields\n    }\n  }\n",
): (typeof documents)["\n  mutation MoveDocument($documentId: ID!, $folderId: ID) {\n    moveDocument(id: $documentId, folderID: $folderId) {\n      ...DocumentFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  subscription DocumentUpdated($documentId: ID!) {\n    documentUpdated(documentId: $documentId) {\n      ...DocumentFields\n    }\n  }\n",
): (typeof documents)["\n  subscription DocumentUpdated($documentId: ID!) {\n    documentUpdated(documentId: $documentId) {\n      ...DocumentFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  subscription DocumentInserted($userId: ID!) {\n    documentInserted(userId: $userId) {\n      ...DocumentFields\n    }\n  }\n",
): (typeof documents)["\n  subscription DocumentInserted($userId: ID!) {\n    documentInserted(userId: $userId) {\n      ...DocumentFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation UploadImage($file: Upload!, $docId: ID!) {\n    uploadImage(file: $file, docId: $docId) {\n      id\n      docId\n      url\n      createdAt\n      mimeType\n      status\n      error\n    }\n  }\n",
): (typeof documents)["\n  mutation UploadImage($file: Upload!, $docId: ID!) {\n    uploadImage(file: $file, docId: $docId) {\n      id\n      docId\n      url\n      createdAt\n      mimeType\n      status\n      error\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query ListDocumentImages($docId: ID!) {\n    listDocumentImages(docId: $docId) {\n      id\n      docId\n      url\n      createdAt\n      mimeType\n      status\n      error\n    }\n  }\n",
): (typeof documents)["\n  query ListDocumentImages($docId: ID!) {\n    listDocumentImages(docId: $docId) {\n      id\n      docId\n      url\n      createdAt\n      mimeType\n      status\n      error\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetImageSignedUrl($docId: ID!, $imageId: ID!) {\n    getImageSignedUrl(docId: $docId, imageId: $imageId) {\n      url\n      expiresAt\n    }\n  }\n",
): (typeof documents)["\n  query GetImageSignedUrl($docId: ID!, $imageId: ID!) {\n    getImageSignedUrl(docId: $docId, imageId: $imageId) {\n      url\n      expiresAt\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetImage($docId: ID!, $imageId: ID!) {\n    getImage(docId: $docId, imageId: $imageId) {\n      id\n      docId\n      url\n      createdAt\n      mimeType\n      status\n      error\n    }\n  }\n",
): (typeof documents)["\n  query GetImage($docId: ID!, $imageId: ID!) {\n    getImage(docId: $docId, imageId: $imageId) {\n      id\n      docId\n      url\n      createdAt\n      mimeType\n      status\n      error\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetAIThreads($documentId: ID!) {\n    getAskAiThreads(documentId: $documentId) {\n      id\n      title\n      updatedAt\n\n      user {\n        id\n        name\n      }\n    }\n  }\n",
): (typeof documents)["\n  query GetAIThreads($documentId: ID!) {\n    getAskAiThreads(documentId: $documentId) {\n      id\n      title\n      updatedAt\n\n      user {\n        id\n        name\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\nfragment MessageFields on Message {\n  id\n  containerId\n  content\n  createdAt\n  lifecycleStage\n  lifecycleReason\n  authorId\n  hidden\n  user {\n    id\n    name\n    picture\n  }\n  aiContent {\n    concludingMessage\n    feedback\n  }\n  metadata {\n    allowDraftEdits\n    contentAddressBefore\n    contentAddress\n    contentAddressAfter\n    contentAddressAfterTimestamp\n    revisionStatus\n  }\n  attachments {\n    __typename\n    ... on Selection {\n      start\n      end\n      content\n    }\n    ... on Revision {\n      start\n      end\n      updated\n      beforeAddress\n      afterAddress\n      appliedOps\n    }\n    ... on Suggestion {\n      content\n    }\n    ... on AttachmentContent {\n      text\n    }\n    ... on AttachmentError {\n      title\n      text\n    }\n    ... on AttachmentFile {\n      id\n      filename\n      contentType\n    }\n    ... on AttachedRevisoDocument {\n      id\n      title\n    }\n  }\n}",
): (typeof documents)["\nfragment MessageFields on Message {\n  id\n  containerId\n  content\n  createdAt\n  lifecycleStage\n  lifecycleReason\n  authorId\n  hidden\n  user {\n    id\n    name\n    picture\n  }\n  aiContent {\n    concludingMessage\n    feedback\n  }\n  metadata {\n    allowDraftEdits\n    contentAddressBefore\n    contentAddress\n    contentAddressAfter\n    contentAddressAfterTimestamp\n    revisionStatus\n  }\n  attachments {\n    __typename\n    ... on Selection {\n      start\n      end\n      content\n    }\n    ... on Revision {\n      start\n      end\n      updated\n      beforeAddress\n      afterAddress\n      appliedOps\n    }\n    ... on Suggestion {\n      content\n    }\n    ... on AttachmentContent {\n      text\n    }\n    ... on AttachmentError {\n      title\n      text\n    }\n    ... on AttachmentFile {\n      id\n      filename\n      contentType\n    }\n    ... on AttachedRevisoDocument {\n      id\n      title\n    }\n  }\n}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetAIThreadMessages($documentId: ID!, $threadId: ID!) {\n    getAskAiThreadMessages(documentId: $documentId, threadId: $threadId) {\n      ...MessageFields\n    }\n  }\n",
): (typeof documents)["\n  query GetAIThreadMessages($documentId: ID!, $threadId: ID!) {\n    getAskAiThreadMessages(documentId: $documentId, threadId: $threadId) {\n      ...MessageFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\nsubscription MessageUpserted($documentId: ID!, $channelId: ID!) {\n  messageUpserted(documentId: $documentId, channelId: $channelId) {\n    ...MessageFields\n  }\n}\n",
): (typeof documents)["\nsubscription MessageUpserted($documentId: ID!, $channelId: ID!) {\n  messageUpserted(documentId: $documentId, channelId: $channelId) {\n    ...MessageFields\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\nfragment ThreadFields on Thread {\n  id\n  title\n  updatedAt\n}",
): (typeof documents)["\nfragment ThreadFields on Thread {\n  id\n  title\n  updatedAt\n}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation UpdateMessageRevisionStatus($containerId: ID!, $messageId: ID!, $status: MessageRevisionStatus!, $contentAddress: String!) {\n    updateMessageRevisionStatus(containerId: $containerId, messageId: $messageId, status: $status, contentAddress: $contentAddress) {\n      ...MessageFields\n    }\n  }\n",
): (typeof documents)["\n  mutation UpdateMessageRevisionStatus($containerId: ID!, $messageId: ID!, $status: MessageRevisionStatus!, $contentAddress: String!) {\n    updateMessageRevisionStatus(containerId: $containerId, messageId: $messageId, status: $status, contentAddress: $contentAddress) {\n      ...MessageFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation CreateAIThread($documentId: ID!) {\n    createAskAiThread(documentId: $documentId) {\n      id\n    }\n  }\n",
): (typeof documents)["\n  mutation CreateAIThread($documentId: ID!) {\n    createAskAiThread(documentId: $documentId) {\n      id\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation CreateAIThreadMessage($documentId: ID!, $threadId: ID!, $input: MessageInput!) {\n    createAskAiThreadMessage(documentId: $documentId, threadId: $threadId, input: $input) {\n      ...MessageFields\n    }\n  }\n",
): (typeof documents)["\n  mutation CreateAIThreadMessage($documentId: ID!, $threadId: ID!, $input: MessageInput!) {\n    createAskAiThreadMessage(documentId: $documentId, threadId: $threadId, input: $input) {\n      ...MessageFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\nsubscription ThreadUpserted($documentId: ID!) {\n  threadUpserted(documentId: $documentId) {\n    ...ThreadFields\n  }\n}\n",
): (typeof documents)["\nsubscription ThreadUpserted($documentId: ID!) {\n  threadUpserted(documentId: $documentId) {\n    ...ThreadFields\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query SubscriptionPlans {\n    subscriptionPlans {\n      id\n      name\n      priceCents\n      currency\n      interval\n    }\n  }\n",
): (typeof documents)["\n  query SubscriptionPlans {\n    subscriptionPlans {\n      id\n      name\n      priceCents\n      currency\n      interval\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation Checkout($id: ID!) {\n    checkoutSubscriptionPlan(id: $id) {\n      url\n    }\n  }\n",
): (typeof documents)["\n  mutation Checkout($id: ID!) {\n    checkoutSubscriptionPlan(id: $id) {\n      url\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation BillingPortalSession {\n    billingPortalSession {\n      url\n    }\n  }\n",
): (typeof documents)["\n  mutation BillingPortalSession {\n    billingPortalSession {\n      url\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\nfragment TimelineEventFields on TimelineEvent {\n  id\n  replyTo\n  createdAt\n  authorId\n  user {\n    id\n    name\n    picture\n  }\n  event {\n    __typename\n    ... on TLJoinV1 {\n      action\n    }\n    ... on TLMarkerV1 {\n      title\n    }\n    ... on TLUpdateV1 {\n      eventId\n      title\n      content\n      startingContentAddress\n      endingContentAddress\n      flaggedVersionName\n      flaggedVersionCreatedAt\n      flaggedVersionID\n      state\n      flaggedByUser {\n        name\n      }\n    }\n    ... on TLEmpty {\n      placeholder \n    }\n    ... on TLAttributeChangeV1 {\n      attribute\n      oldValue\n      newValue\n    }\n    ... on TLAccessChangeV1 {\n      action\n      userIdentifiers\n    }\n    ... on TLPasteV1 {\n      contentAddressBefore\n      contentAddressAfter\n    }\n    ... on TLMessageV1 {\n      eventId \n      content\n      contentAddress\n      selectionStartId\n      selectionEndId\n      selectionMarkdown\n\n      replies {\n        id\n        replyTo\n        createdAt\n        authorId\n        user {\n          id\n          name\n          picture\n        }\n        event {\n          __typename\n          ... on TLMessageV1 {\n            eventId \n            content\n            contentAddress\n            selectionStartId\n            selectionEndId\n            selectionMarkdown\n          }\n          ... on TLMessageResolutionV1 {\n            eventId\n            resolutionSummary\n            resolved\n          }\n        }\n      }\n    }\n  }\n}",
): (typeof documents)["\nfragment TimelineEventFields on TimelineEvent {\n  id\n  replyTo\n  createdAt\n  authorId\n  user {\n    id\n    name\n    picture\n  }\n  event {\n    __typename\n    ... on TLJoinV1 {\n      action\n    }\n    ... on TLMarkerV1 {\n      title\n    }\n    ... on TLUpdateV1 {\n      eventId\n      title\n      content\n      startingContentAddress\n      endingContentAddress\n      flaggedVersionName\n      flaggedVersionCreatedAt\n      flaggedVersionID\n      state\n      flaggedByUser {\n        name\n      }\n    }\n    ... on TLEmpty {\n      placeholder \n    }\n    ... on TLAttributeChangeV1 {\n      attribute\n      oldValue\n      newValue\n    }\n    ... on TLAccessChangeV1 {\n      action\n      userIdentifiers\n    }\n    ... on TLPasteV1 {\n      contentAddressBefore\n      contentAddressAfter\n    }\n    ... on TLMessageV1 {\n      eventId \n      content\n      contentAddress\n      selectionStartId\n      selectionEndId\n      selectionMarkdown\n\n      replies {\n        id\n        replyTo\n        createdAt\n        authorId\n        user {\n          id\n          name\n          picture\n        }\n        event {\n          __typename\n          ... on TLMessageV1 {\n            eventId \n            content\n            contentAddress\n            selectionStartId\n            selectionEndId\n            selectionMarkdown\n          }\n          ... on TLMessageResolutionV1 {\n            eventId\n            resolutionSummary\n            resolved\n          }\n        }\n      }\n    }\n  }\n}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetDocumentTimeline($documentId: ID!, $filter: TimelineEventFilter) {\n    getDocumentTimeline(documentId: $documentId, filter: $filter) {\n      ...TimelineEventFields\n    }\n  }\n",
): (typeof documents)["\n  query GetDocumentTimeline($documentId: ID!, $filter: TimelineEventFilter) {\n    getDocumentTimeline(documentId: $documentId, filter: $filter) {\n      ...TimelineEventFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  subscription TimelineEventInserted($documentId: ID!) {\n    timelineEventInserted(documentId: $documentId) {\n      ...TimelineEventFields\n    }\n  }\n",
): (typeof documents)["\n  subscription TimelineEventInserted($documentId: ID!) {\n    timelineEventInserted(documentId: $documentId) {\n      ...TimelineEventFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  subscription TimelineEventUpdated($documentId: ID!) {\n    timelineEventUpdated(documentId: $documentId) {\n      ...TimelineEventFields\n    }\n  }\n",
): (typeof documents)["\n  subscription TimelineEventUpdated($documentId: ID!) {\n    timelineEventUpdated(documentId: $documentId) {\n      ...TimelineEventFields\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  subscription TimelineEventDeleted($documentId: ID!) {\n    timelineEventDeleted(documentId: $documentId) {\n      id\n      replyTo\n      event {\n        ... on TLMessageV1 {\n          eventId\n        }\n      }\n    }\n  }\n",
): (typeof documents)["\n  subscription TimelineEventDeleted($documentId: ID!) {\n    timelineEventDeleted(documentId: $documentId) {\n      id\n      replyTo\n      event {\n        ... on TLMessageV1 {\n          eventId\n        }\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation CreateTimelineMessage($documentId: ID!, $input: TimelineMessageInput!) {\n    createTimelineMessage(documentId: $documentId, input: $input) {\n      id\n    }\n  }\n",
): (typeof documents)["\n  mutation CreateTimelineMessage($documentId: ID!, $input: TimelineMessageInput!) {\n    createTimelineMessage(documentId: $documentId, input: $input) {\n      id\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation EditTimelineMessage($documentId: ID!, $messageId: ID!, $input: EditTimelineMessageInput!) {\n    editTimelineMessage(documentId: $documentId, messageId: $messageId, input: $input) {\n      id\n    }\n  }\n",
): (typeof documents)["\n  mutation EditTimelineMessage($documentId: ID!, $messageId: ID!, $input: EditTimelineMessageInput!) {\n    editTimelineMessage(documentId: $documentId, messageId: $messageId, input: $input) {\n      id\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation UpdateMessageResolution($documentId: ID!, $messageId: ID!, $input: UpdateMessageResolutionInput!) {\n    updateMessageResolution(documentId: $documentId, messageId: $messageId, input: $input) {\n      id\n    }\n  }\n",
): (typeof documents)["\n  mutation UpdateMessageResolution($documentId: ID!, $messageId: ID!, $input: UpdateMessageResolutionInput!) {\n    updateMessageResolution(documentId: $documentId, messageId: $messageId, input: $input) {\n      id\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation EditTimelineUpdateSummary($documentId: ID!, $updateId: ID!, $summary: String!) {\n    editTimelineUpdateSummary(documentId: $documentId, updateId: $updateId, summary: $summary) {\n      id\n      event {\n        ... on TLUpdateV1 {\n          eventId\n          content\n        }\n      }\n    }\n  }\n",
): (typeof documents)["\n  mutation EditTimelineUpdateSummary($documentId: ID!, $updateId: ID!, $summary: String!) {\n    editTimelineUpdateSummary(documentId: $documentId, updateId: $updateId, summary: $summary) {\n      id\n      event {\n        ... on TLUpdateV1 {\n          eventId\n          content\n        }\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation EditMessageResolutionSummary($documentId: ID!, $messageId: ID!, $summary: String!) {\n    editMessageResolutionSummary(documentId: $documentId, messageId: $messageId, summary: $summary) {\n      id\n      event {\n      ... on TLMessageResolutionV1 {\n        eventId\n          resolutionSummary\n        }\n      }\n    }\n  }\n",
): (typeof documents)["\n  mutation EditMessageResolutionSummary($documentId: ID!, $messageId: ID!, $summary: String!) {\n    editMessageResolutionSummary(documentId: $documentId, messageId: $messageId, summary: $summary) {\n      id\n      event {\n      ... on TLMessageResolutionV1 {\n        eventId\n          resolutionSummary\n        }\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation ForceTimelineUpdateSummary($documentId: ID!, $userId: String!) {\n    forceTimelineUpdateSummary(documentId: $documentId, userId: $userId)      \n  }\n",
): (typeof documents)["\n  mutation ForceTimelineUpdateSummary($documentId: ID!, $userId: String!) {\n    forceTimelineUpdateSummary(documentId: $documentId, userId: $userId)      \n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation DeleteTimelineMessage($documentId: ID!, $messageId: ID!) {\n    deleteTimelineMessage(documentId: $documentId, messageId: $messageId)\n  }\n",
): (typeof documents)["\n  mutation DeleteTimelineMessage($documentId: ID!, $messageId: ID!) {\n    deleteTimelineMessage(documentId: $documentId, messageId: $messageId)\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetMe {\n    me {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n      subscriptionStatus\n    }\n  }\n",
): (typeof documents)["\n  query GetMe {\n    me {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n      subscriptionStatus\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetUser($id: ID!) {\n    user(id: $id) {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n    }\n  }\n",
): (typeof documents)["\n  query GetUser($id: ID!) {\n    user(id: $id) {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetUsers($ids: [ID!]!) {\n    users(ids: $ids) {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n    }\n  }\n",
): (typeof documents)["\n  query GetUsers($ids: [ID!]!) {\n    users(ids: $ids) {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetMyPreference {\n    myPreference {\n      enableActivityNotifications\n    }\n  }\n",
): (typeof documents)["\n  query GetMyPreference {\n    myPreference {\n      enableActivityNotifications\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation UpdateMyPreference($input: UpdateUserPreferenceInput!) {\n    updateMyPreference(input: $input) {\n      enableActivityNotifications\n    }\n  }\n",
): (typeof documents)["\n  mutation UpdateMyPreference($input: UpdateUserPreferenceInput!) {\n    updateMyPreference(input: $input) {\n      enableActivityNotifications\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  mutation UpdateMe($input: UpdateUserInput!) {\n    updateMe(input: $input) {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n    }\n  }\n",
): (typeof documents)["\n  mutation UpdateMe($input: UpdateUserInput!) {\n    updateMe(input: $input) {\n      id\n      email\n      name\n      displayName\n      picture\n      isAdmin\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(
  source: "\n  query GetUsersInMyDomain($includeSelf: Boolean = false) {\n    usersInMyDomain(includeSelf: $includeSelf) {\n      id\n      name\n      displayName\n      email\n      picture\n      isAdmin\n    }\n  }\n",
): (typeof documents)["\n  query GetUsersInMyDomain($includeSelf: Boolean = false) {\n    usersInMyDomain(includeSelf: $includeSelf) {\n      id\n      name\n      displayName\n      email\n      picture\n      isAdmin\n    }\n  }\n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> =
  TDocumentNode extends DocumentNode<infer TType, any> ? TType : never;
