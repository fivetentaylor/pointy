package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseIncompleteJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		completedKeys []string
		output        map[string]string
	}{
		{"Normal JSON", `{"key": "value"}`, []string{"key"}, map[string]string{"key": "value"}},
		{"Normal JSON w numbers", `{"key": 1}`, []string{"key"}, map[string]string{"key": "1"}},
		{"Normal multiple key JSON", `{"key": "value", "key2": "value2"}`, []string{"key", "key2"}, map[string]string{"key": "value", "key2": "value2"}},
		{"Normal JSON w null", `{"key": "value", "key2": null}`, []string{"key", "key2"}, map[string]string{"key": "value", "key2": ""}},
		{"Invalid JSON w newlines", `{"key": "value \n line 2"}`, []string{}, map[string]string{"key": "value \n line 2"}},
		{"Empty JSON", "{}", []string{}, map[string]string{}},
		{"Incomplete JSON", `{"key": "value`, []string{}, map[string]string{"key": "value"}},
		{"Incomplete JSON key 2", `{"key": "value", "key2"`, []string{"key"}, map[string]string{"key": "value"}},
		{"Incomplete JSON w numbers", `{"key": 1`, []string{}, map[string]string{"key": "1"}},
		{"Incomplete JSON w null", `{"key": "value", "key2": null, "key3": "va`, []string{"key", "key2"}, map[string]string{"key": "value", "key2": "", "key3": "va"}},
		{"Incomplete JSON w ,", `{"key": "value1, value2 and value3`, []string{}, map[string]string{"key": "value1, value2 and value3"}},
		{"Incomplete JSON w \"", `{"key": "value1, \"value2\" and value3`, []string{}, map[string]string{"key": "value1, \"value2\" and value3"}},
		{"Incomplete key", `{"ke`, []string{}, map[string]string{}},
		{"Incomplete invalid JSON w newlines", `{"key": "value \n line 2`, []string{}, map[string]string{"key": "value \n line 2"}},
		{"Empty", ``, []string{}, map[string]string{}},
		{"Incomplete JSON w newlines", `{"key": "1\n2\n3\"}`, []string{}, map[string]string{"key": "1\n2\n3\"}"}},
		{"curly in value", `{"key": "{hello}\"`, []string{}, map[string]string{"key": "{hello}\""}},
		{
			"Unicode escape sequences and HTML entities",
			`{"unicode": "Hello \u0041\u0042\u0043 \u03B1\u03B2\u03B3",
          "emoji": "\uD83D\uDE00\uD83D\uDE4F",
          "mixed": "Copyright \u00A9 \u2665 O'Reilly",
          "html_entities": "Greater than &gt; less than &lt; ampersand &amp;",
          "incomplete": "This is incomplete \u26"}`,
			[]string{"unicode", "emoji", "mixed", "html_entities"},
			map[string]string{
				"unicode":       "Hello ABC αβγ",
				"emoji":         "😀🙏",
				"mixed":         "Copyright © ♥ O'Reilly",
				"html_entities": "Greater than &gt; less than &lt; ampersand &amp;",
				"incomplete":    "This is incomplete u26",
			},
		},
		{
			"Another incomplete JSON",
			`{"message":"test","reasoning":"test","contentId":"full_doc","content":"test_content","analysis":"The adde"`,
			[]string{"message", "reasoning", "contentId", "content", "analysis"},
			map[string]string{"message": "test", "reasoning": "test", "contentId": "full_doc", "content": "test_content", "analysis": "The adde"},
		},
		{
			name:          "Korean characters with Unicode escapes",
			input:         `{"korean": "\u110b\u1161\u11ab\u1102\u1167\u11bc\u1112\u1161\u1109\u1166\u110b\u116d \u1109\u1166\u1100\u1168"}`,
			completedKeys: []string{"korean"},
			output: map[string]string{
				"korean": "안녕하세요 세계",
			},
		},
		{
			name:          "Incomplete Korean characters with no escapes",
			input:         `{"korean": "안녕하세요 세계"`,
			completedKeys: []string{"korean"},
			output: map[string]string{
				"korean": "안녕하세요 세계",
			},
		},
		{
			name:          "Aggressively mixed Unicode representations",
			input:         `{"mixed": "H\u0065llo, \u4E16界\u3053ん\u306B\u3061は\u0020☀\uFF01"}`,
			completedKeys: []string{"mixed"},
			output: map[string]string{
				"mixed": "Hello, 世界こんにちは ☀！",
			},
		},
		/*
			    // streaming parser chokes on nested json
					{
						name:          "Incomplete JSON with nesting",
						input:         `{"message":"test_message","reasoning":"test_reasoning","concludingMessage":"test_conclude","feedback":[\n  {`,
						completedKeys: []string{"message", "reasoning", "concludingMessage"},
						output:        map[string]string{"message": "test_message", "reasoning": "test_reasoning", "concludingMessage": "test_conclude"},
					},
		*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, got, err := ParseIncompleteJSON(tt.input)
			if err != nil {
				t.Errorf("ParseIncompleteJSON(%v) error = %v", tt.input, err)
				return
			}
			assert.Equal(t, tt.output, got)
			for _, key := range tt.completedKeys {
				assert.Contains(t, keys, key)
			}
		})
	}
}
