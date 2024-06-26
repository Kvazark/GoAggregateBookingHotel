package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Hotel struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Address         string  `json:"address"`
	HotelStarRating float32 `json:"hotelStarRating "`
}

type Room struct {
	ID           string  `json:"id"`
	NumberRoom   string  `json:"numberRoom"`
	RoomCategory string  `json:"roomCategory "`
	Capacity     int     `json:"capacity"`
	HotelId      *string `json:"hotel"`
}

var hotels = make(map[string]Hotel)
var rooms = make(map[string]Room)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/addNewHotel", addNewHotel)
	r.Post("/updateHotelStarRating/{idHotel}", updateHotelRating)
	r.Post("/addNewRoom/{idHotel}", addNewRoom)
	r.Get("/hotels", listHotels)
	r.Get("/rooms", listRooms)

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}

func listHotels(w http.ResponseWriter, _ *http.Request) {
	hotelList := make([]Hotel, 0, len(hotels))
	for _, hotel := range hotels {
		hotelList = append(hotelList, hotel)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotelList)
}

func addNewHotel(w http.ResponseWriter, r *http.Request) {
	var hotel Hotel
	err := json.NewDecoder(r.Body).Decode(&hotel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if hotel.ID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	_, ok := hotels[hotel.ID]
	if ok {
		http.Error(w, "Hotel already exists", http.StatusConflict)
		return
	}

	hotels[hotel.ID] = hotel

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotel)
}

func updateHotelRating(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var hotel Hotel
	err := json.NewDecoder(r.Body).Decode(&hotel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hotel, ok := hotels[idStr]

	if !ok {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	hotels[idStr] = hotel

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotel)
}

func listRooms(w http.ResponseWriter, _ *http.Request) {
	roomList := make([]Room, 0, len(rooms))
	for _, room := range rooms {
		roomList = append(roomList, room)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roomList)
}

func addNewRoom(w http.ResponseWriter, r *http.Request) {
	var room Room
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if room.ID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	_, ok := rooms[room.ID]
	if ok {
		http.Error(w, "Hotel already exists", http.StatusConflict)
		return
	}

	rooms[room.ID] = room

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}
