package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type WeatherResponse struct {
	Current CurrentData `json:"current"`
}

type CurrentData struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func isValidCEP(cep string) bool {
	matched, _ := regexp.MatchString(`^\d{5}-?\d{3}$`, cep)
	return matched
}

func getViaCepData(baseUrl, cep string) (*ViaCep, error) {
	if cep == "" || !isValidCEP(cep) {
		return nil, fmt.Errorf("invalid zipcode")
	}

	cepScaped := url.QueryEscape(cep)
	urlFormated := fmt.Sprintf("%s/ws/%s/json/", baseUrl, cepScaped)

	resp, err := http.Get(urlFormated)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var viaCepResponse ViaCep
	err = json.Unmarshal(body, &viaCepResponse)
	if err != nil {
		return nil, err
	}

	return &viaCepResponse, nil
}

func handleWeatherApiCall(baseUrl, location string) (*WeatherResponse, error) {
	locationScaped := url.QueryEscape(location)
	apiKey := os.Getenv("API_KEY")
	url := fmt.Sprintf("%s?key=%s&q=%s&aqi=no", baseUrl, apiKey, locationScaped)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	cep := queryValues.Get("cep")

	responseCep, err := getViaCepData("http://viacep.com.br/", cep)

	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	if responseCep.Localidade == "" {
		http.Error(w, "can not find zipcode", http.StatusInternalServerError)
		return
	}

	weather, err := handleWeatherApiCall("http://api.weatherapi.com/v1/current.json", responseCep.Localidade)
	if err != nil {
		http.Error(w, "error fetching weather data", http.StatusInternalServerError)
		return
	}

	response := TemperatureResponse{
		TempC: weather.Current.TempC,
		TempF: weather.Current.TempF,
		TempK: weather.Current.TempC + 273.15,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error processing request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, continuing with environment variables")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	http.HandleFunc("/cep", cepHandler)
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
