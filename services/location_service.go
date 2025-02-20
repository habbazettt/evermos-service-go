package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"strings"
)

// Struct untuk Province
type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Struct untuk City
type City struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

// API Wilayah Indonesia
const apiProvinsi = "https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json"
const apiKota = "https://www.emsifa.com/api-wilayah-indonesia/api/regencies/"

// Get All Provinces
func GetAllProvinces() ([]Province, error) {
	resp, err := http.Get(apiProvinsi)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var provinces []Province
	err = json.NewDecoder(resp.Body).Decode(&provinces)
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

// Get All Cities by Province ID
func GetCitiesByProvince(provinceID string) ([]City, error) {
	if provinceID == "" {
		return nil, errors.New("province ID is required")
	}

	resp, err := http.Get(apiKota + provinceID + ".json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cities []City
	err = json.NewDecoder(resp.Body).Decode(&cities)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

// Get List Provinces with Search, Limit, and Pagination
func GetListProvinces(search string, limit, page int) ([]Province, error) {
	resp, err := http.Get(apiProvinsi)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var provinces []Province
	err = json.NewDecoder(resp.Body).Decode(&provinces)
	if err != nil {
		return nil, err
	}

	// Apply search filter
	if search != "" {
		filtered := []Province{}
		search = strings.ToLower(search)
		for _, p := range provinces {
			if strings.Contains(strings.ToLower(p.Name), search) {
				filtered = append(filtered, p)
			}
		}
		provinces = filtered
	}

	// Apply pagination
	start := (page - 1) * limit
	if start > len(provinces) {
		return []Province{}, nil
	}

	end := start + limit
	if end > len(provinces) {
		end = len(provinces)
	}

	return provinces[start:end], nil
}

// Get Detail Province
func GetProvinceDetail(provinceID string) (*Province, error) {
	resp, err := http.Get(apiProvinsi)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var provinces []Province
	err = json.NewDecoder(resp.Body).Decode(&provinces)
	if err != nil {
		return nil, err
	}

	for _, p := range provinces {
		if p.ID == provinceID {
			return &p, nil
		}
	}

	return nil, errors.New("province not found")
}

// Get List Cities with Search, Limit, and Pagination
func GetListCities(provinceID, search string, limit, page int) ([]City, error) {
	if provinceID == "" {
		return nil, errors.New("province ID is required")
	}

	resp, err := http.Get(apiKota + provinceID + ".json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cities []City
	err = json.NewDecoder(resp.Body).Decode(&cities)
	if err != nil {
		return nil, err
	}

	// Apply search filter
	if search != "" {
		filtered := []City{}
		search = strings.ToLower(search)
		for _, c := range cities {
			if strings.Contains(strings.ToLower(c.Name), search) {
				filtered = append(filtered, c)
			}
		}
		cities = filtered
	}

	// Apply pagination
	start := (page - 1) * limit
	if start > len(cities) {
		return []City{}, nil
	}

	end := start + limit
	if end > len(cities) {
		end = len(cities)
	}

	return cities[start:end], nil
}

// Get Detail City
func GetCityDetail(provinceID, cityID string) (*City, error) {
	if provinceID == "" || cityID == "" {
		return nil, errors.New("province ID and city ID are required")
	}

	resp, err := http.Get(apiKota + provinceID + ".json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cities []City
	err = json.NewDecoder(resp.Body).Decode(&cities)
	if err != nil {
		return nil, err
	}

	for _, c := range cities {
		if c.ID == cityID {
			return &c, nil
		}
	}

	return nil, errors.New("city not found")
}

// Get Detail City by City ID Only
func GetCityDetailByID(cityID string) (*City, error) {
	if cityID == "" {
		return nil, errors.New("city ID is required")
	}

	// Panggil API untuk mendapatkan data kota berdasarkan ID
	url := fmt.Sprintf("https://api.rajaongkir.com/starter/city/%s", cityID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var city City
	err = json.NewDecoder(resp.Body).Decode(&city)
	if err != nil {
		return nil, err
	}

	return &city, nil
}
