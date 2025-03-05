package v3

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"runtime"
	"sort"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
)

func randomRune(rand *rand.Rand) rune {
	// Define some ranges of valid Unicode characters, including emojis
	ranges := [][]int{
		{0x0020, 0x007F},   // Basic Latin
		{0x00A0, 0x00FF},   // Latin-1 Supplement
		{0x0400, 0x04FF},   // Cyrillic
		{0x1F300, 0x1F5FF}, // Miscellaneous Symbols and Pictographs
		{0x1F600, 0x1F64F}, // Emoticons
		{0x1F680, 0x1F6FF}, // Transport and Map Symbols
		{0x1F900, 0x1F9FF}, // Supplemental Symbols and Pictographs
		// More ranges can be added as needed
	}

	// Choose a random range
	chosenRange := ranges[rand.Intn(len(ranges))]

	// Generate a random rune within the chosen range
	return rune(rand.Intn(chosenRange[1]-chosenRange[0]+1) + chosenRange[0])
}

func RandomString(rand *rand.Rand, length int) string {
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteRune(randomRune(rand))
	}
	return result.String()
}

func WithAuthor(r *Rogue, author string, f func(r *Rogue)) *Rogue {
	prevAuthor := r.Author
	r.Author = author
	f(r)
	r.Author = prevAuthor
	return r
}

func ToJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(b)
}

func bisect[T any, K constraints.Ordered](array []T, key K, keyExtractor func(T) (K, error)) (int, error) {
	low, high := 0, len(array)

	for low < high {
		mid := low + (high-low)/2
		midValue, err := keyExtractor(array[mid])
		if err != nil {
			return 0, err
		}

		if midValue <= key {
			low = mid + 1
		} else {
			high = mid
		}
	}

	return low, nil
}

func bisectLeft[T any, K constraints.Ordered](array []T, key K, keyExtractor func(T) K) int {
	low, high := 0, len(array)

	for low < high {
		mid := low + (high-low)/2
		midValue := keyExtractor(array[mid])

		if midValue < key {
			low = mid + 1
		} else {
			high = mid
		}
	}

	return low
}

func sortedInsert[T any](slice []T, element T, comp func(T, T) int) (int, []T) {
	index := findInsertIndex(slice, element, comp)

	// Ensure slice has enough capacity
	slice = append(slice, element) // Append at the end to increase the size

	if index < len(slice)-1 { // Check if not appending at the end
		// Shift elements to the right to make room
		copy(slice[index+1:], slice[index:len(slice)-1])
		slice[index] = element
	}

	return index, slice
}

func findInsertIndex[T any](slice []T, element T, comp func(T, T) int) int {
	low, high := 0, len(slice)
	for low < high {
		mid := low + (high-low)/2
		if comp(element, slice[mid]) < 0 {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func SafeSlice[T any](slice []T, start, end int) []T {
	if start < 0 {
		start = 0
	}
	if end > len(slice) {
		end = len(slice)
	}
	if start > end {
		start = end
	}

	// Create a new slice with the appropriate length
	newSlice := make([]T, end-start)
	// Copy the elements from the original slice to the new slice
	copy(newSlice, slice[start:end])

	return newSlice
}

// max returns the maximum of two values.
func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two values.
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func PrintJson(prefix string, i interface{}) {
	js, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %s\n", prefix, string(js))
}

func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func InsertAt[T any](slice []T, index int, element T) []T {
	if index < 0 || index > len(slice) {
		return slice
	}

	if index == len(slice) {
		return append(slice, element)
	}

	newSlice := make([]T, 0, len(slice)+1)
	newSlice = append(newSlice, slice[:index]...)
	newSlice = append(newSlice, element)
	newSlice = append(newSlice, slice[index:]...)

	return newSlice
}

func InsertSliceAt[T any](slice []T, index int, elements []T) []T {
	if index < 0 || index > len(slice) {
		return slice
	}

	newSlice := make([]T, 0, len(slice)+len(elements))
	newSlice = append(newSlice, slice[:index]...)
	newSlice = append(newSlice, elements...)
	newSlice = append(newSlice, slice[index:]...)

	return newSlice
}

/*func DeleteAt[T any](slice []T, index, length int) []T {
	startIx := min(max(index, 0), len(slice))
	endIx := min(max(index+length, 0), len(slice))
	return append(slice[:startIx], slice[endIx:]...)
}*/

func DeleteAt[T any](slice []T, index, length int) []T {
	startIx := min(max(index, 0), len(slice))
	endIx := min(max(index+length, 0), len(slice))

	if startIx >= len(slice) || endIx <= startIx {
		return slice
	}

	numToMove := len(slice) - endIx

	if numToMove > 0 {
		copy(slice[startIx:], slice[endIx:])
	}

	return slice[:len(slice)-(endIx-startIx)]
}

func ShallowCompareMaps(map1, map2 map[string]interface{}) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, valueMap1 := range map1 {
		valueMap2, ok := map2[key]
		if !ok {
			// Key does not exist in map2
			return false
		}

		// Direct comparison for basic types; you might need specific handling for non-comparable types
		if valueMap1 != valueMap2 {
			return false
		}
	}
	return true
}

func IsSurrogate(u uint16) bool {
	return IsHighSurrogate(u) || IsLowSurrogate(u)
}

// Left or first half of a surrogate pair
func IsHighSurrogate(u uint16) bool {
	return u >= 0xD800 && u <= 0xDBFF
}

// Right or second half of a surrogate pair
func IsLowSurrogate(u uint16) bool {
	return u >= 0xDC00 && u <= 0xDFFF
}

func StrToUint16(s string) []uint16 {
	r := []rune(s)
	return utf16.Encode(r)
}

// uint16SliceToHexString converts a slice of uint16 to a string of hex values prefixed with 0x
func Uint16SliceToHexString(slice []uint16) string {
	var hexStrings []string
	for _, v := range slice {
		hexStrings = append(hexStrings, fmt.Sprintf("0x%X", v))
	}
	return "[" + strings.Join(hexStrings, ", ") + "]"
}

func Uint16ToStr(u []uint16) string {
	r := utf16.Decode(u)
	return string(r)
}

func UTF16Length(s string) int {
	var length int
	for _, r := range s {
		if r >= 0x10000 {
			length += 2
		} else {
			length++
		}
	}
	return length
}

func IsValidUTF16(u []uint16) bool {
	var i int
	for i < len(u) {
		if IsHighSurrogate(u[i]) {
			i++
			if i == len(u) || !IsLowSurrogate(u[i]) {
				return false
			}
		} else if IsLowSurrogate(u[i]) {
			return false
		}
		i++
	}
	return true
}

func Uint16Equal(a, b []uint16) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

/*func UTF16Length(s string) int {
	return len(StrToUint16(s))
}*/

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func abs[T number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Ordered constraint represents types that are ordered,
// such as int, float64, and string.
type Ordered interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr |
		float32 | float64 |
		string
}

// SortMapKeys is a generic function that takes a map with keys of any comparable type
// and returns a sorted slice of those keys.
func MapSortedKeys[K Ordered, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j] // This works because K is constrained to comparable types
	})

	return keys
}

