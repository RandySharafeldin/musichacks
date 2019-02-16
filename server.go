package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)


func main() {
	frontend := http.FileServer(http.Dir("static"))
	musicfs := http.FileServer(http.Dir("music"))
	http.Handle("/static/", http.StripPrefix("/static/", frontend))
	http.Handle("/music/", http.StripPrefix("/music/", musicfs))  
	http.HandleFunc("/upload", handleUpload)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
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

