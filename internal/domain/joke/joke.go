// Package joke contains domain data for jokes
package joke

import "time"

type Joke struct {
	ID        int
	Joke      string
	Nsfw      bool
	CreatedAt time.Time
}
