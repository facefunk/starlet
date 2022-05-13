package starlet

// CSSStatement ...
type CSSStatement struct {
	Property     string
	Value        string
	OriginalFile string
	OriginalLine int
	OriginalName string
}

// SetMapping sets the original starlet file mapping information for this CSSStatement.
func (statement *CSSStatement) SetMapping(file string, line int, name string) {
	statement.OriginalFile = file
	statement.OriginalLine = line
	statement.OriginalName = name
}
