package v3

// This is going to be the new way of rendering html from rogue.
// The current approach is a little too brittle and hard to change.
// The flow now will be rogue -> nos -> ast -> html
// Currently it is rogue -> nos -> html

type AstNodeType int

const (
	AstNodeSpan AstNodeType = iota
	AstNodeLine
	AstNodeBlock
	AstNodeRoot
)

type AstNode struct {
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *AstNode
	Type                                                    AstNodeType
	Text                                                    *string
	Format                                                  *FormatV3
	StartID, EndID                                          *ID
}

/*func ToAst(vis *FugueVis, spanNOS, lineNOS *NOS, includeIDs bool) (*AstNode, error) {
	root := &AstNode{Type: AstNodeRoot}
	lines := lineNOS.tree.AsSlice()

	for _, line := range lines {
		// write the opening tag of the block if there is one
		var curLineFmt FormatV3 = line.Format

		var startID *ID = nil
		var lineID *ID = nil
		if includeIDs {
			startID = &vis.IDs[line.StartIx]
			if line.EndIx < len(vis.IDs) {
				lineID = &vis.IDs[line.EndIx]
			}
		}

	}

	// close the last blocks
}*/
