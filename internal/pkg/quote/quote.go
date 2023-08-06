package quote

import "math/rand"

var quotes = []string{
	"A rose by any other name would smell as sweet. - William Shakespeare",
	"All that glitters is not gold. - William Shakespeare",
	"All the world's a stage, and all the men and women merely players. - William Shakespeare",
	"Ask not what your country can do for you; ask what you can do for your country. - John Kennedy",
}

// Random returns random quote, it is a convenience wrapper for tests
var Random = RandomDefault

func RandomDefault() string {
	return quotes[rand.Intn(len(quotes)-1)] //nolint:gosec // it's just quotes
}
