package processor

// Item is passed to the callback. Child may be nil — callback that derefs it panics (genre: callback/visitor).
type Item struct {
	Name  string
	Child *Item
}
