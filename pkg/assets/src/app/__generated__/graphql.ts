/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from "@graphql-typed-document-node/core";
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
export type MakeEmpty<
  T extends { [key: string]: unknown },
  K extends keyof T,
> = { [_ in K]?: never };
export type Incremental<T> =
  | T
  | {
      [P in keyof T]?: P extends " $fragmentName" | "__typename" ? T[P] : never;
    };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
  JSON: { input: string; output: string };
  Time: { input: string; output: string };
  Upload: { input: File; output: File };
};

export type AiContent = {
  __typename?: "AiContent";
  concludingMessage?: Maybe<Scalars["String"]["output"]>;
  feedback?: Maybe<Scalars["String"]["output"]>;
  notes?: Maybe<Scalars["String"]["output"]>;
};

export type AttachedRevisoDocument = {
  __typename?: "AttachedRevisoDocument";
  id: Scalars["String"]["output"];
  title: Scalars["String"]["output"];
};

export type AttachmentContent = {
  __typename?: "AttachmentContent";
  role: Scalars["String"]["output"];
  text: Scalars["String"]["output"];
};

export type AttachmentError = {
  __typename?: "AttachmentError";
  error: Scalars["String"]["output"];
  text: Scalars["String"]["output"];
  title: Scalars["String"]["output"];
};

export type AttachmentFile = {
  __typename?: "AttachmentFile";
  contentType: Scalars["String"]["output"];
  filename: Scalars["String"]["output"];
  id: Scalars["String"]["output"];
};

export type AttachmentInput = {
  contentType?: InputMaybe<Scalars["String"]["input"]>;
  id: Scalars["ID"]["input"];
  name: Scalars["String"]["input"];
  type: AttachmentInputType;
};

export type AttachmentInputType = "DRAFT" | "FILE" | "UNKNOWN";

export type AttachmentProgressType = "DONE" | "THINKING" | "UNKNOWN";

export type AttachmentValue =
  | AttachedRevisoDocument
  | AttachmentContent
  | AttachmentError
  | AttachmentFile
  | Revision
  | Selection
  | Suggestion;

export type BillingPortalSession = {
  __typename?: "BillingPortalSession";
  url: Scalars["String"]["output"];
};

export type Chain = {
  __typename?: "Chain";
  id: Scalars["ID"]["output"];
  messages: Array<Message>;
};

export type ChanType = "DIRECT" | "GENERAL" | "REVISO" | "UNKNOWN";

export type Checkout = {
  __typename?: "Checkout";
  url: Scalars["String"]["output"];
};

export type CommentNotificationPayloadValue = {
  __typename?: "CommentNotificationPayloadValue";
  author: User;
  authorId: Scalars["ID"]["output"];
  channelId: Scalars["ID"]["output"];
  commentType: CommentNotificationType;
  containerId: Scalars["ID"]["output"];
  documentId: Scalars["ID"]["output"];
  message: Message;
  messageId: Scalars["ID"]["output"];
};

export type CommentNotificationType =
  | "COMMENT"
  | "MENTION"
  | "REPLY"
  | "UNKNOWN";

export type ContentAddress = {
  __typename?: "ContentAddress";
  documentId: Scalars["ID"]["output"];
  id: Scalars["ID"]["output"];
  payload: Scalars["JSON"]["output"];
};

export type Document = {
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
};

export type DocumentAttachment = {
  __typename?: "DocumentAttachment";
  contentType: Scalars["String"]["output"];
  createdAt: Scalars["Time"]["output"];
  filename: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
};

export type DocumentConnection = {
  __typename?: "DocumentConnection";
  edges: Array<DocumentEdge>;
  pageInfo: PageInfo;
  totalCount: Scalars["Int"]["output"];
};

export type DocumentEdge = {
  __typename?: "DocumentEdge";
  cursor: Scalars["String"]["output"];
  node: Document;
};

export type DocumentInput = {
  isPublic?: InputMaybe<Scalars["Boolean"]["input"]>;
  title?: InputMaybe<Scalars["String"]["input"]>;
};

export type DocumentPreference = {
  __typename?: "DocumentPreference";
  enableAllCommentNotifications: Scalars["Boolean"]["output"];
  enableDMNotifications: Scalars["Boolean"]["output"];
  enableFirstOpenNotifications: Scalars["Boolean"]["output"];
  enableMentionNotifications: Scalars["Boolean"]["output"];
};

export type DocumentPreferenceInput = {
  enableAllCommentNotifications: Scalars["Boolean"]["input"];
  enableDMNotifications: Scalars["Boolean"]["input"];
  enableFirstOpenNotifications: Scalars["Boolean"]["input"];
  enableMentionNotifications: Scalars["Boolean"]["input"];
};

export type DocumentScreenshots = {
  __typename?: "DocumentScreenshots";
  darkUrl: Scalars["String"]["output"];
  lightUrl: Scalars["String"]["output"];
};

export type EditTimelineMessageInput = {
  content: Scalars["String"]["input"];
  contentAddress?: InputMaybe<Scalars["String"]["input"]>;
  endID?: InputMaybe<Scalars["String"]["input"]>;
  selectionMarkdown?: InputMaybe<Scalars["String"]["input"]>;
  startID?: InputMaybe<Scalars["String"]["input"]>;
};

export type FlaggedVersionInput = {
  name: Scalars["String"]["input"];
  updateID: Scalars["String"]["input"];
};

export type Image = {
  __typename?: "Image";
  createdAt: Scalars["Time"]["output"];
  docId: Scalars["ID"]["output"];
  error?: Maybe<Scalars["String"]["output"]>;
  id: Scalars["ID"]["output"];
  mimeType: Scalars["String"]["output"];
  status: Status;
  url: Scalars["String"]["output"];
};

export type LifecycleStage =
  | "COMPLETED"
  | "PENDING"
  | "REVISED"
  | "REVISING"
  | "UNKNOWN";

export type Message = {
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
};

export type MessageInput = {
  allowDraftEdits: Scalars["Boolean"]["input"];
  attachments?: InputMaybe<Array<AttachmentInput>>;
  authorId: Scalars["ID"]["input"];
  content: Scalars["String"]["input"];
  contentAddress: Scalars["String"]["input"];
  llm?: InputMaybe<MsgLlm>;
  replyTo?: InputMaybe<Scalars["ID"]["input"]>;
  selection?: InputMaybe<SelectionInput>;
};

export type MessageRevisionStatus = "ACCEPTED" | "DECLINED" | "UNSPECIFIED";

export type MessageUpdateInput = {
  content: Scalars["String"]["input"];
  id: Scalars["ID"]["input"];
};

export type MsgLlm = "CLAUDE" | "GPT4O";

export type MsgMetadata = {
  __typename?: "MsgMetadata";
  allowDraftEdits: Scalars["Boolean"]["output"];
  contentAddress: Scalars["String"]["output"];
  contentAddressAfter: Scalars["String"]["output"];
  contentAddressAfterTimestamp?: Maybe<Scalars["Time"]["output"]>;
  contentAddressBefore: Scalars["String"]["output"];
  llm: MsgLlm;
  revisionStatus: MessageRevisionStatus;
};

export type Mutation = {
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
};

export type MutationCheckoutSubscriptionPlanArgs = {
  id: Scalars["ID"]["input"];
};

export type MutationCopyDocumentArgs = {
  address?: InputMaybe<Scalars["String"]["input"]>;
  id: Scalars["ID"]["input"];
  isBranch?: InputMaybe<Scalars["Boolean"]["input"]>;
};

export type MutationCreateAskAiThreadArgs = {
  documentId: Scalars["ID"]["input"];
};

export type MutationCreateAskAiThreadMessageArgs = {
  documentId: Scalars["ID"]["input"];
  input: MessageInput;
  threadId: Scalars["ID"]["input"];
};

export type MutationCreateFlaggedVersionArgs = {
  documentId: Scalars["ID"]["input"];
  input: FlaggedVersionInput;
};

export type MutationCreateShareLinksArgs = {
  documentID: Scalars["ID"]["input"];
  emails: Array<Scalars["String"]["input"]>;
  message?: InputMaybe<Scalars["String"]["input"]>;
};

export type MutationCreateTimelineMessageArgs = {
  documentId: Scalars["ID"]["input"];
  input: TimelineMessageInput;
};

export type MutationDeleteDocumentArgs = {
  deleteChildren?: InputMaybe<Scalars["Boolean"]["input"]>;
  id: Scalars["ID"]["input"];
};

export type MutationDeleteFlaggedVersionArgs = {
  flaggedVersionId: Scalars["ID"]["input"];
  timelineEventId: Scalars["ID"]["input"];
};

export type MutationDeleteTimelineMessageArgs = {
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
};

export type MutationEditFlaggedVersionArgs = {
  flaggedVersionId: Scalars["ID"]["input"];
  input: FlaggedVersionInput;
};

export type MutationEditMessageResolutionSummaryArgs = {
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  summary: Scalars["String"]["input"];
};

export type MutationEditTimelineMessageArgs = {
  documentId: Scalars["ID"]["input"];
  input: EditTimelineMessageInput;
  messageId: Scalars["ID"]["input"];
};

export type MutationEditTimelineUpdateSummaryArgs = {
  documentId: Scalars["ID"]["input"];
  summary: Scalars["String"]["input"];
  updateId: Scalars["ID"]["input"];
};

export type MutationForceTimelineUpdateSummaryArgs = {
  documentId: Scalars["ID"]["input"];
  userId: Scalars["String"]["input"];
};

