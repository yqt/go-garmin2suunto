package suunto

import "strings"

var (
	NoteActivityIdBeginTag = "_GM_B_"
	NoteActivityIdEndTag   = "_GM_E_"
)

type MoveItem struct {
	ActivityID       int64  `json:"ActivityID"`
	MoveID           int64  `json:"MoveID"`
	Notes            string `json:"Notes"`
	LastModifiedDate string `json:"LastModifiedDate"`
	LocalStartTime   string `json:"LocalStartTime"`
}

func (m *MoveItem) GetStartDate() string {
	if m.LocalStartTime == "" {
		return ""
	}
	return m.LocalStartTime[:10]
}

func (m *MoveItem) GetGarminActivityId() string {
	if m.Notes == "" {
		return ""
	}

	beginPos := strings.Index(m.Notes, NoteActivityIdBeginTag)
	if beginPos == -1 {
		return ""
	}
	rest := m.Notes[beginPos:]
	endPos := strings.Index(rest, NoteActivityIdEndTag)
	if endPos == -1 {
		return ""
	}

	return rest[len(NoteActivityIdBeginTag):endPos]
}
