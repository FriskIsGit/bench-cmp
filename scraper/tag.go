package scraper

import (
	"fmt"
	"strings"
)

// ComparisonStrategy defines the comparison policy for each pair of attributes.
type ComparisonStrategy int

const (
	// Exact matches both key and value.
	Exact ComparisonStrategy = iota

	// MatchKey matches only the key.
	MatchKey

	// MatchValue matches only the value.
	MatchValue

	// KeyPrefix requires the key to be prefixed with the given string.
	KeyPrefix

	// ValuePrefix requires the value to be prefixed with the given string.
	ValuePrefix

	// KeySuffix requires the key to be suffixed with the given string.
	KeySuffix

	// ValueSuffix requires the value to be suffixed with the given string.
	ValueSuffix

	// KeyAndValuePrefix requires the key to match fully and the value to be prefixed with the given string.
	KeyAndValuePrefix

	// KeyPrefixAndValue requires the key to be prefixed with the given string and the value to match fully.
	KeyPrefixAndValue

	// KeyAndValueSuffix requires the key to match fully and the value to be suffixed with the given string.
	KeyAndValueSuffix

	// KeySuffixAndValue requires the key to be suffixed with the given string and the value to match fully.
	KeySuffixAndValue
)

// Compare is an attribute comparison wrapper used as a predicate when searching for tags.
type Compare struct {
	Key      string
	Value    string
	Strategy ComparisonStrategy
}

// ExactCompare creates a Compare that matches both key and value exactly.
func ExactCompare(key, value string) Compare {
	return Compare{Key: key, Value: value, Strategy: Exact}
}

// KeyCompare creates a Compare that matches only the key.
func KeyCompare(key string) Compare {
	return Compare{Key: key, Strategy: MatchKey}
}

// ValueCompare creates a Compare that matches only the value.
func ValueCompare(value string) Compare {
	return Compare{Value: value, Strategy: MatchValue}
}

// KeyPrefixCompare creates a Compare where the key must be prefixed with the given string.
func KeyPrefixCompare(keyPrefix string) Compare {
	return Compare{Key: keyPrefix, Strategy: KeyPrefix}
}

// ValuePrefixCompare creates a Compare where the value must be prefixed with the given string.
func ValuePrefixCompare(valuePrefix string) Compare {
	return Compare{Value: valuePrefix, Strategy: ValuePrefix}
}

// KeySuffixCompare creates a Compare where the key must be suffixed with the given string.
func KeySuffixCompare(keySuffix string) Compare {
	return Compare{Key: keySuffix, Strategy: KeySuffix}
}

// ValueSuffixCompare creates a Compare where the value must be suffixed with the given string.
func ValueSuffixCompare(valueSuffix string) Compare {
	return Compare{Value: valueSuffix, Strategy: ValueSuffix}
}

// KeyAndValuePrefixCompare creates a Compare where the key must match fully and the value must be prefixed.
func KeyAndValuePrefixCompare(key, valuePrefix string) Compare {
	return Compare{Key: key, Value: valuePrefix, Strategy: KeyAndValuePrefix}
}

// KeyPrefixAndValueCompare creates a Compare where the key must be prefixed and the value must match fully.
func KeyPrefixAndValueCompare(keyPrefix, value string) Compare {
	return Compare{Key: keyPrefix, Value: value, Strategy: KeyPrefixAndValue}
}

// KeySuffixAndValueCompare creates a Compare where the key must be suffixed and the value must match fully.
func KeySuffixAndValueCompare(keySuffix, value string) Compare {
	return Compare{Key: keySuffix, Value: value, Strategy: KeySuffixAndValue}
}

// KeyAndValueSuffixCompare creates a Compare where the key must match fully and the value must be suffixed.
func KeyAndValueSuffixCompare(key, valueSuffix string) Compare {
	return Compare{Key: key, Value: valueSuffix, Strategy: KeyAndValueSuffix}
}

// Attribute represents a single key-value pair on an HTML tag.
type Attribute struct {
	Key   string
	Value string
}

// Tag represents an HTML tag which may include attributes in the form of:
//
//	<tagname prop1="val1" prop2="val2">
type Tag struct {
	Name        string
	StartOffset int
	EndOffset   int
	Attributes  []Attribute
}

// NewTag creates a new Tag with the given name and no attributes.
func NewTag(name string) *Tag {
	return &Tag{
		Name:       name,
		EndOffset:  -1,
		Attributes: []Attribute{},
	}
}

// NewTagWithAttributes creates a new Tag with the given name and attributes.
func NewTagWithAttributes(name string, attributes []Attribute) *Tag {
	return &Tag{
		Name:       name,
		EndOffset:  -1,
		Attributes: attributes,
	}
}

// GetAttribute searches for an attribute by key. If another pair with the same key exists,
// the first pair's value is returned. Returns the value and true if found, or empty string
// and false if no match is found.
func (t *Tag) GetAttribute(key string) (string, bool) {
	for _, attr := range t.Attributes {
		if attr.Key == key {
			return attr.Value, true
		}
	}
	return "", false
}

// CompareAttributes checks whether the tag's attributes match all of the given comparison predicates.
// The tag must have at least as many attributes as comparisons. The order of attributes does not matter.
func (t *Tag) CompareAttributes(comparisons []Compare) bool {
	if len(t.Attributes) < len(comparisons) {
		return false
	}

	for _, predicate := range comparisons {
		if !t.attributesMatchPredicate(predicate) {
			return false
		}
	}

	return true
}

// attributesMatchPredicate determines if any of the tag's attributes match the given predicate.
func (t *Tag) attributesMatchPredicate(predicate Compare) bool {
	for _, attr := range t.Attributes {
		switch predicate.Strategy {
		case Exact:
			if attr.Key == predicate.Key && attr.Value == predicate.Value {
				return true
			}
		case MatchKey:
			if attr.Key == predicate.Key {
				return true
			}
		case MatchValue:
			if attr.Value == predicate.Value {
				return true
			}
		case KeyPrefix:
			if strings.HasPrefix(attr.Key, predicate.Key) {
				return true
			}
		case ValuePrefix:
			if strings.HasPrefix(attr.Value, predicate.Value) {
				return true
			}
		case KeySuffix:
			if strings.HasSuffix(attr.Key, predicate.Key) {
				return true
			}
		case ValueSuffix:
			if strings.HasSuffix(attr.Value, predicate.Value) {
				return true
			}
		case KeyAndValuePrefix:
			if attr.Key == predicate.Key && strings.HasPrefix(attr.Value, predicate.Value) {
				return true
			}
		case KeyAndValueSuffix:
			if attr.Key == predicate.Key && strings.HasSuffix(attr.Value, predicate.Value) {
				return true
			}
		case KeyPrefixAndValue:
			if strings.HasPrefix(attr.Key, predicate.Key) && attr.Value == predicate.Value {
				return true
			}
		case KeySuffixAndValue:
			if strings.HasSuffix(attr.Key, predicate.Key) && attr.Value == predicate.Value {
				return true
			}
		}
	}
	return false
}

// EstimateEndOffset estimates the end offset of the opening tag based on the tag name
// and attribute lengths.
func (t *Tag) EstimateEndOffset() int {
	end := t.StartOffset + 1 + len(t.Name)
	for _, attr := range t.Attributes {
		end += 1 + len(attr.Key) + len(attr.Value) + 3
	}
	return end
}

// String returns a string representation of the Tag.
func (t *Tag) String() string {
	return fmt.Sprintf("Tag{name=%s, index=%d}", t.Name, t.StartOffset)
}
