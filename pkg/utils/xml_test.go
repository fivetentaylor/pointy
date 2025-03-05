package utils_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/teamreviso/code/pkg/utils"
)

func Test_ParseIncompleteXML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []*utils.Tag
	}{
		{"Normal XML",
			`<key>value</key>`,
			[]*utils.Tag{
				{
					Key:            "key",
					Value:          "value",
					RawValue:       "value",
					Complete:       true,
					ContentStartIx: 5,
					ContentEndIx:   10,
				},
			},
		},
		{"XML missing closing tag",
			`<key>value`,
			[]*utils.Tag{
				{
					Key:            "key",
					Value:          "value",
					RawValue:       "value",
					ContentStartIx: 5,
					ContentEndIx:   10,
				},
			},
		},
		{"nested XML",
			`<message><name>joe bob</name></message>`,
			[]*utils.Tag{
				{
					Key:            "message",
					Value:          "",
					RawValue:       "<name>joe bob</name>",
					Complete:       true,
					ContentStartIx: 9,
					ContentEndIx:   29,
					Children: []*utils.Tag{
						{
							Key:            "name",
							Value:          "joe bob",
							RawValue:       "joe bob",
							Complete:       true,
							ContentStartIx: 15,
							ContentEndIx:   22,
						},
					},
				},
			},
		},
		{"nested XML with text",
			`<message>Hello <name>joe bob</name> how are you?</message>`,
			[]*utils.Tag{
				{
					Key:            "message",
					Value:          "Hello  how are you?",
					RawValue:       "Hello <name>joe bob</name> how are you?",
					Complete:       true,
					ContentStartIx: 9,
					ContentEndIx:   48,
					Children: []*utils.Tag{
						{
							Key:            "name",
							Value:          "joe bob",
							RawValue:       "joe bob",
							Complete:       true,
							ContentStartIx: 21,
							ContentEndIx:   28,
						},
					},
				},
			},
		},
		{"no parent node",
			`<message>hello</message><name>joe bob</name>`,
			[]*utils.Tag{
				{
					Key:            "message",
					Value:          "hello",
					RawValue:       "hello",
					Complete:       true,
					ContentStartIx: 9,
					ContentEndIx:   14,
				},
				{
					Key:            "name",
					Value:          "joe bob",
					RawValue:       "joe bob",
					Complete:       true,
					ContentStartIx: 30,
					ContentEndIx:   37,
				},
			},
		},
		{"no parent node (partial)",
			`<message>hello</message><name>joe bo`,
			[]*utils.Tag{
				{
					Key:            "message",
					Value:          "hello",
					RawValue:       "hello",
					Complete:       true,
					ContentStartIx: 9,
					ContentEndIx:   14,
				},
				{
					Key:            "name",
					Value:          "joe bo",
					RawValue:       "joe bo",
					Complete:       false,
					ContentStartIx: 30,
					ContentEndIx:   36,
				},
			},
		},
		{"node with attributes",
			`<message id="taco" ref="beans">hello</message>`,
			[]*utils.Tag{
				{
					Key:      "message",
					Value:    "hello",
					RawValue: "hello",
					Complete: true,
					Attributes: map[string]string{
						"id":  "taco",
						"ref": "beans",
					},
					ContentStartIx: 31,
					ContentEndIx:   36,
				},
			},
		},
		{"node with attributes partial",
			`<message id="taco">hello`,
			[]*utils.Tag{
				{
					Key:      "message",
					Value:    "hello",
					RawValue: "hello",
					Complete: false,
					Attributes: map[string]string{
						"id": "taco",
					},
					ContentStartIx: 19,
					ContentEndIx:   24,
				},
			},
		},
		{"node with attributes without quotes",
			`<message id=taco> hello</message>`,
			[]*utils.Tag{
				{
					Key:      "message",
					Value:    "hello",
					RawValue: " hello",
					Attributes: map[string]string{
						"id": "taco",
					},
					Complete:       true,
					ContentStartIx: 17,
					ContentEndIx:   23,
				},
			},
		},
		{"node with partial attribute",
			`<message id="tac`,
			[]*utils.Tag{},
		},
		{"multiple nodes with the same name with attributes",
			`<taco flavor=beef>with cheese</taco><taco flavor=chicken>with sour cream</taco>`,
			[]*utils.Tag{
				{
					Key:      "taco",
					Value:    "with cheese",
					RawValue: "with cheese",
					Attributes: map[string]string{
						"flavor": "beef",
					},
					Complete:       true,
					ContentStartIx: 18,
					ContentEndIx:   29,
				},
				{
					Key:      "taco",
					Value:    "with sour cream",
					RawValue: "with sour cream",
					Attributes: map[string]string{
						"flavor": "chicken",
					},
					Complete:       true,
					ContentStartIx: 57,
					ContentEndIx:   72,
				},
			},
		},
		{"multiple nodes with the same name, nested",
			`<message><name>joe bob</name></message><message><name>jim bob</name></message>`,
			[]*utils.Tag{
				{
					Key:      "message",
					Value:    "",
					RawValue: "<name>joe bob</name>",
					Complete: true,
					Children: []*utils.Tag{
						{
							Key:            "name",
							Value:          "joe bob",
							RawValue:       "joe bob",
							Complete:       true,
							ContentStartIx: 15,
							ContentEndIx:   22,
						},
					},
					ContentStartIx: 9,
					ContentEndIx:   29,
				},
				{
					Key:      "message",
					Value:    "",
					RawValue: "<name>jim bob</name>",
					Complete: true,
					Children: []*utils.Tag{
						{
							Key:            "name",
							Value:          "jim bob",
							RawValue:       "jim bob",
							Complete:       true,
							ContentStartIx: 54,
							ContentEndIx:   61,
						},
					},
					ContentStartIx: 48,
					ContentEndIx:   68,
				},
			},
		},
		{"self closing XML",
			`<message><name value="bob" /></message>`,
			[]*utils.Tag{
				{
					Key:      "message",
					Value:    "",
					RawValue: "<name value=\"bob\" />",
					Complete: true,
					Children: []*utils.Tag{
						{
							Key:      "name",
							Complete: true,
							Attributes: map[string]string{
								"value": "bob",
							},
							ContentStartIx: 28,
							ContentEndIx:   28,
						},
					},
					ContentStartIx: 9,
					ContentEndIx:   29,
				},
			},
		},
		{"incomplete tag in the middle of the XML",
			`<message>Hello <name value="bob" > how are you?</message>`,
			[]*utils.Tag{
				{
					Key:      "message",
					Value:    "Hello <name value=\"bob\" > how are you?",
					RawValue: "Hello <name value=\"bob\" > how are you?",
					Complete: true,
					Children: []*utils.Tag{
						{
							Key:      "name",
							Complete: false,
							Attributes: map[string]string{
								"value": "bob",
							},
							ContentStartIx: 34,
							ContentEndIx:   47,
						},
					},
					ContentStartIx: 9,
					ContentEndIx:   47,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ParseIncompleteXML(tt.input)
			if err != nil {
				t.Errorf("ParseIncompleteXML(%q) error = %v", tt.input, err)
				return
			}

			if !reflect.DeepEqual(got.Tags, tt.want) {
				gotTags := []utils.Tag{}
				for _, tag := range got.Tags {
					gotTags = append(gotTags, *tag)
				}
				wantTags := []utils.Tag{}
				for _, tag := range tt.want {
					wantTags = append(wantTags, *tag)
				}

				gotStr := got.String()
				wantStr := strings.Join(utils.TagsToString(tt.want, ""), "\n")

				t.Errorf("ParseIncompleteXML(%q)\n\tgot:%#v\n\nPrinted:\n%s\n\n\twant:%#v\n\nPrinted:\n%s", tt.input, gotTags, gotStr, wantTags, wantStr)
			}
		})
	}
}

