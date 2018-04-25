package producerwine

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// PushWinesToDB pushes predefined wines to database
func PushWinesToDB(serverURI string) error {

	wineNames := []string{
		"Bin 707 Cabernet Sauvignon",
		"Caymus Vineyards Cabernet Sauvignon",
		"Dom Pérignon",
		"Échezeaux Grand Cru",
		"Gaja Barbaresco DOCG",
		"Bokkereyder Framboos Noyaux",
		"Hill of Grace Shiraz (Henschke)",
		"Haut-Brion (Château)",
		"Insignia (Joseph Phelps Vineyards)",
		"Klein Constantia Vin de Constance",
		"Lafite Rothschild (Château)",
		"Dr. L Riesling (Loosen Bros)",
		"Margaux (Château)",
		"Ornellaia Bolgheri Superiore",
		"Palmer (Château)",
		"Quinta do Noval ‘Nacional’ Vintage Port",
	}

	for _, wineName := range wineNames {
		var values map[string]string
		values = map[string]string{"winename": wineName}
		jsonValue, _ := json.Marshal(values)
		_, err := http.Post(serverURI, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			return err
		}
	}
	return nil
}