export type MutationJoinShareLinkArgs = {
  inviteLink: Scalars["String"]["input"];
};

export type MutationMoveDocumentArgs = {
  folderID?: InputMaybe<Scalars["ID"]["input"]>;
  id: Scalars["ID"]["input"];
};

export type MutationSaveContentAddressArgs = {
  documentId: Scalars["ID"]["input"];
  payload: Scalars["JSON"]["input"];
};

export type MutationSendAccessLinkForInviteArgs = {
  inviteLink: Scalars["String"]["input"];
};

export type MutationShareDocumentArgs = {
  documentID: Scalars["ID"]["input"];
  emails: Array<Scalars["String"]["input"]>;
  message?: InputMaybe<Scalars["String"]["input"]>;
};

export type MutationSoftDeleteDocumentArgs = {
  id: Scalars["ID"]["input"];
};

export type MutationUnshareDocumentArgs = {
  documentID: Scalars["ID"]["input"];
  editorID: Scalars["ID"]["input"];
};

export type MutationUpdateDocumentArgs = {
  id: Scalars["ID"]["input"];
  input: DocumentInput;
};

export type MutationUpdateDocumentPreferenceArgs = {
  id: Scalars["ID"]["input"];
  input: DocumentPreferenceInput;
};

export type MutationUpdateMeArgs = {
  input: UpdateUserInput;
};

export type MutationUpdateMessageResolutionArgs = {
  documentId: Scalars["ID"]["input"];
  input: UpdateMessageResolutionInput;
  messageId: Scalars["ID"]["input"];
};

export type MutationUpdateMessageRevisionStatusArgs = {
  containerId: Scalars["ID"]["input"];
  contentAddress: Scalars["String"]["input"];
  messageId: Scalars["ID"]["input"];
  status: MessageRevisionStatus;
};

export type MutationUpdateMyPreferenceArgs = {
  input: UpdateUserPreferenceInput;
};

export type MutationUpdateShareLinkArgs = {
  inviteLink: Scalars["String"]["input"];
  isActive: Scalars["Boolean"]["input"];
};

export type MutationUploadAttachmentArgs = {
  docId: Scalars["ID"]["input"];
  file: Scalars["Upload"]["input"];
};

export type MutationUploadImageArgs = {
  docId: Scalars["ID"]["input"];
  file: Scalars["Upload"]["input"];
};

export type MutationResponse = {
  __typename?: "MutationResponse";
  id: Scalars["ID"]["output"];
};

export type Notification = {
  __typename?: "Notification";
  createdAt: Scalars["Time"]["output"];
  documentId: Scalars["ID"]["output"];
  id: Scalars["ID"]["output"];
  payload: NotificationPayloadValue;
  read: Scalars["Boolean"]["output"];
};

export type NotificationConnection = {
  __typename?: "NotificationConnection";
  edges: Array<Notification>;
};

export type NotificationPayloadValue = CommentNotificationPayloadValue;

export type PageInfo = {
  __typename?: "PageInfo";
  hasNextPage: Scalars["Boolean"]["output"];
};

export type Query = {
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
};

export type QueryBaseDocumentsArgs = {
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  offset?: InputMaybe<Scalars["Int"]["input"]>;
};

export type QueryBranchesArgs = {
  id: Scalars["ID"]["input"];
};

export type QueryDocumentArgs = {
  id: Scalars["ID"]["input"];
};

export type QueryDocumentsArgs = {
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  offset?: InputMaybe<Scalars["Int"]["input"]>;
};

export type QueryFolderDocumentsArgs = {
  folderID: Scalars["ID"]["input"];
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  offset?: InputMaybe<Scalars["Int"]["input"]>;
};

export type QueryGetAskAiThreadMessagesArgs = {
  documentId: Scalars["ID"]["input"];
  threadId: Scalars["ID"]["input"];
};

export type QueryGetAskAiThreadsArgs = {
  documentId: Scalars["ID"]["input"];
};

export type QueryGetAttachmentSignedUrlArgs = {
  attachmentId: Scalars["ID"]["input"];
};

export type QueryGetContentAddressArgs = {
  addressId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
};

export type QueryGetDocumentTimelineArgs = {
  documentId: Scalars["ID"]["input"];
  filter?: InputMaybe<TimelineEventFilter>;
};

export type QueryGetImageArgs = {
  docId: Scalars["ID"]["input"];
  imageId: Scalars["ID"]["input"];
};

export type QueryGetImageSignedUrlArgs = {
  docId: Scalars["ID"]["input"];
  imageId: Scalars["ID"]["input"];
};

export type QueryListDocumentAttachmentsArgs = {
  docId: Scalars["ID"]["input"];
};

export type QueryListDocumentImagesArgs = {
  docId: Scalars["ID"]["input"];
};

export type QuerySearchDocumentsArgs = {
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  query: Scalars["String"]["input"];
};

export type QuerySharedDocumentsArgs = {
  limit?: InputMaybe<Scalars["Int"]["input"]>;
  offset?: InputMaybe<Scalars["Int"]["input"]>;
};

export type QuerySharedLinkArgs = {
  inviteLink: Scalars["String"]["input"];
};

export type QuerySharedLinksArgs = {
  documentID: Scalars["ID"]["input"];
};

export type QueryUnauthenticatedSharedLinkArgs = {
  inviteLink: Scalars["String"]["input"];
};

export type QueryUserArgs = {
  id: Scalars["ID"]["input"];
};

export type QueryUsersArgs = {
  ids: Array<Scalars["ID"]["input"]>;
};

export type QueryUsersInMyDomainArgs = {
  includeSelf?: InputMaybe<Scalars["Boolean"]["input"]>;
};

export type Revision = {
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
};

export type Selection = {
  __typename?: "Selection";
  content: Scalars["String"]["output"];
  end: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
  start: Scalars["String"]["output"];
};

export type SelectionInput = {
  content: Scalars["String"]["input"];
  end: Scalars["String"]["input"];
  start: Scalars["String"]["input"];
};

export type SharedDocumentLink = {
  __typename?: "SharedDocumentLink";
  createdAt: Scalars["Time"]["output"];
  document: Document;
  inviteLink: Scalars["String"]["output"];
  invitedBy: User;
  inviteeEmail: Scalars["String"]["output"];
  inviteeUser?: Maybe<User>;
  isActive: Scalars["Boolean"]["output"];
  updatedAt: Scalars["Time"]["output"];
};

export type SignedImageUrl = {
  __typename?: "SignedImageUrl";
  expiresAt: Scalars["Time"]["output"];
  url: Scalars["String"]["output"];
};

export type Status = "ERROR" | "LOADING" | "SUCCESS";

export type Subscription = {
  __typename?: "Subscription";
  documentInserted: Document;
  documentUpdated: Document;
  messageUpserted: Message;
  threadUpserted: Thread;
  timelineEventDeleted: TimelineEvent;
  timelineEventInserted: TimelineEvent;
  timelineEventUpdated: TimelineEvent;
};

export type SubscriptionDocumentInsertedArgs = {
  userId: Scalars["ID"]["input"];
};

export type SubscriptionDocumentUpdatedArgs = {
  documentId: Scalars["ID"]["input"];
};

export type SubscriptionMessageUpsertedArgs = {
  channelId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
};

export type SubscriptionThreadUpsertedArgs = {
  documentId: Scalars["ID"]["input"];
};

export type SubscriptionTimelineEventDeletedArgs = {
  documentId: Scalars["ID"]["input"];
};

export type SubscriptionTimelineEventInsertedArgs = {
  documentId: Scalars["ID"]["input"];
};

export type SubscriptionTimelineEventUpdatedArgs = {
  documentId: Scalars["ID"]["input"];
};

export type SubscriptionPlan = {
  __typename?: "SubscriptionPlan";
  currency: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
  interval: Scalars["String"]["output"];
  name: Scalars["String"]["output"];
  priceCents: Scalars["Int"]["output"];
};

export type Suggestion = {
  __typename?: "Suggestion";
  content: Scalars["String"]["output"];
};

export type TlAccessChangeV1 = {
  __typename?: "TLAccessChangeV1";
  action: Scalars["String"]["output"];
  userIdentifiers: Array<Scalars["String"]["output"]>;
};

export type TlAttributeChangeV1 = {
  __typename?: "TLAttributeChangeV1";
  attribute: Scalars["String"]["output"];
  newValue: Scalars["String"]["output"];
  oldValue: Scalars["String"]["output"];
};

export type TlEmpty = {
  __typename?: "TLEmpty";
  placeholder: Scalars["String"]["output"];
};

export type TlEventPayload =
  | TlAccessChangeV1
  | TlAttributeChangeV1
  | TlEmpty
  | TlJoinV1
  | TlMarkerV1
  | TlMessageResolutionV1
  | TlMessageV1
  | TlPasteV1
  | TlUpdateV1;

export type TlJoinV1 = {
  __typename?: "TLJoinV1";
  action: Scalars["String"]["output"];
};

export type TlMarkerV1 = {
  __typename?: "TLMarkerV1";
  title: Scalars["String"]["output"];
};

export type TlMessageResolutionV1 = {
  __typename?: "TLMessageResolutionV1";
  eventId: Scalars["String"]["output"];
  resolutionSummary: Scalars["String"]["output"];
  resolved: Scalars["Boolean"]["output"];
};

