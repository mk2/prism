package prism

import "strings"

const (
	PathSeparator = "/"
)

type Path struct {
	PathInterface
	resourceType string
	ID           string
}

type PathInterface interface {
	hasID() bool
}

func (p *Path) hasID() bool {

	return len(p.ID) > 0
}

func NewPath(str string) *Path {

	str = strings.Trim(str, PathSeparator)
	tokens := strings.Split(str, PathSeparator)

	var resourceType, ID string

	if len(tokens) > 1 {
		resourceType, ID = tokens[0], tokens[1]
	} else {
		resourceType, ID = tokens[0], ""
	}

	return &Path{
		resourceType: resourceType,
		ID:           ID,
	}
}
