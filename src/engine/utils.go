package engine

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
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

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func SaveToFile(r *http.Request, fromfile, tofile string) error {
	file, _, err := r.FormFile(fromfile)
	if err != nil {
		return err
	}
	defer file.Close()
	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	return err
}

func CountChar(str string) int {
	cnt := 0
	for i := 0; i < len(str); i++ {
		if str[i] == ' ' {
			cnt++
		} else {
			break
		}
	}
	return cnt
}
