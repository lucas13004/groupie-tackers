package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/dgkg/project/model"
)

type ServicesHandler struct {
	tmpl *template.Template
}

type HMap map[string]interface{}

func New(tmpl *template.Template) *ServicesHandler {
	return &ServicesHandler{
		tmpl: tmpl,
	}
}

var artistsPaternID = regexp.MustCompile(`artists/([0-9]+)`) //
var artistsPatern = regexp.MustCompile(`artists`)            //

func (sh *ServicesHandler) Route(w http.ResponseWriter, r *http.Request) {
	switch {
	case artistsPaternID.MatchString(r.URL.Path):
		sh.GetArtist(w, r)
	case artistsPatern.MatchString(r.URL.Path):
		sh.GetAllArtist(w, r)
	default:
		w.Write([]byte("Unknown Pattern"))
	}
}

func (sh *ServicesHandler) GetAllArtist(w http.ResponseWriter, r *http.Request) {

	var listArtists []model.Artist
	url := "https://groupietrackers.herokuapp.com/api/artists"
	err := RequestGet(url, &listArtists)
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	sh.responseOk(w, "list-artist", HMap{
		"title":   "List of artists",
		"artists": listArtists,
	})

}

const (
	apiURL = "https://groupietrackers.herokuapp.com/api/artists/"
)

func (sh *ServicesHandler) GetArtist(w http.ResponseWriter, r *http.Request) {

	url := apiURL + strings.ReplaceAll(r.URL.Path, "/artists/", "")
	var artist model.Artist
	err := RequestGet(url, &artist)
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	var location model.Location
	if len(artist.Locations) != 0 {
		err = RequestGet(artist.Locations, &location)
		if err != nil {
			responseError(w, http.StatusBadRequest, err)
			return
		}
	}

	var date model.Date
	if len(artist.Concertdates) != 0 {
		err = RequestGet(artist.Concertdates, &date)
		if err != nil {
			responseError(w, http.StatusBadRequest, err)
			return
		}
	}

	var relation model.Relation
	if len(artist.Relations) != 0 {
		err = RequestGet(artist.Relations, &relation)
		if err != nil {
			responseError(w, http.StatusBadRequest, err)
			return
		}
	}

	// render location
	renderLocation := "<ul>"
	for _, v := range location.Locations {
		if len(v) != 0 {
			renderLocation += "<li>" + v + "</li>"
		}
	}
	renderLocation += "</ul>"

	// render relation
	renderRelation := "<ul>"
	for k, v := range relation.Dateslocations {
		res := ""
		for i := 0; i < len(v); i++ {
			if len(v[i]) != 0 {
				res += v[i] + ", "
			}
		}
		if len(res) != 0 {
			res = res[:len(res)-2]
		}
		res = res
		renderRelation += fmt.Sprintf("<li>%v - %v </li>", k, res)
	}
	renderRelation += "</ul>"

	sh.responseOk(w, "artist", HMap{
		"title":        artist.Name,
		"image":        artist.Image,
		"name":         artist.Name,
		"members":      artist.Members,
		"creationDate": artist.Creationdate,
		"firstAlbum":   artist.Firstalbum,
		"locations":    template.HTML(renderLocation),
		"concertDates": date.Dates,
		"relations":    template.HTML(renderRelation),
	})
}

func RequestGet(url string, bind interface{}) error {

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, bind)
}

func responseError(w http.ResponseWriter, statusCode int, err error) {
	log.Print(err)
	data, err := json.Marshal(err)
	if err != nil {
		log.Print(err)
		return
	}
	w.Write(data)
	w.WriteHeader(statusCode)
}

func (sh *ServicesHandler) responseOk(w http.ResponseWriter, templateName string, payodad interface{}) {
	w.WriteHeader(http.StatusOK)

	err := sh.tmpl.ExecuteTemplate(w, templateName, payodad)
	if err != nil {
		log.Print(err)
		return
	}
	//(wr io.Writer, name string, data interface{})
}
