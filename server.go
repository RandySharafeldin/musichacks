package musichacks

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"google.golang.org/appengine"
)


func main() {

	http.HandleFunc("/", handleUpload)
	appengine.Main()
}


func handleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("music")
	if err != nil {
		fmt.Print(err)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "%v", handler.Header)

	f, err := os.OpenFile("./music/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print(err)
		return
	}

	defer f.Close()

	io.Copy(f, file)

}
