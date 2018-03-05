package client

import (
    "github.com/omgitsotis/backend-challenge/dao"
    "github.com/gorilla/mux"
    "net/http"
    "fmt"
    "math"
    "strconv"
)

type Client struct {
    dao dao.DAO
}

func NewClient(dao dao.DAO) *Client {
    return &Client{dao}
}

func ServeAPI(dao dao.DAO) error {
    client := NewClient(dao)
    r := mux.NewRouter()
    r.Methods("GET").
        Path("/search").
        Queries("searchTerm", "{searchTerm}", "lat", "{lat}", "lng", "{lng}").
        HandlerFunc(client.search)

	return http.ListenAndServe(":4000", r)
}

func (c *Client) search(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    searchTerm := vars["searchTerm"]
    lng := vars["lng"]
    lat := vars["lat"]

    fmt.Printf("%s | %s | %s\n", searchTerm, lng, lat)

    longitude, err := strconv.ParseFloat(lng, 64)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    latitude, err := strconv.ParseFloat(lat, 64)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    for i:= 1; i < 10; i++{
        minLat := latitude - (0.01 * math.Pow(2.0, float64(i)))
        maxLat := latitude + (0.01 * math.Pow(2.0, float64(i)))

        minLong := longitude - (0.01 * math.Pow(2.0, float64(i)))
        maxLong := longitude + (0.01 * math.Pow(2.0, float64(i)))

        r := dao.Radius{minLat, maxLat, minLong, maxLong}

        count, err := c.dao.GetItemCount(searchTerm, r)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        fmt.Println(count)
    }

    fmt.Fprintln(w, "done")
}
