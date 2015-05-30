package prism

import "testing"

func TestNewPath(t *testing.T) {
	test_data := map[string]*Path{
		"articles/": &Path{
			resourceType: "articles",
			ID:           "",
		},
	}

	for pathStr, path := range test_data {
		convPath := NewPath(pathStr)
		if convPath.resourceType != path.resourceType {
			t.Errorf("not match: %v : %v", convPath.resourceType, path.resourceType)
		}
	}
}
