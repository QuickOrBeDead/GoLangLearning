package datastructures

type LinkSet map[string]struct{}

func (s LinkSet) Add(link string) {
	s[link] = struct{}{}
}

func (s LinkSet) Remove(link string) {
	delete(s, link)
}

func (s LinkSet) Contains(link string) bool {
	_, ok := s[link]
	return ok
}
