package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/omgitsotis/backend-challenge/dao"
)

type MockDAO struct {
	Rows dao.Rows
}

func (m *MockDAO) GetItems(term string, r dao.Radius) ([]dao.Row, error) {
	if term == "camera" {
		if r.MaxLatitude > 7 {
			return m.Rows[:20], nil
		}

		return m.Rows[:15], nil
	}

	if term == "canon" {
		return m.Rows[0:9], nil
	}

	if term == "nikon" {
		return m.Rows[10:21], nil
	}

	return []dao.Row{}, errors.New("Bad request")
}

func (sql *MockDAO) CountItems(term string, r dao.Radius) (int, error) {
	if term == "camera" {
		if r.MaxLatitude > 7 {
			return 21, nil
		}

		return 15, nil
	}

	if term == "nikon" {
		return 12, nil
	}

	if term == "canon" {
		return 9, nil
	}

	return 0, errors.New("Bad request")
}

func TestIncorrectQueryParams(t *testing.T) {
	c := setupClient()

	req, err := http.NewRequest("GET", "/search?searchTerm=camera&lat=51.412&lng=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	NewRouter(c).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Incorrect status code: expected [%d], got [%d]",
			http.StatusBadRequest, rr.Code)
	}

	req, err = http.NewRequest("GET", "/search?searchTerm=camera&lat=test&lng=51.412", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	NewRouter(c).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Incorrect status code: expected [%d], got [%d]",
			http.StatusBadRequest, rr.Code)
	}
}

func TestSearch(t *testing.T) {
	c := setupClient()

	req, err := http.NewRequest("GET", "/search?searchTerm=camera&lat=6.82&lng=6.43", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	NewRouter(c).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Incorrect status code: expected [%d], got [%d]",
			http.StatusOK, rr.Code)
	}

	var results dao.Rows
	if err = json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}

	if len(results) != 20 {
		t.Fatalf("Incorrect number of items: expected [%d], got [%d]",
			20, len(results))
	}

	if results[0].ItemID != "canon camera 1" {
		t.Fatalf("Incorrect item: expected [%s], got [%s]",
			"canon camera 1", results[0].ItemID)
	}

	if results[19].ItemID != "nikon camera 4" {
		t.Fatalf("Incorrect item: expected [%s], got [%s]",
			"nikon camera 4", results[19].ItemID)
	}
}

func TestSearchBelowTwenty(t *testing.T) {
	c := setupClient()

	req, err := http.NewRequest("GET", "/search?searchTerm=canon&lat=6.00&lng=6.00", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	NewRouter(c).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Incorrect status code: expected [%d], got [%d]",
			http.StatusOK, rr.Code)
	}

	var results dao.Rows
	if err = json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}

	if len(results) != 9 {
		t.Fatalf("Incorrect number of items: expected [%d], got [%d]",
			9, len(results))
	}

	if results[0].ItemID != "canon camera 1" {
		t.Fatalf("Incorrect item: expected [%s], got [%s]",
			"canon camera 1", results[0].ItemID)
	}

	if results[8].ItemID != "canon camera 5" {
		t.Fatalf("Incorrect item: expected [%s], got [%s]",
			"nikon camera 4", results[8].ItemID)
	}
}

func TestSearchOutOfRange(t *testing.T) {
	c := setupClient()

	req, err := http.NewRequest("GET", "/search?searchTerm=camera&lat=0.00&lng=0.00", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	NewRouter(c).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Incorrect status code: expected [%d], got [%d]",
			http.StatusOK, rr.Code)
	}

	var results dao.Rows
	if err = json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}

	if len(results) != 15 {
		t.Fatalf("Incorrect number of items: expected [%d], got [%d]",
			15, len(results))
	}

	// The closest out of the first 15 items in the array is canon camera 4
	if results[0].ItemID != "canon camera 1" {
		t.Fatalf("Incorrect item: expected [%s], got [%s]",
			"canon camera 1", results[0].ItemID)
	}

	// The furthest away out of the first 15 items in the array is nikon camera 4
	if results[14].ItemID != "nikon camera 4" {
		t.Fatalf("Incorrect item: expected [%s], got [%s]",
			"nikon camera 4", results[14].ItemID)
	}
}

func TestSearchBadRequest(t *testing.T) {
	c := setupClient()

	req, err := http.NewRequest("GET", "/search?searchTerm=television&lat=0.00&lng=0.00", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	NewRouter(c).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Incorrect status code: expected [%d], got [%d]",
			http.StatusBadRequest, rr.Code)
	}
}

func setupClient() *Client {
	var mock MockDAO
	rows := []dao.Row{
		dao.Row{"canon camera 1", 5.00, 5.00, "canon.camera", []string{"camera.jpg"}, 2},
		dao.Row{"canon camera 2", 5.50, 5.50, "canon.camera", []string{"camera.jpg"}, 2.5},
		dao.Row{"canon camera 3", 6.00, 6.00, "canon.camera", []string{"camera.jpg"}, 3},
		dao.Row{"canon camera 4", 6.10, 6.10, "canon.camera", []string{"camera.jpg"}, 3.1},
		dao.Row{"canon camera 5", 6.60, 6.60, "canon.camera", []string{"camera.jpg"}, 3.6},
		dao.Row{"canon camera 6", 6.20, 6.20, "canon.camera", []string{"camera.jpg"}, 3.2},
		dao.Row{"canon camera 7", 6.30, 6.30, "canon.camera", []string{"camera.jpg"}, 3.3},
		dao.Row{"canon camera 8", 6.40, 6.40, "canon.camera", []string{"camera.jpg"}, 3.4},
		dao.Row{"canon camera 9", 6.50, 6.50, "canon.camera", []string{"camera.jpg"}, 3.5},
		dao.Row{"nikon camera 1", 6.70, 6.70, "nikon.camera", []string{"camera.jpg"}, 3.7},
		dao.Row{"nikon camera 2", 6.80, 6.80, "nikon.camera", []string{"camera.jpg"}, 3.8},
		dao.Row{"nikon camera 3", 6.90, 6.90, "nikon.camera", []string{"camera.jpg"}, 3.9},
		dao.Row{"nikon camera 4", 6.91, 6.91, "nikon.camera", []string{"camera.jpg"}, 3.91},
		dao.Row{"nikon camera 5", 5.90, 5.90, "nikon.camera", []string{"camera.jpg"}, 2.9},
		dao.Row{"nikon camera 6", 5.91, 5.91, "nikon.camera", []string{"camera.jpg"}, 2.91},
		dao.Row{"nikon camera 7", 5.92, 5.92, "nikon.camera", []string{"camera.jpg"}, 2.92},
		dao.Row{"nikon camera 8", 5.93, 5.93, "nikon.camera", []string{"camera.jpg"}, 2.93},
		dao.Row{"nikon camera 9", 5.94, 5.94, "nikon.camera", []string{"camera.jpg"}, 2.94},
		dao.Row{"nikon camera 10", 5.95, 5.95, "nikon.camera", []string{"camera.jpg"}, 2.95},
		dao.Row{"nikon camera 11", 5.96, 5.96, "nikon.camera", []string{"camera.jpg"}, 2.96},
		dao.Row{"nikon camera 12", 5.97, 5.97, "nikon.camera", []string{"camera.jpg"}, 2.97},
		dao.Row{"samsung television", 6.40, 6.40, "samsung.tv", []string{"tv.jpg"}, 3.40},
	}
	mock.Rows = rows
	return NewClient(&mock)
}
