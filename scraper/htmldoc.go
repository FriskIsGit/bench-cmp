package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:130.0) Gecko/20100101 Firefox/130"

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// HtmlDoc is an HTML parser that supports finding tags by name and attributes,
// extracting nested tags, and extracting text content.
type HtmlDoc struct {
	html        string
	len         int
	concatChar  byte
	delimitTags bool
	brToNewline bool
}

// NewHtmlDoc creates a new HtmlDoc from the given HTML string.
func NewHtmlDoc(contents string) *HtmlDoc {
	return &HtmlDoc{
		html:        contents,
		len:         len(contents),
		concatChar:  '\n',
		delimitTags: true,
		brToNewline: false,
	}
}

// SetConcatenatingChar sets the character used to delimit text extracted from adjacent tags.
func (d *HtmlDoc) SetConcatenatingChar(ch byte) {
	d.concatChar = ch
}

// ReplaceLineBreakWithNewLine controls whether <br> tags are replaced with newline characters
// during text extraction.
func (d *HtmlDoc) ReplaceLineBreakWithNewLine(enabled bool) {
	d.brToNewline = enabled
}

// DelimitTags controls whether a delimiter is inserted when concatenating text from nested tags.
// In HTML, text content within tags can be combined with other elements to create a formatted
// nested tree. Therefore it's important to specify a delimiter so that extracting the text can
// resemble the original. See SetConcatenatingChar.
func (d *HtmlDoc) DelimitTags(enabled bool) {
	d.delimitTags = enabled
}

// Find finds a tag by name matching at least the given attributes.
// The tag returned will have at least the number of attributes specified.
// The order in which attributes are provided does not matter.
// Returns the Tag or nil if not found.
func (d *HtmlDoc) Find(tag string, comparisons ...Compare) *Tag {
	return d.FindFrom(tag, 0, comparisons...)
}

// FindAll finds all tags by name and attributes from the beginning of the document.
// Returns all Tag objects which matched, possibly an empty slice.
func (d *HtmlDoc) FindAll(tag string, comparisons ...Compare) []*Tag {
	return d.FindAllFrom(tag, 0, comparisons...)
}

// FindAllFrom finds all tags by name and attributes starting at the given index until the
// end of the document. Returns all Tag objects which matched, possibly an empty slice.
func (d *HtmlDoc) FindAllFrom(tag string, from int, comparisons ...Compare) []*Tag {
	var tags []*Tag
	cursor := from
	for cursor < d.len {
		found := d.FindFrom(tag, cursor, comparisons...)
		if found == nil {
			break
		}
		tags = append(tags, found)
		if found.EndOffset == -1 {
			cursor = found.StartOffset + 1
		} else {
			cursor = found.EndOffset + 1
		}
	}
	return tags
}

// FindFrom linearly searches for a tag in the HtmlDoc returning the first Tag that matches
// predicates, starting from the given index. Returns the Tag or nil if not found.
func (d *HtmlDoc) FindFrom(tag string, from int, comparisons ...Compare) *Tag {
	for i := from; i < d.len; i++ {
		chr := d.html[i]
		if chr != '<' {
			continue
		}

		if d.isCommentAhead(i) {
			i = d.skipComment(i + 1)
			continue
		}

		closing := i+1 < d.len && d.html[i+1] == '/'
		if closing {
			tagNameSt := i + 2
			i = d.indexOf(">", tagNameSt)
			continue
		}

		hasAttributes := false
		j := i + 1
		for ; j < d.len; j++ {
			if isHtmlWhitespace(d.html[j]) {
				hasAttributes = true
				break
			}
			if d.html[j] == '>' {
				break
			}
		}

		tagName := strings.TrimSuffix(d.html[i+1:j], "/")
		parsedTag := NewTag(tagName)
		end := -1
		if hasAttributes {
			end = d.parseAttributes(parsedTag, j+1)
		}

		if parsedTag.Name == tag && parsedTag.CompareAttributes(comparisons) {
			parsedTag.StartOffset = i
			return parsedTag
		}

		lowerName := strings.ToLower(parsedTag.Name)
		if lowerName == "script" || lowerName == "style" {
			i = d.skipRawText(lowerName, i)
			continue
		}

		if hasAttributes {
			if end == -1 {
				break
			}
			i = end
		}
	}

	return nil
}

