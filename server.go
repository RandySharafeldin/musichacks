package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
)

func main() {
	react := http.FileServer(http.Dir("app/build"))
	storage := http.FileServer(http.Dir("music"))
	http.Handle("/", http.StripPrefix("/", react))
	http.Handle("/storage/", http.StripPrefix("/storage/", storage))

	http.HandleFunc("/albums", handleAlbum)
	http.HandleFunc("/songs", handleMusic)
	http.HandleFunc("/upload", handleUpload)
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleAlbum(w http.ResponseWriter, r *http.Request) {
	albums, _ := ioutil.ReadDir("storage/albums")
	albumNames := list.List()
	for _, album := range albums {
		albumNames.PushBack(album.Name());
	}
	j, err := json.Marshal(albumNames)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Fprintf(w, "%v", string(j))
}

func handleMusic(w http.ResponseWriter, r *http.Request) {

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

