package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsValidCEP(t *testing.T) {
	var tests = []struct {
		cep    string
		wanted bool
	}{
		{"12345-678", true},
		{"12345678", true},
		{"1234-567", false},
		{"12345-6789", false},
		{"abcd-efgh", false},
	}

	for _, tt := range tests {
		t.Run(tt.cep, func(t *testing.T) {
			got := isValidCEP(tt.cep)
			if got != tt.wanted {
				t.Errorf("isValidCEP(%s) = %v, want %v", tt.cep, got, tt.wanted)
			}
		})
	}
}

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return m.DoFunc(req)
}

func TestGetViaCepData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "01001000") {
			mockResponse := ViaCep{
				Cep:        "01001000",
				Logradouro: "Praça da Sé",
				Localidade: "São Paulo",
			}
			json.NewEncoder(w).Encode(mockResponse)
		} else {
			http.Error(w, "CEP not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	tests := []struct {
		name    string
		cep     string
		wantErr bool
		wantCep string
	}{
		{"Valid CEP", "01001000", false, "01001000"},
		{"Invalid CEP Format", "invalid", true, ""},
		{"CEP Not Found", "00000000", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getViaCepData(server.URL, tt.cep)

			if (err != nil) != tt.wantErr {
				t.Errorf("getViaCepData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.Cep != tt.wantCep {
				t.Errorf("getViaCepData() got = %v, want %v", got.Cep, tt.wantCep)
			}
		})
	}
}

func TestHandleWeatherApiCall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "valid-location") {
			mockResponse := WeatherResponse{
				Current: CurrentData{
					TempC: 25.0,
					TempF: 77.0,
				},
			}
			json.NewEncoder(w).Encode(mockResponse)
		} else {
			http.Error(w, "Location not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	tests := []struct {
		name     string
		location string
		wantErr  bool
		wantTemp float64
	}{
		{"Valid Location", "valid-location", false, 25.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleWeatherApiCall(server.URL, tt.location)

			if (err != nil) != tt.wantErr {
				t.Errorf("handleWeatherApiCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.Current.TempC != tt.wantTemp {
				t.Errorf("handleWeatherApiCall() got = %v, want %v", got.Current.TempC, tt.wantTemp)
			}
		})
	}
}
