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
}

interface AiContent {
  __typename?: "AiContent";
  concludingMessage?: Maybe<Scalars["String"]["output"]>;
  feedback?: Maybe<Scalars["String"]["output"]>;
  notes?: Maybe<Scalars["String"]["output"]>;
}

type AttachmentValue =
  | Revision
  | Selection
  | Suggestion
  | { __typename?: "%other" };

interface Chain {
  __typename?: "Chain";
  id: Scalars["ID"]["output"];
  messages: Array<Message>;
}

type ChanType = "DIRECT" | "GENERAL" | "REVISO" | "UNKNOWN";

interface Channel {
  __typename?: "Channel";
  channelType: ChanType;
  documentId: Scalars["ID"]["output"];
  id: Scalars["ID"]["output"];
  isActive: Scalars["Boolean"]["output"];
  messages: Array<Message>;
  unreadMentionCount: Scalars["Int"]["output"];
  unreadMessageCount: Scalars["Int"]["output"];
  updatedAt: Scalars["Time"]["output"];
  users: Array<UserWithAccess>;
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
  createdAt: Scalars["Time"]["output"];
  editors: Array<User>;
  hasUnreadNotifications: Scalars["Boolean"]["output"];
  id: Scalars["ID"]["output"];
  isPublic: Scalars["Boolean"]["output"];
  ownedBy: User;
  preferences: DocumentPreference;
  screenshots?: Maybe<DocumentScreenshots>;
  title: Scalars["String"]["output"];
  updatedAt: Scalars["Time"]["output"];
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

type LifecycleStage = "COMPLETED" | "PENDING" | "REVISED" | "UNKNOWN";

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
  id: Scalars["ID"]["output"];
  lifecycleStage: LifecycleStage;
  metadata: MsgMetadata;
  parentContainerId?: Maybe<Scalars["ID"]["output"]>;
  parentMessageId?: Maybe<Scalars["ID"]["output"]>;
  read: Scalars["Boolean"]["output"];
  replies: Array<Message>;
  replyCount: Scalars["Int"]["output"];
  replyingUserIds: Array<Scalars["ID"]["output"]>;
  replyingUsers: Array<User>;
  unreadReplyCount: Scalars["Int"]["output"];
  user: User;
  userId: Scalars["ID"]["output"];
}

interface MessageInput {
  allowDraftEdits: Scalars["Boolean"]["input"];
  authorId: Scalars["ID"]["input"];
  content: Scalars["String"]["input"];
  replyTo?: InputMaybe<Scalars["ID"]["input"]>;
  selection?: InputMaybe<SelectionInput>;
}

type MessageRevisionStatus = "ACCEPTED" | "DECLINED" | "UNSPECIFIED";

interface MessageUpdateInput {
  content: Scalars["String"]["input"];
  id: Scalars["ID"]["input"];
}

interface MsgMetadata {
  __typename?: "MsgMetadata";
  contentAddress: Scalars["String"]["output"];
  contentAddressAfter: Scalars["String"]["output"];
  contentAddressAfterTimestamp?: Maybe<Scalars["Time"]["output"]>;
  contentAddressBefore: Scalars["String"]["output"];
  revisionStatus: MessageRevisionStatus;
}

interface Mutation {
  __typename?: "Mutation";
  createAskAiThread: Thread;
  createAskAiThreadMessage: Message;
  createChannel: Channel;
  createDocument: Document;
  createMessage: Message;
  createMessageToReviso: Message;
  createReplyMessage: Message;
  createShareLinks: Array<SharedDocumentLink>;
  createTimelineMessage: TimelineEvent;
  deleteDocument?: Maybe<Scalars["Boolean"]["output"]>;
  forceTimelineUpdateSummary: Scalars["Boolean"]["output"];
  joinShareLink: Document;
  markAllNotificationsAsRead: Scalars["Boolean"]["output"];
  markChannelAsRead: Scalars["Boolean"]["output"];
  markMessageAsRead: Scalars["Boolean"]["output"];
  markMessageAsUnread: Scalars["Boolean"]["output"];
  markNotificationAsUnread: Scalars["Boolean"]["output"];
  markNotificationsAsRead: Scalars["Boolean"]["output"];
  regenerateMessage: Message;
  saveContentAddress?: Maybe<MutationResponse>;
  sendAccessLinkForInvite?: Maybe<Scalars["Boolean"]["output"]>;
  shareDocument: Array<SharedDocumentLink>;
  softDeleteDocument?: Maybe<Scalars["Boolean"]["output"]>;
  unshareDocument: Document;
  updateDocument?: Maybe<Document>;
  updateDocumentPreference?: Maybe<DocumentPreference>;
  updateMe?: Maybe<User>;
  updateMessageRevisionStatus: Message;
  updateMyPreference: UserPreference;
  updateShareLink: SharedDocumentLink;
}

