package statusline

// RegisterBuiltinProviders adds all built-in data providers to the registry.
func RegisterBuiltinProviders(registry *ProviderRegistry) {
	registry.Register(&pwdProvider{})
	registry.Register(&gitProvider{})
	registry.Register(&contextProvider{})
	registry.Register(&modelProvider{})
	registry.Register(&costProvider{})
	registry.Register(&sessionProvider{})
}