export type TlMessageV1 = {
  __typename?: "TLMessageV1";
  content: Scalars["String"]["output"];
  contentAddress: Scalars["String"]["output"];
  documentId: Scalars["ID"]["output"];
  eventId: Scalars["String"]["output"];
  replies: Array<TimelineEvent>;
  selectionEndId: Scalars["String"]["output"];
  selectionMarkdown: Scalars["String"]["output"];
  selectionStartId: Scalars["String"]["output"];
};

export type TlPasteV1 = {
  __typename?: "TLPasteV1";
  contentAddressAfter: Scalars["String"]["output"];
  contentAddressBefore: Scalars["String"]["output"];
};

export type TlUpdateState = "COMPLETE" | "SUMMARIZING";

export type TlUpdateV1 = {
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
};

export type Thread = {
  __typename?: "Thread";
  documentId: Scalars["ID"]["output"];
  id: Scalars["ID"]["output"];
  messages: Array<Message>;
  title: Scalars["String"]["output"];
  updatedAt: Scalars["Time"]["output"];
  user: User;
  userId: Scalars["ID"]["output"];
};

export type TimelineEvent = {
  __typename?: "TimelineEvent";
  authorId: Scalars["String"]["output"];
  createdAt: Scalars["Time"]["output"];
  documentId: Scalars["ID"]["output"];
  event: TlEventPayload;
  id: Scalars["ID"]["output"];
  replyTo: Scalars["ID"]["output"];
  user: User;
};

export type TimelineEventFilter = "ALL" | "COMMENTS" | "EDITS";

export type TimelineMessageInput = {
  authorId: Scalars["String"]["input"];
  content: Scalars["String"]["input"];
  contentAddress: Scalars["String"]["input"];
  endID?: InputMaybe<Scalars["String"]["input"]>;
  replyTo?: InputMaybe<Scalars["String"]["input"]>;
  selectionMarkdown?: InputMaybe<Scalars["String"]["input"]>;
  startID?: InputMaybe<Scalars["String"]["input"]>;
};

export type UnauthenticatedSharedLink = {
  __typename?: "UnauthenticatedSharedLink";
  documentTitle: Scalars["String"]["output"];
  inviteLink: Scalars["String"]["output"];
  invitedByEmail: Scalars["String"]["output"];
  invitedByName: Scalars["String"]["output"];
};

export type UpdateMessageResolutionInput = {
  authorID: Scalars["String"]["input"];
  resolved: Scalars["Boolean"]["input"];
};

export type UpdateUserInput = {
  displayName: Scalars["String"]["input"];
  name: Scalars["String"]["input"];
};

export type UpdateUserPreferenceInput = {
  enableActivityNotifications?: InputMaybe<Scalars["Boolean"]["input"]>;
};

export type User = {
  __typename?: "User";
  displayName: Scalars["String"]["output"];
  email: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
  isAdmin: Scalars["Boolean"]["output"];
  name: Scalars["String"]["output"];
  picture?: Maybe<Scalars["String"]["output"]>;
  subscriptionStatus: Scalars["String"]["output"];
};

export type UserPreference = {
  __typename?: "UserPreference";
  enableActivityNotifications: Scalars["Boolean"]["output"];
};

export type UploadAttachmentMutationVariables = Exact<{
  file: Scalars["Upload"]["input"];
  docId: Scalars["ID"]["input"];
}>;

export type UploadAttachmentMutation = {
  __typename?: "Mutation";
  uploadAttachment: {
    __typename?: "DocumentAttachment";
    id: string;
    filename: string;
    contentType: string;
    createdAt: string;
  };
};

export type ListDocumentAttachmentsQueryVariables = Exact<{
  docId: Scalars["ID"]["input"];
}>;

export type ListDocumentAttachmentsQuery = {
  __typename?: "Query";
  listDocumentAttachments: Array<{
    __typename?: "DocumentAttachment";
    id: string;
    filename: string;
    contentType: string;
    createdAt: string;
  }>;
};

export type ListUsersAttachmentsQueryVariables = Exact<{
  [key: string]: never;
}>;

export type ListUsersAttachmentsQuery = {
  __typename?: "Query";
  listUsersAttachments: Array<{
    __typename?: "DocumentAttachment";
    id: string;
    filename: string;
    contentType: string;
    createdAt: string;
  }>;
};

export type DocumentFieldsFragment = {
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
} & { " $fragmentName"?: "DocumentFieldsFragment" };

export type GetDocumentQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

export type GetDocumentQuery = {
  __typename?: "Query";
  document?:
    | ({ __typename?: "Document" } & {
        " $fragmentRefs"?: { DocumentFieldsFragment: DocumentFieldsFragment };
      })
    | null;
};

export type UpdateDocumentTitleMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
  input: DocumentInput;
}>;

export type UpdateDocumentTitleMutation = {
  __typename?: "Mutation";
  updateDocument?: {
    __typename?: "Document";
    id: string;
    title: string;
  } | null;
};

export type GetDocumentsQueryVariables = Exact<{
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
}>;

export type GetDocumentsQuery = {
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

export type GetBaseDocumentsQueryVariables = Exact<{
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
}>;

export type GetBaseDocumentsQuery = {
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

export type GetSharedDocumentsQueryVariables = Exact<{
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
}>;

export type GetSharedDocumentsQuery = {
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

export type GetFolderDocumentsQueryVariables = Exact<{
  folderId: Scalars["ID"]["input"];
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
}>;

export type GetFolderDocumentsQuery = {
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

export type SearchDocumentsQueryVariables = Exact<{
  query: Scalars["String"]["input"];
  offset?: InputMaybe<Scalars["Int"]["input"]>;
  limit?: InputMaybe<Scalars["Int"]["input"]>;
}>;

export type SearchDocumentsQuery = {
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

export type CreateDocumentMutationVariables = Exact<{ [key: string]: never }>;

export type CreateDocumentMutation = {
  __typename?: "Mutation";
  createDocument: { __typename?: "Document"; id: string };
};

export type CreateFolderMutationVariables = Exact<{ [key: string]: never }>;

export type CreateFolderMutation = {
  __typename?: "Mutation";
  createFolder: { __typename?: "Document" } & {
    " $fragmentRefs"?: { DocumentFieldsFragment: DocumentFieldsFragment };
  };
};

export type DeleteDocumentMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
  deleteChildren?: InputMaybe<Scalars["Boolean"]["input"]>;
}>;

export type DeleteDocumentMutation = {
  __typename?: "Mutation";
  deleteDocument?: boolean | null;
};

export type ShareDocumentMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
  emails: Array<Scalars["String"]["input"]> | Scalars["String"]["input"];
  message?: InputMaybe<Scalars["String"]["input"]>;
}>;

export type ShareDocumentMutation = {
  __typename?: "Mutation";
  shareDocument: Array<{
    __typename?: "SharedDocumentLink";
    inviteLink: string;
  }>;
};

export type UnshareDocumentMutationVariables = Exact<{
  docId: Scalars["ID"]["input"];
  editorId: Scalars["ID"]["input"];
}>;

export type UnshareDocumentMutation = {
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

export type UpdateSharedLinkMutationVariables = Exact<{
  inviteLink: Scalars["String"]["input"];
  isActive: Scalars["Boolean"]["input"];
}>;

export type UpdateSharedLinkMutation = {
  __typename?: "Mutation";
  updateShareLink: { __typename?: "SharedDocumentLink"; inviteLink: string };
};

export type UpdateDocumentVisibilityMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
  input: DocumentInput;
}>;

export type UpdateDocumentVisibilityMutation = {
  __typename?: "Mutation";
  updateDocument?: {
    __typename?: "Document";
    id: string;
    isPublic: boolean;
  } | null;
};

export type SharedDocumentLinksQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

export type SharedDocumentLinksQuery = {
  __typename?: "Query";
  sharedLinks: Array<{
    __typename?: "SharedDocumentLink";
    inviteLink: string;
    inviteeEmail: string;
    isActive: boolean;
    invitedBy: { __typename?: "User"; name: string };
  }>;
};

export type UpdateDocumentPreferencesMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  input: DocumentPreferenceInput;
}>;

export type UpdateDocumentPreferencesMutation = {
  __typename?: "Mutation";
  updateDocumentPreference?: {
    __typename?: "DocumentPreference";
    enableFirstOpenNotifications: boolean;
    enableMentionNotifications: boolean;
    enableDMNotifications: boolean;
    enableAllCommentNotifications: boolean;
  } | null;
};

export type CreateFlaggedVersionMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  input: FlaggedVersionInput;
}>;

export type CreateFlaggedVersionMutation = {
  __typename?: "Mutation";
  createFlaggedVersion: boolean;
};

export type EditFlaggedVersionMutationVariables = Exact<{
  flaggedVersionId: Scalars["ID"]["input"];
  input: FlaggedVersionInput;
}>;

export type EditFlaggedVersionMutation = {
  __typename?: "Mutation";
  editFlaggedVersion: boolean;
};

export type DeleteFlaggedVersionMutationVariables = Exact<{
  flaggedVersionId: Scalars["ID"]["input"];
  timelineEventId: Scalars["ID"]["input"];
}>;

export type DeleteFlaggedVersionMutation = {
  __typename?: "Mutation";
  deleteFlaggedVersion: boolean;
};

export type MoveDocumentMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  folderId?: InputMaybe<Scalars["ID"]["input"]>;
}>;

export type MoveDocumentMutation = {
  __typename?: "Mutation";
  moveDocument?:
    | ({ __typename?: "Document" } & {
        " $fragmentRefs"?: { DocumentFieldsFragment: DocumentFieldsFragment };
      })
    | null;
};

