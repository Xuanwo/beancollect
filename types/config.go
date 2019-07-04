package types

// Config is the config for beancollect transform.
type Config struct {
	Account map[string]string `yaml:"account"`
	Rules   []Rule            `yaml:"rules"`
}

// Rule is the config for roles.
type Rule struct {
	Type      string
	Condition map[string]string
	Value     string
}
