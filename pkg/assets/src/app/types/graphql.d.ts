type Maybe<T> = T | null;
type InputMaybe<T> = Maybe<T>;
type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = {
  [_ in K]?: never;
};
type Incremental<T> =
  | T
  | {
      [P in keyof T]?: P extends " $fragmentName" | "__typename" ? T[P] : never;
    };
/** All built-in and custom scalars, mapped to their actual values */
interface Scalars {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
  JSON: { input: string; output: string };
  Time: { input: string; output: string };
  Upload: { input: File; output: File };
}

interface AiContent {
  __typename?: "AiContent";
  concludingMessage?: Maybe<Scalars["String"]["output"]>;
  feedback?: Maybe<Scalars["String"]["output"]>;
  notes?: Maybe<Scalars["String"]["output"]>;
}

interface AttachedRevisoDocument {
  __typename?: "AttachedRevisoDocument";
  id: Scalars["String"]["output"];
  title: Scalars["String"]["output"];
}

interface AttachmentContent {
  __typename?: "AttachmentContent";
  role: Scalars["String"]["output"];
  text: Scalars["String"]["output"];
}

interface AttachmentError {
  __typename?: "AttachmentError";
  error: Scalars["String"]["output"];
  text: Scalars["String"]["output"];
  title: Scalars["String"]["output"];
}

interface AttachmentFile {
  __typename?: "AttachmentFile";
  contentType: Scalars["String"]["output"];
  filename: Scalars["String"]["output"];
  id: Scalars["String"]["output"];
}

interface AttachmentInput {
  contentType?: InputMaybe<Scalars["String"]["input"]>;
  id: Scalars["ID"]["input"];
  name: Scalars["String"]["input"];
  type: AttachmentInputType;
}

type AttachmentInputType = "DRAFT" | "FILE" | "UNKNOWN";

type AttachmentProgressType = "DONE" | "THINKING" | "UNKNOWN";

type AttachmentValue =
  | AttachedRevisoDocument
  | AttachmentContent
  | AttachmentError
  | AttachmentFile
  | Revision
  | Selection
  | Suggestion
  | { __typename?: "%other" };

interface BillingPortalSession {
  __typename?: "BillingPortalSession";
  url: Scalars["String"]["output"];
}

interface Chain {
  __typename?: "Chain";
  id: Scalars["ID"]["output"];
  messages: Array<Message>;
}

type ChanType = "DIRECT" | "GENERAL" | "REVISO" | "UNKNOWN";

interface Checkout {
  __typename?: "Checkout";
  url: Scalars["String"]["output"];
}

interface CommentNotificationPayloadValue {
  __typename?: "CommentNotificationPayloadValue";
  author: User;
  authorId: Scalars["ID"]["output"];
  channelId: Scalars["ID"]["output"];
  commentType: CommentNotificationType;
  containerId: Scalars["ID"]["output"];
  documentId: Scalars["ID"]["output"];
  message: Message;
  messageId: Scalars["ID"]["output"];
}

type CommentNotificationType = "COMMENT" | "MENTION" | "REPLY" | "UNKNOWN";

interface ContentAddress {
  __typename?: "ContentAddress";
  documentId: Scalars["ID"]["output"];
  id: Scalars["ID"]["output"];
  payload: Scalars["JSON"]["output"];
}

interface Document {
  __typename?: "Document";
  access: Scalars["String"]["output"];
  branchCopies: Array<Document>;
  createdAt: Scalars["Time"]["output"];
  editors: Array<User>;
  folderID?: Maybe<Scalars["ID"]["output"]>;
  hasUnreadNotifications: Scalars["Boolean"]["output"];
  id: Scalars["ID"]["output"];
  isFolder: Scalars["Boolean"]["output"];
  isPublic: Scalars["Boolean"]["output"];
  ownedBy: User;
  parentAddress?: Maybe<Scalars["String"]["output"]>;
  parentID?: Maybe<Scalars["ID"]["output"]>;
  preferences: DocumentPreference;
  rootParentID: Scalars["ID"]["output"];
  screenshots?: Maybe<DocumentScreenshots>;
  title: Scalars["String"]["output"];
  updatedAt: Scalars["Time"]["output"];
}

interface DocumentAttachment {
  __typename?: "DocumentAttachment";
  contentType: Scalars["String"]["output"];
  createdAt: Scalars["Time"]["output"];
  filename: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
}

interface DocumentConnection {
  __typename?: "DocumentConnection";
  edges: Array<DocumentEdge>;
  pageInfo: PageInfo;
  totalCount: Scalars["Int"]["output"];
}

interface DocumentEdge {
  __typename?: "DocumentEdge";
  cursor: Scalars["String"]["output"];
  node: Document;
}

interface DocumentInput {
  isPublic?: InputMaybe<Scalars["Boolean"]["input"]>;
  title?: InputMaybe<Scalars["String"]["input"]>;
}

interface DocumentPreference {
  __typename?: "DocumentPreference";
  enableAllCommentNotifications: Scalars["Boolean"]["output"];
  enableDMNotifications: Scalars["Boolean"]["output"];
  enableFirstOpenNotifications: Scalars["Boolean"]["output"];
  enableMentionNotifications: Scalars["Boolean"]["output"];
}

interface DocumentPreferenceInput {
  enableAllCommentNotifications: Scalars["Boolean"]["input"];
  enableDMNotifications: Scalars["Boolean"]["input"];
  enableFirstOpenNotifications: Scalars["Boolean"]["input"];
  enableMentionNotifications: Scalars["Boolean"]["input"];
}

interface DocumentScreenshots {
  __typename?: "DocumentScreenshots";
  darkUrl: Scalars["String"]["output"];
  lightUrl: Scalars["String"]["output"];
}

interface EditTimelineMessageInput {
  content: Scalars["String"]["input"];
  contentAddress?: InputMaybe<Scalars["String"]["input"]>;
  endID?: InputMaybe<Scalars["String"]["input"]>;
  selectionMarkdown?: InputMaybe<Scalars["String"]["input"]>;
  startID?: InputMaybe<Scalars["String"]["input"]>;
}

interface FlaggedVersionInput {
  name: Scalars["String"]["input"];
  updateID: Scalars["String"]["input"];
}

interface Image {
  __typename?: "Image";
  createdAt: Scalars["Time"]["output"];
  docId: Scalars["ID"]["output"];
  error?: Maybe<Scalars["String"]["output"]>;
  id: Scalars["ID"]["output"];
  mimeType: Scalars["String"]["output"];
  status: Status;
  url: Scalars["String"]["output"];
}

type LifecycleStage =
  | "COMPLETED"
  | "PENDING"
  | "REVISED"
  | "REVISING"
  | "UNKNOWN";

interface Message {
  __typename?: "Message";
  aiContent?: Maybe<AiContent>;
  attachments: Array<AttachmentValue>;
  authorId: Scalars["ID"]["output"];
  chain: Scalars["String"]["output"];
  channelId: Scalars["ID"]["output"];
  containerId: Scalars["ID"]["output"];
  content: Scalars["String"]["output"];
  createdAt: Scalars["Time"]["output"];
  forkedMessageIds: Array<Scalars["ID"]["output"]>;
  hidden: Scalars["Boolean"]["output"];
  id: Scalars["ID"]["output"];
  lifecycleReason: Scalars["String"]["output"];
  lifecycleStage: LifecycleStage;
  metadata: MsgMetadata;
  parentContainerId?: Maybe<Scalars["ID"]["output"]>;
  parentMessageId?: Maybe<Scalars["ID"]["output"]>;
  replies: Array<Message>;
  replyCount: Scalars["Int"]["output"];
  replyingUserIds: Array<Scalars["ID"]["output"]>;
  replyingUsers: Array<User>;
  user: User;
  userId: Scalars["ID"]["output"];
}

interface MessageInput {
  allowDraftEdits: Scalars["Boolean"]["input"];
  attachments?: InputMaybe<Array<AttachmentInput>>;
  authorId: Scalars["ID"]["input"];
  content: Scalars["String"]["input"];
  contentAddress: Scalars["String"]["input"];
  llm?: InputMaybe<MsgLlm>;
  replyTo?: InputMaybe<Scalars["ID"]["input"]>;
  selection?: InputMaybe<SelectionInput>;
}

type MessageRevisionStatus = "ACCEPTED" | "DECLINED" | "UNSPECIFIED";

