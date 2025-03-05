package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

type xmlState int

const (
	stateText xmlState = iota
	stateOpenTag
	stateCloseTag
	stateAttribute
)

func (s xmlState) String() string {
	return [...]string{
		"text",
		"openTag",
		"closeTag",
		"attribute",
	}[s]
}

type xmlAttrState int

func (s xmlAttrState) String() string {
	return [...]string{
		"name",
		"equal",
		"value",
	}[s]
}

const (
	attrStateName xmlAttrState = iota
	attrStateEqual
	attrStateValue
)

type XMLDocument struct {
	Tags []*Tag
}

func (doc *XMLDocument) String() string {
	return strings.Join(TagsToString(doc.Tags, ""), "\n")
}

func TagsToString(tags []*Tag, prefix string) []string {
	var result []string
	for _, tag := range tags {
		result = append(result, prefix+tag.String())
		result = append(result, TagsToString(tag.Children, prefix+"  ")...)
	}
	return result
}

// Find returns the first tag with the given key at the top level.
func (doc *XMLDocument) Find(key string) *Tag {
	for _, tag := range doc.Tags {
		if tag.Key == key {
			return tag
		}
	}
	return nil
}

// FindAll returns all tags with the given key at the top level.
func (doc *XMLDocument) FindAll(key string) []*Tag {
	var result []*Tag
	for _, tag := range doc.Tags {
		if tag.Key == key {
			result = append(result, tag)
		}
	}
	return result
}

// FindDeep returns the first tag with the given key at any level.
func (doc *XMLDocument) FindDeep(key string) *Tag {
	for _, tag := range doc.Tags {
		if found := tag.FindDeep(key); found != nil {
			return found
		}
	}
	return nil
}

// FindAllDeep returns all tags with the given key at any level.
func (doc *XMLDocument) FindAllDeep(key string) []*Tag {
	var result []*Tag
	for _, tag := range doc.Tags {
		result = append(result, tag.FindAllDeep(key)...)
	}
	return result
}

type Tag struct {
	Key        string
	Value      string
	RawValue   string
	Attributes map[string]string
	Complete   bool
	Children   []*Tag

	ContentStartIx int
	ContentEndIx   int
}

func (t *Tag) String() string {
	icon := "ⓧ"
	if t.Complete {
		icon = "✓"
	}
	trucated := ""
	if len(t.RawValue) > 100 {
		trucated = "..."
	}
	return fmt.Sprintf(
		"[%s %d-%d] %s: %q attrs: %+v",
		icon, t.ContentStartIx, t.ContentEndIx, t.Key, t.RawValue[:min(100, len(t.RawValue))]+trucated, t.Attributes,
	)
}

func (t *Tag) Attr(key string) string {
	if t.Attributes == nil {
		return ""
	}
	return t.Attributes[key]
}

func (t *Tag) Find(key string) *Tag {
	for _, child := range t.Children {
		if child.Key == key {
			return child
		}
	}
	return nil
}

func (t *Tag) FindDeep(key string) *Tag {
	if t.Key == key {
		return t
	}
	for _, child := range t.Children {
		if found := child.FindDeep(key); found != nil {
			return found
		}
	}
	return nil
}

func (t *Tag) FindAllDeep(key string) []*Tag {
	var result []*Tag
	if t.Key == key {
		result = append(result, t)
	}
	for _, child := range t.Children {
		result = append(result, child.FindAllDeep(key)...)
	}
	return result
}

