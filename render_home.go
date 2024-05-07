package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, data interface{}) {
	htmlString := `
	<html>
    <head>
        <h1>
            Hi this is MuhammadDev Backend.
          
        </h1>
        
    </head>



</p>
</html>
	`
	t, err := template.New("index").Parse(htmlString)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error passing html file err : %v", err))
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error exucting html  err : %v", err))

	}
}

func renderHome(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "")
}