interface Mutation_CreateAskAiThreadArgs {
  documentId: Scalars["ID"]["input"];
}

interface Mutation_CreateAskAiThreadMessageArgs {
  documentId: Scalars["ID"]["input"];
  input: MessageInput;
  threadId: Scalars["ID"]["input"];
}

interface Mutation_CreateChannelArgs {
  documentId: Scalars["ID"]["input"];
}

interface Mutation_CreateMessageArgs {
  channelId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
  input: MessageInput;
}

interface Mutation_CreateMessageToRevisoArgs {
  documentId: Scalars["ID"]["input"];
  promptKey: Scalars["String"]["input"];
  replyMessageId?: InputMaybe<Scalars["ID"]["input"]>;
  selection?: InputMaybe<SelectionInput>;
}

interface Mutation_CreateReplyMessageArgs {
  channelId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
  input: MessageInput;
  messageId: Scalars["ID"]["input"];
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
  id: Scalars["ID"]["input"];
}

interface Mutation_ForceTimelineUpdateSummaryArgs {
  contentAddress: Scalars["String"]["input"];
  documentId: Scalars["ID"]["input"];
  userId: Scalars["String"]["input"];
}

interface Mutation_JoinShareLinkArgs {
  inviteLink: Scalars["String"]["input"];
}

interface Mutation_MarkAllNotificationsAsReadArgs {
  documentId: Scalars["ID"]["input"];
}

interface Mutation_MarkChannelAsReadArgs {
  channelId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
}

interface Mutation_MarkMessageAsReadArgs {
  channelId: Scalars["ID"]["input"];
  containerId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}

interface Mutation_MarkMessageAsUnreadArgs {
  channelId: Scalars["ID"]["input"];
  containerId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}

interface Mutation_MarkNotificationAsUnreadArgs {
  documentId: Scalars["ID"]["input"];
  id: Scalars["ID"]["input"];
}

interface Mutation_MarkNotificationsAsReadArgs {
  documentId: Scalars["ID"]["input"];
  ids: Array<Scalars["ID"]["input"]>;
}

interface Mutation_RegenerateMessageArgs {
  containerId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
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
  document?: Maybe<Document>;
  documents: DocumentConnection;
  getAskAiThreadMessages: Array<Message>;
  getAskAiThreads: Array<Thread>;
  getChannel: Channel;
  getChannelChains: Array<Chain>;
  getChannelMessage: Message;
  getChannelReplyMessages: Array<Message>;
  getChannels: Array<Channel>;
  getContentAddress?: Maybe<ContentAddress>;
  getDocumentTimeline: Array<TimelineEvent>;
  getNotificationsForDocument: NotificationConnection;
  getUnreadMessageCountForDocument: Scalars["Int"]["output"];
  getUnreadNotificationsCount: Scalars["Int"]["output"];
  getUnreadNotificationsCountForDocument: Scalars["Int"]["output"];
  me?: Maybe<User>;
  myPreference: UserPreference;
  sharedLink?: Maybe<SharedDocumentLink>;
  sharedLinks: Array<SharedDocumentLink>;
  unauthenticatedSharedLink?: Maybe<UnauthenticatedSharedLink>;
  user?: Maybe<User>;
  users?: Maybe<Array<Maybe<User>>>;
}

interface Query_DocumentArgs {
  id: Scalars["ID"]["input"];
}

