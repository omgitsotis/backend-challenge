package client

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/omgitsotis/backend-challenge/dao"
)

type Client struct {
	dao dao.DAO
}

type ErrorResponse struct {
	Message string `msg`
}

func NewClient(dao dao.DAO) *Client {
	return &Client{dao}
}

func NewRouter(c *Client) *mux.Router {
	r := mux.NewRouter()
	r.Methods("GET").
		Path("/search").
		Queries("searchTerm", "{searchTerm}", "lat", "{lat}", "lng", "{lng}").
		HandlerFunc(c.search)

	return r
}

func ServeAPI(dao dao.DAO) error {
	client := NewClient(dao)
	return http.ListenAndServe(":4000", NewRouter(client))
}

func (c *Client) search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	searchTerm, ok := vars["searchTerm"]
	if !ok {
		log.Println("no search term provided")
		c.writeErrorResponse(w, http.StatusBadRequest, "no search term provided")
		return
	}

	lat, ok := vars["lat"]
	if !ok {
		log.Println("no latitude provided")
		c.writeErrorResponse(w, http.StatusBadRequest, "no latitude provided")
		return
	}

	lng, ok := vars["lng"]
	if !ok {
		log.Println("no longitude provided")
		c.writeErrorResponse(w, http.StatusBadRequest, "no longitude provided")
		return
	}

	log.Printf("%s | %s | %s", searchTerm, lng, lat)

	longitude, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		c.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		c.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var rad dao.Radius

	for i := 1; i < 10; i++ {
		minLat := latitude - (0.01 * math.Pow(2.0, float64(i)))
		maxLat := latitude + (0.01 * math.Pow(2.0, float64(i)))

		minLong := longitude - (0.01 * math.Pow(2.0, float64(i)))
		maxLong := longitude + (0.01 * math.Pow(2.0, float64(i)))

		log.Printf(
			"Range: (%f, %f) - (%f, %f)",
			minLat, minLong, maxLat, maxLong,
		)

		rad = dao.Radius{
			minLat, maxLat,
			minLong, maxLong,
			latitude, longitude,
		}

		count, err := c.dao.CountItems(searchTerm, rad)
		if err != nil {
			c.writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Printf("Found %d item(s)", count)
		if count >= 20 {
			break
		}
	}

	results, err := c.dao.GetItems(searchTerm, rad)
	if err != nil {
		c.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	sort.Sort(dao.Rows(results))

	if len(results) > 20 {
		results = results[:20]
	}

	b, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}

func (c *Client) writeErrorResponse(w http.ResponseWriter, code int, msg string) {
	errorMsg := ErrorResponse{msg}
	b, err := json.Marshal(errorMsg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(b)
}
