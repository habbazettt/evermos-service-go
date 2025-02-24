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
const apiKota = "https://www.emsifa.com/api-wilayah-indonesia/api/regencies.json"
const apiKotaBase = "https://www.emsifa.com/api-wilayah-indonesia/api/regencies/"

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

// Get List of Cities by Province ID
func GetCitiesByProvince(provinceID string) ([]City, error) {
	if provinceID == "" {
		return nil, errors.New("province ID is required")
	}

	// Panggil API berdasarkan provinsi
	url := fmt.Sprintf("%s%s.json", apiKotaBase, provinceID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cities []City
	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return nil, err
	}

	// Jika tidak ada data kota ditemukan
	if len(cities) == 0 {
		return nil, errors.New("no cities found for this province")
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

// Get City Detail by City ID
func GetCityDetailByID(cityID string) (*City, error) {
	if cityID == "" {
		return nil, errors.New("city ID is required")
	}

	// Karena API tidak menyediakan endpoint langsung untuk `cityID`,
	// kita harus memeriksa semua kota dari berbagai provinsi.
	provinces, err := GetAllProvinces()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provinces: %w", err)
	}

	// Loop ke semua provinsi untuk mencari kota
	for _, province := range provinces {
		url := fmt.Sprintf("%s%s.json", apiKotaBase, province.ID)
		resp, err := http.Get(url)
		if err != nil {
			continue // Lewati jika gagal
		}
		defer resp.Body.Close()

		// Cek jika API tidak mengembalikan status 200
		if resp.StatusCode != http.StatusOK {
			continue
		}

		// Decode JSON daftar kota dari provinsi ini
		var cities []City
		if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
			continue
		}

		// Cari kota yang sesuai dengan cityID
		for _, city := range cities {
			if city.ID == cityID {
				return &city, nil
			}
		}
	}

	return nil, errors.New("city not found")
}