// ExtractTags searches for all tags matching the target name within the given tag.
// Returns a slice of Tag objects found within the parent tag's scope.
func (d *HtmlDoc) ExtractTags(tag *Tag, target string) []*Tag {
	if tag.StartOffset < 0 || tag.StartOffset >= d.len {
		return nil
	}

	var extractedTags []*Tag
	from := tag.EndOffset
	if from == -1 {
		from = tag.EstimateEndOffset()
	}

	tagStack := []string{tag.Name}

	for i := from; i < d.len; i++ {
		chr := d.html[i]
		if chr != '<' {
			continue
		}

		if d.isCommentAhead(i) {
			i = d.skipComment(i + 1)
			continue
		}

		closing := i+1 < d.len && d.html[i+1] == '/'
		if closing {
			tagNameSt := i + 2
			i = d.indexOf(">", tagNameSt)
			closedTagName := strings.TrimSpace(d.html[tagNameSt:i])
			// Pop all void tags
			for len(tagStack) > 0 {
				popped := tagStack[len(tagStack)-1]
				tagStack = tagStack[:len(tagStack)-1]
				if popped == closedTagName {
					break
				}
			}
			if len(tagStack) == 0 {
				return extractedTags
			}
			continue
		}

		hasAttributes := false
		j := i + 1
		for ; j < d.len; j++ {
			if isHtmlWhitespace(d.html[j]) {
				hasAttributes = true
				break
			}
			if d.html[j] == '>' {
				break
			}
		}

		tagName := strings.TrimSuffix(d.html[i+1:j], "/")
		parsedTag := NewTag(tagName)
		end := -1
		if hasAttributes {
			end = d.parseAttributes(parsedTag, j+1)
		}

		tagStack = append(tagStack, parsedTag.Name)
		if end != -1 && d.html[end-1] == '/' {
			tagStack = tagStack[:len(tagStack)-1]
		}
		// It could be that the tag given is a void tag in which case it has no tags
		if len(tagStack) == 0 {
			return extractedTags
		}

		if parsedTag.Name == target {
			parsedTag.StartOffset = i
			extractedTags = append(extractedTags, parsedTag)
		}

		if hasAttributes {
			if end == -1 {
				break
			}
			i = end
		}
	}

	return extractedTags
}

// ExtractText extracts text from the given tag and all its sub-tags beginning at tag.StartOffset.
// The extraction will begin at the top level tag and will continue until either the tag is closed
// or the end of document is reached. Text extracted from sub-tags will be concatenated using the
// specified delimiter that can be set by calling SetConcatenatingChar. The EndOffset field of the
// given tag will be set to the index where parsing finished if its value was -1.
func (d *HtmlDoc) ExtractText(tag *Tag) string {
	if tag.StartOffset < 0 || tag.StartOffset >= d.len {
		return ""
	}

	appendText := false
	concatenate := false
	var stack []string
	var text strings.Builder

	for i := tag.StartOffset; i < d.len; i++ {
		chr := d.html[i]
		switch chr {
		case ' ':
			if appendText {
				if d.delimitTags && concatenate && text.Len() > 0 {
					text.WriteByte(d.concatChar)
				}
				text.WriteByte(' ')
				concatenate = false
			}

		case '"':
			if appendText {
				text.WriteByte('"')
			}

		case '<':
			if d.isCommentAhead(i) {
				i = d.skipComment(i + 1)
				continue
			}

			closing := i+1 < d.len && d.html[i+1] == '/'
			if closing {
				i++
			}

			hasAttributes := false
			tagEnd := i + 1
			for j := tagEnd; j < d.len; j++ {
				c := d.html[j]
				if c == ' ' {
					if closing {
						continue
					}
					hasAttributes = true
					tagEnd = j
					break
				}
				if c == '>' {
					tagEnd = j
					break
				}
			}

			anyTag := d.html[i+1 : tagEnd]
			if d.brToNewline && strings.HasPrefix(anyTag, "br") {
				text.WriteByte('\n')
			}

			i = tagEnd
			voidTag := len(anyTag) > 0 && anyTag[len(anyTag)-1] == '/'

			if closing {
				for len(stack) > 0 {
					popped := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					if popped == anyTag {
						break
					}
				}
				if len(stack) == 0 {
					if tag.EndOffset == -1 {
						tag.EndOffset = i
					}
					return text.String()
				}
			} else if !voidTag {
				stack = append(stack, anyTag)
			}

			appendText = !hasAttributes
			concatenate = true

		case '>':
			concatenate = true
			appendText = true

		default:
			if appendText {
				if d.delimitTags && concatenate && text.Len() > 0 {
					text.WriteByte(d.concatChar)
				}
				concatenate = false
				text.WriteByte(chr)
			}
		}
	}

	if tag.EndOffset == -1 {
		tag.EndOffset = d.len
	}
	return text.String()
}

