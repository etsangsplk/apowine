package producerwine

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// PushWinesToDB pushes predefined wines to database
func PushWinesToDB(serverURI string) error {
	zap.L().Info("Reading wines from file")
	zap.L().Info("ServerURI", zap.String("URI", serverURI))

	data, err := ioutil.ReadFile("/apowine/producerwine.txt")
	if err != nil {
		return err
	}

	strData := string(data)

	wineNames := strings.SplitAfter(strData, ",")

	for _, wineNameWithSymbol := range wineNames {
		wineName := strings.Replace(wineNameWithSymbol, ",", "", -1)
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
