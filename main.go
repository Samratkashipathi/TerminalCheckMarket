package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"encoding/json"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Config struct {
	CryptoStock             []string `yaml:"CRYPTO_STOCK"`
	CoinApiKey              string   `yaml:"COIN_API_KEY"`
	CoinApiExchangeCurrency string   `yaml:"COIN_API_EXCHANGE_CURRENCY"`
}

type SpecificRate struct {
	Time string
	Rate float64
}

type TimeSeriesData struct {
	TimePeriodStart string
	TimePeriodEnd   string
	RateClose       float64
}

func (c *Config) getConfig() *Config {

	os.Setenv("TERMINAL_CHECK_MARKET_CONFIG_PATH", "config.yaml")

	configFilePath := os.Getenv("TERMINAL_CHECK_MARKET_CONFIG_PATH")
	if len(configFilePath) == 0 {
		os.Exit(1)
	}

	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Error in getting yaml file path", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println("Error in reading yaml file")
		os.Exit(1)
	}

	return c
}

func getCurrentPrice(config Config, crypto string) (*SpecificRate, error) {
	client := &http.Client{}

	url := "https://rest.coinapi.io/v1/exchangerate/" + crypto + "/" + config.CoinApiExchangeCurrency
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	request.Header.Set("X-CoinAPI-Key", config.CoinApiKey)

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	bodyString := string(bodyBytes)
	jsonData := SpecificRate{}

	err = json.Unmarshal([]byte(bodyString), &jsonData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &jsonData, nil
}

func getOneDayData(config Config, crypto string) (*[]TimeSeriesData, error) {
	client := &http.Client{}

	url := "https://rest.coinapi.io/v1/exchangerate/" + crypto + "/" + config.CoinApiExchangeCurrency +
		"/history?period_id=1HRS&time_start=" + fmt.Sprintf(time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05")) + "&time_end" + fmt.Sprintf(time.Now().Format("2006-01-02T15:04:05"))

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	request.Header.Set("X-CoinAPI-Key", config.CoinApiKey)

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	bodyString := string(bodyBytes)
	jsonData := []TimeSeriesData{}

	err = json.Unmarshal([]byte(bodyString), &jsonData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &jsonData, nil
}

func main() {

	var config Config
	config.getConfig()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	cryptoWidget := widgets.NewList()
	cryptoWidget.Title = "Crypto Watchlist"
	cryptoWidget.Rows = config.CryptoStock
	cryptoWidget.TextStyle = ui.NewStyle(ui.ColorYellow)
	cryptoWidget.WrapText = false
	cryptoWidget.SetRect(110, 0, 140, 17)

	headingWidget := widgets.NewParagraph()
	headingWidget.Title = "Name of the Crypto/Stock"
	headingWidget.Text = ""
	headingWidget.SetRect(0, 0, 100, 5)

	sinData := func() [][]float64 {
		n := 220
		data := make([][]float64, 2)
		data[0] = make([]float64, n)
		data[1] = make([]float64, n)
		for i := 0; i < n; i++ {
			data[0][i] = 1 + math.Sin(float64(i)/5)
			data[1][i] = 1 + math.Cos(float64(i)/5)
		}
		return data
	}()

	graphWidget := widgets.NewPlot()
	graphWidget.Title = "1 Day Change"
	graphWidget.Marker = widgets.MarkerDot
	graphWidget.Data = make([][]float64, 2)
	graphWidget.Data[0] = []float64{1, 2, 3, 4, 5}
	graphWidget.Data[1] = sinData[1][4:]
	graphWidget.SetRect(0, 12, 100, 40)
	graphWidget.AxesColor = ui.ColorWhite
	graphWidget.LineColors[0] = ui.ColorCyan
	graphWidget.PlotType = widgets.ScatterPlot

	errorWidget := widgets.NewParagraph()
	errorWidget.Title = "Error"
	errorWidget.Text = ""
	errorWidget.SetRect(0, 5, 100, 10)

	ui.Render(headingWidget, graphWidget, cryptoWidget, errorWidget)

	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "Down":
			cryptoWidget.ScrollDown()
		case "k", "<Up>":
			cryptoWidget.ScrollUp()
		}

		selectedCrypto := cryptoWidget.Rows[cryptoWidget.SelectedRow]

		currentPrice, err := getCurrentPrice(config, selectedCrypto)

		if err != nil {
			errorWidget.Text = err.Error()
			ui.Render(headingWidget, graphWidget, cryptoWidget, errorWidget)
			continue
		}

		timeSeriesData, err := getOneDayData(config, selectedCrypto)

		if err != nil {
			errorWidget.Text = err.Error()
			ui.Render(headingWidget, graphWidget, cryptoWidget, errorWidget)
			continue
		}

		fmt.Println(timeSeriesData)

		headingWidget.Text = fmt.Sprintln(selectedCrypto, "\n Current Price:", fmt.Sprintf("%f", currentPrice.Rate), "\n Time:", currentPrice.Time)
		ui.Render(headingWidget, graphWidget, cryptoWidget, errorWidget)
	}

}