interface MessageUpdateInput {
  content: Scalars["String"]["input"];
  id: Scalars["ID"]["input"];
}

type MsgLlm = "CLAUDE" | "GPT4O";

interface MsgMetadata {
  __typename?: "MsgMetadata";
  allowDraftEdits: Scalars["Boolean"]["output"];
  contentAddress: Scalars["String"]["output"];
  contentAddressAfter: Scalars["String"]["output"];
  contentAddressAfterTimestamp?: Maybe<Scalars["Time"]["output"]>;
  contentAddressBefore: Scalars["String"]["output"];
  llm: MsgLlm;
  revisionStatus: MessageRevisionStatus;
}

interface Mutation {
  __typename?: "Mutation";
  billingPortalSession: BillingPortalSession;
  checkoutSubscriptionPlan: Checkout;
  copyDocument?: Maybe<Document>;
  createAskAiThread: Thread;
  createAskAiThreadMessage: Message;
  createDocument: Document;
  createFlaggedVersion: Scalars["Boolean"]["output"];
  createFolder: Document;
  createShareLinks: Array<SharedDocumentLink>;
  createTimelineMessage: TimelineEvent;
  deleteDocument?: Maybe<Scalars["Boolean"]["output"]>;
  deleteFlaggedVersion: Scalars["Boolean"]["output"];
  deleteTimelineMessage: Scalars["Boolean"]["output"];
  editFlaggedVersion: Scalars["Boolean"]["output"];
  editMessageResolutionSummary: TimelineEvent;
  editTimelineMessage: TimelineEvent;
  editTimelineUpdateSummary: TimelineEvent;
  forceTimelineUpdateSummary: Scalars["Boolean"]["output"];
  joinShareLink: Document;
  moveDocument?: Maybe<Document>;
  saveContentAddress?: Maybe<MutationResponse>;
  sendAccessLinkForInvite?: Maybe<Scalars["Boolean"]["output"]>;
  shareDocument: Array<SharedDocumentLink>;
  softDeleteDocument?: Maybe<Scalars["Boolean"]["output"]>;
  unshareDocument: Document;
  updateDocument?: Maybe<Document>;
  updateDocumentPreference?: Maybe<DocumentPreference>;
  updateMe?: Maybe<User>;
  updateMessageResolution: TimelineEvent;
  updateMessageRevisionStatus: Message;
  updateMyPreference: UserPreference;
  updateShareLink: SharedDocumentLink;
  uploadAttachment: DocumentAttachment;
  uploadImage: Image;
}

interface Mutation_CheckoutSubscriptionPlanArgs {
  id: Scalars["ID"]["input"];
}

interface Mutation_CopyDocumentArgs {
  address?: InputMaybe<Scalars["String"]["input"]>;
  id: Scalars["ID"]["input"];
  isBranch?: InputMaybe<Scalars["Boolean"]["input"]>;
}

interface Mutation_CreateAskAiThreadArgs {
  documentId: Scalars["ID"]["input"];
}

interface Mutation_CreateAskAiThreadMessageArgs {
  documentId: Scalars["ID"]["input"];
  input: MessageInput;
  threadId: Scalars["ID"]["input"];
}

interface Mutation_CreateFlaggedVersionArgs {
  documentId: Scalars["ID"]["input"];
  input: FlaggedVersionInput;
}

interface Mutation_CreateShareLinksArgs {
  documentID: Scalars["ID"]["input"];
  emails: Array<Scalars["String"]["input"]>;
  message?: InputMaybe<Scalars["String"]["input"]>;
}

interface Mutation_CreateTimelineMessageArgs {
  documentId: Scalars["ID"]["input"];
  input: TimelineMessageInput;
}

interface Mutation_DeleteDocumentArgs {
  deleteChildren?: InputMaybe<Scalars["Boolean"]["input"]>;
  id: Scalars["ID"]["input"];
}

interface Mutation_DeleteFlaggedVersionArgs {
  flaggedVersionId: Scalars["ID"]["input"];
  timelineEventId: Scalars["ID"]["input"];
}

interface Mutation_DeleteTimelineMessageArgs {
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}

interface Mutation_EditFlaggedVersionArgs {
  flaggedVersionId: Scalars["ID"]["input"];
  input: FlaggedVersionInput;
}

interface Mutation_EditMessageResolutionSummaryArgs {
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  summary: Scalars["String"]["input"];
}

interface Mutation_EditTimelineMessageArgs {
  documentId: Scalars["ID"]["input"];
  input: EditTimelineMessageInput;
  messageId: Scalars["ID"]["input"];
}

interface Mutation_EditTimelineUpdateSummaryArgs {
  documentId: Scalars["ID"]["input"];
  summary: Scalars["String"]["input"];
  updateId: Scalars["ID"]["input"];
}

interface Mutation_ForceTimelineUpdateSummaryArgs {
  documentId: Scalars["ID"]["input"];
  userId: Scalars["String"]["input"];
}

interface Mutation_JoinShareLinkArgs {
  inviteLink: Scalars["String"]["input"];
}

interface Mutation_MoveDocumentArgs {
  folderID?: InputMaybe<Scalars["ID"]["input"]>;
  id: Scalars["ID"]["input"];
}

interface Mutation_SaveContentAddressArgs {
  documentId: Scalars["ID"]["input"];
  payload: Scalars["JSON"]["input"];
}

interface Mutation_SendAccessLinkForInviteArgs {
  inviteLink: Scalars["String"]["input"];
}

interface Mutation_ShareDocumentArgs {
  documentID: Scalars["ID"]["input"];
  emails: Array<Scalars["String"]["input"]>;
  message?: InputMaybe<Scalars["String"]["input"]>;
}

interface Mutation_SoftDeleteDocumentArgs {
  id: Scalars["ID"]["input"];
}

interface Mutation_UnshareDocumentArgs {
  documentID: Scalars["ID"]["input"];
  editorID: Scalars["ID"]["input"];
}

interface Mutation_UpdateDocumentArgs {
  id: Scalars["ID"]["input"];
  input: DocumentInput;
}

interface Mutation_UpdateDocumentPreferenceArgs {
  id: Scalars["ID"]["input"];
  input: DocumentPreferenceInput;
}

interface Mutation_UpdateMeArgs {
  input: UpdateUserInput;
}

interface Mutation_UpdateMessageResolutionArgs {
  documentId: Scalars["ID"]["input"];
  input: UpdateMessageResolutionInput;
  messageId: Scalars["ID"]["input"];
}

interface Mutation_UpdateMessageRevisionStatusArgs {
  containerId: Scalars["ID"]["input"];
  contentAddress: Scalars["String"]["input"];
  messageId: Scalars["ID"]["input"];
  status: MessageRevisionStatus;
}

interface Mutation_UpdateMyPreferenceArgs {
  input: UpdateUserPreferenceInput;
}

interface Mutation_UpdateShareLinkArgs {
  inviteLink: Scalars["String"]["input"];
  isActive: Scalars["Boolean"]["input"];
}

interface Mutation_UploadAttachmentArgs {
  docId: Scalars["ID"]["input"];
  file: Scalars["Upload"]["input"];
}

interface Mutation_UploadImageArgs {
  docId: Scalars["ID"]["input"];
  file: Scalars["Upload"]["input"];
}

interface MutationResponse {
  __typename?: "MutationResponse";
  id: Scalars["ID"]["output"];
}

interface Notification {
  __typename?: "Notification";
  createdAt: Scalars["Time"]["output"];
  documentId: Scalars["ID"]["output"];
  id: Scalars["ID"]["output"];
  payload: NotificationPayloadValue;
  read: Scalars["Boolean"]["output"];
}

interface NotificationConnection {
  __typename?: "NotificationConnection";
  edges: Array<Notification>;
}

type NotificationPayloadValue =
  | CommentNotificationPayloadValue
  | { __typename?: "%other" };

interface PageInfo {
  __typename?: "PageInfo";
  hasNextPage: Scalars["Boolean"]["output"];
}

