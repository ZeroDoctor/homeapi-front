package model

// Label :
type Label struct {
	Name     string
	ID       string
	Folder   bool
	Open     bool
	Index    int
	Depth    int
	Children []Label
}

// NewLabel :
func NewLabel(name, id string, folder, open bool, index, depth int, children []Label) Label {
	return Label{
		Name:     name,
		ID:       id,
		Folder:   folder,
		Open:     open,
		Index:    index,
		Depth:    depth,
		Children: children,
	}
}

// can't convert []T to []interface{}.
// I would if I could but: https://golang.org/doc/faq#convert_slice_of_interface

// Insert :
func Insert(a []Label, c Label, i int) []Label {
	a = append(a, NewLabel("", "", false, false, -1, -1, nil))
	copy(a[i+1:], a[i:])
	a[i] = c
	return a
}

// RemoveFrom :
func RemoveFrom(a []Label, i, j int) []Label {
	return append(a[:i], a[j:]...)
}

// Remove :
func Remove(a []Label, i int) []Label {
	return append(a[:i], a[i+1:]...)
}
