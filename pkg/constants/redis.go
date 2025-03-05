package constants

// PubSub keys
const (
	ChannelNewDocsFormat             = "chanNewDocs:%s"              // userID (new docs)
	DocUpsertChanFormat              = "chanDoc:%s"                  // docID (model updates)
	UserDocUpdatesChanFormat         = "chanUserDoc:%s:%s"           // docID, userID
	MsgUpsertChanFormat              = "chanChanMsgs:%s"             // channelID
	ChannelUpsertChanFormat          = "chanDocChans:%s"             // docID
	DocUpdateChanFormat              = "chanDocUpdates:%s"           // docID (rogue updates)
	ActivityUnreadChanFormat         = "chanUnreadActivity:%s"       // userID
	UnreadChannelUpdateChanFormat    = "chanUnreadChannels:%s:%s"    // docID, userID
	UnreadMessageUpdateChanFormat    = "chanUnreadMessages:%s:%s"    // docID, userID
	ThreadUpdateChanFormat           = "chanThreads:%s:%s"           // docID, userID
	ChannelTimelineEventUpdateFormat = "chanTimelineEventsUpdate:%s" // docID
	ChannelTimelineEventInsertFormat = "chanTimelineEventsInsert:%s" // docID
	ChannelTimelineEventDeleteFormat = "chanTimelineEventsDelete:%s" // docID
)

const (
	DocCounterKeyFormat     = "doc:%s:counter"           // docID
	DocEventsKeyFormat      = "doc:%s:events"            // docID
	DocActiveConnectionsKey = "doc:%s:connections"       // docID
	DocUserConnectionKey    = "doc:%s:user:%s:author:%s" // docID, userID, authorID
	DocUserLastMessageKey   = "doc:%s:user:%s:message"   // docID, userID
)
