package constants

// Restricted at the DB level, update document_access table if you change these
const (
	AccessLevelNone  = "none" // No access, does not exist in DB
	AccessLevelRead  = "read"
	AccessLevelWrite = "write"
	AccessLevelOwner = "owner"
	AccessLevelAdmin = "admin"
)

var AccessLevels = []string{AccessLevelRead, AccessLevelWrite, AccessLevelOwner, AccessLevelAdmin}

var AccessLevelsWithEdit = []string{AccessLevelWrite, AccessLevelOwner, AccessLevelAdmin}
