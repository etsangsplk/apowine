package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/aporeto-inc/apowine/source/mongodb-lib"
)

// Client holds data to connect to the serverß
type Client struct {
	serverAddress string
	drinkName     string
	beer          mongodb.Beer
	realm         string
	validity      string
	midgardToken  string
	wine          mongodb.Wine
}

// GenerateClientPage generates HTML to manipulate data
func GenerateLoginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOGIN PAGE")
	t, err := template.New("login.html").ParseFiles("/Users/sibi/apomux/workspace/code/go/src/github.com/aporeto-inc/apowine/source/frontend-ui/templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 3)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// NewClient creates new client handler
func NewClient(serverAddress string, realm, validity string) *Client {

	return &Client{
		serverAddress: serverAddress,
		validity:      validity,
		realm:         realm,
	}
}

func (c *Client) CatchToken(w http.ResponseWriter, r *http.Request) {

	googleJWT := r.FormValue("idtoken")
	fmt.Println(googleJWT)

	url := "https://api.console.aporeto.com/issue"

	var jsonStr = []byte(fmt.Sprintf(`{"data":"%s","realm":"%s","validity":"%s"}`, googleJWT, c.realm, c.validity))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	dumpReq, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(dumpReq))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var midgardResponse map[string]interface{}
	json.Unmarshal(body, &midgardResponse)
	if midgardResponse["token"] != nil {
		c.midgardToken = midgardResponse["token"].(string)
	} else {
		http.Error(w, "Error from midgard issuing token", http.StatusInternalServerError)
	}
}

// GenerateClientPage generates HTML to manipulate data
func (c *Client) GenerateClientPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("homepage.html").ParseFiles("/Users/sibi/apomux/workspace/code/go/src/github.com/aporeto-inc/apowine/source/frontend-ui/templates/homepage.html")
	if err != nil {
		fmt.Println(err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(c.beer)
		if err != nil {
			http.Error(w, err.Error(), 2)
		}
	} else if c.drinkName == mongodb.WINE {
		c.drinkName = mongodb.WINE
		operation := r.URL.Query().Get("operationType")
		err := c.manipulateData(operation, r, &c.wine, mongodb.WINE)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("id_token", "token")
		err = json.NewEncoder(w).Encode(c.wine)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	reader := bytes.NewReader(data)
	err = json.NewDecoder(reader).Decode(&beer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewDecoder(reader).Decode(&wine)
	w.Header().Set("Content-Type", "application/json")
	if beer.BeerName != "" {
		err = json.NewEncoder(w).Encode(beer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if wine.WineName != "" {
		err = json.NewEncoder(w).Encode(wine)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
