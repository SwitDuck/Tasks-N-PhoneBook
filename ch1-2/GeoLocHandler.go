package main

import (
	"fmt"
	"log"
	"net/http"
	//"strings"
	"encoding/json"
	"encoding/csv"
	"path/filepath"
	"os"
	"time"
)

type GPSData struct {
    User	string	`json:"user1"`
    CurrentTime	string	`json:"current_time"`
    GPGGA	string	`json:"gpgga"`
    GPRMC	string	`json:"gprmc"`
}

var last_data GPSData

func processGPS(data GPSData) error {
    currentDir, err := os.Getwd()
    if err != nil {
        fmt.Printf("Ошибка получения текущего каталога: %v\n", err)
        return err
    }
    
    csvPath := filepath.Join(currentDir, "track_data.csv")
    
    file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Ошибка открытия файла: %v\n", err)
        return err
    }
    defer file.Close()
    
    writer := csv.NewWriter(file)
    defer writer.Flush()
    
    fileInfo, _ := file.Stat()
    if fileInfo.Size() == 0 {
        headers := []string{"user", "current_time", "gpgga", "gprmc"}
        if err := writer.Write(headers); err != nil {
            return err
        }
        writer.Flush()
    }
    
    row := []string{data.User, data.CurrentTime, data.GPGGA, data.GPRMC}
    if err := writer.Write(row); err != nil {
        return err
    }
    
    fmt.Printf("Данные сохранены в CSV: %s\n", csvPath)
    return nil
}
func handleGPS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var data GPSData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Получены данные - Пользователь: %s, Время: %s", data.User, data.CurrentTime)
	log.Printf("GPGGA: %s", data.GPGGA[:min(30, len(data.GPGGA))])
	log.Printf("GPRMC: %s", data.GPRMC[:min(30, len(data.GPRMC))])

	if err := processGPS(data); err != nil {
		log.Printf("Ошибка сохранения CSV: %v", err)
		http.Error(w, "Failed to save to CSV", http.StatusInternalServerError)
		return
	}

	last_data = data

	response := map[string]string{"status": "processed", "message": "GPS data saved"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is:"
	fmt.Fprintf(w, "<h1 align=\"center\">%s</h1>", Body)
	fmt.Fprintf(w, "<h2 align=\"center\">%s</h2>\n", t)
	fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Printf("Served time for: %s\n", r.Host)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
}