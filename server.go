package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"encoding/json"
)

func main() {
	react := http.FileServer(http.Dir("app/build"))
	storage := http.FileServer(http.Dir("music"))
	http.Handle("/", http.StripPrefix("/", react))
	http.Handle("/storage/", http.StripPrefix("/storage/", storage))

	http.HandleFunc("/albums", handleAlbums)
	http.HandleFunc("/songs", handleMusic)
	http.HandleFunc("/upload", handleUpload)
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleAlbums(w http.ResponseWriter, r *http.Request) {
	albums, _ := ioutil.ReadDir("storage/albums")
	albumNames := make([]string, 0)
	for _, album := range albums {
		if album.Name() != "index.html" {
			albumNames = append(albumNames, album.Name());
		}
	}
	j, err := json.Marshal(albumNames)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Fprintf(w, "%v", string(j))
}

func handleMusic(w http.ResponseWriter, r *http.Request) {
	q, _ := url.ParseQuery(r.URL.RawQuery)
	albumName := q["album"][0]
	songs, _ := ioutil.ReadDir("storage/albums/" + albumName)
	songNames := make([]string, 0)
	for _, song := range songs {
		if song.Name() != "index.html"  {
			songNames = append(songNames, song.Name());
		}
	}
	j, err := json.Marshal(songNames)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Fprintf(w, "%v", string(j))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Uploading")
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("music")
	if err != nil {
		fmt.Print(err)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "%v", handler.Header)

	f, err := os.OpenFile("./music/" + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print(err)
		return
	}

	defer f.Close()
	io.Copy(f, file)
}

