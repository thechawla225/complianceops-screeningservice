package screening

type Entry struct {
	UID            int      `json:"uid"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Programs       []string `json:"programs"`
	AlternateNames []string `json:"alternateNames,omitempty"`
}

type Result struct {
	Verdict       string  `json:"verdict"`
	MatchedEntity *string `json:"matchedEntity"`
	ScreeningRef  string  `json:"screeningRef"`
}

const WatchlistKey = "watchlist:entries"
