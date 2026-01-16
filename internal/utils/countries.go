package utils

import countries "github.com/mrz1836/go-countries"

var countriesMap map[string]string

func init() {
	countriesMap = map[string]string{}
	for _, country := range countries.GetAll() {
		countriesMap[country.Alpha2] = country.Name
	}
}

// GetCountries returns a map of country codes to country names
// it returns original map, and it is not intended to be modified
func GetCountries() map[string]string {
	return countriesMap
}

// CountryExists checks if a country code exists in the countries map
func CountryExists(code string) bool {
	_, exists := countriesMap[code]
	return exists
}