interface Query {
  __typename?: "Query";
  baseDocuments: DocumentConnection;
  branches: Array<Document>;
  document?: Maybe<Document>;
  documents: DocumentConnection;
  folderDocuments: DocumentConnection;
  getAskAiThreadMessages: Array<Message>;
  getAskAiThreads: Array<Thread>;
  getAttachmentSignedUrl: SignedImageUrl;
  getContentAddress?: Maybe<ContentAddress>;
  getDocumentTimeline: Array<TimelineEvent>;
  getImage: Image;
  getImageSignedUrl: SignedImageUrl;
  listDocumentAttachments: Array<DocumentAttachment>;
  listDocumentImages: Array<Image>;
  listUsersAttachments: Array<DocumentAttachment>;
  me?: Maybe<User>;
  myPreference: UserPreference;
  searchDocuments: DocumentConnection;
  sharedDocuments: DocumentConnection;
  sharedLink?: Maybe<SharedDocumentLink>;
  sharedLinks: Array<SharedDocumentLink>;
  subscriptionPlans: Array<SubscriptionPlan>;
  unauthenticatedSharedLink?: Maybe<UnauthenticatedSharedLink>;
  user?: Maybe<User>;
  users?: Maybe<Array<Maybe<User>>>;
  usersInMyDomain: Array<User>;
}

interface Query_BaseDocumentsArgs {
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  offset?: InputMaybe<Scalars["Int"]["input"]>;
}

interface Query_BranchesArgs {
  id: Scalars["ID"]["input"];
}

interface Query_DocumentArgs {
  id: Scalars["ID"]["input"];
}

interface Query_DocumentsArgs {
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  offset?: InputMaybe<Scalars["Int"]["input"]>;
}

interface Query_FolderDocumentsArgs {
  folderID: Scalars["ID"]["input"];
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  offset?: InputMaybe<Scalars["Int"]["input"]>;
}

interface Query_GetAskAiThreadMessagesArgs {
  documentId: Scalars["ID"]["input"];
  threadId: Scalars["ID"]["input"];
}

interface Query_GetAskAiThreadsArgs {
  documentId: Scalars["ID"]["input"];
}

interface Query_GetAttachmentSignedUrlArgs {
  attachmentId: Scalars["ID"]["input"];
}

interface Query_GetContentAddressArgs {
  addressId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
}

interface Query_GetDocumentTimelineArgs {
  documentId: Scalars["ID"]["input"];
  filter?: InputMaybe<TimelineEventFilter>;
}

interface Query_GetImageArgs {
  docId: Scalars["ID"]["input"];
  imageId: Scalars["ID"]["input"];
}

interface Query_GetImageSignedUrlArgs {
  docId: Scalars["ID"]["input"];
  imageId: Scalars["ID"]["input"];
}

interface Query_ListDocumentAttachmentsArgs {
  docId: Scalars["ID"]["input"];
}

interface Query_ListDocumentImagesArgs {
  docId: Scalars["ID"]["input"];
}

interface Query_SearchDocumentsArgs {
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  query: Scalars["String"]["input"];
}

interface Query_SharedDocumentsArgs {
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  offset?: InputMaybe<Scalars["Int"]["input"]>;
}

interface Query_SharedLinkArgs {
  inviteLink: Scalars["String"]["input"];
}

interface Query_SharedLinksArgs {
  documentID: Scalars["ID"]["input"];
}

interface Query_UnauthenticatedSharedLinkArgs {
  inviteLink: Scalars["String"]["input"];
}

interface Query_UserArgs {
  id: Scalars["ID"]["input"];
}

interface Query_UsersArgs {
  ids: Array<Scalars["ID"]["input"]>;
}

interface Query_UsersInMyDomainArgs {
  includeSelf?: InputMaybe<Scalars["Boolean"]["input"]>;
}

interface Revision {
  __typename?: "Revision";
  afterAddress?: Maybe<Scalars["String"]["output"]>;
  appliedOps?: Maybe<Scalars["String"]["output"]>;
  beforeAddress?: Maybe<Scalars["String"]["output"]>;
  end: Scalars["String"]["output"];
  explanation?: Maybe<Scalars["String"]["output"]>;
  followUps?: Maybe<Scalars["String"]["output"]>;
  marshalledOperations: Scalars["String"]["output"];
  start: Scalars["String"]["output"];
  updated: Scalars["String"]["output"];
}

interface Selection {
  __typename?: "Selection";
  content: Scalars["String"]["output"];
  end: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
  start: Scalars["String"]["output"];
}

interface SelectionInput {
  content: Scalars["String"]["input"];
  end: Scalars["String"]["input"];
  start: Scalars["String"]["input"];
}

interface SharedDocumentLink {
  __typename?: "SharedDocumentLink";
  createdAt: Scalars["Time"]["output"];
  document: Document;
  inviteLink: Scalars["String"]["output"];
  invitedBy: User;
  inviteeEmail: Scalars["String"]["output"];
  inviteeUser?: Maybe<User>;
  isActive: Scalars["Boolean"]["output"];
  updatedAt: Scalars["Time"]["output"];
}

interface SignedImageUrl {
  __typename?: "SignedImageUrl";
  expiresAt: Scalars["Time"]["output"];
  url: Scalars["String"]["output"];
}

type Status = "ERROR" | "LOADING" | "SUCCESS";

interface Subscription {
  __typename?: "Subscription";
  documentInserted: Document;
  documentUpdated: Document;
  messageUpserted: Message;
  threadUpserted: Thread;
  timelineEventDeleted: TimelineEvent;
  timelineEventInserted: TimelineEvent;
  timelineEventUpdated: TimelineEvent;
}

interface Subscription_DocumentInsertedArgs {
  userId: Scalars["ID"]["input"];
}

interface Subscription_DocumentUpdatedArgs {
  documentId: Scalars["ID"]["input"];
}

interface Subscription_MessageUpsertedArgs {
  channelId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
}

interface Subscription_ThreadUpsertedArgs {
  documentId: Scalars["ID"]["input"];
}

interface Subscription_TimelineEventDeletedArgs {
  documentId: Scalars["ID"]["input"];
}

interface Subscription_TimelineEventInsertedArgs {
  documentId: Scalars["ID"]["input"];
}

interface Subscription_TimelineEventUpdatedArgs {
  documentId: Scalars["ID"]["input"];
}

interface SubscriptionPlan {
  __typename?: "SubscriptionPlan";
  currency: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
  interval: Scalars["String"]["output"];
  name: Scalars["String"]["output"];
  priceCents: Scalars["Int"]["output"];
}

interface Suggestion {
  __typename?: "Suggestion";
  content: Scalars["String"]["output"];
}

interface TlAccessChangeV1 {
  __typename?: "TLAccessChangeV1";
  action: Scalars["String"]["output"];
  userIdentifiers: Array<Scalars["String"]["output"]>;
}

interface TlAttributeChangeV1 {
  __typename?: "TLAttributeChangeV1";
  attribute: Scalars["String"]["output"];
  newValue: Scalars["String"]["output"];
  oldValue: Scalars["String"]["output"];
}

interface TlEmpty {
  __typename?: "TLEmpty";
  placeholder: Scalars["String"]["output"];
}

type TlEventPayload =
  | TlAccessChangeV1
  | TlAttributeChangeV1
  | TlEmpty
  | TlJoinV1
  | TlMarkerV1
  | TlMessageResolutionV1
  | TlMessageV1
  | TlPasteV1
  | TlUpdateV1
  | { __typename?: "%other" };

interface TlJoinV1 {
  __typename?: "TLJoinV1";
  action: Scalars["String"]["output"];
}

interface TlMarkerV1 {
  __typename?: "TLMarkerV1";
  title: Scalars["String"]["output"];
}

interface TlMessageResolutionV1 {
  __typename?: "TLMessageResolutionV1";
  eventId: Scalars["String"]["output"];
  resolutionSummary: Scalars["String"]["output"];
  resolved: Scalars["Boolean"]["output"];
}

interface TlMessageV1 {
  __typename?: "TLMessageV1";
  content: Scalars["String"]["output"];
  contentAddress: Scalars["String"]["output"];
  documentId: Scalars["ID"]["output"];
  eventId: Scalars["String"]["output"];
  replies: Array<TimelineEvent>;
  selectionEndId: Scalars["String"]["output"];
  selectionMarkdown: Scalars["String"]["output"];
  selectionStartId: Scalars["String"]["output"];
}

interface TlPasteV1 {
  __typename?: "TLPasteV1";
  contentAddressAfter: Scalars["String"]["output"];
  contentAddressBefore: Scalars["String"]["output"];
}

type TlUpdateState = "COMPLETE" | "SUMMARIZING";

