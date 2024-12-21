package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		http.Error(w, "Failed to open file for appending", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Failed to stat file:", err)
		return
	}
	size := fileInfo.Size()

	if size == 0 {
		file.Write([]byte("["))
	} else {
		_, err = file.Seek(size-1, 0)
		if err != nil {
			fmt.Println("Failed to seek to last byte:", err)
			return
		}
		if _, err = file.Write([]byte(",")); err != nil {
			fmt.Println("Failed to write new character:", err)
			return
		}
	}

	if _, err := file.Write(append(body, ']')); err != nil {
		http.Error(w, "Failed to write to file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Data appended successfully.")
}

func handleView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("view.html")

	if err != nil {
		http.Error(w, "Failed to load HTML template", http.StatusInternalServerError)
		return
	}

	var aspirasi []Aspirasi

	if _, err := os.Stat(fileName); err == nil {
		file, err := os.Open(fileName)
		if err != nil {
			http.Error(w, "Failed to open existing file", http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(file).Decode(&aspirasi); err != nil && err != io.EOF {
			http.Error(w, "Failed to parse existing JSON file", http.StatusInternalServerError)
			file.Close()
			return
		}
	}

	t.Execute(w, aspirasi)
}

func main() {
	http.HandleFunc("/Aspirasi/post", handlePost)
	http.HandleFunc("/Aspirasi/view", handleView)

	fmt.Println("Listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
