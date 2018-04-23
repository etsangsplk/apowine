package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aporeto-inc/apowine/source/frontend-ui/templates"
	"github.com/aporeto-inc/apowine/source/mongodb-lib"
)

// GenerateClientPage generates HTML to manipulate data
func GenerateClientPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("template").Parse(templates.UITemplate)
	if err != nil {
		fmt.Println(err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 1)
	}

	w.Header().Set("Content-Type", "text/html")
}

// AllDrinks returns drinks based on type in JSON format
func GenerateBeerManipulator(w http.ResponseWriter, r *http.Request, serverIP string, serverPort string) {
	defer r.Body.Close()
	var beer mongodb.Beer
	endpoint := strings.SplitAfter(r.URL.RequestURI(), "/")
	param := strings.SplitAfter(endpoint[1], "?")
	drinkName := strings.Replace(param[0], "?", "", -1)
	switch r.URL.Query().Get("type") {
	case "random":
		response, _ := http.Get(serverIP + serverPort + "/" + drinkName + "/random")
		ioioi, _ := ioutil.ReadAll(response.Body)
		re := bytes.NewReader(ioioi)
		json.NewDecoder(re).Decode(&beer)
	case "create":
		//response, _ := http.Post(serverIP + serverPort + "/" + drinkName)

	case "read":
		response, _ := http.Get(serverIP + serverPort + "/" + drinkName + "/random")
		ioioi, _ := ioutil.ReadAll(response.Body)
		re := bytes.NewReader(ioioi)
		json.NewDecoder(re).Decode(&beer)
	case "update":
		response, _ := http.Get(serverIP + serverPort + "/" + drinkName + "/random")
		ioioi, _ := ioutil.ReadAll(response.Body)
		re := bytes.NewReader(ioioi)
		json.NewDecoder(re).Decode(&beer)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(beer)
}

func GenerateWineManipulator(w http.ResponseWriter, r *http.Request, serverIP string, serverPort string) {
	defer r.Body.Close()

	defer r.Body.Close()
	var wine mongodb.Wine
	endpoint := strings.SplitAfter(r.URL.RequestURI(), "/")
	param := strings.SplitAfter(endpoint[1], "?")
	drinkName := strings.Replace(param[0], "?", "", -1)
	switch r.URL.Query().Get("type") {
	case "random":
		response, _ := http.Get(serverIP + serverPort + "/" + drinkName + "/random")
		ioioi, _ := ioutil.ReadAll(response.Body)
		re := bytes.NewReader(ioioi)
		json.NewDecoder(re).Decode(&wine)
	}
}
