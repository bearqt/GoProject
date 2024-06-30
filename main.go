package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Game struct {
	Id          int     `json:"id"`
	BetAmount   float64 `json:"bet_amount"`
	Title       string  `json:"title"`
	RoundNumber int     `json:"round_number"`
	UsersCount  int     `json:"users_count"`
}

var games = make(map[int]Game)

func getGames(w http.ResponseWriter, r *http.Request) {
	gameList := make([]Game, 0, len(games))
	for _, game := range games {
		gameList = append(gameList, game)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameList)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "ID is not present!", http.StatusBadRequest)
		return
	}
	game, ok := games[id]
	if !ok {
		http.Error(w, "Game not found!", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func createGame(w http.ResponseWriter, r *http.Request) {
	var newGame Game
	err := json.NewDecoder(r.Body).Decode(&newGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newGame.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
	}
	_, ok := games[newGame.Id]
	if ok {
		http.Error(w, "Game Id is already present!", http.StatusConflict)
		return
	}
	games[newGame.Id] = newGame

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newGame)
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/games", getGames)
	router.Get("/games/{id}", getGame)
	router.Post("/games", createGame)

	fmt.Println("Starting server at port 8000")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		panic(err)
	}
}
