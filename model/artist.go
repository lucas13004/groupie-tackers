package model

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	Creationdate int      `json:"creationDate"`
	Firstalbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	Concertdates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}