// ParseIncompleteXML parses an XML string and returns a list of completed tags, their values with associated attributes and any errors.
func ParseIncompleteXML(input string) (*XMLDocument, error) {
	document := &XMLDocument{
		Tags: make([]*Tag, 0),
	}
	state := stateText
	var currentTag string
	var currentTagContext Tag
	tagStack := []*Tag{}

	var attrState xmlAttrState
	var attrName string
	var attrValue string
	var attrQuoteChar rune
	var collectingUnquotedValue bool

	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])

		var tagNames []string
		for _, tag := range tagStack {
			tagNames = append(tagNames, tag.Key)
		}

		// fmt.Printf("%d: (%s) %s [%s]\n", i, state.String(), string(r), strings.Join(tagNames, ", "))

		switch state {
		case stateText:
			if r == '<' {
				if i+1 < len(input) && input[i+1] == '/' {
					if len(tagStack) > 0 {
						// When we find a closing tag, get raw content from stored start position to current
						current := tagStack[len(tagStack)-1]
						current.RawValue = input[current.ContentStartIx:i]
						current.ContentEndIx = i // Set the content end index before the closing tag
					}
					state = stateCloseTag
					i++ // Skip the '/'
					currentTag = ""
					currentTagContext = Tag{}
				} else {
					state = stateOpenTag
					currentTag = ""
					currentTagContext = Tag{}
				}
			} else {
				if len(tagStack) > 0 {
					tagStack[len(tagStack)-1].Value += string(r)
				}
			}
		case stateOpenTag:
			if r == '>' {
				state = stateText
				currentTagContext.Key = currentTag
				currentTagContext.Value = ""
				currentTagContext.ContentStartIx = i + size // Store start position after '>'
				currentTagContext.ContentEndIx = i + size   // Initialize ContentEndIx, will be updated when closing tag is found
				newTag := currentTagContext
				tagStack = append(tagStack, &newTag)
			} else if r == '/' && i+1 < len(input) && input[i+1] == '>' {
				currentTagContext.Key = norm.NFC.String(currentTag)
				currentTagContext.Complete = true
				currentTagContext.RawValue = ""
				currentTagContext.ContentStartIx = i + size // Set content start for self-closing tag
				currentTagContext.ContentEndIx = i + size   // Set content end for self-closing tag (same as start for empty tags)
				state = stateText
				i++ // Skip the next '>'
			} else if IsSpace(r) {
				// Start parsing attributes
				state = stateAttribute
				attrState = attrStateName
				attrName = ""
				attrValue = ""
				attrQuoteChar = 0
				collectingUnquotedValue = false
				// Initialize currentTagContext
				currentTagContext.Key = currentTag
				currentTagContext.Value = ""
			} else {
				currentTag += string(r)
			}
		case stateCloseTag:
			if r == '>' {
				// Close current tag
				state = stateText
				if len(tagStack) > 0 {
					var (
						closedTag    *Tag
						orphanedTags []*Tag
					)

					// Find the closed tag (there might be tags that are not closed in the middle of the XML)
					for j := len(tagStack) - 1; j >= 0; j-- {
						if tagStack[j].Key == currentTag {
							closedTag = tagStack[j]
							tagStack = tagStack[:j]
							break
						}
						orphanedTags = append(orphanedTags, tagStack[j])
					}

					if len(orphanedTags) > 0 {
						// If we have orphaned tags, we need to update the value of the closed tag since the orphaned tags
						// no longer own the values
						minOrphanedStartIx := min(closedTag.ContentStartIx, orphanedTags[0].ContentStartIx)
						maxOrphanedEndIx := orphanedTags[len(orphanedTags)-1].ContentEndIx
						closedTag.Value = norm.NFC.String(strings.TrimSpace(input[minOrphanedStartIx:maxOrphanedEndIx]))
						closedTag.RawValue = input[minOrphanedStartIx:maxOrphanedEndIx]
						closedTag.ContentEndIx = maxOrphanedEndIx
					} else {
						closedTag.Value = norm.NFC.String(strings.TrimSpace(closedTag.Value))
					}

					// We've found unclosed tags between the opening tag and closing tag. The values of these tags belong to the closing tag
					// so we set the orphaned tags as children of the closed tag and set their values to ""
					// Note we keep the ContentStartIx of the orphaned tags (this can be changed later if it's an issue)
					for _, orphanedTag := range orphanedTags {
						orphanedTag.RawValue = ""
						orphanedTag.Value = ""
						closedTag.Children = append(closedTag.Children, orphanedTag)
					}

					closedTag.Key = norm.NFC.String(closedTag.Key)
					closedTag.Complete = true

					if len(tagStack) > 0 {
						// Update parent's raw value to include this entire tag
						parent := tagStack[len(tagStack)-1]
						parent.RawValue = input[parent.ContentStartIx : i+size]
						parent.ContentEndIx = i + size // Update parent's content end index
						// Add nested tag
						parent.Children = append(parent.Children, closedTag)
					} else {
						document.Tags = append(document.Tags, closedTag)
					}
				}
				currentTag = ""
				currentTagContext = Tag{}
			} else {
				currentTag += string(r)
			}
		case stateAttribute:
			switch attrState {
			case attrStateName:
				if IsSpace(r) {
					if attrName != "" {
						// We've read the attribute name, waiting for '='
						attrState = attrStateEqual
					}
					// Else skip extra spaces
				} else if r == '=' {
					attrState = attrStateValue
				} else if r == '>' {
					// End of tag
					state = stateText
					currentTagContext.ContentStartIx = i + size
					currentTagContext.ContentEndIx = i + size // Initialize ContentEndIx
					newTag := currentTagContext
					tagStack = append(tagStack, &newTag)
				} else if r == '/' && i+1 < len(input) && input[i+1] == '>' {
					// Self-closing tag
					currentTagContext.Key = norm.NFC.String(currentTag)
					currentTagContext.Complete = true
					currentTagContext.RawValue = ""
					currentTagContext.ContentStartIx = i + size
					currentTagContext.ContentEndIx = i + size // Set ContentEndIx for self-closing tag
					state = stateText

					selfClosingTag := currentTagContext
					currentTagContext = Tag{}

					if len(tagStack) > 0 {
						// Add nested tag
						tagStack[len(tagStack)-1].Children = append(tagStack[len(tagStack)-1].Children, &selfClosingTag)
					} else {
						document.Tags = append(document.Tags, &selfClosingTag)
					}

					i++ // Skip the next '>'
				} else {
					attrName += string(r)
				}
			case attrStateEqual:
				if IsSpace(r) {
					// Skip whitespace
				} else if r == '=' {
					attrState = attrStateValue
				} else {
					// In XML, '=' is required after attribute name
					attrState = attrStateValue
					i -= size // Re-process this character in attrStateValue
					continue
				}
			case attrStateValue:
				if attrQuoteChar != 0 {
					// Collecting quoted value
					if r == attrQuoteChar {
						// End of quoted attribute value
						if currentTagContext.Attributes == nil {
							currentTagContext.Attributes = make(map[string]string)
						}
						currentTagContext.Attributes[strings.TrimSpace(attrName)] = attrValue

						// Reset attribute parsing variables
						attrName = ""
						attrValue = ""
						attrState = attrStateName
						attrQuoteChar = 0
						collectingUnquotedValue = false
					} else {
						attrValue += string(r)
					}
				} else if collectingUnquotedValue {
					// Collecting unquoted value
					if IsSpace(r) || r == '>' || (r == '/' && i+1 < len(input) && input[i+1] == '>') {
						// End of unquoted attribute value
						if currentTagContext.Attributes == nil {
							currentTagContext.Attributes = make(map[string]string)
						}
						currentTagContext.Attributes[strings.TrimSpace(attrName)] = attrValue

						// Reset attribute parsing variables
						attrName = ""
						attrValue = ""
						attrState = attrStateName
						collectingUnquotedValue = false

						// Handle '>' or '/>'
						if r == '>' {
							state = stateText
							currentTagContext.ContentStartIx = i + size
							currentTagContext.ContentEndIx = i + size // Initialize ContentEndIx
							newTag := currentTagContext
							tagStack = append(tagStack, &newTag)
							i++ // Skip the '>'
						} else if r == '/' && i+1 < len(input) && input[i+1] == '>' {
							// Self-closing tag
							currentTagContext.Key = norm.NFC.String(currentTag)
							currentTagContext.Complete = true
							currentTagContext.RawValue = ""
							currentTagContext.ContentStartIx = i + size
							currentTagContext.ContentEndIx = i + size // Set ContentEndIx for self-closing tag
							state = stateText
							i++ // Skip the next '>'
						}
						continue // Important to skip the increment of i in this iteration
					} else {
						attrValue += string(r)
					}
				} else {
					// Not yet started collecting value
					if IsSpace(r) {
						// Skip leading whitespace
					} else if r == '"' || r == '\'' {
						// Start collecting quoted value
						attrQuoteChar = r
					} else {
						// Start collecting unquoted value
						attrValue += string(r)
						collectingUnquotedValue = true
					}
				}
			}
		}

		i += size
	}

	// Process any remaining open tags
	if len(tagStack) > 0 {
		// Set the raw value and content end index for the final tag
		current := tagStack[len(tagStack)-1]
		current.RawValue = input[current.ContentStartIx:len(input)]
		current.ContentEndIx = len(input) // Set ContentEndIx to end of input for unclosed tags
	}

	for len(tagStack) > 0 {
		openTag := tagStack[len(tagStack)-1]
		tagStack = tagStack[:len(tagStack)-1]
		openTag.Key = norm.NFC.String(openTag.Key)
		openTag.Value = norm.NFC.String(strings.TrimSpace(openTag.Value))
		if len(tagStack) > 0 {
			// Add nested tag
			tagStack[len(tagStack)-1].Children = append(tagStack[len(tagStack)-1].Children, openTag)
		} else {
			document.Tags = append(document.Tags, openTag)
		}
	}

	return document, nil
}

func IsSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
