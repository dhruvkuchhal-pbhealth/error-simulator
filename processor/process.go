package processor

// Process invokes fn for each item. Panic can happen inside fn (caller's code) = callback genre.
func Process(items []Item, fn func(Item)) {
	for _, it := range items {
		invoke(it, fn)
	}
}
