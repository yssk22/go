package config

// a list of []*ServiceConfig that contains global configurations whose 'Value' as default value.
// isGlobal in each element should be true
var globalDefaults []*ServiceConfig

// Global register a service configuration shared in whole applications, and it can be overwritten by a service.
func Global(key string, defaultValue string, description string) {
	sc := newServiceConfig(key, defaultValue, description)
	sc.isGlobal = true
	globalDefaults = append(globalDefaults, sc)
}
