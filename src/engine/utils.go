package engine

import (
	"html/template"
	"log"
	"net/http"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Template(w http.ResponseWriter, name string, data map[interface{}]interface{}) {
	t, err := template.ParseFiles("src/views/T.header.tpl", "src/views/T.navbar.tpl", "src/views/T.footer.tpl", "src/views/"+name)
	CheckError(err)
	err = t.ExecuteTemplate(w, name, data)
	CheckError(err)
}
