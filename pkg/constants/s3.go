package constants

const (
	S3Prefix                       = "v3"
	LogsPrefix                     = "logs/%s"                         // docID
	ConvosPrefix                   = "convos/%s"                       // docID
	ConvosFileKeyFormat            = "convos/%s/%s"                    // docID, filepath
	DagsDir                        = "dags"                            //
	DagsDirByParentId              = "dags/%s"                         // dagName, parentId (usually the docID)
	DagsDirNamedFileKeyFormat      = "dags/%s/%s"                      // dagName, filepath
	DocumentSnapshotPrefix         = "%s/%s/snapshots"                 // S3Prefix, docID
	UserAvatarKeyFormat            = "images/users/%s/avatar.png"      // userID
	DocumentImageKeyFormat         = "images/documents/%s/%s.png"      // docID, imageID
	OriginalDocumentImageKeyFormat = "images/documents/%s/%s-original" // docID, imageID
	DocumentAttachmentOriginalKey  = "attachments/%s/original"         // fileID
	DocumentAttachmentFileKey      = "attachments/%s/%s"               // fileID, filename
)
