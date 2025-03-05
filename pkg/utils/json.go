package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

type state int

const (
	stateUnknown state = iota
	stateLookingForKey
	stateLookingForValue

	stateInKey

	stateInStringValue
	stateInNumberValue
	stateInUnknownValue
)

func (s state) String() string {
	return [...]string{
		"unknown",
		"lookingFor: Key",
		"lookingFor: Value",
		"in: Key",
		"in: StringValue",
		"in: NumberValue",
		"in: UnknownValue",
	}[s]
}

var digitRegexp = regexp.MustCompile(`^\d$`)

func PrettyPrintJSON(input interface{}) ([]byte, error) {
	bytes, err := json.MarshalIndent(input, "", "  ") // Indent with 2 spaces for pretty printing
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func ParseIncompleteJSON(input string) ([]string, map[string]string, error) {
	var result map[string]string

	err := json.Unmarshal([]byte(input), &result)
	if err == nil {
		keys := make([]string, 0, len(result))
		normMap := map[string]string{}
		for k, v := range result {
			nk := norm.NFC.String(k)
			keys = append(keys, nk)
			normMap[nk] = norm.NFC.String(v)
		}
		return keys, normMap, nil
	}

	completedKeys := []string{}
	result = make(map[string]string)
	escape := -1
	var key, value string
	var state = stateUnknown

	writeValue := func(keyComplete bool) {
		// log.Info("writeValue", "keyComplete", keyComplete, "key", key, "value", value, "result", result)
		if keyComplete {
			completedKeys = append(completedKeys, key)
		}
		result[norm.NFC.String(key)] = norm.NFC.String(strings.TrimRight(value, " "))
		key, value = "", ""
		state = stateLookingForKey
	}

	i := 0
	for i < len(input) {
		// log.Infof("char[%s]: %d %q", state, i, char)
		char, size := utf8.DecodeRuneInString(input[i:])

		if escape == i {
			v := string(char)
			switch char {
			case 'n':
				v = "\n"
			case 't':
				v = "\t"
			case 'r':
				v = "\r"
			case 'b':
				v = "\b"
			case 'f':
				v = "\f"
			case 'v':
				v = "\v"
			case 'u':
				if i+5 < len(input) {
					unicodeSeq := input[i+1 : i+5]
					if r, err := strconv.ParseUint(unicodeSeq, 16, 32); err == nil {
						if utf16.IsSurrogate(rune(r)) && i+11 < len(input) && input[i+5:i+7] == "\\u" {
							// This is a surrogate pair
							secondSeq := input[i+7 : i+11]
							if r2, err := strconv.ParseUint(secondSeq, 16, 32); err == nil {
								v = string(utf16.DecodeRune(rune(r), rune(r2)))
								i += 10 // Skip both \u sequences
							}
						} else {
							v = string(rune(r))
							i += 4
						}
					}
				}
			}

			switch state {
			case stateInKey:
				key = key + v
			case stateInStringValue:
				value = value + v
			}

			i += size
			continue
		}

		switch state {
		case stateInKey:
			switch char {
			case '{', '}', ':', ',':
				key = key + string(char)
			case '\\':
				escape = i + size
			case '"':
				state = stateLookingForValue
			case ' ':
				// Do nothing, ignore spaces in key
			default:
				key = key + string(char)
			}
		case stateInUnknownValue:
			switch char {
			case '"':
				if len(value) > 0 {
					return completedKeys, result, fmt.Errorf("unexpected '\"' while in value %q", value)
				}
				state = stateInStringValue
			case ' ':
				if len(value) > 0 {
					value = value + string(char)
				}
			case '{':
				return completedKeys, result, fmt.Errorf("unexpected '{' while in value %d", i)
			case '}', ':', ',':
				return completedKeys, result, fmt.Errorf("unexpected '%c' while in value %d", char, i)
			case '\\':
				escape = i + size
			default:
				value = value + string(char)
				if value == "null" {
					value = ""
					writeValue(true)
					state = stateLookingForKey
				}
			}
		case stateInStringValue:
			switch char {
			case '{', '}', ':', ',':
				value = value + string(char)
			case '\\':
				escape = i + size
			case '"':
				writeValue(true)
			case ' ':
				value = value + string(char)
			default:
				value = value + string(char)
			}
		case stateInNumberValue:
			switch char {
			case '{':
				return completedKeys, result, fmt.Errorf("unexpected '{' while in number value %d", i)
			case '}', ':', ',', '\\', '"':
				writeValue(true)
			case ' ':
				writeValue(true)
			default:
				if digitRegexp.MatchString(string(char)) || char == '.' {
					value = value + string(char)
				} else {
					writeValue(true)
				}
			}
		case stateLookingForKey:
			switch char {
			case '{':
				state = stateLookingForKey
			case '}':
				if key != "" {
					writeValue(true)
				}
			case '\\':
				escape = i + size
			case '"':
				state = stateInKey
			case ':':
				return completedKeys, result, fmt.Errorf("unexpected '%c' while looking for key %d", char, i)
			case ' ', ',':
				// Do nothing, ignore spaces and commas while looking for key
			default:
				i += size
				continue
			}
		case stateLookingForValue:
			switch char {
			case '{', '}', ',':
				return completedKeys, result, fmt.Errorf("unexpected '%c' while looking for value %d", char, i)
			case '\\':
				escape = i + size
			case '"':
				state = stateInStringValue
			case ' ', ':':
				// Do nothing, ignore spaces and colons while looking for value
			default:
				if digitRegexp.MatchString(string(char)) {
					state = stateInNumberValue
					value = value + string(char)
				} else {
					state = stateInUnknownValue
					value = value + string(char)
				}
			}
		default:
			switch char {
			case '{':
				state = stateLookingForKey
			case '}':
				if key != "" {
					writeValue(true)
				}
			case '\\':
				escape = i + 1
			case '"':
				state = stateInKey
			case ':', ',':
				return completedKeys, result, fmt.Errorf("unexpected '%c' while in unknown state %d", char, i)
			case ' ':
				// Do nothing, ignore spaces in unknown state
			default:
				i += size
				continue
			}
		}

		i += size
	}

	if key != "" && value != "" {
		writeValue(false)
	}

	return completedKeys, result, nil
}
