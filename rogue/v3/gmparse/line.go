package gmparse

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Line represents a single line within a paragraph or text block.
type Line struct {
	ast.BaseInline
}

// Dump implements Node.Dump.
func (n *Line) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindLine is a NodeKind for Line nodes.
var KindLine = ast.NewNodeKind("Line")

// Kind implements Node.Kind.
func (n *Line) Kind() ast.NodeKind {
	return KindLine
}

// NewLine returns a new Line node.
func NewLine() *Line {
	return &Line{
		BaseInline: ast.BaseInline{},
	}
}

// lineTransformer is an AST transformer that breaks paragraphs and text blocks into lines
type lineTransformer struct{}

func (t *lineTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n := n.(type) {
		case *ast.Paragraph:
			t.transform(n, reader)
		case *ast.TextBlock:
			t.transform(n, reader)
		}

		return ast.WalkContinue, nil
	})
}

func maxStop(n ast.Node) int {
	switch n := n.(type) {
	case *ast.Text:
		return n.Segment.Stop
	default:
		ms := 0
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			ms = max(ms, maxStop(c))
		}
		return ms
	}
}

func (t *lineTransformer) transform(para ast.Node, reader text.Reader) {
	var curLines []text.Segment
	segments := para.Lines()
	for i := 0; i < segments.Len(); i++ {
		line := segments.At(i)
		curLines = append(curLines, line)
	}

	if len(curLines) == 0 {
		return
	}

	paragraphLines := make([][]ast.Node, len(curLines))

	lineIx := 0
	for c := para.FirstChild(); c != nil; c = c.NextSibling() {
		if maxStop(c) > curLines[lineIx].Stop {
			lineIx = min(lineIx+1, len(curLines)-1)
		}

		paragraphLines[lineIx] = append(paragraphLines[lineIx], c)
	}

	para.RemoveChildren(nil)
	for _, line := range paragraphLines {
		paragraphLine := NewLine()
		for _, n := range line {
			paragraphLine.AppendChild(paragraphLine, n)
		}
		para.AppendChild(para, paragraphLine)
	}
}