func Test_ParseIncompleteXML_Realworld(t *testing.T) {
	text := `<thinking>The user's request is to make the first and last paragraphs of the document shorter. Therefore, I need to identify the first and last paragraphs in the document provided.

The document begins with the paragraph about the Battle of Glasgow being fought on October 15, 1864, and ends with the paragraph detailing the Confederate retreat reaching Texas. These paragraphs are the targets for the user's requested changes.

I will identify these two paragraphs as the relevant sections for the user's request.

<relevant_sections>
1. The first paragraph:
"The Battle of Glasgow was fought on October 15, 1864 in and near Glasgow, Missouri as part of Price's Missouri Expedition during the American Civil War. The battle resulted in the capture of needed weapons and improved Confederate morale, which had been dented after a defeat in the Battle of Pilot Knob."

2. The last paragraph:
"On October 23, the Confederates were decisively defeated at the Battle of Westport. Price's men retreated, but were harried for much of the way by Union pursuit. The retreat eventually reached Texas."

These two sections will be identified for updating based on the user's request.

I will proceed to include these relevant sections in the output.</thinking>

<initial_response>
<relevant_sections>
<relevant_section>
The Battle of Glasgow was fought on October 15, 1864 in and near Glasgow, Missouri as part of Price's Missouri Expedition during the American Civil War. The battle resulted in the capture of needed weapons and improved Confederate morale, which had been dented after a defeat in the Battle of Pilot Knob.
</relevant_section>
<relevant_section>
On October 23, the Confederates were decisively defeated at the Battle of Westport. Price's men retreated, but were harried for much of the way by Union pursuit. The retreat eventually reached Texas.
</relevant_section>
</relevant_sections>
</initial_response>

<reflect>I have identified the correct sections based on the user's request to make the first and last paragraphs shorter. The first and last paragraphs of the document have been correctly selected as the relevant sections for the user's requested changes. My response aligns with the instructions and examples provided. I am confident in my selection.</reflect>

<response>
<relevant_sections>
<relevant_section>
The Battle of Glasgow was fought on October 15, 1864 in and near Glasgow, Missouri as part of Price's Missouri Expedition during the American Civil War. The battle resulted in the capture of needed weapons and improved Confederate morale, which had been dented after a defeat in the Battle of Pilot Knob.
</relevant_section>
<relevant_section>
On October 23, the Confederates were decisively defeated at the Battle of Westport. Price's men retreated, but were harried for much of the way by Union pursuit. The retreat eventually reached Texas.
</relevant_section>
</relevant_sections>
</response>`

	doc, err := utils.ParseIncompleteXML(text)
	if err != nil {
		t.Error(err)
	}

	// Notice there is a <relevant_section> tag in the thinking tag, that is unclosed
	// normally we could use Find("response"), but in this case we need to use FindDeep to get the response out of the nested tags
	resp := doc.FindDeep("response")
	fmt.Println(doc.String())
	if resp == nil {
		t.Errorf("response not found in\n\t%+v", doc)
	}
}

