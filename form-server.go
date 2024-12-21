package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Aspirasi struct {
	Nama       string `json:"Nama"`
	Kelas      string `json:"Kelas"`
	FasilitasA string `json:"FasilitasA"`
	FasilitasB string `json:"FasilitasB"`
	FasilitasC string `json:"FasilitasC"`
	KBMA       string `json:"KBMA"`
	KBMB       string `json:"KBMB"`
	KinerjaA   string `json:"KinerjaA"`
	KinerjaB   string `json:"KinerjaB"`
	KinerjaC   string `json:"KinerjaC"`
	EkskulA    string `json:"EkskulA"`
	EkskulB    string `json:"EkskulB"`
}

const fileName = "Aspirasi.json"

var fileMutex sync.Mutex

func handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Post Request")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Read the request body.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 2. Parse the incoming JSON into a struct (optional if you don’t need validation).
	var newAspirasi Aspirasi
	if err := json.Unmarshal(body, &newAspirasi); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fileMutex.Lock()
	defer fileMutex.Unlock()

	// 3. Read the existing file (if any) and unmarshal into a slice.
	existingData := []Aspirasi{}
	if _, err := os.Stat(fileName); err == nil {
		// File exists, read it
		file, err := os.Open(fileName)
		if err != nil {
			http.Error(w, "Failed to open existing file", http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(file).Decode(&existingData); err != nil && err != io.EOF {
			http.Error(w, "Failed to parse existing JSON file", http.StatusInternalServerError)
			file.Close()
			return
		}
		file.Close()
	}

	// 4. Append the new object to the slice.
	existingData = append(existingData, newAspirasi)

	// 5. Write the updated array back to the file.
	file, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Failed to create output file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print
	if err := encoder.Encode(existingData); err != nil {
		http.Error(w, "Failed to write JSON to file", http.StatusInternalServerError)
		return
	}

	// Respond with success
	fmt.Fprintln(w, "Data appended successfully.")
}

func handleView(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/Aspirasi/post", handlePost)
	http.HandleFunc("/Aspirasi/view", handleView)

	fmt.Println("Listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
