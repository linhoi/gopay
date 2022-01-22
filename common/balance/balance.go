package balance

// Weight is a interface that implement a weighted round robin algorithm.
type Weight interface {
	// Next gets next selected item.
	// Next is  goroutine-safe. You don't use the snchronization primitive to protect it in concurrent cases.
	Next() (item interface{})
	// Add adds a weighted item for selection. if adds an exit item ,weight  plussed
	Add(item interface{}, weight int)
	//All returns all items.
	All() map[interface{}]int

	// Reset resets the balancing algorithm.
	Reset()

	// RemoveAll removes all weighted items. Next will return nil if remove all
	RemoveAll()

	OnChange(c *Config)
}

type Config struct {
	Items []Item `yaml:"items"`
}

type Item struct {
	Name   string `yaml:"name"`
	Weight int    `yaml:"weight"`
}