func TestFindAll(t *testing.T) {
	tests := []struct {
		Name     string
		InputXML string
		Key      string
		Want     []*utils.Tag
	}{
		{
			Name: "xml with html",
			InputXML: `<updated_subsection id="84aee0fd-6b86-420b-abf0-43d0320b5066" action="replace">
<figure><img src="https://app.reviso.dev:9090/drafts/76399538-0aa4-4a30-9d96-449681595b9d/images/5c762946-5401-4ebb-bd49-5091c0ff7c2c" /><figcaption>My kitty</figcaption></figure>

## Cat Virtues

- Independence: Masters of self-reliance
- Grace: Elegant in movement and demeanor
- Cleanliness: Meticulous self-groomers
- Curiosity: Explorers of the unknown
- Patience: Skilled hunters and observers
- Affection: Loving companions on their own terms
- Playfulness: Eternally youthful spirits
- Adaptability: Comfortable in various environments
- Intuition: Keen sense of their surroundings
- Serenity: Zen-like ability to relax and unwind
</updated_subsection>

<conclusion>

I've added a list of cat virtues below the image as requested. The list includes ten virtues that are commonly associated with cats, highlighting their unique qualities and characteristics. This addition complements the existing content about felines and adds depth to the document. Let me know if you'd like any further adjustments or additions!
</conclusion>`,
			Key: "updated_subsection",
			Want: []*utils.Tag{
				{
					Key:      "updated_subsection",
					Value:    "## Cat Virtues\n\n- Independence: Masters of self-reliance\n- Grace: Elegant in movement and demeanor\n- Cleanliness: Meticulous self-groomers\n- Curiosity: Explorers of the unknown\n- Patience: Skilled hunters and observers\n- Affection: Loving companions on their own terms\n- Playfulness: Eternally youthful spirits\n- Adaptability: Comfortable in various environments\n- Intuition: Keen sense of their surroundings\n- Serenity: Zen-like ability to relax and unwind",
					RawValue: "\n<figure><img src=\"https://app.reviso.dev:9090/drafts/76399538-0aa4-4a30-9d96-449681595b9d/images/5c762946-5401-4ebb-bd49-5091c0ff7c2c\" /><figcaption>My kitty</figcaption></figure>\n\n## Cat Virtues\n\n- Independence: Masters of self-reliance\n- Grace: Elegant in movement and demeanor\n- Cleanliness: Meticulous self-groomers\n- Curiosity: Explorers of the unknown\n- Patience: Skilled hunters and observers\n- Affection: Loving companions on their own terms\n- Playfulness: Eternally youthful spirits\n- Adaptability: Comfortable in various environments\n- Intuition: Keen sense of their surroundings\n- Serenity: Zen-like ability to relax and unwind\n",
					Attributes: map[string]string{
						"action": "replace",
						"id":     "84aee0fd-6b86-420b-abf0-43d0320b5066",
					},
					Complete: true,
					Children: []*utils.Tag{
						{
							Key:        "figure",
							Value:      "",
							RawValue:   "<img src=\"https://app.reviso.dev:9090/drafts/76399538-0aa4-4a30-9d96-449681595b9d/images/5c762946-5401-4ebb-bd49-5091c0ff7c2c\" /><figcaption>My kitty</figcaption>",
							Attributes: nil,
							Complete:   true,
							Children: []*utils.Tag{
								{
									Key:      "img",
									Value:    "",
									RawValue: "",
									Attributes: map[string]string{
										"src": "https://app.reviso.dev:9090/drafts/76399538-0aa4-4a30-9d96-449681595b9d/images/5c762946-5401-4ebb-bd49-5091c0ff7c2c",
									},
									Complete:       true,
									Children:       nil,
									ContentStartIx: 216,
									ContentEndIx:   216,
								},
								{
									Key:            "figcaption",
									Value:          "My kitty",
									RawValue:       "My kitty",
									Attributes:     nil,
									Complete:       true,
									Children:       nil,
									ContentStartIx: 229,
									ContentEndIx:   237,
								},
							},
							ContentStartIx: 88,
							ContentEndIx:   250,
						},
					},
					ContentStartIx: 79,
					ContentEndIx:   719,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			doc, err := utils.ParseIncompleteXML(tt.InputXML)
			if err != nil {
				t.Error(err)
			}

			got := doc.FindAll(tt.Key)
			require.Equal(t, tt.Want, got)
		})
	}
}