interface TlUpdateV1 {
  __typename?: "TLUpdateV1";
  content: Scalars["String"]["output"];
  endingContentAddress: Scalars["String"]["output"];
  eventId: Scalars["String"]["output"];
  flaggedByUser?: Maybe<User>;
  flaggedVersionCreatedAt?: Maybe<Scalars["Time"]["output"]>;
  flaggedVersionID?: Maybe<Scalars["String"]["output"]>;
  flaggedVersionName?: Maybe<Scalars["String"]["output"]>;
  startingContentAddress: Scalars["String"]["output"];
  state: TlUpdateState;
  title: Scalars["String"]["output"];
}

interface Thread {
  __typename?: "Thread";
  documentId: Scalars["ID"]["output"];
  id: Scalars["ID"]["output"];
  messages: Array<Message>;
  title: Scalars["String"]["output"];
  updatedAt: Scalars["Time"]["output"];
  user: User;
  userId: Scalars["ID"]["output"];
}

interface TimelineEvent {
  __typename?: "TimelineEvent";
  authorId: Scalars["String"]["output"];
  createdAt: Scalars["Time"]["output"];
  documentId: Scalars["ID"]["output"];
  event: TlEventPayload;
  id: Scalars["ID"]["output"];
  replyTo: Scalars["ID"]["output"];
  user: User;
}

type TimelineEventFilter = "ALL" | "COMMENTS" | "EDITS";

interface TimelineMessageInput {
  authorId: Scalars["String"]["input"];
  content: Scalars["String"]["input"];
  contentAddress: Scalars["String"]["input"];
  endID?: InputMaybe<Scalars["String"]["input"]>;
  replyTo?: InputMaybe<Scalars["String"]["input"]>;
  selectionMarkdown?: InputMaybe<Scalars["String"]["input"]>;
  startID?: InputMaybe<Scalars["String"]["input"]>;
}

interface UnauthenticatedSharedLink {
  __typename?: "UnauthenticatedSharedLink";
  documentTitle: Scalars["String"]["output"];
  inviteLink: Scalars["String"]["output"];
  invitedByEmail: Scalars["String"]["output"];
  invitedByName: Scalars["String"]["output"];
}

interface UpdateMessageResolutionInput {
  authorID: Scalars["String"]["input"];
  resolved: Scalars["Boolean"]["input"];
}

interface UpdateUserInput {
  displayName: Scalars["String"]["input"];
  name: Scalars["String"]["input"];
}

interface UpdateUserPreferenceInput {
  enableActivityNotifications?: InputMaybe<Scalars["Boolean"]["input"]>;
}

interface User {
  __typename?: "User";
  displayName: Scalars["String"]["output"];
  email: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
  isAdmin: Scalars["Boolean"]["output"];
  name: Scalars["String"]["output"];
  picture?: Maybe<Scalars["String"]["output"]>;
  subscriptionStatus: Scalars["String"]["output"];
}

interface UserPreference {
  __typename?: "UserPreference";
  enableActivityNotifications: Scalars["Boolean"]["output"];
}

type UploadAttachmentMutationVariables = Exact<{
  file: Scalars["Upload"]["input"];
  docId: Scalars["ID"]["input"];
}>;

type UploadAttachmentMutation = {
  __typename?: "Mutation";
  uploadAttachment: {
    __typename?: "DocumentAttachment";
    id: string;
    filename: string;
    contentType: string;
    createdAt: string;
  };
};

type ListDocumentAttachmentsQueryVariables = Exact<{
  docId: Scalars["ID"]["input"];
}>;

type ListDocumentAttachmentsQuery = {
  __typename?: "Query";
  listDocumentAttachments: Array<{
    __typename?: "DocumentAttachment";
    id: string;
    filename: string;
    contentType: string;
    createdAt: string;
  }>;
};

type ListUsersAttachmentsQueryVariables = Exact<{ [key: string]: never }>;

type ListUsersAttachmentsQuery = {
  __typename?: "Query";
  listUsersAttachments: Array<{
    __typename?: "DocumentAttachment";
    id: string;
    filename: string;
    contentType: string;
    createdAt: string;
  }>;
};

type DocumentFieldsFragment = {
  __typename?: "Document";
  id: string;
  title: string;
  isPublic: boolean;
  isFolder: boolean;
  folderID?: string | null;
  updatedAt: string;
  access: string;
  ownedBy: { __typename?: "User"; id: string };
  editors: Array<{
    __typename?: "User";
    id: string;
    name: string;
    displayName: string;
    email: string;
    picture?: string | null;
  }>;
  preferences: {
    __typename?: "DocumentPreference";
    enableFirstOpenNotifications: boolean;
    enableMentionNotifications: boolean;
    enableDMNotifications: boolean;
    enableAllCommentNotifications: boolean;
  };
};

type GetDocumentQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

type GetDocumentQuery = {
  __typename?: "Query";
  document?: {
    __typename?: "Document";
    id: string;
    title: string;
    isPublic: boolean;
    isFolder: boolean;
    folderID?: string | null;
    updatedAt: string;
    access: string;
    ownedBy: { __typename?: "User"; id: string };
    editors: Array<{
      __typename?: "User";
      id: string;
      name: string;
      displayName: string;
      email: string;
      picture?: string | null;
    }>;
    preferences: {
      __typename?: "DocumentPreference";
      enableFirstOpenNotifications: boolean;
      enableMentionNotifications: boolean;
      enableDMNotifications: boolean;
      enableAllCommentNotifications: boolean;
    };
  } | null;
};

type UpdateDocumentTitleMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
  input: DocumentInput;
}>;

type UpdateDocumentTitleMutation = {
  __typename?: "Mutation";
  updateDocument?: {
    __typename?: "Document";
    id: string;
    title: string;
  } | null;
};

type GetDocumentsQueryVariables = Exact<{
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
}>;

type GetDocumentsQuery = {
  __typename?: "Query";
  documents: {
    __typename?: "DocumentConnection";
    totalCount: number;
    edges: Array<{
      __typename?: "DocumentEdge";
      node: {
        __typename?: "Document";
        id: string;
        title: string;
        isFolder: boolean;
        folderID?: string | null;
        updatedAt: string;
        ownedBy: { __typename?: "User"; id: string };
      };
    }>;
    pageInfo: { __typename?: "PageInfo"; hasNextPage: boolean };
  };
};

type GetBaseDocumentsQueryVariables = Exact<{
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
}>;

type GetBaseDocumentsQuery = {
  __typename?: "Query";
  baseDocuments: {
    __typename?: "DocumentConnection";
    totalCount: number;
    edges: Array<{
      __typename?: "DocumentEdge";
      node: {
        __typename?: "Document";
        id: string;
        title: string;
        isFolder: boolean;
        folderID?: string | null;
        updatedAt: string;
        ownedBy: { __typename?: "User"; id: string };
      };
    }>;
    pageInfo: { __typename?: "PageInfo"; hasNextPage: boolean };
  };
};

type GetSharedDocumentsQueryVariables = Exact<{
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
}>;

type GetSharedDocumentsQuery = {
  __typename?: "Query";
  sharedDocuments: {
    __typename?: "DocumentConnection";
    totalCount: number;
    edges: Array<{
      __typename?: "DocumentEdge";
      node: {
        __typename?: "Document";
        id: string;
        title: string;
        isFolder: boolean;
        folderID?: string | null;
        updatedAt: string;
        ownedBy: { __typename?: "User"; id: string };
      };
    }>;
    pageInfo: { __typename?: "PageInfo"; hasNextPage: boolean };
  };
};

type GetFolderDocumentsQueryVariables = Exact<{
  folderId: Scalars["ID"]["input"];
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
}>;

type GetFolderDocumentsQuery = {
  __typename?: "Query";
  folderDocuments: {
    __typename?: "DocumentConnection";
    totalCount: number;
    edges: Array<{
      __typename?: "DocumentEdge";
      node: {
        __typename?: "Document";
        id: string;
        title: string;
        isFolder: boolean;
        folderID?: string | null;
        updatedAt: string;
        ownedBy: { __typename?: "User"; id: string };
      };
    }>;
    pageInfo: { __typename?: "PageInfo"; hasNextPage: boolean };
  };
};

type SearchDocumentsQueryVariables = Exact<{
  query: Scalars["String"]["input"];
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
}>;

