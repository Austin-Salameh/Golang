package client

import (
	"encoding/json"
	"log"
	"net/http"

	"git.mobilestar.us/MobileStar/msstatus/html-css/crypto/model"
)

func FetchCrypto(fiat string, crypto string) (string, error) {
	URL := "https://api.nomics.com/v1/currencies/ticker?key=3990ec554a414b59dd85d29b2286dd85&interval=1d&ids=" + crypto + "&convert=" + fiat

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var cResp model.Cryptoresponse

	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Fatalln(err)
	}
	return cResp.TextOutput(), nil
}
