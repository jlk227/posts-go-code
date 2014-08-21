package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

//db table post
var posts map[int]string

type Page struct {
	Title string
	Body  []byte
	Data  map[int]string
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	// if it is a number
	if postInx, err := strconv.Atoi(title); err == nil {
		post := posts[postInx]
		if post != "" {
			return &Page{Title: title, Body: []byte(post)}, nil
		}
		return &Page{Title: title, Body: []byte("No Post related to this id")}, nil
	}

	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* Home page Handle URI start with /edit/ */
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

/* Home page Handle URI start with /view/ */
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

/* Home page Handle URI start with /save/ */
func saveHandler(w http.ResponseWriter, r *http.Request) {

	postInx, _ := strconv.Atoi(r.URL.Path[len("/save/"):])
	if posts[postInx] == "" {
		postInx = len(posts) + 1
	}

	body := r.FormValue("body")

	posts[postInx] = body
	http.Redirect(w, r, "/", http.StatusFound)
}

/* Home page Handle URI start with / */
func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", &Page{Title: "home", Data: posts})
}

func main() {
	posts = make(map[int]string)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.ListenAndServe(":8080", nil)
}
