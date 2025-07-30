package templates

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/notioncodes/types"
)

func main() {
	f, err := os.Open("data/page.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var page types.Page
	err = json.NewDecoder(f).Decode(&page)
	if err != nil {
		log.Fatal(err)
	}

	tmpls, _ := template.ParseFiles(filepath.Join("templates", "post.gotmpl"))
	tmpl := tmpls.Lookup("post.gotmpl")
	tmpl.Execute(os.Stdout, page)
}