export type DocumentUpdatedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

export type DocumentUpdatedSubscription = {
  __typename?: "Subscription";
  documentUpdated: { __typename?: "Document" } & {
    " $fragmentRefs"?: { DocumentFieldsFragment: DocumentFieldsFragment };
  };
};

export type DocumentInsertedSubscriptionVariables = Exact<{
  userId: Scalars["ID"]["input"];
}>;

export type DocumentInsertedSubscription = {
  __typename?: "Subscription";
  documentInserted: { __typename?: "Document" } & {
    " $fragmentRefs"?: { DocumentFieldsFragment: DocumentFieldsFragment };
  };
};

export type UploadImageMutationVariables = Exact<{
  file: Scalars["Upload"]["input"];
  docId: Scalars["ID"]["input"];
}>;

export type UploadImageMutation = {
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

export type ListDocumentImagesQueryVariables = Exact<{
  docId: Scalars["ID"]["input"];
}>;

export type ListDocumentImagesQuery = {
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

export type GetImageSignedUrlQueryVariables = Exact<{
  docId: Scalars["ID"]["input"];
  imageId: Scalars["ID"]["input"];
}>;

export type GetImageSignedUrlQuery = {
  __typename?: "Query";
  getImageSignedUrl: {
    __typename?: "SignedImageUrl";
    url: string;
    expiresAt: string;
  };
};

export type GetImageQueryVariables = Exact<{
  docId: Scalars["ID"]["input"];
  imageId: Scalars["ID"]["input"];
}>;

export type GetImageQuery = {
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

export type GetAiThreadsQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

export type GetAiThreadsQuery = {
  __typename?: "Query";
  getAskAiThreads: Array<{
    __typename?: "Thread";
    id: string;
    title: string;
    updatedAt: string;
    user: { __typename?: "User"; id: string; name: string };
  }>;
};

export type MessageFieldsFragment = {
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
} & { " $fragmentName"?: "MessageFieldsFragment" };

export type GetAiThreadMessagesQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  threadId: Scalars["ID"]["input"];
}>;

export type GetAiThreadMessagesQuery = {
  __typename?: "Query";
  getAskAiThreadMessages: Array<
    { __typename?: "Message" } & {
      " $fragmentRefs"?: { MessageFieldsFragment: MessageFieldsFragment };
    }
  >;
};

export type MessageUpsertedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  channelId: Scalars["ID"]["input"];
}>;

export type MessageUpsertedSubscription = {
  __typename?: "Subscription";
  messageUpserted: { __typename?: "Message" } & {
    " $fragmentRefs"?: { MessageFieldsFragment: MessageFieldsFragment };
  };
};

export type ThreadFieldsFragment = {
  __typename?: "Thread";
  id: string;
  title: string;
  updatedAt: string;
} & { " $fragmentName"?: "ThreadFieldsFragment" };

export type UpdateMessageRevisionStatusMutationVariables = Exact<{
  containerId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  status: MessageRevisionStatus;
  contentAddress: Scalars["String"]["input"];
}>;

export type UpdateMessageRevisionStatusMutation = {
  __typename?: "Mutation";
  updateMessageRevisionStatus: { __typename?: "Message" } & {
    " $fragmentRefs"?: { MessageFieldsFragment: MessageFieldsFragment };
  };
};

export type CreateAiThreadMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

export type CreateAiThreadMutation = {
  __typename?: "Mutation";
  createAskAiThread: { __typename?: "Thread"; id: string };
};

export type CreateAiThreadMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  threadId: Scalars["ID"]["input"];
  input: MessageInput;
}>;

export type CreateAiThreadMessageMutation = {
  __typename?: "Mutation";
  createAskAiThreadMessage: { __typename?: "Message" } & {
    " $fragmentRefs"?: { MessageFieldsFragment: MessageFieldsFragment };
  };
};

export type ThreadUpsertedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

export type ThreadUpsertedSubscription = {
  __typename?: "Subscription";
  threadUpserted: { __typename?: "Thread" } & {
    " $fragmentRefs"?: { ThreadFieldsFragment: ThreadFieldsFragment };
  };
};

export type SubscriptionPlansQueryVariables = Exact<{ [key: string]: never }>;

export type SubscriptionPlansQuery = {
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

export type CheckoutMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

export type CheckoutMutation = {
  __typename?: "Mutation";
  checkoutSubscriptionPlan: { __typename?: "Checkout"; url: string };
};

export type BillingPortalSessionMutationVariables = Exact<{
  [key: string]: never;
}>;

export type BillingPortalSessionMutation = {
  __typename?: "Mutation";
  billingPortalSession: { __typename?: "BillingPortalSession"; url: string };
};

export type TimelineEventFieldsFragment = {
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
} & { " $fragmentName"?: "TimelineEventFieldsFragment" };

export type GetDocumentTimelineQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  filter?: InputMaybe<TimelineEventFilter>;
}>;

export type GetDocumentTimelineQuery = {
  __typename?: "Query";
  getDocumentTimeline: Array<
    { __typename?: "TimelineEvent" } & {
      " $fragmentRefs"?: {
        TimelineEventFieldsFragment: TimelineEventFieldsFragment;
      };
    }
  >;
};

export type TimelineEventInsertedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

export type TimelineEventInsertedSubscription = {
  __typename?: "Subscription";
  timelineEventInserted: { __typename?: "TimelineEvent" } & {
    " $fragmentRefs"?: {
      TimelineEventFieldsFragment: TimelineEventFieldsFragment;
    };
  };
};

export type TimelineEventUpdatedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

export type TimelineEventUpdatedSubscription = {
  __typename?: "Subscription";
  timelineEventUpdated: { __typename?: "TimelineEvent" } & {
    " $fragmentRefs"?: {
      TimelineEventFieldsFragment: TimelineEventFieldsFragment;
    };
  };
};

export type TimelineEventDeletedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

export type TimelineEventDeletedSubscription = {
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

export type CreateTimelineMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  input: TimelineMessageInput;
}>;

export type CreateTimelineMessageMutation = {
  __typename?: "Mutation";
  createTimelineMessage: { __typename?: "TimelineEvent"; id: string };
};

export type EditTimelineMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  input: EditTimelineMessageInput;
}>;

export type EditTimelineMessageMutation = {
  __typename?: "Mutation";
  editTimelineMessage: { __typename?: "TimelineEvent"; id: string };
};

export type UpdateMessageResolutionMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  input: UpdateMessageResolutionInput;
}>;

export type UpdateMessageResolutionMutation = {
  __typename?: "Mutation";
  updateMessageResolution: { __typename?: "TimelineEvent"; id: string };
};

export type EditTimelineUpdateSummaryMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  updateId: Scalars["ID"]["input"];
  summary: Scalars["String"]["input"];
}>;

export type EditTimelineUpdateSummaryMutation = {
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

export type EditMessageResolutionSummaryMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  summary: Scalars["String"]["input"];
}>;

export type EditMessageResolutionSummaryMutation = {
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

export type ForceTimelineUpdateSummaryMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  userId: Scalars["String"]["input"];
}>;

export type ForceTimelineUpdateSummaryMutation = {
  __typename?: "Mutation";
  forceTimelineUpdateSummary: boolean;
};

export type DeleteTimelineMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}>;

export type DeleteTimelineMessageMutation = {
  __typename?: "Mutation";
  deleteTimelineMessage: boolean;
};

export type GetMeQueryVariables = Exact<{ [key: string]: never }>;

