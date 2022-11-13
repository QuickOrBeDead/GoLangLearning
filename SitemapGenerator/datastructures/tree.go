package datastructures

type SitemapNode struct {
	Url      string
	Level    int
	Children []*SitemapNode
	Parent   *SitemapNode
}

func (t *SitemapNode) AddChild(url string, parent *SitemapNode) *SitemapNode {
	if t.hasUrl(url) {
		return nil
	}

	child := &SitemapNode{Url: url, Children: make([]*SitemapNode, 0), Parent: t, Level: t.Level + 1}
	t.Children = append(t.Children, child)

	return child
}

func (t *SitemapNode) hasUrl(url string) bool {
	if t.Url == url {
		return true
	}

	for _, c := range t.Children {
		if c.Url == url {
			return true
		}
	}

	for p := t.Parent; p != nil; p = p.Parent {
		if p.Url == url {
			return true
		}
	}

	return false
}
