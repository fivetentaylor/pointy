package graph_test

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/testutils"
)

type ListDocumentsResponse struct {
	Documents struct {
		TotalCount int `json:"totalCount"`
		Edges      []struct {
			Node struct {
				ID           string    `json:"id"`
				ParentID     *string   `json:"parentID"` // Using pointer as it might be null
				RootParentID string    `json:"rootParentID"`
				Title        string    `json:"title"`
				CreatedAt    time.Time `json:"createdAt"`
				UpdatedAt    time.Time `json:"updatedAt"`
				BranchCopies []struct {
					ID       string  `json:"id"`
					ParentID *string `json:"parentID"` // Using pointer as it might be null
					Title    string  `json:"title"`
				} `json:"branchCopies"`
			} `json:"node"`
			Cursor string `json:"cursor"`
		} `json:"edges"`
		PageInfo struct {
			HasNextPage bool `json:"hasNextPage"`
		} `json:"pageInfo"`
	} `json:"documents"`
}

var listDocumentsQuery = `query GetDocuments($limit: Int, $offset: Int) {
  documents(limit: $limit, offset: $offset) {
    totalCount
    edges {
      node {
        id
        parentID
        rootParentID
        title
        createdAt
        updatedAt
        branchCopies {
          id
          parentID
          title
        }
      }
      cursor
    }
    pageInfo {
      hasNextPage
    }
  }
}`

type CopyDocumentResponse struct {
	CopyDocument models.Document
}

var copyDocumentMutation = `mutation CopyDocument($id: ID!, $isBranch: Boolean, $address: String) {
  copyDocument(id: $id, isBranch: $isBranch, address: $address) {
    id
    parentID
    title
    ownedBy {
      id
      name  # Assuming User type has a name field
    }
    editors {
      id
      name  # Assuming User type has a name field
    }
    createdAt
    updatedAt
    isPublic
    hasUnreadNotifications
    preferences {
      enableFirstOpenNotifications
      enableMentionNotifications
      enableDMNotifications
      enableAllCommentNotifications
    }
    screenshots {
      lightUrl
      darkUrl
    }
  }
}`

var branchesQuery = `query GetBranches($id: ID!) {
  branches(id: $id) {
    id
    parentID
    title
    ownedBy {
      id
      name  # Assuming User type has a name field
    }
    editors {
      id
      name  # Assuming User type has a name field
    }
    createdAt
    updatedAt
    isPublic
    hasUnreadNotifications
    preferences {
      enableFirstOpenNotifications
      enableMentionNotifications
      enableDMNotifications
      enableAllCommentNotifications
    }
    screenshots {
      lightUrl
      darkUrl
    }
  }
}`

type DocumentResponse struct {
	Document models.Document
}

var documentQuery = `query GetDocument($id: ID!) {
  document(id: $id) {
    id
    parentID
    title
    ownedBy {
      id
      name  # Assuming User type has a name field
    }
    editors {
      id
      name  # Assuming User type has a name field
    }
    createdAt
    updatedAt
    isPublic
    hasUnreadNotifications
    preferences {
      enableFirstOpenNotifications
      enableMentionNotifications
      enableDMNotifications
      enableAllCommentNotifications
    }
    screenshots {
      lightUrl
      darkUrl
    }
  }
}`

var deleteDocumentMutation = `mutation DeleteDocument($id: ID!) {
  deleteDocument(id: $id)
}`

var softDeleteDocumentMutation = `mutation SoftDeleteDocument($id: ID!) {
  softDeleteDocument(id: $id)
}`

func (suite *GraphTestSuite) TestListDocuments() {
	ctx := testutils.TestContext()
	user := testutils.CreateAdmin(suite.T(), ctx)
	ctx = testutils.WithUserClaimForUser(user)(ctx)

	result, err := RunQuery[ListDocumentsResponse](ctx, suite.handler, makeQuery(listDocumentsQuery, nil))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 0, result.Documents.TotalCount)
}

func (suite *GraphTestSuite) TestCopyDocumentBasic() {
	ctx := testutils.TestContext()
	user := testutils.CreateAdmin(suite.T(), ctx)
	ctx = testutils.WithUserClaimForUser(user)(ctx)

	docID := uuid.NewString()
	testutils.CreateTestDocument(suite.T(), ctx, docID, "")
	testutils.AddOwnerToDocument(suite.T(), ctx, docID, user.ID)

	result, err := RunQuery[ListDocumentsResponse](ctx, suite.handler, makeQuery(listDocumentsQuery, nil))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 1, result.Documents.TotalCount)
	require.Equal(suite.T(), 0, len(result.Documents.Edges[0].Node.BranchCopies))
	require.Equal(suite.T(), docID, result.Documents.Edges[0].Node.RootParentID)

	_, err = suite.RunQuery(ctx, makeQuery(copyDocumentMutation, map[string]interface{}{
		"id": docID,
	}))
	require.NoError(suite.T(), err)

	result, err = RunQuery[ListDocumentsResponse](ctx, suite.handler, makeQuery(listDocumentsQuery, nil))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 2, result.Documents.TotalCount)
	require.Equal(suite.T(), 0, len(result.Documents.Edges[0].Node.BranchCopies))

	for _, edge := range result.Documents.Edges {
		require.Equal(suite.T(), edge.Node.ID, edge.Node.RootParentID)
	}
}