export type GetMeQuery = {
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

export type GetUserQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

export type GetUserQuery = {
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

export type GetUsersQueryVariables = Exact<{
  ids: Array<Scalars["ID"]["input"]> | Scalars["ID"]["input"];
}>;

export type GetUsersQuery = {
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

export type GetMyPreferenceQueryVariables = Exact<{ [key: string]: never }>;

export type GetMyPreferenceQuery = {
  __typename?: "Query";
  myPreference: {
    __typename?: "UserPreference";
    enableActivityNotifications: boolean;
  };
};

export type UpdateMyPreferenceMutationVariables = Exact<{
  input: UpdateUserPreferenceInput;
}>;

export type UpdateMyPreferenceMutation = {
  __typename?: "Mutation";
  updateMyPreference: {
    __typename?: "UserPreference";
    enableActivityNotifications: boolean;
  };
};

export type UpdateMeMutationVariables = Exact<{
  input: UpdateUserInput;
}>;

export type UpdateMeMutation = {
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

export type GetUsersInMyDomainQueryVariables = Exact<{
  includeSelf?: InputMaybe<Scalars["Boolean"]["input"]>;
}>;

export type GetUsersInMyDomainQuery = {
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

export const DocumentFieldsFragmentDoc = {
  kind: "Document",
  definitions: [
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "DocumentFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Document" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "title" } },
          { kind: "Field", name: { kind: "Name", value: "isPublic" } },
          { kind: "Field", name: { kind: "Name", value: "isFolder" } },
          { kind: "Field", name: { kind: "Name", value: "folderID" } },
          { kind: "Field", name: { kind: "Name", value: "updatedAt" } },
          { kind: "Field", name: { kind: "Name", value: "access" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "ownedBy" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "editors" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "preferences" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableFirstOpenNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableMentionNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableDMNotifications" },
                },
                {
                  kind: "Field",
                  name: {
                    kind: "Name",
                    value: "enableAllCommentNotifications",
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<DocumentFieldsFragment, unknown>;
export const MessageFieldsFragmentDoc = {
  kind: "Document",
  definitions: [
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "MessageFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Message" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "containerId" } },
          { kind: "Field", name: { kind: "Name", value: "content" } },
          { kind: "Field", name: { kind: "Name", value: "createdAt" } },
          { kind: "Field", name: { kind: "Name", value: "lifecycleStage" } },
          { kind: "Field", name: { kind: "Name", value: "lifecycleReason" } },
          { kind: "Field", name: { kind: "Name", value: "authorId" } },
          { kind: "Field", name: { kind: "Name", value: "hidden" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "user" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "aiContent" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "concludingMessage" },
                },
                { kind: "Field", name: { kind: "Name", value: "feedback" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "metadata" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "allowDraftEdits" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressBefore" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddress" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressAfter" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressAfterTimestamp" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "revisionStatus" },
                },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "attachments" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "__typename" } },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Selection" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "start" } },
                      { kind: "Field", name: { kind: "Name", value: "end" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Revision" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "start" } },
                      { kind: "Field", name: { kind: "Name", value: "end" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "updated" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "beforeAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "afterAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "appliedOps" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Suggestion" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentContent" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "text" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentError" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                      { kind: "Field", name: { kind: "Name", value: "text" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentFile" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "filename" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentType" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachedRevisoDocument" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<MessageFieldsFragment, unknown>;
export const ThreadFieldsFragmentDoc = {
  kind: "Document",
  definitions: [
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "ThreadFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Thread" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "title" } },
          { kind: "Field", name: { kind: "Name", value: "updatedAt" } },
        ],
      },
    },
  ],
} as unknown as DocumentNode<ThreadFieldsFragment, unknown>;
export const TimelineEventFieldsFragmentDoc = {
  kind: "Document",
  definitions: [
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "TimelineEventFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "TimelineEvent" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "replyTo" } },
          { kind: "Field", name: { kind: "Name", value: "createdAt" } },
          { kind: "Field", name: { kind: "Name", value: "authorId" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "user" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "event" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "__typename" } },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLJoinV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "action" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLMarkerV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLUpdateV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "eventId" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "startingContentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "endingContentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedVersionName" },
                      },
                      {
                        kind: "Field",
                        name: {
                          kind: "Name",
                          value: "flaggedVersionCreatedAt",
                        },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedVersionID" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "state" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedByUser" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "name" },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLEmpty" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "placeholder" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLAttributeChangeV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "attribute" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "oldValue" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "newValue" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLAccessChangeV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "action" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "userIdentifiers" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLPasteV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddressBefore" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddressAfter" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLMessageV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "eventId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionStartId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionEndId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionMarkdown" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "replies" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "id" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "replyTo" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "createdAt" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "authorId" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "user" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "id" },
                                  },
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "name" },
                                  },
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "picture" },
                                  },
                                ],
                              },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "event" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "__typename" },
                                  },
                                  {
                                    kind: "InlineFragment",
                                    typeCondition: {
                                      kind: "NamedType",
                                      name: {
                                        kind: "Name",
                                        value: "TLMessageV1",
                                      },
                                    },
                                    selectionSet: {
                                      kind: "SelectionSet",
                                      selections: [
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "eventId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "content",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "contentAddress",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionStartId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionEndId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionMarkdown",
                                          },
                                        },
                                      ],
                                    },
                                  },
                                  {
                                    kind: "InlineFragment",
                                    typeCondition: {
                                      kind: "NamedType",
                                      name: {
                                        kind: "Name",
                                        value: "TLMessageResolutionV1",
                                      },
                                    },
                                    selectionSet: {
                                      kind: "SelectionSet",
                                      selections: [
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "eventId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "resolutionSummary",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "resolved",
                                          },
                                        },
                                      ],
                                    },
                                  },
                                ],
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<TimelineEventFieldsFragment, unknown>;
export const UploadAttachmentDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UploadAttachment" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "file" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "Upload" },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "docId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "uploadAttachment" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "file" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "file" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "docId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "docId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "filename" } },
                { kind: "Field", name: { kind: "Name", value: "contentType" } },
                { kind: "Field", name: { kind: "Name", value: "createdAt" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UploadAttachmentMutation,
  UploadAttachmentMutationVariables
>;
export const ListDocumentAttachmentsDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "ListDocumentAttachments" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "docId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "listDocumentAttachments" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "docId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "docId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "filename" } },
                { kind: "Field", name: { kind: "Name", value: "contentType" } },
                { kind: "Field", name: { kind: "Name", value: "createdAt" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  ListDocumentAttachmentsQuery,
  ListDocumentAttachmentsQueryVariables
>;
export const ListUsersAttachmentsDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "ListUsersAttachments" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "listUsersAttachments" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "filename" } },
                { kind: "Field", name: { kind: "Name", value: "contentType" } },
                { kind: "Field", name: { kind: "Name", value: "createdAt" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  ListUsersAttachmentsQuery,
  ListUsersAttachmentsQueryVariables
>;
export const GetDocumentDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetDocument" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "document" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "id" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "DocumentFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "DocumentFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Document" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "title" } },
          { kind: "Field", name: { kind: "Name", value: "isPublic" } },
          { kind: "Field", name: { kind: "Name", value: "isFolder" } },
          { kind: "Field", name: { kind: "Name", value: "folderID" } },
          { kind: "Field", name: { kind: "Name", value: "updatedAt" } },
          { kind: "Field", name: { kind: "Name", value: "access" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "ownedBy" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "editors" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "preferences" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableFirstOpenNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableMentionNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableDMNotifications" },
                },
                {
                  kind: "Field",
                  name: {
                    kind: "Name",
                    value: "enableAllCommentNotifications",
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetDocumentQuery, GetDocumentQueryVariables>;
export const UpdateDocumentTitleDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateDocumentTitle" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "DocumentInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateDocument" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "id" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "title" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateDocumentTitleMutation,
  UpdateDocumentTitleMutationVariables
>;
export const GetDocumentsDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetDocuments" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "offset" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Int" } },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "limit" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Int" } },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "documents" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "offset" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "offset" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "limit" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "limit" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "totalCount" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "edges" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "node" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "id" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "title" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "isFolder" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "folderID" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "updatedAt" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "ownedBy" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "id" },
                                  },
                                ],
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "pageInfo" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "hasNextPage" },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetDocumentsQuery, GetDocumentsQueryVariables>;
export const GetBaseDocumentsDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetBaseDocuments" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "offset" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Int" } },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "limit" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Int" } },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "baseDocuments" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "offset" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "offset" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "limit" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "limit" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "totalCount" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "edges" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "node" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "id" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "title" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "isFolder" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "folderID" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "updatedAt" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "ownedBy" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "id" },
                                  },
                                ],
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "pageInfo" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "hasNextPage" },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  GetBaseDocumentsQuery,
  GetBaseDocumentsQueryVariables
>;
export const GetSharedDocumentsDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetSharedDocuments" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "offset" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Int" } },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "limit" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Int" } },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "sharedDocuments" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "offset" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "offset" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "limit" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "limit" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "totalCount" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "edges" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "node" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "id" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "title" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "isFolder" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "folderID" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "updatedAt" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "ownedBy" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "id" },
                                  },
                                ],
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "pageInfo" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "hasNextPage" },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  GetSharedDocumentsQuery,
  GetSharedDocumentsQueryVariables
>;
export const GetFolderDocumentsDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetFolderDocuments" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "folderId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "offset" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Int" } },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "limit" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Int" } },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "folderDocuments" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "folderID" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "folderId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "offset" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "offset" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "limit" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "limit" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "totalCount" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "edges" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "node" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "id" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "title" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "isFolder" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "folderID" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "updatedAt" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "ownedBy" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "id" },
                                  },
                                ],
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "pageInfo" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "hasNextPage" },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  GetFolderDocumentsQuery,
  GetFolderDocumentsQueryVariables
>;
export const SearchDocumentsDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "SearchDocuments" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "query" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "offset" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Int" } },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "limit" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Int" } },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "searchDocuments" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "query" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "query" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "offset" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "offset" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "limit" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "limit" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "totalCount" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "edges" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "node" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "id" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "title" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "isFolder" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "folderID" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "updatedAt" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "ownedBy" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "id" },
                                  },
                                ],
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "pageInfo" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "hasNextPage" },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  SearchDocumentsQuery,
  SearchDocumentsQueryVariables
>;
export const CreateDocumentDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CreateDocument" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "createDocument" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CreateDocumentMutation,
  CreateDocumentMutationVariables
