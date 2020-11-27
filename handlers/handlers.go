package handlers

import (
	"golang-fifa-world-cup-web-service/data"
	"net/http"
)

// RootHandler returns an empty body status code
func RootHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNoContent)
}

// ListWinners returns winners from the list
func ListWinners(res http.ResponseWriter, req *http.Request) {
	// Set the header.
	res.Header().Set("Content-Type", "application/json")

	// Get the year.
	year := req.URL.Query().Get("year")

	// If no year, list all.
	if year == "" {
		// Get all winners.
		winners, err := data.ListAllJSON()

		// Handle errors.
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Respond with winners.
		res.Write(winners)
	} else {
		// Get winners by year.
		filteredWinners, err := data.ListAllByYear(year)

		// Handle errors.
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// Respond with winners by year.
		res.Write(filteredWinners)
	}
}

// AddNewWinner adds new winner to the list
func AddNewWinner(res http.ResponseWriter, req *http.Request) {
	// Get access token and determine its validity.
	accessToken := req.Header.Get("X-ACCESS-TOKEN")
	isTokenValid := data.IsAccessTokenValid(accessToken)

	// Stop if the token is invalid.
	if !isTokenValid {
		res.WriteHeader(http.StatusUnauthorized)
	} else {
		// Add a new winner.
		err := data.AddNewWinner(req.Body)

		// Handle errors.
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		// Respond with success.
		res.WriteHeader(http.StatusCreated)
	}
}

// WinnersHandler is the dispatcher for all /winners URL
func WinnersHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		ListWinners(res, req)
	case http.MethodPost:
		AddNewWinner(res, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}
