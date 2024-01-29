package utilities

import (
	"bytes"
	"os"

	"howett.net/plist"
)

// plistCustomType represents a custom data structure for plist encoding and decoding.
type plistCustomType struct {
	Name  string
	Value int
}

// EncodeMapToPlist encodes a map[string]interface{} to a plist format and writes it to standard output.
// The `format` parameter determines the plist format (e.g., XML, Binary).
// Returns an error if encoding fails.
func EncodeMapToPlist(data map[string]interface{}, format int) error {
	encoder := plist.NewEncoderForFormat(os.Stdout, format)
	return encoder.Encode(data)
}

// DecodePlistToMap decodes plist data from a byte slice into a map[string]interface{}.
// Returns the decoded map, the format of the plist, and an error if decoding fails.
func DecodePlistToMap(data []byte) (map[string]interface{}, int, error) {
	var result map[string]interface{}
	decoder := plist.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&result)
	return result, decoder.Format, err
}

// EncodeCustomTypeToPlist encodes an instance of plistCustomType to a plist format and writes it to standard output.
// The `format` parameter determines the plist format (e.g., XML, Binary).
// Returns an error if encoding fails.
func EncodeCustomTypeToPlist(customType plistCustomType, format int) error {
	encoder := plist.NewEncoderForFormat(os.Stdout, format)
	return encoder.Encode(customType)
}

// DecodePlistToCustomType decodes plist data from a byte slice into a provided plistCustomType pointer.
// Updates the value pointed by `customType` with the decoded data.
// Returns the format of the plist and an error if decoding fails.
func DecodePlistToCustomType(data []byte, customType *plistCustomType) (int, error) {
	decoder := plist.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(customType)
	return decoder.Format, err
}