>;
export const CreateFolderDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CreateFolder" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "createFolder" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "DocumentFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "DocumentFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Document" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "title" } },
          { kind: "Field", name: { kind: "Name", value: "isPublic" } },
          { kind: "Field", name: { kind: "Name", value: "isFolder" } },
          { kind: "Field", name: { kind: "Name", value: "folderID" } },
          { kind: "Field", name: { kind: "Name", value: "updatedAt" } },
          { kind: "Field", name: { kind: "Name", value: "access" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "ownedBy" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "editors" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "preferences" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableFirstOpenNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableMentionNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableDMNotifications" },
                },
                {
                  kind: "Field",
                  name: {
                    kind: "Name",
                    value: "enableAllCommentNotifications",
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CreateFolderMutation,
  CreateFolderMutationVariables
>;
export const DeleteDocumentDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "DeleteDocument" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "deleteChildren" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Boolean" } },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "deleteDocument" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "id" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "deleteChildren" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "deleteChildren" },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  DeleteDocumentMutation,
  DeleteDocumentMutationVariables
>;
export const ShareDocumentDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "ShareDocument" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "emails" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "ListType",
              type: {
                kind: "NonNullType",
                type: {
                  kind: "NamedType",
                  name: { kind: "Name", value: "String" },
                },
              },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "message" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "String" } },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "shareDocument" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentID" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "emails" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "emails" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "message" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "message" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "inviteLink" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  ShareDocumentMutation,
  ShareDocumentMutationVariables
>;
export const UnshareDocumentDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UnshareDocument" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "docId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "editorId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "unshareDocument" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentID" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "docId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "editorID" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "editorId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "editors" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      { kind: "Field", name: { kind: "Name", value: "name" } },
                      { kind: "Field", name: { kind: "Name", value: "email" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "picture" },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UnshareDocumentMutation,
  UnshareDocumentMutationVariables
>;
export const UpdateSharedLinkDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateSharedLink" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "inviteLink" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "isActive" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "Boolean" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateShareLink" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "inviteLink" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "inviteLink" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "isActive" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "isActive" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "inviteLink" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateSharedLinkMutation,
  UpdateSharedLinkMutationVariables
>;
export const UpdateDocumentVisibilityDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateDocumentVisibility" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "DocumentInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateDocument" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "id" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "isPublic" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateDocumentVisibilityMutation,
  UpdateDocumentVisibilityMutationVariables
>;
export const SharedDocumentLinksDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "SharedDocumentLinks" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "sharedLinks" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentID" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "inviteLink" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "inviteeEmail" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "invitedBy" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "name" } },
                    ],
                  },
                },
                { kind: "Field", name: { kind: "Name", value: "isActive" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  SharedDocumentLinksQuery,
  SharedDocumentLinksQueryVariables
>;
export const UpdateDocumentPreferencesDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateDocumentPreferences" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "DocumentPreferenceInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateDocumentPreference" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "id" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableFirstOpenNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableMentionNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableDMNotifications" },
                },
                {
                  kind: "Field",
                  name: {
                    kind: "Name",
                    value: "enableAllCommentNotifications",
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateDocumentPreferencesMutation,
  UpdateDocumentPreferencesMutationVariables
>;
export const CreateFlaggedVersionDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CreateFlaggedVersion" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "FlaggedVersionInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "createFlaggedVersion" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CreateFlaggedVersionMutation,
  CreateFlaggedVersionMutationVariables
>;
export const EditFlaggedVersionDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "EditFlaggedVersion" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "flaggedVersionId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "FlaggedVersionInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "editFlaggedVersion" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "flaggedVersionId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "flaggedVersionId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  EditFlaggedVersionMutation,
  EditFlaggedVersionMutationVariables
>;
export const DeleteFlaggedVersionDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "DeleteFlaggedVersion" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "flaggedVersionId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "timelineEventId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "deleteFlaggedVersion" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "flaggedVersionId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "flaggedVersionId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "timelineEventId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "timelineEventId" },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  DeleteFlaggedVersionMutation,
  DeleteFlaggedVersionMutationVariables
>;
export const MoveDocumentDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "MoveDocument" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "folderId" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "moveDocument" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "id" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "folderID" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "folderId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "DocumentFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "DocumentFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Document" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "title" } },
          { kind: "Field", name: { kind: "Name", value: "isPublic" } },
          { kind: "Field", name: { kind: "Name", value: "isFolder" } },
          { kind: "Field", name: { kind: "Name", value: "folderID" } },
          { kind: "Field", name: { kind: "Name", value: "updatedAt" } },
          { kind: "Field", name: { kind: "Name", value: "access" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "ownedBy" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "editors" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "preferences" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableFirstOpenNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableMentionNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableDMNotifications" },
                },
                {
                  kind: "Field",
                  name: {
                    kind: "Name",
                    value: "enableAllCommentNotifications",
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  MoveDocumentMutation,
  MoveDocumentMutationVariables
>;
export const DocumentUpdatedDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "subscription",
      name: { kind: "Name", value: "DocumentUpdated" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "documentUpdated" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "DocumentFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "DocumentFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Document" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "title" } },
          { kind: "Field", name: { kind: "Name", value: "isPublic" } },
          { kind: "Field", name: { kind: "Name", value: "isFolder" } },
          { kind: "Field", name: { kind: "Name", value: "folderID" } },
          { kind: "Field", name: { kind: "Name", value: "updatedAt" } },
          { kind: "Field", name: { kind: "Name", value: "access" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "ownedBy" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "editors" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "preferences" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableFirstOpenNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableMentionNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableDMNotifications" },
                },
                {
                  kind: "Field",
                  name: {
                    kind: "Name",
                    value: "enableAllCommentNotifications",
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  DocumentUpdatedSubscription,
  DocumentUpdatedSubscriptionVariables
>;
export const DocumentInsertedDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "subscription",
      name: { kind: "Name", value: "DocumentInserted" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "userId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "documentInserted" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "userId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "userId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "DocumentFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "DocumentFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Document" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "title" } },
          { kind: "Field", name: { kind: "Name", value: "isPublic" } },
          { kind: "Field", name: { kind: "Name", value: "isFolder" } },
          { kind: "Field", name: { kind: "Name", value: "folderID" } },
          { kind: "Field", name: { kind: "Name", value: "updatedAt" } },
          { kind: "Field", name: { kind: "Name", value: "access" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "ownedBy" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "editors" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "preferences" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableFirstOpenNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableMentionNotifications" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableDMNotifications" },
                },
                {
                  kind: "Field",
                  name: {
                    kind: "Name",
                    value: "enableAllCommentNotifications",
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  DocumentInsertedSubscription,
  DocumentInsertedSubscriptionVariables
>;
export const UploadImageDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UploadImage" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "file" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "Upload" },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "docId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "uploadImage" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "file" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "file" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "docId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "docId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "docId" } },
                { kind: "Field", name: { kind: "Name", value: "url" } },
                { kind: "Field", name: { kind: "Name", value: "createdAt" } },
                { kind: "Field", name: { kind: "Name", value: "mimeType" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
                { kind: "Field", name: { kind: "Name", value: "error" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<UploadImageMutation, UploadImageMutationVariables>;
export const ListDocumentImagesDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "ListDocumentImages" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "docId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "listDocumentImages" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "docId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "docId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "docId" } },
                { kind: "Field", name: { kind: "Name", value: "url" } },
                { kind: "Field", name: { kind: "Name", value: "createdAt" } },
                { kind: "Field", name: { kind: "Name", value: "mimeType" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
                { kind: "Field", name: { kind: "Name", value: "error" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  ListDocumentImagesQuery,
  ListDocumentImagesQueryVariables
>;
export const GetImageSignedUrlDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetImageSignedUrl" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "docId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "imageId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "getImageSignedUrl" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "docId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "docId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "imageId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "imageId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "url" } },
                { kind: "Field", name: { kind: "Name", value: "expiresAt" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  GetImageSignedUrlQuery,
  GetImageSignedUrlQueryVariables
>;
export const GetImageDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetImage" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "docId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "imageId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "getImage" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "docId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "docId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "imageId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "imageId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "docId" } },
                { kind: "Field", name: { kind: "Name", value: "url" } },
                { kind: "Field", name: { kind: "Name", value: "createdAt" } },
                { kind: "Field", name: { kind: "Name", value: "mimeType" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
                { kind: "Field", name: { kind: "Name", value: "error" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetImageQuery, GetImageQueryVariables>;
export const GetAiThreadsDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetAIThreads" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "getAskAiThreads" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "title" } },
                { kind: "Field", name: { kind: "Name", value: "updatedAt" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "user" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      { kind: "Field", name: { kind: "Name", value: "name" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetAiThreadsQuery, GetAiThreadsQueryVariables>;
export const GetAiThreadMessagesDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetAIThreadMessages" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "threadId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "getAskAiThreadMessages" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "threadId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "threadId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "MessageFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "MessageFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Message" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "containerId" } },
          { kind: "Field", name: { kind: "Name", value: "content" } },
          { kind: "Field", name: { kind: "Name", value: "createdAt" } },
          { kind: "Field", name: { kind: "Name", value: "lifecycleStage" } },
          { kind: "Field", name: { kind: "Name", value: "lifecycleReason" } },
          { kind: "Field", name: { kind: "Name", value: "authorId" } },
          { kind: "Field", name: { kind: "Name", value: "hidden" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "user" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "aiContent" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "concludingMessage" },
                },
                { kind: "Field", name: { kind: "Name", value: "feedback" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "metadata" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "allowDraftEdits" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressBefore" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddress" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressAfter" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressAfterTimestamp" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "revisionStatus" },
                },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "attachments" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "__typename" } },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Selection" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "start" } },
                      { kind: "Field", name: { kind: "Name", value: "end" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Revision" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "start" } },
                      { kind: "Field", name: { kind: "Name", value: "end" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "updated" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "beforeAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "afterAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "appliedOps" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Suggestion" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentContent" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "text" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentError" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                      { kind: "Field", name: { kind: "Name", value: "text" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentFile" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "filename" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentType" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachedRevisoDocument" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  GetAiThreadMessagesQuery,
  GetAiThreadMessagesQueryVariables
>;
export const MessageUpsertedDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "subscription",
      name: { kind: "Name", value: "MessageUpserted" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "channelId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "messageUpserted" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "channelId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "channelId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "MessageFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "MessageFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Message" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "containerId" } },
          { kind: "Field", name: { kind: "Name", value: "content" } },
          { kind: "Field", name: { kind: "Name", value: "createdAt" } },
          { kind: "Field", name: { kind: "Name", value: "lifecycleStage" } },
          { kind: "Field", name: { kind: "Name", value: "lifecycleReason" } },
          { kind: "Field", name: { kind: "Name", value: "authorId" } },
          { kind: "Field", name: { kind: "Name", value: "hidden" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "user" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "aiContent" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "concludingMessage" },
                },
                { kind: "Field", name: { kind: "Name", value: "feedback" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "metadata" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "allowDraftEdits" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressBefore" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddress" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressAfter" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressAfterTimestamp" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "revisionStatus" },
                },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "attachments" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "__typename" } },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Selection" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "start" } },
                      { kind: "Field", name: { kind: "Name", value: "end" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Revision" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "start" } },
                      { kind: "Field", name: { kind: "Name", value: "end" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "updated" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "beforeAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "afterAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "appliedOps" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Suggestion" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentContent" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "text" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentError" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                      { kind: "Field", name: { kind: "Name", value: "text" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentFile" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "filename" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentType" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachedRevisoDocument" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  MessageUpsertedSubscription,
  MessageUpsertedSubscriptionVariables
>;
export const UpdateMessageRevisionStatusDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateMessageRevisionStatus" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "containerId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "messageId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "status" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "MessageRevisionStatus" },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "contentAddress" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateMessageRevisionStatus" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "containerId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "containerId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "messageId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "messageId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "status" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "status" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "contentAddress" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "contentAddress" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "MessageFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "MessageFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Message" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "containerId" } },
          { kind: "Field", name: { kind: "Name", value: "content" } },
          { kind: "Field", name: { kind: "Name", value: "createdAt" } },
          { kind: "Field", name: { kind: "Name", value: "lifecycleStage" } },
          { kind: "Field", name: { kind: "Name", value: "lifecycleReason" } },
          { kind: "Field", name: { kind: "Name", value: "authorId" } },
          { kind: "Field", name: { kind: "Name", value: "hidden" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "user" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "aiContent" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "concludingMessage" },
                },
                { kind: "Field", name: { kind: "Name", value: "feedback" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "metadata" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "allowDraftEdits" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressBefore" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddress" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressAfter" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressAfterTimestamp" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "revisionStatus" },
                },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "attachments" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "__typename" } },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Selection" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "start" } },
                      { kind: "Field", name: { kind: "Name", value: "end" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Revision" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "start" } },
                      { kind: "Field", name: { kind: "Name", value: "end" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "updated" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "beforeAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "afterAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "appliedOps" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Suggestion" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentContent" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "text" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentError" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                      { kind: "Field", name: { kind: "Name", value: "text" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentFile" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "filename" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentType" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachedRevisoDocument" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateMessageRevisionStatusMutation,
  UpdateMessageRevisionStatusMutationVariables
>;
export const CreateAiThreadDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CreateAIThread" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "createAskAiThread" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CreateAiThreadMutation,
  CreateAiThreadMutationVariables
>;
export const CreateAiThreadMessageDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CreateAIThreadMessage" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "threadId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "MessageInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "createAskAiThreadMessage" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "threadId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "threadId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "MessageFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "MessageFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Message" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "containerId" } },
          { kind: "Field", name: { kind: "Name", value: "content" } },
          { kind: "Field", name: { kind: "Name", value: "createdAt" } },
          { kind: "Field", name: { kind: "Name", value: "lifecycleStage" } },
          { kind: "Field", name: { kind: "Name", value: "lifecycleReason" } },
          { kind: "Field", name: { kind: "Name", value: "authorId" } },
          { kind: "Field", name: { kind: "Name", value: "hidden" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "user" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "aiContent" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "concludingMessage" },
                },
                { kind: "Field", name: { kind: "Name", value: "feedback" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "metadata" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "allowDraftEdits" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressBefore" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddress" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressAfter" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "contentAddressAfterTimestamp" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "revisionStatus" },
                },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "attachments" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "__typename" } },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Selection" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "start" } },
                      { kind: "Field", name: { kind: "Name", value: "end" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Revision" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "start" } },
                      { kind: "Field", name: { kind: "Name", value: "end" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "updated" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "beforeAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "afterAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "appliedOps" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "Suggestion" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentContent" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "text" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentError" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                      { kind: "Field", name: { kind: "Name", value: "text" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachmentFile" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "filename" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentType" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "AttachedRevisoDocument" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CreateAiThreadMessageMutation,
  CreateAiThreadMessageMutationVariables
>;
export const ThreadUpsertedDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "subscription",
      name: { kind: "Name", value: "ThreadUpserted" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "threadUpserted" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "ThreadFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "ThreadFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "Thread" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "title" } },
          { kind: "Field", name: { kind: "Name", value: "updatedAt" } },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  ThreadUpsertedSubscription,
  ThreadUpsertedSubscriptionVariables
>;
export const SubscriptionPlansDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "SubscriptionPlans" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "subscriptionPlans" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "priceCents" } },
                { kind: "Field", name: { kind: "Name", value: "currency" } },
                { kind: "Field", name: { kind: "Name", value: "interval" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  SubscriptionPlansQuery,
  SubscriptionPlansQueryVariables
>;
export const CheckoutDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "Checkout" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "checkoutSubscriptionPlan" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "id" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "url" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<CheckoutMutation, CheckoutMutationVariables>;
export const BillingPortalSessionDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "BillingPortalSession" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "billingPortalSession" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "url" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  BillingPortalSessionMutation,
  BillingPortalSessionMutationVariables
>;
export const GetDocumentTimelineDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetDocumentTimeline" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "filter" },
          },
          type: {
            kind: "NamedType",
            name: { kind: "Name", value: "TimelineEventFilter" },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "getDocumentTimeline" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "filter" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "filter" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "TimelineEventFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "TimelineEventFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "TimelineEvent" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "replyTo" } },
          { kind: "Field", name: { kind: "Name", value: "createdAt" } },
          { kind: "Field", name: { kind: "Name", value: "authorId" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "user" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "event" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "__typename" } },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLJoinV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "action" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLMarkerV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLUpdateV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "eventId" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "startingContentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "endingContentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedVersionName" },
                      },
                      {
                        kind: "Field",
                        name: {
                          kind: "Name",
                          value: "flaggedVersionCreatedAt",
                        },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedVersionID" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "state" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedByUser" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "name" },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLEmpty" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "placeholder" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLAttributeChangeV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "attribute" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "oldValue" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "newValue" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLAccessChangeV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "action" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "userIdentifiers" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLPasteV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddressBefore" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddressAfter" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLMessageV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "eventId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionStartId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionEndId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionMarkdown" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "replies" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "id" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "replyTo" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "createdAt" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "authorId" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "user" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "id" },
                                  },
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "name" },
                                  },
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "picture" },
                                  },
                                ],
                              },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "event" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "__typename" },
                                  },
                                  {
                                    kind: "InlineFragment",
                                    typeCondition: {
                                      kind: "NamedType",
                                      name: {
                                        kind: "Name",
                                        value: "TLMessageV1",
                                      },
                                    },
                                    selectionSet: {
                                      kind: "SelectionSet",
                                      selections: [
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "eventId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "content",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "contentAddress",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionStartId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionEndId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionMarkdown",
                                          },
                                        },
                                      ],
                                    },
                                  },
                                  {
                                    kind: "InlineFragment",
                                    typeCondition: {
                                      kind: "NamedType",
                                      name: {
                                        kind: "Name",
                                        value: "TLMessageResolutionV1",
                                      },
                                    },
                                    selectionSet: {
                                      kind: "SelectionSet",
                                      selections: [
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "eventId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "resolutionSummary",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "resolved",
                                          },
                                        },
                                      ],
                                    },
                                  },
                                ],
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  GetDocumentTimelineQuery,
  GetDocumentTimelineQueryVariables
>;
export const TimelineEventInsertedDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "subscription",
      name: { kind: "Name", value: "TimelineEventInserted" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "timelineEventInserted" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "TimelineEventFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "TimelineEventFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "TimelineEvent" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "replyTo" } },
          { kind: "Field", name: { kind: "Name", value: "createdAt" } },
          { kind: "Field", name: { kind: "Name", value: "authorId" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "user" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "event" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "__typename" } },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLJoinV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "action" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLMarkerV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLUpdateV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "eventId" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "startingContentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "endingContentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedVersionName" },
                      },
                      {
                        kind: "Field",
                        name: {
                          kind: "Name",
                          value: "flaggedVersionCreatedAt",
                        },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedVersionID" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "state" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedByUser" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "name" },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLEmpty" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "placeholder" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLAttributeChangeV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "attribute" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "oldValue" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "newValue" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLAccessChangeV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "action" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "userIdentifiers" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLPasteV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddressBefore" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddressAfter" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLMessageV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "eventId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionStartId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionEndId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionMarkdown" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "replies" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "id" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "replyTo" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "createdAt" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "authorId" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "user" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "id" },
                                  },
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "name" },
                                  },
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "picture" },
                                  },
                                ],
                              },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "event" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "__typename" },
                                  },
                                  {
                                    kind: "InlineFragment",
                                    typeCondition: {
                                      kind: "NamedType",
                                      name: {
                                        kind: "Name",
                                        value: "TLMessageV1",
                                      },
                                    },
                                    selectionSet: {
                                      kind: "SelectionSet",
                                      selections: [
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "eventId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "content",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "contentAddress",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionStartId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionEndId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionMarkdown",
                                          },
                                        },
                                      ],
                                    },
                                  },
                                  {
                                    kind: "InlineFragment",
                                    typeCondition: {
                                      kind: "NamedType",
                                      name: {
                                        kind: "Name",
                                        value: "TLMessageResolutionV1",
                                      },
                                    },
                                    selectionSet: {
                                      kind: "SelectionSet",
                                      selections: [
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "eventId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "resolutionSummary",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "resolved",
                                          },
                                        },
                                      ],
                                    },
                                  },
                                ],
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  TimelineEventInsertedSubscription,
  TimelineEventInsertedSubscriptionVariables
