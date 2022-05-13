package starlet

import (
	"strings"
)

var (
	firstIDChars      = "abcdefghijklmnopqrstuvwxyz"
	subsequentIDChars = "abcdefghijklmnopqrstuvwxyz0123456789"
	firstBase         int
	subsequentBase    int
)

func init() {
	firstBase = len(firstIDChars)
	subsequentBase = len(subsequentIDChars)
}

type RenamingMap struct {
	Map map[string]string
	Len int
}

func NewRenamingMap() *RenamingMap {
	return &RenamingMap{
		Map: make(map[string]string),
		Len: 0,
	}
}

func (r *RenamingMap) Assign(longID string) string {
	shortID, ok := r.Map[longID]
	if ok {
		return shortID
	}

	base := 0
	iDChars := ""
	value := r.Len
	shortID = ""

	for {
		if value < subsequentBase {
			base = firstBase
			iDChars = firstIDChars
			if shortID != "" {
				value -= 1
			}
		} else {
			base = subsequentBase
			iDChars = subsequentIDChars
		}

		d := value % base
		value /= base
		shortID = string(iDChars[d]) + shortID

		if value == 0 {
			break
		}
	}

	r.Map[longID] = shortID
	r.Len++
	return shortID
}

func renameClasses(rules []*CSSRule, renamingMap *RenamingMap) {
	for _, rule := range rules {
		renameClasses(rule.Duplicates, renamingMap)
		for p, part := range rule.Selector {
			if part.Type != ClassSelector {
				continue
			}
			parts := strings.Split(part.Name, "-")
			for p, part := range parts {
				parts[p] = renamingMap.Assign(part)
			}
			rule.Selector[p].Name = strings.Join(parts, "-")
		}
	}
}
