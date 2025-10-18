// yo this is so freaking cool

package main

import (
	"log"
	"net/http"
	"html/template"
	"sync"
	"strconv"
)

// globals for our form 
var formResult string
var mu sync.Mutex // thread safety 

func main() {
	// server a static http server from files here 
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// server index page 
	// woa woa you can just write a function in a function 
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		mu.Lock()
		data := map[string]string { "Result": formResult }
		mu.Unlock()
		tmpl.Execute(w, data)
	})

	// handle submissions
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return 
		}

		// get values from the form 
		numberStr := r.FormValue("numberInput")
		optionStr := r.FormValue("optionSelect")
		
		var number int
		var to_base int 
		var result string

		// make our option (or base number) a int
		to_base, err = strconv.Atoi(optionStr)
		if err != nil {
			http.Error(w, "Error atoi for base option", http.StatusInternalServerError)
			return
		}

		number, err  = strconv.Atoi(numberStr)
		if err != nil {
			http.Error(w, "Error atoi for number str", http.StatusInternalServerError)
			return
		}

		// convert our number back to a string 
		// ok cool type casting isnt too different
		result = strconv.FormatInt(int64(number), to_base)
		

		log.Println("Number: ", number, " Result: ", result)

		mu.Lock()
		formResult = "Base 10: " + numberStr + " -> Base " + optionStr + ": " + result
		mu.Unlock()


		// load and render result 
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		mu.Lock()
		data := map[string]string { "Result": formResult }
		mu.Unlock()
		tmpl.Execute(w, data)

	})



	// start it on port 2047 
	log.Println("Server running at http://localhost:2047")
	err := http.ListenAndServe(":2047", nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
