package producerbeer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// PushBeersToDB pushes predefined beers to database
func PushBeersToDB(serverURI string) error {
	zap.L().Info("Reading beers from file")
	zap.L().Info("ServerURI", zap.String("URI", serverURI))

	data, err := ioutil.ReadFile("/apowine/producerbeer.txt")
	if err != nil {
		return err
	}

	strData := string(data)

	beerNames := strings.SplitAfter(strData, ",")

	for _, beerNameWithSymbol := range beerNames {
		beerName := strings.Replace(beerNameWithSymbol, ",", "", -1)
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
