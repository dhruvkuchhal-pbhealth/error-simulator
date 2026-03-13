package processor

// invoke calls the visitor for one item. Stack: handler (callback) → invoke.go → process.go.
func invoke(it Item, fn func(Item)) {
	fn(it)
}
