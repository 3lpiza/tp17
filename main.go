package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type PriceResponse struct {
	Time       TimeData `json:"time"`
	Disclaimer string   `json:"disclaimer"`
	BPI        BPIData  `json:"bpi"`
}

type TimeData struct {
	Updated    string `json:"updated"`
	UpdatedISO string `json:"updatedISO"`
	UpdatedUK  string `json:"updateduk"`
}

type BPIData struct {
	USD Currency `json:"USD"`
	BTC Currency `json:"BTC"`
}

type Currency struct {
	Code        string  `json:"code"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	RateFloat   float64 `json:"rate_float"`
}

func main() {
	url := "https://api.coindesk.com/v1/bpi/currentprice/BTC.json"

	var maxPrice, minPrice float64
	maxPrice, minPrice = loadMaxMinFromCSV("btc_prices.csv")

	for {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error al hacer la solicitud HTTP:", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error al leer la respuesta HTTP:", err)
			return
		}

		var priceData PriceResponse
		err = json.Unmarshal(body, &priceData)
		if err != nil {
			fmt.Println("Error al decodificar el JSON:", err)
			return
		}

		btcPrice := priceData.BPI.USD.RateFloat
		timestamp := time.Now()

		//TODO Imprimir el precio del Bitcon

		if btcPrice > maxPrice {
			//TODO si el precio supera al maximo registrado
			//actualizar el maximo
			//e imprimir por pantalla una alerta
		}

		if btcPrice < minPrice {
			//TODO si el precio es inferior al minimo registrado
			//actualizar el minimo
			//e imprimir por pantalla una alerta
		}

		if err := updateCSV(btcPrice, timestamp); err != nil {
			fmt.Println("Error al actualizar el archivo CSV:", err)
		}

		time.Sleep(time.Minute) // Esperar un minuto antes de verificar nuevamente
	}
}

func loadMaxMinFromCSV(fileName string) (float64, float64) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 1e20
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil || len(records) <= 1 {
		return 0, 1e20
	}

	maxPrice, _ := strconv.ParseFloat(strings.TrimSpace(records[len(records)-1][1]), 64)
	minPrice, _ := strconv.ParseFloat(strings.TrimSpace(records[len(records)-1][1]), 64)

	for _, record := range records {
		price, _ := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		if price > maxPrice {
			maxPrice = price
		}
		if price < minPrice {
			minPrice = price
		}
	}

	return maxPrice, minPrice
}

func updateCSV(price float64, timestamp time.Time) error {
	fileName := "btc_prices.csv"

	// Si el archivo no existe, creamos el encabezado
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		writer.Write([]string{"timestamp", "price"})
		writer.Flush()
	}

	// Abrimos el archivo en modo append
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{timestamp.Format(time.RFC3339), fmt.Sprintf("%.4f", price)})
	if err != nil {
		return err
	}

	return nil
}
