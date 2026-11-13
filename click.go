package goatcounter

import "strings"

const (
	clickSegmentMaxLength = 64
	clickElementMaxLength = 128
)

// ClickPath is the four-level identifier encoded in a click event path.
// Click events use: click/<project>.<page>.<type>.<element>.
type ClickPath struct {
	Project string
	Page    string
	Type    string
	Element string
}

// ParseClickPath parses a canonical click event path. Event paths are stored
// without a leading slash by Hit.cleanPath.
func ParseClickPath(path string) (ClickPath, bool) {
	var click ClickPath
	if !strings.HasPrefix(path, "click/") {
		return click, false
	}

	parts := strings.Split(strings.TrimPrefix(path, "click/"), ".")
	if len(parts) != 4 ||
		!validClickSegment(parts[0], clickSegmentMaxLength) ||
		!validClickSegment(parts[1], clickSegmentMaxLength) ||
		!validClickSegment(parts[2], clickSegmentMaxLength) ||
		!validClickSegment(parts[3], clickElementMaxLength) {
		return click, false
	}

	return ClickPath{
		Project: parts[0],
		Page:    parts[1],
		Type:    parts[2],
		Element: parts[3],
	}, true
}

func validClickSegment(segment string, maxLength int) bool {
	if len(segment) == 0 || len(segment) > maxLength {
		return false
	}
	for _, r := range segment {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return true
}
