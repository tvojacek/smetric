package smetric

import (
	"bytes"
	"errors"
	"strings"
	"sync"
)

var InvalidCharacterError = errors.New("invalid character")

const toLower = 'a' - 'A'

// GetSnakeMetricName convert input string to snake case metric names function return error if input contain invalid character
// valid characters are [a-z,A-Z,0-9,_]
func GetSnakeMetricName(input string) (string, error) {
	var output strings.Builder
	output.Grow(len(input) + 5)
	lastUpper := true
	for _, char := range input {
		switch {
		case char >= 'A' && char <= 'Z':
			if !lastUpper {
				output.WriteByte('_')
				lastUpper = true
			}
			output.WriteRune(char + toLower)
		case char >= '0' && char <= '9':
			fallthrough
		case char >= 'a' && char <= 'z':
			output.WriteRune(char)
			lastUpper = false
		case char == '_':
			output.WriteByte('_')
			lastUpper = true
		default:
			return output.String(), InvalidCharacterError
		}
	}
	return output.String(), nil
}

var builderPool = sync.Pool{
	New: func() any {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		var builder []byte
		return builder
	},
}

// GetSnakeMetricNameSync convert input string to snake case metric names function return error if input contain invalid character
// valid characters are [a-z,A-Z,0-9,_]
func GetSnakeMetricNameSync(input string) (string, error) {
	buf := builderPool.Get().([]byte)
	output := bytes.NewBuffer(buf)
	defer func() {
		output.Reset()
		builderPool.Put(buf[:0])
	}()
	output.Grow(len(input) + 5)
	lastUpper := true
	for _, char := range input {
		switch {
		case char >= 'A' && char <= 'Z':
			if !lastUpper {
				output.WriteByte('_')
				lastUpper = true
			}
			output.WriteRune(char + toLower)
		case char >= '0' && char <= '9':
			fallthrough
		case char >= 'a' && char <= 'z':
			output.WriteRune(char)
			lastUpper = false
		case char == '_':
			output.WriteByte('_')
			lastUpper = true
		default:
			return output.String(), InvalidCharacterError
		}
	}
	return output.String(), nil
}
