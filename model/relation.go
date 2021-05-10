package model

type Relation struct {
	ID             int                 `json:"id"`
	Dateslocations map[string][]string `json:"datesLocations"`
}
