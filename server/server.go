package server

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fetch-challenge/models"

	"github.com/google/uuid"
)

var (
	validator           = newReciptValidator()
	reciptsById         = map[uuid.UUID]*models.Receipt{}
	alphanumericPattern = regexp.MustCompile(`[a-zA-Z0-9]`)
)

func Serve() {
	http.HandleFunc("/receipts/process", processReciptsHandler)
	http.HandleFunc("/receipts/{id}/points", getPointsHandler)

	fmt.Println("Server is listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func processReciptsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		badRequest(w)
		return
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w)
		return
	}

	recipt := &models.Receipt{}
	err = json.Unmarshal(buf, recipt)
	if err != nil {
		badRequest(w)
		return
	}

	err = validator.Struct(recipt)
	if err != nil {
		badRequest(w)
		return
	}

	uuid := generateUUID()

	reciptsById[uuid] = recipt

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.RecieptResponse{Id: uuid.String()})
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		notFound(w)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		notFound(w)
		return
	}

	intId, err := uuid.Parse(id)
	if err != nil {
		notFound(w)
		return
	}

	recipt, ok := reciptsById[intId]
	if !ok || recipt == nil {
		notFound(w)
		return
	}

	points := computePoints(recipt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.PointsResponse{Points: points})
}

func generateUUID() uuid.UUID {
	id := uuid.New()
	for _, ok := reciptsById[id]; ok; {
		id = uuid.New()
	}
	return id
}

func computePoints(receipt *models.Receipt) int {

	points := getRetailerPoints(receipt.Retailer)
	points += getTotalIsRoundPoints(receipt.Total)
	points += getTotalIsMultipleTwentyFiveCentsPoints(receipt.Total)
	points += getPointsForEveryTwoItems(receipt.Items)
	points += getDescriptionPoints(receipt.Items)
	points += getTimePoints(receipt.PurchaseTime)
	points += getDatePoints(receipt.PurchaseDate)
	return points
}

func getRetailerPoints(retailer string) int {
	return len(alphanumericPattern.FindAllString(retailer, -1))
}

func getTotalIsRoundPoints(t string) int {
	total, err := strconv.ParseFloat(t, 64)
	if err == nil && math.Floor(total)-total == 0 {
		return 50
	}
	return 0
}

func getTotalIsMultipleTwentyFiveCentsPoints(t string) int {
	total, err := strconv.ParseFloat(t, 64)
	if err == nil && math.Mod(total, 0.25) == 0 {
		return 25
	}
	return 0
}

func getPointsForEveryTwoItems(items []*models.Item) int {
	return 5 * (len(items) / 2)
}

func getDescriptionPoints(items []*models.Item) int {
	points := 0
	for _, item := range items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				itemPoints := math.Ceil(price * 0.2)
				points += int(itemPoints)
			}
		}
	}
	return points
}

func getDatePoints(date string) int {
	if date, err := time.Parse("2006-01-02", date); err == nil && date.Day()%2 != 0 {
		return 6
	}
	return 0
}

func getTimePoints(t string) int {
	if parsedTime, err := time.Parse("15:04", t); err == nil {
		if parsedTime.Hour() > 14 && parsedTime.Hour() < 16 {
			return 10
		}
		if parsedTime.Hour() == 14 && parsedTime.Minute() > 0 {
			return 10
		}
	}
	return 0
}

func badRequest(w http.ResponseWriter) {
	http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
}

func notFound(w http.ResponseWriter) {
	http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
}