type SearchDocumentsQuery = {
  __typename?: "Query";
  searchDocuments: {
    __typename?: "DocumentConnection";
    totalCount: number;
    edges: Array<{
      __typename?: "DocumentEdge";
      node: {
        __typename?: "Document";
        id: string;
        title: string;
        isFolder: boolean;
        folderID?: string | null;
        updatedAt: string;
        ownedBy: { __typename?: "User"; id: string };
      };
    }>;
    pageInfo: { __typename?: "PageInfo"; hasNextPage: boolean };
  };
};

type CreateDocumentMutationVariables = Exact<{ [key: string]: never }>;

type CreateDocumentMutation = {
  __typename?: "Mutation";
  createDocument: { __typename?: "Document"; id: string };
};

type CreateFolderMutationVariables = Exact<{ [key: string]: never }>;

type CreateFolderMutation = {
  __typename?: "Mutation";
  createFolder: {
    __typename?: "Document";
    id: string;
    title: string;
    isPublic: boolean;
    isFolder: boolean;
    folderID?: string | null;
    updatedAt: string;
    access: string;
    ownedBy: { __typename?: "User"; id: string };
    editors: Array<{
      __typename?: "User";
      id: string;
      name: string;
      displayName: string;
      email: string;
      picture?: string | null;
    }>;
    preferences: {
      __typename?: "DocumentPreference";
      enableFirstOpenNotifications: boolean;
      enableMentionNotifications: boolean;
      enableDMNotifications: boolean;
      enableAllCommentNotifications: boolean;
    };
  };
};

type DeleteDocumentMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
  deleteChildren?: InputMaybe<Scalars["Boolean"]["input"]>;
}>;

type DeleteDocumentMutation = {
  __typename?: "Mutation";
  deleteDocument?: boolean | null;
};

type ShareDocumentMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
  emails: Array<Scalars["String"]["input"]> | Scalars["String"]["input"];
  message?: InputMaybe<Scalars["String"]["input"]>;
}>;

type ShareDocumentMutation = {
  __typename?: "Mutation";
  shareDocument: Array<{
    __typename?: "SharedDocumentLink";
    inviteLink: string;
  }>;
};

type UnshareDocumentMutationVariables = Exact<{
  docId: Scalars["ID"]["input"];
  editorId: Scalars["ID"]["input"];
}>;

type UnshareDocumentMutation = {
  __typename?: "Mutation";
  unshareDocument: {
    __typename?: "Document";
    id: string;
    editors: Array<{
      __typename?: "User";
      id: string;
      name: string;
      email: string;
      picture?: string | null;
    }>;
  };
};

type UpdateSharedLinkMutationVariables = Exact<{
  inviteLink: Scalars["String"]["input"];
  isActive: Scalars["Boolean"]["input"];
}>;

type UpdateSharedLinkMutation = {
  __typename?: "Mutation";
  updateShareLink: { __typename?: "SharedDocumentLink"; inviteLink: string };
};

type UpdateDocumentVisibilityMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
  input: DocumentInput;
}>;

type UpdateDocumentVisibilityMutation = {
  __typename?: "Mutation";
  updateDocument?: {
    __typename?: "Document";
    id: string;
    isPublic: boolean;
  } | null;
};

type SharedDocumentLinksQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

type SharedDocumentLinksQuery = {
  __typename?: "Query";
  sharedLinks: Array<{
    __typename?: "SharedDocumentLink";
    inviteLink: string;
    inviteeEmail: string;
    isActive: boolean;
    invitedBy: { __typename?: "User"; name: string };
  }>;
};

type UpdateDocumentPreferencesMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  input: DocumentPreferenceInput;
}>;

type UpdateDocumentPreferencesMutation = {
  __typename?: "Mutation";
  updateDocumentPreference?: {
    __typename?: "DocumentPreference";
    enableFirstOpenNotifications: boolean;
    enableMentionNotifications: boolean;
    enableDMNotifications: boolean;
    enableAllCommentNotifications: boolean;
  } | null;
};

type CreateFlaggedVersionMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  input: FlaggedVersionInput;
}>;

type CreateFlaggedVersionMutation = {
  __typename?: "Mutation";
  createFlaggedVersion: boolean;
};

type EditFlaggedVersionMutationVariables = Exact<{
  flaggedVersionId: Scalars["ID"]["input"];
  input: FlaggedVersionInput;
}>;

type EditFlaggedVersionMutation = {
  __typename?: "Mutation";
  editFlaggedVersion: boolean;
};

type DeleteFlaggedVersionMutationVariables = Exact<{
  flaggedVersionId: Scalars["ID"]["input"];
  timelineEventId: Scalars["ID"]["input"];
}>;

type DeleteFlaggedVersionMutation = {
  __typename?: "Mutation";
  deleteFlaggedVersion: boolean;
};

type MoveDocumentMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  folderId?: InputMaybe<Scalars["ID"]["input"]>;
}>;

type MoveDocumentMutation = {
  __typename?: "Mutation";
  moveDocument?: {
    __typename?: "Document";
    id: string;
    title: string;
    isPublic: boolean;
    isFolder: boolean;
    folderID?: string | null;
    updatedAt: string;
    access: string;
    ownedBy: { __typename?: "User"; id: string };
    editors: Array<{
      __typename?: "User";
      id: string;
      name: string;
      displayName: string;
      email: string;
      picture?: string | null;
    }>;
    preferences: {
      __typename?: "DocumentPreference";
      enableFirstOpenNotifications: boolean;
      enableMentionNotifications: boolean;
      enableDMNotifications: boolean;
      enableAllCommentNotifications: boolean;
    };
  } | null;
};

type DocumentUpdatedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type DocumentUpdatedSubscription = {
  __typename?: "Subscription";
  documentUpdated: {
    __typename?: "Document";
    id: string;
    title: string;
    isPublic: boolean;
    isFolder: boolean;
    folderID?: string | null;
    updatedAt: string;
    access: string;
    ownedBy: { __typename?: "User"; id: string };
    editors: Array<{
      __typename?: "User";
      id: string;
      name: string;
      displayName: string;
      email: string;
      picture?: string | null;
    }>;
    preferences: {
      __typename?: "DocumentPreference";
      enableFirstOpenNotifications: boolean;
      enableMentionNotifications: boolean;
      enableDMNotifications: boolean;
      enableAllCommentNotifications: boolean;
    };
  };
};

type DocumentInsertedSubscriptionVariables = Exact<{
  userId: Scalars["ID"]["input"];
}>;

type DocumentInsertedSubscription = {
  __typename?: "Subscription";
  documentInserted: {
    __typename?: "Document";
    id: string;
    title: string;
    isPublic: boolean;
    isFolder: boolean;
    folderID?: string | null;
    updatedAt: string;
    access: string;
    ownedBy: { __typename?: "User"; id: string };
    editors: Array<{
      __typename?: "User";
      id: string;
      name: string;
      displayName: string;
      email: string;
      picture?: string | null;
    }>;
    preferences: {
      __typename?: "DocumentPreference";
      enableFirstOpenNotifications: boolean;
      enableMentionNotifications: boolean;
      enableDMNotifications: boolean;
      enableAllCommentNotifications: boolean;
    };
  };
};

type UploadImageMutationVariables = Exact<{
  file: Scalars["Upload"]["input"];
  docId: Scalars["ID"]["input"];
}>;

type UploadImageMutation = {
  __typename?: "Mutation";
  uploadImage: {
    __typename?: "Image";
    id: string;
    docId: string;
    url: string;
    createdAt: string;
    mimeType: string;
    status: Status;
    error?: string | null;
  };
};

type ListDocumentImagesQueryVariables = Exact<{
  docId: Scalars["ID"]["input"];
}>;

type ListDocumentImagesQuery = {
  __typename?: "Query";
  listDocumentImages: Array<{
    __typename?: "Image";
    id: string;
    docId: string;
    url: string;
    createdAt: string;
    mimeType: string;
    status: Status;
    error?: string | null;
  }>;
};

type GetImageSignedUrlQueryVariables = Exact<{
  docId: Scalars["ID"]["input"];
  imageId: Scalars["ID"]["input"];
}>;

type GetImageSignedUrlQuery = {
  __typename?: "Query";
  getImageSignedUrl: {
    __typename?: "SignedImageUrl";
    url: string;
    expiresAt: string;
  };
};

