package search

import "github.com/fruitsco/goji/conf"

type SearchDriver string

const (
	Typesense SearchDriver = "typesense"
)

type Config struct {
	Driver SearchDriver `conf:"driver"`

	Typesense *TypesenseConfig `conf:"typesense"`
}

var DefaultConfig = conf.DefaultConfig{
	"tasks.driver": "typesense",
}
