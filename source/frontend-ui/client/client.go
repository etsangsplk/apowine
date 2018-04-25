package client

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/aporeto-inc/apowine/source/frontend-ui/templates"
	"github.com/aporeto-inc/apowine/source/mongodb-lib"
)

// Client holds data to connect to the serverß
type Client struct {
	serverAddress string
	drinkName     string
	beer          mongodb.Beer
	wine          mongodb.Wine
}

// GenerateClientPage generates HTML to manipulate data
func GenerateClientPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("template").Parse(templates.UITemplate)
	if err != nil {
		http.Error(w, err.Error(), 2)
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 3)
	}

	w.Header().Set("Content-Type", "text/html")
}

// NewClient creates new client handler
func NewClient(serverAddress string) *Client {

	return &Client{
		serverAddress: serverAddress,
	}
}

// GenerateDrinkManipulator returns drinks based on type in JSON format
func (c *Client) GenerateDrinkManipulator(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	c.drinkName = r.URL.Query().Get("drinkType")
	if c.drinkName == mongodb.BEER {
		c.drinkName = mongodb.BEER
		operation := r.URL.Query().Get("operationType")
		err := c.manipulateData(operation, r, &c.beer, mongodb.BEER)
		if err != nil {
			http.Error(w, err.Error(), 2)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(c.beer)
		if err != nil {
			http.Error(w, err.Error(), 2)
		}
	} else if c.drinkName == mongodb.WINE {
		c.drinkName = mongodb.WINE
		operation := r.URL.Query().Get("operationType")
		err := c.manipulateData(operation, r, &c.wine, mongodb.WINE)
		if err != nil {
			http.Error(w, err.Error(), 2)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(c.wine)
		if err != nil {
			http.Error(w, err.Error(), 3)
		}
	}
}

func (c *Client) manipulateData(operation string, r *http.Request, drinkTypeData interface{}, drinkType string) error {

	switch operation {
	case "random":
		response, err := http.Get(c.serverAddress + "/" + c.drinkName + "/random")
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		reader := bytes.NewReader(data)
		err = json.NewDecoder(reader).Decode(drinkTypeData)
		if err != nil {
			return err
		}
	case "create":
		var values map[string]string
		name := r.URL.Query().Get("name")
		if drinkType == mongodb.BEER {
			values = map[string]string{"beername": name}
		} else {
			values = map[string]string{"winename": name}
		}
		jsonValue, err := json.Marshal(values)
		if err != nil {
			return err
		}
		_, err = http.Post(c.serverAddress+"/"+c.drinkName, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			return err
		}
	case "read":
		id := r.URL.Query().Get("id")
		response, err := http.Get(c.serverAddress + "/" + c.drinkName + "/" + id)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		reader := bytes.NewReader(data)
		err = json.NewDecoder(reader).Decode(drinkTypeData)
		if err != nil {
			return err
		}
	case "update":
		var values map[string]string
		name := r.URL.Query().Get("name")
		id := r.URL.Query().Get("id")
		if drinkType == mongodb.BEER {
			values = map[string]string{"id": id, "beername": name}
		} else {
			values = map[string]string{"id": id, "winename": name}
		}
		jsonValue, err := json.Marshal(values)
		if err != nil {
			return err
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, c.serverAddress+"/"+c.drinkName, bytes.NewBuffer(jsonValue))
		if err != nil {
			return err
		}
		_, err = client.Do(req)
		if err != nil {
			return err
		}
	case "delete":
		id := r.URL.Query().Get("id")
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodDelete, c.serverAddress+"/"+c.drinkName+"/"+id, nil)
		if err != nil {
			return err
		}
		_, err = client.Do(req)
		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateRandomDrinkManipulator generates random drinks
func (c *Client) GenerateRandomDrinkManipulator(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var beer mongodb.Beer
	var wine mongodb.Wine

	response, err := http.Get(c.serverAddress + "/random")
	if err != nil {
		http.Error(w, err.Error(), 2)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), 3)
	}
	reader := bytes.NewReader(data)
	err = json.NewDecoder(reader).Decode(&beer)
	if err != nil {
		http.Error(w, err.Error(), 4)
	}

	json.NewDecoder(reader).Decode(&wine)
	w.Header().Set("Content-Type", "application/json")
	if beer.BeerName != "" {
		err = json.NewEncoder(w).Encode(beer)
		if err != nil {
			http.Error(w, err.Error(), 5)
		}
	} else if wine.WineName != "" {
		err = json.NewEncoder(w).Encode(wine)
		if err != nil {
			http.Error(w, err.Error(), 6)
		}
	}
}
