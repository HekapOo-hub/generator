package model

import "encoding/json"

type GeneratedPrice struct {
	Ask    float64 `json:"ask"`
	Bid    float64 `json:"bid"`
	Symbol string  `json:"symbol"`
}

func DecodePrice(data []byte) (GeneratedPrice, error) {
	var msg GeneratedPrice
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return GeneratedPrice{}, err
	}
	return msg, nil
}

func (gp GeneratedPrice) MarshalBinary() ([]byte, error) {
	return json.Marshal(gp)
}