// isHtmlWhitespace checks if a character is an HTML whitespace character.
// See: https://www.w3.org/TR/html4/struct/text.html#h-9.1
func isHtmlWhitespace(c byte) bool {
	return c == ' ' || c == '\r' || c == '\n' || c == '\t' || c == '\f'
}

// parseAttributes parses tag attributes starting from the given index (after the tag name).
// Returns the index of '>' where the tag ends, or -1 if EOF is reached.
func (d *HtmlDoc) parseAttributes(parsedTag *Tag, from int) int {
	var name, value strings.Builder
	afterEqual := false
	var quote byte

	appendAttribute := func() {
		if name.Len() == 0 {
			return
		}
		parsedTag.Attributes = append(parsedTag.Attributes,
			Attribute{Key: name.String(), Value: value.String()})
		name.Reset()
		value.Reset()
		afterEqual = false
	}

	for i := from; i < d.len; i++ {
		chr := d.html[i]

		if quote != 0 {
			if chr == quote {
				appendAttribute()
				quote = 0
			} else {
				value.WriteByte(chr)
			}
			continue
		}

		switch chr {
		case '=':
			afterEqual = true

		case ' ', '\n', '\r', '\t', '\f':
			if !afterEqual {
				next := i + 1
				for next < d.len && isHtmlWhitespace(d.html[next]) {
					next++
				}
				if next >= d.len || d.html[next] != '=' {
					appendAttribute()
				}
				break
			}
			if value.Len() > 0 {
				appendAttribute()
			}

		case '"', '\'':
			if afterEqual {
				quote = chr
			} else {
				name.WriteByte(chr)
			}

		case '/':
			if afterEqual {
				value.WriteByte(chr)
				break
			}
			endOffset := i + 1
			for ; endOffset < d.len; endOffset++ {
				if d.html[endOffset] == '>' {
					break
				}
			}
			appendAttribute()
			return endOffset

		case '>':
			appendAttribute()
			return i

		default:
			if afterEqual {
				value.WriteByte(chr)
			} else {
				name.WriteByte(chr)
			}
		}
	}

	return -1
}

func (d *HtmlDoc) skipComment(from int) int {
	return d.indexOf("-->", from)
}

func (d *HtmlDoc) isCommentAhead(lt int) bool {
	return lt+1 < d.len && d.html[lt+1] == '!' &&
		lt+2 < d.len && d.html[lt+2] == '-' &&
		lt+3 < d.len && d.html[lt+3] == '-'
}

// indexOf returns the index of the first occurrence of text starting from the given position.
// Always returns indices within bounds.
func (d *HtmlDoc) indexOf(text string, from int) int {
	if from >= d.len {
		return d.len - 1
	}
	index := strings.Index(d.html[from:], text)
	if index == -1 {
		return d.len - 1
	}
	return from + index
}

// skipRawText returns the index of the closing tag's '>' for a raw text element
// such as script or style, or the end of the document if no closing tag exists.
func (d *HtmlDoc) skipRawText(tagName string, from int) int {
	for i := from; i+1 < d.len; i++ {
		if d.html[i] != '<' || d.html[i+1] != '/' {
			continue
		}
		nameStart := i + 2
		nameEnd := nameStart
		for nameEnd < d.len && !isHtmlWhitespace(d.html[nameEnd]) && d.html[nameEnd] != '>' {
			nameEnd++
		}
		if strings.ToLower(d.html[nameStart:nameEnd]) == tagName {
			return d.indexOf(">", nameEnd)
		}
	}
	return d.len - 1
}

// FetchHtml performs an HTTP GET request to the given URL and returns the response body as a string.
func FetchHtml(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Accept-Language", "en-US;q=0.7")
	req.Header.Set("Set-GPC", "1")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return string(body), fmt.Errorf("response code: %d", resp.StatusCode)
	}

	return string(body), nil
}