func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func containsOnlyKeysFrom(myMap map[string]interface{}, keys []string) bool {
	allowedKeys := make(map[string]bool, len(keys))
	for _, key := range keys {
		allowedKeys[key] = true
	}

	for key := range myMap {
		if !allowedKeys[key] {
			return false
		}
	}

	return true
}

func isCharInString(s string, c rune) bool {
	for _, r := range s {
		if r == c {
			return true
		}
	}
	return false
}

func (r *Rogue) RandInsert(rng *rand.Rand, maxLength int) (InsertAction, InsertOp, error) {
	visIx := rng.Intn(r.VisSize + 1)
	strLen := rng.Intn(maxLength) + 1
	s := RandomString(rng, strLen)

	if visIx < r.VisSize {
		c, err := r.GetChar(visIx)
		if err != nil {
			return InsertAction{}, InsertOp{}, fmt.Errorf("GetChar(%d): %w", visIx, err)
		}

		// special handling for surrogate pairs
		if visIx < r.VisSize && IsLowSurrogate(c) {
			visIx += 1
		}
	}

	op, err := r.Insert(visIx, s)
	if err != nil {
		return InsertAction{}, InsertOp{}, fmt.Errorf("Insert(%d, %q): %w", visIx, s, err)
	}

	return InsertAction{Index: visIx, Text: s}, op, nil

}

func (r *Rogue) RandDelete(rng *rand.Rand, maxLength int) (das DeleteAction, op Op, err error) {
	startIx := rng.Intn(r.VisSize)
	endIx := min(rng.Intn(maxLength)+startIx+1, r.VisSize-1)

	startChar, err := r.GetChar(startIx)
	if err != nil {
		return das, nil, fmt.Errorf("GetChar(%d): %w", startIx, err)
	}

	endChar, err := r.GetChar(endIx)
	if err != nil {
		return das, nil, fmt.Errorf("GetChar(%d): %w", endIx, err)
	}

	if IsLowSurrogate(startChar) {
		startIx--
	}

	if IsHighSurrogate(endChar) {
		endIx++
	}

	length := endIx - startIx + 1

	dop, err := r.Delete(startIx, length)
	if err != nil {
		return das, nil, fmt.Errorf("Delete(%d, %d): %w", startIx, length, err)
	}

	return DeleteAction{Index: startIx, Count: length}, dop, nil
}

func _randFormat(rng *rand.Rand) (f FormatV3) {
	t := rng.Float64()

	spans := []string{"b", "i", "u", "s", "e"}
	langs := []string{"python", "javascript", "java", "c", "go"}

	if t < 0.7 {
		k := spans[rng.Intn(len(spans))]
		return FormatV3Span{k: "true"}
	} else if t < 0.75 {
		return FormatV3Line{}
	} else if t < 0.8 {
		return FormatV3BlockQuote{}
	} else if t < 0.85 {
		lang := langs[rng.Intn(len(langs))]
		return FormatV3CodeBlock(lang)
	} else if t < 0.9 {
		indent := rng.Intn(4)
		return FormatV3BulletList(indent)
	} else if t < 0.95 {
		indent := rng.Intn(4)
		return FormatV3OrderedList(indent)
	} else {
		indent := rng.Intn(4) + 1
		return FormatV3Header(indent)
	}
}