>;
export const TimelineEventUpdatedDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "subscription",
      name: { kind: "Name", value: "TimelineEventUpdated" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "timelineEventUpdated" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "FragmentSpread",
                  name: { kind: "Name", value: "TimelineEventFields" },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: "FragmentDefinition",
      name: { kind: "Name", value: "TimelineEventFields" },
      typeCondition: {
        kind: "NamedType",
        name: { kind: "Name", value: "TimelineEvent" },
      },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          { kind: "Field", name: { kind: "Name", value: "id" } },
          { kind: "Field", name: { kind: "Name", value: "replyTo" } },
          { kind: "Field", name: { kind: "Name", value: "createdAt" } },
          { kind: "Field", name: { kind: "Name", value: "authorId" } },
          {
            kind: "Field",
            name: { kind: "Name", value: "user" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
              ],
            },
          },
          {
            kind: "Field",
            name: { kind: "Name", value: "event" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "__typename" } },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLJoinV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "action" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLMarkerV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLUpdateV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "eventId" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "title" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "startingContentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "endingContentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedVersionName" },
                      },
                      {
                        kind: "Field",
                        name: {
                          kind: "Name",
                          value: "flaggedVersionCreatedAt",
                        },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedVersionID" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "state" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "flaggedByUser" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "name" },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLEmpty" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "placeholder" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLAttributeChangeV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "attribute" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "oldValue" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "newValue" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLAccessChangeV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "action" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "userIdentifiers" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLPasteV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddressBefore" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddressAfter" },
                      },
                    ],
                  },
                },
                {
                  kind: "InlineFragment",
                  typeCondition: {
                    kind: "NamedType",
                    name: { kind: "Name", value: "TLMessageV1" },
                  },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "eventId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "content" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "contentAddress" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionStartId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionEndId" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "selectionMarkdown" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "replies" },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "id" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "replyTo" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "createdAt" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "authorId" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "user" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "id" },
                                  },
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "name" },
                                  },
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "picture" },
                                  },
                                ],
                              },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "event" },
                              selectionSet: {
                                kind: "SelectionSet",
                                selections: [
                                  {
                                    kind: "Field",
                                    name: { kind: "Name", value: "__typename" },
                                  },
                                  {
                                    kind: "InlineFragment",
                                    typeCondition: {
                                      kind: "NamedType",
                                      name: {
                                        kind: "Name",
                                        value: "TLMessageV1",
                                      },
                                    },
                                    selectionSet: {
                                      kind: "SelectionSet",
                                      selections: [
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "eventId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "content",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "contentAddress",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionStartId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionEndId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "selectionMarkdown",
                                          },
                                        },
                                      ],
                                    },
                                  },
                                  {
                                    kind: "InlineFragment",
                                    typeCondition: {
                                      kind: "NamedType",
                                      name: {
                                        kind: "Name",
                                        value: "TLMessageResolutionV1",
                                      },
                                    },
                                    selectionSet: {
                                      kind: "SelectionSet",
                                      selections: [
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "eventId",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "resolutionSummary",
                                          },
                                        },
                                        {
                                          kind: "Field",
                                          name: {
                                            kind: "Name",
                                            value: "resolved",
                                          },
                                        },
                                      ],
                                    },
                                  },
                                ],
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  TimelineEventUpdatedSubscription,
  TimelineEventUpdatedSubscriptionVariables