func (suite *GraphTestSuite) TestCopyDocumentBranch() {
	ctx := testutils.TestContext()
	user0 := testutils.CreateAdmin(suite.T(), ctx)
	user1 := testutils.CreateAdmin(suite.T(), ctx)
	ctxUser0 := testutils.WithUserClaimForUser(user0)(ctx)
	ctxUser1 := testutils.WithUserClaimForUser(user1)(ctx)

	docID := uuid.NewString()
	testutils.CreateTestDocument(suite.T(), ctxUser0, docID, "")
	testutils.AddOwnerToDocument(suite.T(), ctxUser0, docID, user0.ID)

	// user0 has access to the document
	result, err := RunQuery[ListDocumentsResponse](ctxUser0, suite.handler, makeQuery(listDocumentsQuery, nil))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 1, result.Documents.TotalCount)
	node := result.Documents.Edges[0].Node
	require.Equal(suite.T(), node.ID, node.RootParentID)

	// user0 makes a branch copy of the document
	branch0, err := RunQuery[CopyDocumentResponse](ctxUser0, suite.handler, makeQuery(copyDocumentMutation, map[string]interface{}{
		"id":       docID,
		"isBranch": true,
	}))
	require.NoError(suite.T(), err)

	// user0 now has access to the branch copy
	result, err = RunQuery[ListDocumentsResponse](ctxUser0, suite.handler, makeQuery(listDocumentsQuery, nil))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 1, result.Documents.TotalCount)
	require.Equal(suite.T(), 1, len(result.Documents.Edges[0].Node.BranchCopies))

	// user0 makes a branch copy of the branch copy
	branch1, err := RunQuery[CopyDocumentResponse](ctxUser0, suite.handler, makeQuery(copyDocumentMutation, map[string]interface{}{
		"id":       branch0.CopyDocument.ID,
		"isBranch": true,
	}))
	require.NoError(suite.T(), err)

	// user0 now has access to two branch copies
	result, err = RunQuery[ListDocumentsResponse](ctxUser0, suite.handler, makeQuery(listDocumentsQuery, nil))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 1, result.Documents.TotalCount)
	require.Equal(suite.T(), 2, len(result.Documents.Edges[0].Node.BranchCopies))

	// grant user1 access to branch0 and branch1
	testutils.AddOwnerToDocument(suite.T(), ctxUser0, branch0.CopyDocument.ID, user1.ID)
	testutils.AddOwnerToDocument(suite.T(), ctxUser0, branch1.CopyDocument.ID, user1.ID)

	// user1 now has access to branch0 and branch1
	result, err = RunQuery[ListDocumentsResponse](ctxUser1, suite.handler, makeQuery(listDocumentsQuery, nil))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 1, result.Documents.TotalCount)
	require.Equal(suite.T(), 1, len(result.Documents.Edges[0].Node.BranchCopies))

	// user0 deletes the document
	_, err = suite.RunQuery(ctxUser0, makeQuery(deleteDocumentMutation, map[string]interface{}{
		"id": docID,
	}))

	// user0 and user1 now have no documents
	result, err = RunQuery[ListDocumentsResponse](ctxUser0, suite.handler, makeQuery(listDocumentsQuery, nil))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 0, result.Documents.TotalCount)

	result, err = RunQuery[ListDocumentsResponse](ctxUser1, suite.handler, makeQuery(listDocumentsQuery, nil))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 0, result.Documents.TotalCount)
}

func (suite *GraphTestSuite) TestTriggerExecution() {
	ctx := testutils.TestContext()
	db := env.RawDB(ctx)
	query := env.Query(ctx)
	documentTbl := query.Document

	var current_user string
	db.Raw(`select current_user`).Scan(&current_user)
	fmt.Println("CURRENT USER", current_user)

	var has_function_privilege bool
	db.Raw(`SELECT has_table_privilege(current_user, 'documents', 'TRIGGER')`).Scan(&has_function_privilege)
	fmt.Println("HAS FUNCTION PRIVILEGE", has_function_privilege)

	var sessionReplicationRole string
	db.Raw(`SHOW session_replication_role`).Scan(&sessionReplicationRole)
	fmt.Println("SESSION REPLICATION ROLE", sessionReplicationRole)
	var err error

	var triggerCount int64
	err = db.Raw(`SELECT COUNT(*)
        FROM pg_trigger
        WHERE tgrelid = 'documents'::regclass::oid
        AND tgname = 'set_root_parent_id_trigger'`).Scan(&triggerCount).Error
	require.NoError(suite.T(), err)
	fmt.Printf("TRIGGER COUNT: %v\n", triggerCount)

	var relHasTriggers bool
	db.Raw(`SELECT relhastriggers
          FROM pg_class
          WHERE relname = 'documents';`).Scan(&relHasTriggers)
	fmt.Println("REL HAS TRIGGERS", relHasTriggers)

	var triggerExists bool
	err = db.Raw(`
        SELECT EXISTS (
            SELECT 1
            FROM pg_trigger
            WHERE tgname = 'set_root_parent_id_trigger'
            AND tgrelid = 'documents'::regclass::oid
        )
    `).Scan(&triggerExists).Error
	require.NoError(suite.T(), err)
	require.True(suite.T(), triggerExists, "Trigger does not exist")

	// Create a document without setting root_parent_id
	doc := &models.Document{
		Title: "Test Trigger Execution",
		// Don't set RootParentID
	}

	err = documentTbl.Omit(documentTbl.RootParentID).Create(doc)
	require.NoError(suite.T(), err)

	// Fetch the document directly from the database
	var fetchedDoc models.Document
	err = db.Raw("SELECT * FROM documents WHERE id = ?", doc.ID).Scan(&fetchedDoc).Error
	require.NoError(suite.T(), err)

	// Check if root_parent_id was set by the trigger
	if fetchedDoc.RootParentID == "" {
		suite.T().Errorf("Trigger did not set root_parent_id. Fetched document: %+v", fetchedDoc)
	} else {
		suite.T().Logf("Trigger successfully set root_parent_id. Fetched document: %+v", fetchedDoc)
	}
}