func (r *Rogue) RandFormat(rng *rand.Rand, maxLength int) (act FormatAction, op Op, err error) {
	if r.VisSize < 3 {
		return act, nil, fmt.Errorf("VisSize < 3")
	}

	startIx := max(0, rng.Intn(r.VisSize-3))
	endIx := min(rng.Intn(maxLength)+startIx+1, r.VisSize-3)

	startChar, err := r.GetChar(startIx)
	if err != nil {
		return act, nil, fmt.Errorf("GetChar(%d): %w", startIx, err)
	}

	endChar, err := r.GetChar(endIx)
	if err != nil {
		return act, nil, fmt.Errorf("GetChar(%d): %w", endIx, err)
	}

	if IsLowSurrogate(startChar) {
		startIx--
	}

	if IsHighSurrogate(endChar) {
		endIx++
	}

	length := endIx - startIx + 1

	format := _randFormat(rng)
	fop, err := r.Format(startIx, length, format)
	if err != nil {
		return act, nil, fmt.Errorf("Format(%d, %d, %v): %w", startIx, length, format, err)
	}

	return FormatAction{Index: startIx, Length: length}, fop, nil
}

func intPow(base, exp int) int {
	result := 1
	for exp > 0 {
		if exp%2 == 1 {
			result *= base
		}
		base *= base
		exp /= 2
	}
	return result
}

func Interleave[T any](rng *rand.Rand, slices ...[]T) []T {
	nonEmptySlices := make([][]T, 0, len(slices))
	totalLength := 0

	for _, s := range slices {
		totalLength += len(s)
		if len(s) > 0 {
			nonEmptySlices = append(nonEmptySlices, s)
		}
	}

	result := make([]T, 0, totalLength)

	for len(nonEmptySlices) > 1 {
		index := rng.Intn(len(nonEmptySlices))
		result = append(result, nonEmptySlices[index][0])
		nonEmptySlices[index] = nonEmptySlices[index][1:]

		if len(nonEmptySlices[index]) == 0 {
			nonEmptySlices = append(nonEmptySlices[:index], nonEmptySlices[index+1:]...)
		}
	}

	if len(nonEmptySlices) == 1 {
		result = append(result, nonEmptySlices[0]...)
	}

	return result
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func IntersectMaps[K comparable, V comparable](map1, map2 map[K]V) map[K]V {
	result := make(map[K]V)

	for key, value := range map1 {
		if val, ok := map2[key]; ok && value == val {
			result[key] = value
		}
	}

	return result
}

func MergeMaps[K comparable, V any](map1, map2 map[K]V) map[K]V {
	result := make(map[K]V, len(map1)+len(map2))

	for key, value := range map1 {
		result[key] = value
	}

	for key, value := range map2 {
		result[key] = value
	}

	return result
}

func Identity[T any](t T) T {
	return t
}

func CopyMap[K comparable, V any](original map[K]V) map[K]V {
	newMap := make(map[K]V, len(original))
	for key, value := range original {
		newMap[key] = value
	}
	return newMap
}

func IsValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if u.Scheme == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}

	if u.Host == "" {
		return false
	}

	parts := strings.Split(u.Host, ".")
	if len(parts) < 2 {
		return false
	}

	return true
}

func CountLeadingChar(s string, char rune) int {
	count := 0
	for _, c := range s {
		if c != char {
			break
		}
		count++
	}
	return count
}

func Utf8ToUtf16Ix(s string, utf8Ix int) int {
	if utf8Ix < 0 || utf8Ix >= len(s) {
		return -1
	}

	utf16Ix := 0
	for i := 0; i < utf8Ix; {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError {
			return -1 // Invalid UTF-8 encoding
		}

		if utf8Ix-i < size {
			return utf16Ix
		}

		utf16Ix += len(utf16.Encode([]rune{r}))
		i += size
	}

	return utf16Ix
}

func MapSetMax[K comparable, V constraints.Ordered](m map[K]V, key K, newValue V) {
	if currentValue, exists := m[key]; !exists || newValue > currentValue {
		m[key] = newValue
	}
}

func PtrTo[T any](v T) *T {
	return &v
}

func SliceMinIx[T constraints.Ordered](slice []T) (ix int) {
	if len(slice) == 0 {
		return -1
	}

	minVal := slice[0]
	for i, v := range slice {
		if v < minVal {
			minVal = v
			ix = i
		}
	}

	return ix
}

func SliceMaxIx[T constraints.Ordered](slice []T) (ix int) {
	if len(slice) == 0 {
		return -1
	}

	maxVal := slice[0]
	for i, v := range slice {
		if v > maxVal {
			maxVal = v
			ix = i
		}
	}

	return ix
}
