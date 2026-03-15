package statusline

import "sync"

// SegmentRegistry maps segment type names to their implementations.
type SegmentRegistry struct {
	segments map[string]Segment
}

// NewSegmentRegistry creates an empty segment registry.
func NewSegmentRegistry() *SegmentRegistry {
	return &SegmentRegistry{segments: make(map[string]Segment)}
}

// Register adds a segment implementation.
func (r *SegmentRegistry) Register(seg Segment) {
	r.segments[seg.Name()] = seg
}

// Get returns the segment for the given type name, or nil.
func (r *SegmentRegistry) Get(name string) Segment {
	return r.segments[name]
}

// ProviderRegistry maps provider names to their implementations.
type ProviderRegistry struct {
	providers map[string]DataProvider
}

// NewProviderRegistry creates an empty provider registry.
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{providers: make(map[string]DataProvider)}
}

// Register adds a data provider implementation.
func (r *ProviderRegistry) Register(p DataProvider) {
	r.providers[p.Name()] = p
}

// CollectProviderNames walks the tree and returns the set of provider
// names needed for rendering (skipping disabled nodes).
func (r *ProviderRegistry) CollectProviderNames(tree []SegmentNode) map[string]bool {
	names := make(map[string]bool)
	collectNames(tree, names)
	return names
}

func collectNames(nodes []SegmentNode, names map[string]bool) {
	for _, node := range nodes {
		if node.Enabled != nil && !*node.Enabled {
			continue
		}
		if node.Provider != "" {
			names[node.Provider] = true
		}
		if len(node.Children) > 0 {
			collectNames(node.Children, names)
		}
	}
}

// ResolveAll resolves all named providers concurrently and returns
// a map of provider name → resolved data.
func (r *ProviderRegistry) ResolveAll(names map[string]bool, session *SessionData) map[string]any {
	results := make(map[string]any)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for name := range names {
		p, ok := r.providers[name]
		if !ok {
			continue
		}
		wg.Add(1)
		go func(provider DataProvider) {
			defer wg.Done()
			data, err := provider.Resolve(session)
			if err != nil {
				return
			}
			mu.Lock()
			results[provider.Name()] = data
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	return results
}
