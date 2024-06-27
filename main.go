package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Hotel struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Address         string  `json:"address"`
	HotelStarRating float64 `json:"hotelStarRating"`
}

type Room struct {
	ID            string  `json:"id"`
	NumberRoom    string  `json:"numberRoom"`
	RoomCategory  string  `json:"roomCategory"`
	Capacity      int     `json:"capacity"`
	BaseRoomPrice float64 `json:"baseRoomPrice"`
	HotelInfo     Hotel   `json:"hotelInfo,omitempty"`
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

	fmt.Println(hotel)
	hotels[hotel.ID] = hotel

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotel)
}

func updateHotelRating(w http.ResponseWriter, r *http.Request) {
	hotelId := chi.URLParam(r, "idHotel")
	var hotel Hotel
	err := json.NewDecoder(r.Body).Decode(&hotel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingHotel, ok := hotels[hotelId]
	if !ok {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	existingHotel.HotelStarRating = hotel.HotelStarRating
	hotels[hotelId] = existingHotel

	for _, room := range rooms {
		if room.HotelInfo.ID == hotelId {
			room.HotelInfo.HotelStarRating = hotel.HotelStarRating
			rooms[room.ID] = room
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingHotel)
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

	if _, ok := rooms[room.ID]; ok {
		http.Error(w, "Room already exists", http.StatusConflict)
		return
	}

	hotelId := chi.URLParam(r, "idHotel")

	if hotelId == "" {
		http.Error(w, "HotelId is required", http.StatusBadRequest)
		return
	}

	hotel, ok := hotels[hotelId]
	if !ok {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	randomCategory, randomCapacity := getRandomRoomDetails()
	room.RoomCategory = randomCategory
	room.Capacity = randomCapacity

	room.BaseRoomPrice = calculateBaseRoomPrice(room.RoomCategory, hotel.HotelStarRating, room.Capacity, w)
	room.HotelInfo = hotel
	rooms[room.ID] = room

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func calculateBaseRoomPrice(roomCategory string, hotelStarRating float64, capacity int, w http.ResponseWriter) float64 {
	var basePrice float64
	switch roomCategory {
	case "Standard":
		basePrice = 50 * hotelStarRating
	case "JuniorSuite":
		basePrice = 100 * hotelStarRating
	case "Deluxe":
		basePrice = 150 * hotelStarRating
	case "Suite":
		basePrice = 200 * hotelStarRating
	default:
		http.Error(w, "Category is required", http.StatusNotFound)
	}

	if capacity > 2 {
		basePrice += 15 * float64(capacity)
	}

	return math.Ceil(basePrice)
}

func getRandomRoomDetails() (string, int) {
	roomCategories := []string{"Standard", "JuniorSuite", "Deluxe", "Suite"}
	capacities := []int{1, 2, 3, 4, 6}

	randomCategory := roomCategories[rand.Intn(len(roomCategories))]
	randomCapacity := capacities[rand.Intn(len(capacities))]

	return randomCategory, randomCapacity
}
