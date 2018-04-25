package producerbeer

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// PushBeersToDB pushes predefined beers to database
func PushBeersToDB(serverURI string) error {

	beerNames := []string{
		"Westvleteren 12 (XII)",
		"Rochefort Trappistes 10",
		"Westvleteren Extra 8",
		"Struise Pannepot (10%)",
		"Cantillon Soleil de Minuit",
		"Bokkereyder Framboos Noyaux",
		"3 Fonteinen Hommage",
		"Cantillon Blåbær Lambik",
		"3 Fonteinen Oude Geuze (Cuvée Armand & Gaston)",
		"Struise Pannepot Reserva",
		"St. Bernardus Abt 12",
		"Struise Black Albert",
		"3 Fonteinen Oude Geuze Vintage",
		"Struise Cuvée Delphine",
		"Cantillon Lou Pepe Pure Kriek",
		"Cantillon La Vie est Belge",
		"Cantillon Lambic d’Aunis",
		"Struise Black Damnation V - Double Black",
		"Cantillon Carignan",
		"3 Fonteinen Schaarbeekse Kriek",
		"Rodenbach Alexander",
		"Bokkereyder Framboos Vanille",
		"Struise Pannepot Grand Reserva",
		"3 Fonteinen Oude Geuze Golden Blend",
		"Bokkereyder Perzik",
		"Rodenbach Caractère Rouge",
		"Rochefort Trappistes 8",
		"Goedele’s Bloesem",
		"Cantillon Lou Pepe Framboise",
		"Cantillon Zelige",
		"Bokkereyder Framboos Cognac",
		"Struise Black Damnation I - Black Berry Albert",
		"Cantillon Zwanze (2016) Framboise",
		"3 Fonteinen Oude Geuze Honing",
		"Cantillon Fou’ Foune",
		"Tilquin Oude Pinot Noir à l’Ancienne",
		"De Dolle Oerbier Special Reserva",
		"Bokkereyder Framboos Puur (2016)",
		"Struise Black Damnation IV - Coffee Club",
		"Abbaye des Rocs Brune",
	}

	for _, beerName := range beerNames {
		var values map[string]string
		values = map[string]string{"beername": beerName}
		jsonValue, _ := json.Marshal(values)
		_, err := http.Post(serverURI, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			return err
		}
	}

	return nil
}