type GetImageQueryVariables = Exact<{
  docId: Scalars["ID"]["input"];
  imageId: Scalars["ID"]["input"];
}>;

type GetImageQuery = {
  __typename?: "Query";
  getImage: {
    __typename?: "Image";
    id: string;
    docId: string;
    url: string;
    createdAt: string;
    mimeType: string;
    status: Status;
    error?: string | null;
  };
};

type GetAiThreadsQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type GetAiThreadsQuery = {
  __typename?: "Query";
  getAskAiThreads: Array<{
    __typename?: "Thread";
    id: string;
    title: string;
    updatedAt: string;
    user: { __typename?: "User"; id: string; name: string };
  }>;
};

type MessageFieldsFragment = {
  __typename?: "Message";
  id: string;
  containerId: string;
  content: string;
  createdAt: string;
  lifecycleStage: LifecycleStage;
  lifecycleReason: string;
  authorId: string;
  hidden: boolean;
  user: {
    __typename?: "User";
    id: string;
    name: string;
    picture?: string | null;
  };
  aiContent?: {
    __typename?: "AiContent";
    concludingMessage?: string | null;
    feedback?: string | null;
  } | null;
  metadata: {
    __typename?: "MsgMetadata";
    allowDraftEdits: boolean;
    contentAddressBefore: string;
    contentAddress: string;
    contentAddressAfter: string;
    contentAddressAfterTimestamp?: string | null;
    revisionStatus: MessageRevisionStatus;
  };
  attachments: Array<
    | { __typename: "AttachedRevisoDocument"; id: string; title: string }
    | { __typename: "AttachmentContent"; text: string }
    | { __typename: "AttachmentError"; title: string; text: string }
    | {
        __typename: "AttachmentFile";
        id: string;
        filename: string;
        contentType: string;
      }
    | {
        __typename: "Revision";
        start: string;
        end: string;
        updated: string;
        beforeAddress?: string | null;
        afterAddress?: string | null;
        appliedOps?: string | null;
      }
    | { __typename: "Selection"; start: string; end: string; content: string }
    | { __typename: "Suggestion"; content: string }
  >;
};

type GetAiThreadMessagesQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  threadId: Scalars["ID"]["input"];
}>;

type GetAiThreadMessagesQuery = {
  __typename?: "Query";
  getAskAiThreadMessages: Array<{
    __typename?: "Message";
    id: string;
    containerId: string;
    content: string;
    createdAt: string;
    lifecycleStage: LifecycleStage;
    lifecycleReason: string;
    authorId: string;
    hidden: boolean;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    aiContent?: {
      __typename?: "AiContent";
      concludingMessage?: string | null;
      feedback?: string | null;
    } | null;
    metadata: {
      __typename?: "MsgMetadata";
      allowDraftEdits: boolean;
      contentAddressBefore: string;
      contentAddress: string;
      contentAddressAfter: string;
      contentAddressAfterTimestamp?: string | null;
      revisionStatus: MessageRevisionStatus;
    };
    attachments: Array<
      | { __typename: "AttachedRevisoDocument"; id: string; title: string }
      | { __typename: "AttachmentContent"; text: string }
      | { __typename: "AttachmentError"; title: string; text: string }
      | {
          __typename: "AttachmentFile";
          id: string;
          filename: string;
          contentType: string;
        }
      | {
          __typename: "Revision";
          start: string;
          end: string;
          updated: string;
          beforeAddress?: string | null;
          afterAddress?: string | null;
          appliedOps?: string | null;
        }
      | { __typename: "Selection"; start: string; end: string; content: string }
      | { __typename: "Suggestion"; content: string }
    >;
  }>;
};

type MessageUpsertedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  channelId: Scalars["ID"]["input"];
}>;

type MessageUpsertedSubscription = {
  __typename?: "Subscription";
  messageUpserted: {
    __typename?: "Message";
    id: string;
    containerId: string;
    content: string;
    createdAt: string;
    lifecycleStage: LifecycleStage;
    lifecycleReason: string;
    authorId: string;
    hidden: boolean;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    aiContent?: {
      __typename?: "AiContent";
      concludingMessage?: string | null;
      feedback?: string | null;
    } | null;
    metadata: {
      __typename?: "MsgMetadata";
      allowDraftEdits: boolean;
      contentAddressBefore: string;
      contentAddress: string;
      contentAddressAfter: string;
      contentAddressAfterTimestamp?: string | null;
      revisionStatus: MessageRevisionStatus;
    };
    attachments: Array<
      | { __typename: "AttachedRevisoDocument"; id: string; title: string }
      | { __typename: "AttachmentContent"; text: string }
      | { __typename: "AttachmentError"; title: string; text: string }
      | {
          __typename: "AttachmentFile";
          id: string;
          filename: string;
          contentType: string;
        }
      | {
          __typename: "Revision";
          start: string;
          end: string;
          updated: string;
          beforeAddress?: string | null;
          afterAddress?: string | null;
          appliedOps?: string | null;
        }
      | { __typename: "Selection"; start: string; end: string; content: string }
      | { __typename: "Suggestion"; content: string }
    >;
  };
};

type ThreadFieldsFragment = {
  __typename?: "Thread";
  id: string;
  title: string;
  updatedAt: string;
};

type UpdateMessageRevisionStatusMutationVariables = Exact<{
  containerId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  status: MessageRevisionStatus;
  contentAddress: Scalars["String"]["input"];
}>;

type UpdateMessageRevisionStatusMutation = {
  __typename?: "Mutation";
  updateMessageRevisionStatus: {
    __typename?: "Message";
    id: string;
    containerId: string;
    content: string;
    createdAt: string;
    lifecycleStage: LifecycleStage;
    lifecycleReason: string;
    authorId: string;
    hidden: boolean;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    aiContent?: {
      __typename?: "AiContent";
      concludingMessage?: string | null;
      feedback?: string | null;
    } | null;
    metadata: {
      __typename?: "MsgMetadata";
      allowDraftEdits: boolean;
      contentAddressBefore: string;
      contentAddress: string;
      contentAddressAfter: string;
      contentAddressAfterTimestamp?: string | null;
      revisionStatus: MessageRevisionStatus;
    };
    attachments: Array<
      | { __typename: "AttachedRevisoDocument"; id: string; title: string }
      | { __typename: "AttachmentContent"; text: string }
      | { __typename: "AttachmentError"; title: string; text: string }
      | {
          __typename: "AttachmentFile";
          id: string;
          filename: string;
          contentType: string;
        }
      | {
          __typename: "Revision";
          start: string;
          end: string;
          updated: string;
          beforeAddress?: string | null;
          afterAddress?: string | null;
          appliedOps?: string | null;
        }
      | { __typename: "Selection"; start: string; end: string; content: string }
      | { __typename: "Suggestion"; content: string }
    >;
  };
};

type CreateAiThreadMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type CreateAiThreadMutation = {
  __typename?: "Mutation";
  createAskAiThread: { __typename?: "Thread"; id: string };
};

type CreateAiThreadMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  threadId: Scalars["ID"]["input"];
  input: MessageInput;
}>;

type CreateAiThreadMessageMutation = {
  __typename?: "Mutation";
  createAskAiThreadMessage: {
    __typename?: "Message";
    id: string;
    containerId: string;
    content: string;
    createdAt: string;
    lifecycleStage: LifecycleStage;
    lifecycleReason: string;
    authorId: string;
    hidden: boolean;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    aiContent?: {
      __typename?: "AiContent";
      concludingMessage?: string | null;
      feedback?: string | null;
    } | null;
    metadata: {
      __typename?: "MsgMetadata";
      allowDraftEdits: boolean;
      contentAddressBefore: string;
      contentAddress: string;
      contentAddressAfter: string;
      contentAddressAfterTimestamp?: string | null;
      revisionStatus: MessageRevisionStatus;
    };
    attachments: Array<
      | { __typename: "AttachedRevisoDocument"; id: string; title: string }
      | { __typename: "AttachmentContent"; text: string }
      | { __typename: "AttachmentError"; title: string; text: string }
      | {
          __typename: "AttachmentFile";
          id: string;
          filename: string;
          contentType: string;
        }
      | {
          __typename: "Revision";
          start: string;
          end: string;
          updated: string;
          beforeAddress?: string | null;
          afterAddress?: string | null;
          appliedOps?: string | null;
        }
      | { __typename: "Selection"; start: string; end: string; content: string }
      | { __typename: "Suggestion"; content: string }
    >;
  };
};

type ThreadUpsertedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type ThreadUpsertedSubscription = {
  __typename?: "Subscription";
  threadUpserted: {
    __typename?: "Thread";
    id: string;
    title: string;
    updatedAt: string;
  };
};

type SubscriptionPlansQueryVariables = Exact<{ [key: string]: never }>;

type SubscriptionPlansQuery = {
  __typename?: "Query";
  subscriptionPlans: Array<{
    __typename?: "SubscriptionPlan";
    id: string;
    name: string;
    priceCents: number;
    currency: string;
    interval: string;
  }>;
};

type CheckoutMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

type CheckoutMutation = {
  __typename?: "Mutation";
  checkoutSubscriptionPlan: { __typename?: "Checkout"; url: string };
};

type BillingPortalSessionMutationVariables = Exact<{ [key: string]: never }>;

type BillingPortalSessionMutation = {
  __typename?: "Mutation";
  billingPortalSession: { __typename?: "BillingPortalSession"; url: string };
};

type TimelineEventFieldsFragment = {
  __typename?: "TimelineEvent";
  id: string;
  replyTo: string;
  createdAt: string;
  authorId: string;
  user: {
    __typename?: "User";
    id: string;
    name: string;
    picture?: string | null;
  };
  event:
    | {
        __typename: "TLAccessChangeV1";
        action: string;
        userIdentifiers: Array<string>;
      }
    | {
        __typename: "TLAttributeChangeV1";
        attribute: string;
        oldValue: string;
        newValue: string;
      }
    | { __typename: "TLEmpty"; placeholder: string }
    | { __typename: "TLJoinV1"; action: string }
    | { __typename: "TLMarkerV1"; title: string }
    | { __typename: "TLMessageResolutionV1" }
    | {
        __typename: "TLMessageV1";
        eventId: string;
        content: string;
        contentAddress: string;
        selectionStartId: string;
        selectionEndId: string;
        selectionMarkdown: string;
        replies: Array<{
          __typename?: "TimelineEvent";
          id: string;
          replyTo: string;
          createdAt: string;
          authorId: string;
          user: {
            __typename?: "User";
            id: string;
            name: string;
            picture?: string | null;
          };
          event:
            | { __typename: "TLAccessChangeV1" }
            | { __typename: "TLAttributeChangeV1" }
            | { __typename: "TLEmpty" }
            | { __typename: "TLJoinV1" }
            | { __typename: "TLMarkerV1" }
            | {
                __typename: "TLMessageResolutionV1";
                eventId: string;
                resolutionSummary: string;
                resolved: boolean;
              }
            | {
                __typename: "TLMessageV1";
                eventId: string;
                content: string;
                contentAddress: string;
                selectionStartId: string;
                selectionEndId: string;
                selectionMarkdown: string;
              }
            | { __typename: "TLPasteV1" }
            | { __typename: "TLUpdateV1" };
        }>;
      }
    | {
        __typename: "TLPasteV1";
        contentAddressBefore: string;
        contentAddressAfter: string;
      }
    | {
        __typename: "TLUpdateV1";
        eventId: string;
        title: string;
        content: string;
        startingContentAddress: string;
        endingContentAddress: string;
        flaggedVersionName?: string | null;
        flaggedVersionCreatedAt?: string | null;
        flaggedVersionID?: string | null;
        state: TlUpdateState;
        flaggedByUser?: { __typename?: "User"; name: string } | null;
      };
};

type GetDocumentTimelineQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  filter?: InputMaybe<TimelineEventFilter>;
}>;

type GetDocumentTimelineQuery = {
  __typename?: "Query";
  getDocumentTimeline: Array<{
    __typename?: "TimelineEvent";
    id: string;
    replyTo: string;
    createdAt: string;
    authorId: string;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    event:
      | {
          __typename: "TLAccessChangeV1";
          action: string;
          userIdentifiers: Array<string>;
        }
      | {
          __typename: "TLAttributeChangeV1";
          attribute: string;
          oldValue: string;
          newValue: string;
        }
      | { __typename: "TLEmpty"; placeholder: string }
      | { __typename: "TLJoinV1"; action: string }
      | { __typename: "TLMarkerV1"; title: string }
      | { __typename: "TLMessageResolutionV1" }
      | {
          __typename: "TLMessageV1";
          eventId: string;
          content: string;
          contentAddress: string;
          selectionStartId: string;
          selectionEndId: string;
          selectionMarkdown: string;
          replies: Array<{
            __typename?: "TimelineEvent";
            id: string;
            replyTo: string;
            createdAt: string;
            authorId: string;
            user: {
              __typename?: "User";
              id: string;
              name: string;
              picture?: string | null;
            };
            event:
              | { __typename: "TLAccessChangeV1" }
              | { __typename: "TLAttributeChangeV1" }
              | { __typename: "TLEmpty" }
              | { __typename: "TLJoinV1" }
              | { __typename: "TLMarkerV1" }
              | {
                  __typename: "TLMessageResolutionV1";
                  eventId: string;
                  resolutionSummary: string;
                  resolved: boolean;
                }
              | {
                  __typename: "TLMessageV1";
                  eventId: string;
                  content: string;
                  contentAddress: string;
                  selectionStartId: string;
                  selectionEndId: string;
                  selectionMarkdown: string;
                }
              | { __typename: "TLPasteV1" }
              | { __typename: "TLUpdateV1" };
          }>;
        }
      | {
          __typename: "TLPasteV1";
          contentAddressBefore: string;
          contentAddressAfter: string;
        }
      | {
          __typename: "TLUpdateV1";
          eventId: string;
          title: string;
          content: string;
          startingContentAddress: string;
          endingContentAddress: string;
          flaggedVersionName?: string | null;
          flaggedVersionCreatedAt?: string | null;
          flaggedVersionID?: string | null;
          state: TlUpdateState;
          flaggedByUser?: { __typename?: "User"; name: string } | null;
        };
  }>;
};

type TimelineEventInsertedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type TimelineEventInsertedSubscription = {
  __typename?: "Subscription";
  timelineEventInserted: {
    __typename?: "TimelineEvent";
    id: string;
    replyTo: string;
    createdAt: string;
    authorId: string;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    event:
      | {
          __typename: "TLAccessChangeV1";
          action: string;
          userIdentifiers: Array<string>;
        }
      | {
          __typename: "TLAttributeChangeV1";
          attribute: string;
          oldValue: string;
          newValue: string;
        }
      | { __typename: "TLEmpty"; placeholder: string }
      | { __typename: "TLJoinV1"; action: string }
      | { __typename: "TLMarkerV1"; title: string }
      | { __typename: "TLMessageResolutionV1" }
      | {
          __typename: "TLMessageV1";
          eventId: string;
          content: string;
          contentAddress: string;
          selectionStartId: string;
          selectionEndId: string;
          selectionMarkdown: string;
          replies: Array<{
            __typename?: "TimelineEvent";
            id: string;
            replyTo: string;
            createdAt: string;
            authorId: string;
            user: {
              __typename?: "User";
              id: string;
              name: string;
              picture?: string | null;
            };
            event:
              | { __typename: "TLAccessChangeV1" }
              | { __typename: "TLAttributeChangeV1" }
              | { __typename: "TLEmpty" }
              | { __typename: "TLJoinV1" }
              | { __typename: "TLMarkerV1" }
              | {
                  __typename: "TLMessageResolutionV1";
                  eventId: string;
                  resolutionSummary: string;
                  resolved: boolean;
                }
              | {
                  __typename: "TLMessageV1";
                  eventId: string;
                  content: string;
                  contentAddress: string;
                  selectionStartId: string;
                  selectionEndId: string;
                  selectionMarkdown: string;
                }
              | { __typename: "TLPasteV1" }
              | { __typename: "TLUpdateV1" };
          }>;
        }
      | {
          __typename: "TLPasteV1";
          contentAddressBefore: string;
          contentAddressAfter: string;
        }
      | {
          __typename: "TLUpdateV1";
          eventId: string;
          title: string;
          content: string;
          startingContentAddress: string;
          endingContentAddress: string;
          flaggedVersionName?: string | null;
          flaggedVersionCreatedAt?: string | null;
          flaggedVersionID?: string | null;
          state: TlUpdateState;
          flaggedByUser?: { __typename?: "User"; name: string } | null;
        };
  };
};

type TimelineEventUpdatedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type TimelineEventUpdatedSubscription = {
  __typename?: "Subscription";
  timelineEventUpdated: {
    __typename?: "TimelineEvent";
    id: string;
    replyTo: string;
    createdAt: string;
    authorId: string;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    event:
      | {
          __typename: "TLAccessChangeV1";
          action: string;
          userIdentifiers: Array<string>;
        }
      | {
          __typename: "TLAttributeChangeV1";
          attribute: string;
          oldValue: string;
          newValue: string;
        }
      | { __typename: "TLEmpty"; placeholder: string }
      | { __typename: "TLJoinV1"; action: string }
      | { __typename: "TLMarkerV1"; title: string }
      | { __typename: "TLMessageResolutionV1" }
      | {
          __typename: "TLMessageV1";
          eventId: string;
          content: string;
          contentAddress: string;
          selectionStartId: string;
          selectionEndId: string;
          selectionMarkdown: string;
          replies: Array<{
            __typename?: "TimelineEvent";
            id: string;
            replyTo: string;
            createdAt: string;
            authorId: string;
            user: {
              __typename?: "User";
              id: string;
              name: string;
              picture?: string | null;
            };
            event:
              | { __typename: "TLAccessChangeV1" }
              | { __typename: "TLAttributeChangeV1" }
              | { __typename: "TLEmpty" }
              | { __typename: "TLJoinV1" }
              | { __typename: "TLMarkerV1" }
              | {
                  __typename: "TLMessageResolutionV1";
                  eventId: string;
                  resolutionSummary: string;
                  resolved: boolean;
                }
              | {
                  __typename: "TLMessageV1";
                  eventId: string;
                  content: string;
                  contentAddress: string;
                  selectionStartId: string;
                  selectionEndId: string;
                  selectionMarkdown: string;
                }
              | { __typename: "TLPasteV1" }
              | { __typename: "TLUpdateV1" };
          }>;
        }
      | {
          __typename: "TLPasteV1";
          contentAddressBefore: string;
          contentAddressAfter: string;
        }
      | {
          __typename: "TLUpdateV1";
          eventId: string;
          title: string;
          content: string;
          startingContentAddress: string;
          endingContentAddress: string;
          flaggedVersionName?: string | null;
          flaggedVersionCreatedAt?: string | null;
          flaggedVersionID?: string | null;
          state: TlUpdateState;
          flaggedByUser?: { __typename?: "User"; name: string } | null;
        };
  };
};

type TimelineEventDeletedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type TimelineEventDeletedSubscription = {
  __typename?: "Subscription";
  timelineEventDeleted: {
    __typename?: "TimelineEvent";
    id: string;
    replyTo: string;
    event:
      | { __typename?: "TLAccessChangeV1" }
      | { __typename?: "TLAttributeChangeV1" }
      | { __typename?: "TLEmpty" }
      | { __typename?: "TLJoinV1" }
      | { __typename?: "TLMarkerV1" }
      | { __typename?: "TLMessageResolutionV1" }
      | { __typename?: "TLMessageV1"; eventId: string }
      | { __typename?: "TLPasteV1" }
      | { __typename?: "TLUpdateV1" };
  };
};

type CreateTimelineMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  input: TimelineMessageInput;
}>;

type CreateTimelineMessageMutation = {
  __typename?: "Mutation";
  createTimelineMessage: { __typename?: "TimelineEvent"; id: string };
};

type EditTimelineMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  input: EditTimelineMessageInput;
}>;

type EditTimelineMessageMutation = {
  __typename?: "Mutation";
  editTimelineMessage: { __typename?: "TimelineEvent"; id: string };
};

type UpdateMessageResolutionMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  input: UpdateMessageResolutionInput;
}>;

type UpdateMessageResolutionMutation = {
  __typename?: "Mutation";
  updateMessageResolution: { __typename?: "TimelineEvent"; id: string };
};

type EditTimelineUpdateSummaryMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  updateId: Scalars["ID"]["input"];
  summary: Scalars["String"]["input"];
}>;

type EditTimelineUpdateSummaryMutation = {
  __typename?: "Mutation";
  editTimelineUpdateSummary: {
    __typename?: "TimelineEvent";
    id: string;
    event:
      | { __typename?: "TLAccessChangeV1" }
      | { __typename?: "TLAttributeChangeV1" }
      | { __typename?: "TLEmpty" }
      | { __typename?: "TLJoinV1" }
      | { __typename?: "TLMarkerV1" }
      | { __typename?: "TLMessageResolutionV1" }
      | { __typename?: "TLMessageV1" }
      | { __typename?: "TLPasteV1" }
      | { __typename?: "TLUpdateV1"; eventId: string; content: string };
  };
};

type EditMessageResolutionSummaryMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  summary: Scalars["String"]["input"];
}>;

type EditMessageResolutionSummaryMutation = {
  __typename?: "Mutation";
  editMessageResolutionSummary: {
    __typename?: "TimelineEvent";
    id: string;
    event:
      | { __typename?: "TLAccessChangeV1" }
      | { __typename?: "TLAttributeChangeV1" }
      | { __typename?: "TLEmpty" }
      | { __typename?: "TLJoinV1" }
      | { __typename?: "TLMarkerV1" }
      | {
          __typename?: "TLMessageResolutionV1";
          eventId: string;
          resolutionSummary: string;
        }
      | { __typename?: "TLMessageV1" }
      | { __typename?: "TLPasteV1" }
      | { __typename?: "TLUpdateV1" };
  };
};

type ForceTimelineUpdateSummaryMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  userId: Scalars["String"]["input"];
}>;

type ForceTimelineUpdateSummaryMutation = {
  __typename?: "Mutation";
  forceTimelineUpdateSummary: boolean;
};

type DeleteTimelineMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}>;

type DeleteTimelineMessageMutation = {
  __typename?: "Mutation";
  deleteTimelineMessage: boolean;
};

type GetMeQueryVariables = Exact<{ [key: string]: never }>;

type GetMeQuery = {
  __typename?: "Query";
  me?: {
    __typename?: "User";
    id: string;
    email: string;
    name: string;
    displayName: string;
    picture?: string | null;
    isAdmin: boolean;
    subscriptionStatus: string;
  } | null;
};

type GetUserQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

type GetUserQuery = {
  __typename?: "Query";
  user?: {
    __typename?: "User";
    id: string;
    email: string;
    name: string;
    displayName: string;
    picture?: string | null;
    isAdmin: boolean;
  } | null;
};

type GetUsersQueryVariables = Exact<{
  ids: Array<Scalars["ID"]["input"]> | Scalars["ID"]["input"];
}>;

type GetUsersQuery = {
  __typename?: "Query";
  users?: Array<{
    __typename?: "User";
    id: string;
    email: string;
    name: string;
    displayName: string;
    picture?: string | null;
    isAdmin: boolean;
  } | null> | null;
};

type GetMyPreferenceQueryVariables = Exact<{ [key: string]: never }>;

type GetMyPreferenceQuery = {
  __typename?: "Query";
  myPreference: {
    __typename?: "UserPreference";
    enableActivityNotifications: boolean;
  };
};

type UpdateMyPreferenceMutationVariables = Exact<{
  input: UpdateUserPreferenceInput;
}>;

type UpdateMyPreferenceMutation = {
  __typename?: "Mutation";
  updateMyPreference: {
    __typename?: "UserPreference";
    enableActivityNotifications: boolean;
  };
};

type UpdateMeMutationVariables = Exact<{
  input: UpdateUserInput;
}>;

type UpdateMeMutation = {
  __typename?: "Mutation";
  updateMe?: {
    __typename?: "User";
    id: string;
    email: string;
    name: string;
    displayName: string;
    picture?: string | null;
    isAdmin: boolean;
  } | null;
};

type GetUsersInMyDomainQueryVariables = Exact<{
  includeSelf?: InputMaybe<Scalars["Boolean"]["input"]>;
}>;

type GetUsersInMyDomainQuery = {
  __typename?: "Query";
  usersInMyDomain: Array<{
    __typename?: "User";
    id: string;
    name: string;
    displayName: string;
    email: string;
    picture?: string | null;
    isAdmin: boolean;
  }>;
};
