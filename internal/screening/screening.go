type Result struct {
	Verdict       string  `json:"verdict"`
	MatchedEntity *string `json:"matchedEntity"` 
	ScreeningRef  string  `json:"screeningRef"`
}

//Returnning a stubbed response for now
func Verdict(debtorName, creditorName string) Result {
	return Result{
		Verdict:      "clear",
		ScreeningRef: "SCR-STUB-000001",
	}
}