interface Query_DocumentsArgs {
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

interface Query_GetChannelArgs {
  documentId: Scalars["ID"]["input"];
  id: Scalars["ID"]["input"];
}

interface Query_GetChannelChainsArgs {
  chains: Array<Scalars["String"]["input"]>;
  channelId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}

interface Query_GetChannelMessageArgs {
  channelId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}

interface Query_GetChannelReplyMessagesArgs {
  channelId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}

interface Query_GetChannelsArgs {
  documentId: Scalars["ID"]["input"];
}

interface Query_GetContentAddressArgs {
  addressId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
}

interface Query_GetDocumentTimelineArgs {
  documentId: Scalars["ID"]["input"];
  filter?: InputMaybe<TimelineEventFilter>;
}

interface Query_GetNotificationsForDocumentArgs {
  documentId: Scalars["ID"]["input"];
  read: Scalars["Boolean"]["input"];
}

interface Query_GetUnreadMessageCountForDocumentArgs {
  documentId: Scalars["ID"]["input"];
}

interface Query_GetUnreadNotificationsCountForDocumentArgs {
  documentId: Scalars["ID"]["input"];
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

interface Subscription {
  __typename?: "Subscription";
  channelUpserted: Channel;
  documentUpserted: Document;
  messageUpserted: Message;
  threadUpserted: Thread;
  timelineEventUpserted: TimelineEvent;
}

interface Subscription_ChannelUpsertedArgs {
  documentId: Scalars["ID"]["input"];
}

interface Subscription_DocumentUpsertedArgs {
  documentId: Scalars["ID"]["input"];
}

interface Subscription_MessageUpsertedArgs {
  channelId: Scalars["ID"]["input"];
  documentId: Scalars["ID"]["input"];
}

interface Subscription_ThreadUpsertedArgs {
  documentId: Scalars["ID"]["input"];
}

interface Subscription_TimelineEventUpsertedArgs {
  documentId: Scalars["ID"]["input"];
}

interface Suggestion {
  __typename?: "Suggestion";
  content: Scalars["String"]["output"];
}

type TlEventPayload =
  | TlJoinV1
  | TlMarkerV1
  | TlMessageV1
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

type TlUpdateState = "COMPLETE" | "SUMMARIZING";

interface TlUpdateV1 {
  __typename?: "TLUpdateV1";
  content: Scalars["String"]["output"];
  endingContentAddress: Scalars["String"]["output"];
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
}

interface TimelineEvent {
  __typename?: "TimelineEvent";
  authorId: Scalars["String"]["output"];
  createdAt: Scalars["Time"]["output"];
  documentId: Scalars["ID"]["output"];
  event: TlEventPayload;
  id: Scalars["ID"]["output"];
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
}

interface UserPreference {
  __typename?: "UserPreference";
  enableActivityNotifications: Scalars["Boolean"]["output"];
}

interface UserWithAccess {
  __typename?: "UserWithAccess";
  hasAccess: Scalars["Boolean"]["output"];
  user: User;
}

type ChannelFieldsFragment = {
  __typename?: "Channel";
  id: string;
  channelType: ChanType;
  isActive: boolean;
  unreadMessageCount: number;
  unreadMentionCount: number;
  users: Array<{
    __typename?: "UserWithAccess";
    hasAccess: boolean;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
  }>;
};

type GetChannelsQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type GetChannelsQuery = {
  __typename?: "Query";
  getChannels: Array<{
    __typename?: "Channel";
    id: string;
    channelType: ChanType;
    isActive: boolean;
    unreadMessageCount: number;
    unreadMentionCount: number;
    users: Array<{
      __typename?: "UserWithAccess";
      hasAccess: boolean;
      user: {
        __typename?: "User";
        id: string;
        name: string;
        picture?: string | null;
      };
    }>;
  }>;
};

type CreateChannelMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type CreateChannelMutation = {
  __typename?: "Mutation";
  createChannel: {
    __typename?: "Channel";
    id: string;
    channelType: ChanType;
    isActive: boolean;
    unreadMessageCount: number;
    unreadMentionCount: number;
    users: Array<{
      __typename?: "UserWithAccess";
      hasAccess: boolean;
      user: {
        __typename?: "User";
        id: string;
        name: string;
        picture?: string | null;
      };
    }>;
  };
};

type MessageFieldsFragment = {
  __typename?: "Message";
  id: string;
  chain: string;
  channelId: string;
  containerId: string;
  content: string;
  createdAt: string;
  replyCount: number;
  unreadReplyCount: number;
  lifecycleStage: LifecycleStage;
  parentContainerId?: string | null;
  forkedMessageIds: Array<string>;
  read: boolean;
  user: {
    __typename?: "User";
    id: string;
    name: string;
    picture?: string | null;
  };
  replyingUsers: Array<{
    __typename?: "User";
    id: string;
    name: string;
    picture?: string | null;
  }>;
  aiContent?: {
    __typename?: "AiContent";
    concludingMessage?: string | null;
  } | null;
  attachments: Array<
    | {
        __typename: "Revision";
        start: string;
        end: string;
        updated: string;
        marshalledOperations: string;
        followUps?: string | null;
      }
    | { __typename: "Selection"; start: string; end: string; content: string }
    | { __typename: "Suggestion"; content: string }
  >;
};

type ChainFieldsFragment = {
  __typename?: "Chain";
  id: string;
  messages: Array<{
    __typename?: "Message";
    id: string;
    chain: string;
    channelId: string;
    containerId: string;
    content: string;
    createdAt: string;
    replyCount: number;
    unreadReplyCount: number;
    lifecycleStage: LifecycleStage;
    parentContainerId?: string | null;
    forkedMessageIds: Array<string>;
    read: boolean;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    replyingUsers: Array<{
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    }>;
    aiContent?: {
      __typename?: "AiContent";
      concludingMessage?: string | null;
    } | null;
    attachments: Array<
      | {
          __typename: "Revision";
          start: string;
          end: string;
          updated: string;
          marshalledOperations: string;
          followUps?: string | null;
        }
      | { __typename: "Selection"; start: string; end: string; content: string }
      | { __typename: "Suggestion"; content: string }
    >;
  }>;
};

type GetChannelMessageQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  channelId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}>;

type GetChannelMessageQuery = {
  __typename?: "Query";
  getChannelMessage: {
    __typename?: "Message";
    id: string;
    chain: string;
    channelId: string;
    containerId: string;
    content: string;
    createdAt: string;
    replyCount: number;
    unreadReplyCount: number;
    lifecycleStage: LifecycleStage;
    parentContainerId?: string | null;
    forkedMessageIds: Array<string>;
    read: boolean;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    replyingUsers: Array<{
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    }>;
    aiContent?: {
      __typename?: "AiContent";
      concludingMessage?: string | null;
    } | null;
    attachments: Array<
      | {
          __typename: "Revision";
          start: string;
          end: string;
          updated: string;
          marshalledOperations: string;
          followUps?: string | null;
        }
      | { __typename: "Selection"; start: string; end: string; content: string }
      | { __typename: "Suggestion"; content: string }
    >;
  };
};

type GetChannelReplyMessagesQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  channelId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}>;