>;
export const TimelineEventDeletedDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "subscription",
      name: { kind: "Name", value: "TimelineEventDeleted" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "timelineEventDeleted" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "replyTo" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "event" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: { kind: "Name", value: "TLMessageV1" },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "eventId" },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  TimelineEventDeletedSubscription,
  TimelineEventDeletedSubscriptionVariables
>;
export const CreateTimelineMessageDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CreateTimelineMessage" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "TimelineMessageInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "createTimelineMessage" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CreateTimelineMessageMutation,
  CreateTimelineMessageMutationVariables
>;
export const EditTimelineMessageDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "EditTimelineMessage" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "messageId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "EditTimelineMessageInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "editTimelineMessage" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "messageId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "messageId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  EditTimelineMessageMutation,
  EditTimelineMessageMutationVariables
>;
export const UpdateMessageResolutionDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateMessageResolution" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "messageId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "UpdateMessageResolutionInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateMessageResolution" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "messageId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "messageId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateMessageResolutionMutation,
  UpdateMessageResolutionMutationVariables
>;
export const EditTimelineUpdateSummaryDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "EditTimelineUpdateSummary" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "updateId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "summary" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "editTimelineUpdateSummary" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "updateId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "updateId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "summary" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "summary" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "event" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: { kind: "Name", value: "TLUpdateV1" },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "eventId" },
                            },
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "content" },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  EditTimelineUpdateSummaryMutation,
  EditTimelineUpdateSummaryMutationVariables
>;
export const EditMessageResolutionSummaryDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "EditMessageResolutionSummary" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "messageId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "summary" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "editMessageResolutionSummary" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "messageId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "messageId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "summary" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "summary" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "event" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      {
                        kind: "InlineFragment",
                        typeCondition: {
                          kind: "NamedType",
                          name: {
                            kind: "Name",
                            value: "TLMessageResolutionV1",
                          },
                        },
                        selectionSet: {
                          kind: "SelectionSet",
                          selections: [
                            {
                              kind: "Field",
                              name: { kind: "Name", value: "eventId" },
                            },
                            {
                              kind: "Field",
                              name: {
                                kind: "Name",
                                value: "resolutionSummary",
                              },
                            },
                          ],
                        },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  EditMessageResolutionSummaryMutation,
  EditMessageResolutionSummaryMutationVariables
>;
export const ForceTimelineUpdateSummaryDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "ForceTimelineUpdateSummary" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "userId" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "forceTimelineUpdateSummary" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "userId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "userId" },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  ForceTimelineUpdateSummaryMutation,
  ForceTimelineUpdateSummaryMutationVariables
>;
export const DeleteTimelineMessageDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "DeleteTimelineMessage" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "documentId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "messageId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "deleteTimelineMessage" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "documentId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "documentId" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "messageId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "messageId" },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  DeleteTimelineMessageMutation,
  DeleteTimelineMessageMutationVariables
>;
export const GetMeDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetMe" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "me" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
                { kind: "Field", name: { kind: "Name", value: "isAdmin" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "subscriptionStatus" },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetMeQuery, GetMeQueryVariables>;
export const GetUserDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetUser" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "user" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "id" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
                { kind: "Field", name: { kind: "Name", value: "isAdmin" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetUserQuery, GetUserQueryVariables>;
export const GetUsersDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetUsers" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "ids" } },
          type: {
            kind: "NonNullType",
            type: {
              kind: "ListType",
              type: {
                kind: "NonNullType",
                type: {
                  kind: "NamedType",
                  name: { kind: "Name", value: "ID" },
                },
              },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "users" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "ids" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "ids" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
                { kind: "Field", name: { kind: "Name", value: "isAdmin" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetUsersQuery, GetUsersQueryVariables>;
export const GetMyPreferenceDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetMyPreference" },
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "myPreference" },
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableActivityNotifications" },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  GetMyPreferenceQuery,
  GetMyPreferenceQueryVariables
>;
export const UpdateMyPreferenceDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateMyPreference" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "UpdateUserPreferenceInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateMyPreference" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                {
                  kind: "Field",
                  name: { kind: "Name", value: "enableActivityNotifications" },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  UpdateMyPreferenceMutation,
  UpdateMyPreferenceMutationVariables
>;
export const UpdateMeDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "UpdateMe" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "UpdateUserInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "updateMe" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
                { kind: "Field", name: { kind: "Name", value: "isAdmin" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<UpdateMeMutation, UpdateMeMutationVariables>;
export const GetUsersInMyDomainDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetUsersInMyDomain" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "includeSelf" },
          },
          type: { kind: "NamedType", name: { kind: "Name", value: "Boolean" } },
          defaultValue: { kind: "BooleanValue", value: false },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "usersInMyDomain" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "includeSelf" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "includeSelf" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "name" } },
                { kind: "Field", name: { kind: "Name", value: "displayName" } },
                { kind: "Field", name: { kind: "Name", value: "email" } },
                { kind: "Field", name: { kind: "Name", value: "picture" } },
                { kind: "Field", name: { kind: "Name", value: "isAdmin" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  GetUsersInMyDomainQuery,
  GetUsersInMyDomainQueryVariables
>;
