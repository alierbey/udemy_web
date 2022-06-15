package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Books struct {
	Id              int
	Name            string
	Author          string
	PublicationDate int
	IsStock         bool
	Category        int
}

var tmpl *template.Template

type Book struct {
	Id              int
	Name            string
	Author          string
	PublicationDate int
}

func booksLoad() []Book {
	return []Book{
		Book{1, "Cehennme Övgü", "Gündüz Vassaf", 1992},
		Book{2, "Dublörün Dilemması", "Murat Menteş", 2005},
		Book{3, "Saatleri Ayarlama Enstitüsü", "Ahmet Hamdi Tanpınar", 1961},
	}
}

type Category struct {
	CategoryId   int
	CategoryName string
}

func categoryLoad() []Category {
	return []Category{
		Category{1, "Roman"},
		Category{2, "Felsefe"},
	}
}

func getFiles(p string) []string {

	var files []string

	dizin, err := os.Open(p)
	if err != nil {
		fmt.Println("Dizine ulaşılamıyor..")
		os.Exit(1)
	}

	liste, _ := dizin.Readdirnames(0)

	for _, fileName := range liste {
		files = append(files, "templates/part/"+fileName)
	}

	files = append(files, "templates/base.html")

	return files

}

func main() {

	r := mux.NewRouter()
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	r.PathPrefix("/static/").Handler(s)
	r.HandleFunc("/", index)
	r.HandleFunc("/about", about)
	r.HandleFunc("/books", books)
	r.HandleFunc("/books/{process}/{id}", books)
	r.HandleFunc("/bookadd", bookAdd)
	r.HandleFunc("/detay/{id}", detay)
	r.HandleFunc("/edit/{id}", bookEdit)
	r.HandleFunc("/delete/{id}", bookDelete)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func templateReady(file string) template.Template {

	var files []string

	part := getFiles("templates/part")

	files = append(files, part...)
	files = append(files, "templates/"+file+".html")

	templates := template.Must(template.ParseFiles(files...))

	return *templates

}

func index(w http.ResponseWriter, r *http.Request) {

	tmp := templateReady("index")
	tmp.ExecuteTemplate(w, "base", nil)
}
func about(w http.ResponseWriter, r *http.Request) {

	tmp := templateReady("about")
	tmp.ExecuteTemplate(w, "base", nil)
}

func books(w http.ResponseWriter, r *http.Request) {

	type Data struct {
		Books      []Book
		Categories []Category
	}

	var data Data
	data.Books = booksLoad()
	data.Categories = categoryLoad()

	tmp := templateReady("books")
	tmp.ExecuteTemplate(w, "base", data)

}

func bookAdd(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		name := r.FormValue("name")
		author := r.FormValue("author")
		publication_date := r.FormValue("publication_date")

		fmt.Println("Veriyi Ekle : ", name, author, publication_date)
	}

	tmp := templateReady("bookadd")
	tmp.ExecuteTemplate(w, "base", nil)

}

func detay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	books := booksLoad()

	var k Book
	for _, v := range books {
		if strconv.Itoa(v.Id) == id {
			fmt.Println("Eşleşti :", v)
			k = Book{v.Id, v.Name, v.Author, v.PublicationDate}
		}
	}

	tmp := templateReady("detay")
	tmp.ExecuteTemplate(w, "base", k)

}

func bookEdit(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		vars := mux.Vars(r)
		id := vars["id"]

		books := booksLoad()

		var k Book
		for _, v := range books {
			if strconv.Itoa(v.Id) == id {
				fmt.Println("Eşleşti :", v)
				k = Book{v.Id, v.Name, v.Author, v.PublicationDate}
			}
		}

		tmp := templateReady("bookedit")
		tmp.ExecuteTemplate(w, "base", k)
	} else if r.Method == "POST" {

		id := r.FormValue("id")
		name := r.FormValue("name")
		author := r.FormValue("author")
		publication_date := r.FormValue("publication_date")

		fmt.Println("Veriyi Guncelle : ", id, name, author, publication_date)

		http.Redirect(w, r, "/books", http.StatusSeeOther)

	}

}

func bookDelete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Println(id + "id silindi")

	http.Redirect(w, r, "/books", http.StatusSeeOther)

	tmp := templateReady("bookdelete")
	tmp.ExecuteTemplate(w, "base", nil)

}