type GetChannelReplyMessagesQuery = {
  __typename?: "Query";
  getChannelReplyMessages: Array<{
    __typename?: "Message";
    id: string;
    chain: string;
    channelId: string;
    containerId: string;
    content: string;
    createdAt: string;
    replyCount: number;
    unreadReplyCount: number;
    lifecycleStage: LifecycleStage;
    parentContainerId?: string | null;
    forkedMessageIds: Array<string>;
    read: boolean;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    replyingUsers: Array<{
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    }>;
    aiContent?: {
      __typename?: "AiContent";
      concludingMessage?: string | null;
    } | null;
    attachments: Array<
      | {
          __typename: "Revision";
          start: string;
          end: string;
          updated: string;
          marshalledOperations: string;
          followUps?: string | null;
        }
      | { __typename: "Selection"; start: string; end: string; content: string }
      | { __typename: "Suggestion"; content: string }
    >;
  }>;
};

type GetChannelChainsQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  channelId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  chains: Array<Scalars["String"]["input"]> | Scalars["String"]["input"];
}>;

type GetChannelChainsQuery = {
  __typename?: "Query";
  getChannelChains: Array<{
    __typename?: "Chain";
    id: string;
    messages: Array<{
      __typename?: "Message";
      id: string;
      chain: string;
      channelId: string;
      containerId: string;
      content: string;
      createdAt: string;
      replyCount: number;
      unreadReplyCount: number;
      lifecycleStage: LifecycleStage;
      parentContainerId?: string | null;
      forkedMessageIds: Array<string>;
      read: boolean;
      user: {
        __typename?: "User";
        id: string;
        name: string;
        picture?: string | null;
      };
      replyingUsers: Array<{
        __typename?: "User";
        id: string;
        name: string;
        picture?: string | null;
      }>;
      aiContent?: {
        __typename?: "AiContent";
        concludingMessage?: string | null;
      } | null;
      attachments: Array<
        | {
            __typename: "Revision";
            start: string;
            end: string;
            updated: string;
            marshalledOperations: string;
            followUps?: string | null;
          }
        | {
            __typename: "Selection";
            start: string;
            end: string;
            content: string;
          }
        | { __typename: "Suggestion"; content: string }
      >;
    }>;
  }>;
};

type GetChannelQueryVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  channelId: Scalars["ID"]["input"];
}>;

type GetChannelQuery = {
  __typename?: "Query";
  getChannel: {
    __typename?: "Channel";
    id: string;
    isActive: boolean;
    messages: Array<{
      __typename?: "Message";
      id: string;
      chain: string;
      channelId: string;
      containerId: string;
      content: string;
      createdAt: string;
      replyCount: number;
      unreadReplyCount: number;
      lifecycleStage: LifecycleStage;
      parentContainerId?: string | null;
      forkedMessageIds: Array<string>;
      read: boolean;
      user: {
        __typename?: "User";
        id: string;
        name: string;
        picture?: string | null;
      };
      replyingUsers: Array<{
        __typename?: "User";
        id: string;
        name: string;
        picture?: string | null;
      }>;
      aiContent?: {
        __typename?: "AiContent";
        concludingMessage?: string | null;
      } | null;
      attachments: Array<
        | {
            __typename: "Revision";
            start: string;
            end: string;
            updated: string;
            marshalledOperations: string;
            followUps?: string | null;
          }
        | {
            __typename: "Selection";
            start: string;
            end: string;
            content: string;
          }
        | { __typename: "Suggestion"; content: string }
      >;
    }>;
  };
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
    chain: string;
    channelId: string;
    containerId: string;
    content: string;
    createdAt: string;
    replyCount: number;
    unreadReplyCount: number;
    lifecycleStage: LifecycleStage;
    parentContainerId?: string | null;
    forkedMessageIds: Array<string>;
    read: boolean;
    user: {
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    };
    replyingUsers: Array<{
      __typename?: "User";
      id: string;
      name: string;
      picture?: string | null;
    }>;
    aiContent?: {
      __typename?: "AiContent";
      concludingMessage?: string | null;
    } | null;
    attachments: Array<
      | {
          __typename: "Revision";
          start: string;
          end: string;
          updated: string;
          marshalledOperations: string;
          followUps?: string | null;
        }
      | { __typename: "Selection"; start: string; end: string; content: string }
      | { __typename: "Suggestion"; content: string }
    >;
  };
};

type ChannelUpsertedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type ChannelUpsertedSubscription = {
  __typename?: "Subscription";
  channelUpserted: {
    __typename?: "Channel";
    id: string;
    channelType: ChanType;
    isActive: boolean;
    unreadMessageCount: number;
    unreadMentionCount: number;
    users: Array<{
      __typename?: "UserWithAccess";
      hasAccess: boolean;
      user: {
        __typename?: "User";
        id: string;
        name: string;
        picture?: string | null;
      };
    }>;
  };
};

type CreateMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  channelId: Scalars["ID"]["input"];
  input: MessageInput;
}>;

type CreateMessageMutation = {
  __typename?: "Mutation";
  createMessage: { __typename?: "Message"; id: string };
};

type CreateReplyMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  channelId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
  input: MessageInput;
}>;

type CreateReplyMessageMutation = {
  __typename?: "Mutation";
  createReplyMessage: { __typename?: "Message"; id: string };
};

type ReadMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  channelId: Scalars["ID"]["input"];
  containerId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}>;

type ReadMessageMutation = {
  __typename?: "Mutation";
  markMessageAsRead: boolean;
};

type UnreadMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  channelId: Scalars["ID"]["input"];
  containerId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}>;

type UnreadMessageMutation = {
  __typename?: "Mutation";
  markMessageAsUnread: boolean;
};

type RegenerateMessageMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  containerId: Scalars["ID"]["input"];
  messageId: Scalars["ID"]["input"];
}>;

type RegenerateMessageMutation = {
  __typename?: "Mutation";
  regenerateMessage: { __typename?: "Message"; id: string };
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
    updatedAt: string;
    hasUnreadNotifications: boolean;
    editors: Array<{ __typename?: "User"; id: string }>;
    ownedBy: { __typename?: "User"; id: string; name: string };
  } | null;
};

type GetBasicDocumentQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

type GetBasicDocumentQuery = {
  __typename?: "Query";
  document?: {
    __typename?: "Document";
    id: string;
    title: string;
    isPublic: boolean;
    updatedAt: string;
  } | null;
};

type DocumentUpsertedSubscriptionVariables = Exact<{
  documentId: Scalars["ID"]["input"];
}>;

type DocumentUpsertedSubscription = {
  __typename?: "Subscription";
  documentUpserted: {
    __typename?: "Document";
    id: string;
    title: string;
    isPublic: boolean;
    updatedAt: string;
    hasUnreadNotifications: boolean;
    preferences: {
      __typename?: "DocumentPreference";
      enableFirstOpenNotifications: boolean;
      enableMentionNotifications: boolean;
      enableDMNotifications: boolean;
      enableAllCommentNotifications: boolean;
    };
    editors: Array<{ __typename?: "User"; id: string }>;
    ownedBy: { __typename?: "User"; id: string; name: string };
  };
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
        updatedAt: string;
        screenshots?: {
          __typename?: "DocumentScreenshots";
          lightUrl: string;
          darkUrl: string;
        } | null;
        ownedBy: { __typename?: "User"; id: string; name: string };
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

type DeleteDocumentMutationVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

type DeleteDocumentMutation = {
  __typename?: "Mutation";
  deleteDocument?: boolean | null;
};

type GetShareLinkQueryVariables = Exact<{
  inviteLink: Scalars["String"]["input"];
}>;

type GetShareLinkQuery = {
  __typename?: "Query";
  sharedLink?: {
    __typename?: "SharedDocumentLink";
    inviteeEmail: string;
    invitedBy: { __typename?: "User"; name: string; email: string };
    document: { __typename?: "Document"; title: string; id: string };
  } | null;
};

type GetUnauthenticatedShareLinkQueryVariables = Exact<{
  inviteLink: Scalars["String"]["input"];
}>;

type GetUnauthenticatedShareLinkQuery = {
  __typename?: "Query";
  unauthenticatedSharedLink?: {
    __typename?: "UnauthenticatedSharedLink";
    inviteLink: string;
    invitedByEmail: string;
    invitedByName: string;
    documentTitle: string;
  } | null;
};

type JoinDocumentMutationVariables = Exact<{
  inviteLink: Scalars["String"]["input"];
}>;

type JoinDocumentMutation = {
  __typename?: "Mutation";
  joinShareLink: { __typename?: "Document"; id: string };
};

type CloseShareLinkMutationVariables = Exact<{
  inviteLink: Scalars["String"]["input"];
}>;

type CloseShareLinkMutation = {
  __typename?: "Mutation";
  updateShareLink: { __typename?: "SharedDocumentLink"; inviteLink: string };
};

type CreateMessageToRevisoMutationVariables = Exact<{
  documentId: Scalars["ID"]["input"];
  replyMessageId?: InputMaybe<Scalars["ID"]["input"]>;
  promptKey: Scalars["String"]["input"];
  selection?: InputMaybe<SelectionInput>;
}>;

type CreateMessageToRevisoMutation = {
  __typename?: "Mutation";
  createMessageToReviso: {
    __typename?: "Message";
    id: string;
    channelId: string;
    containerId: string;
  };
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
  } | null;
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

type GetDocumentWithEditorsQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

type GetDocumentWithEditorsQuery = {
  __typename?: "Query";
  document?: {
    __typename?: "Document";
    id: string;
    title: string;
    isPublic: boolean;
    updatedAt: string;
    editors: Array<{
      __typename?: "User";
      id: string;
      name: string;
      displayName: string;
      email: string;
      picture?: string | null;
    }>;
  } | null;
};
