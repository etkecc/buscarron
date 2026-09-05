package utils

import countries "github.com/mrz1836/go-countries"

var countriesMap map[string]string

func init() {
	countriesMap = map[string]string{}
	for _, country := range countries.GetAll() {
		countriesMap[country.Alpha2] = country.Name
	}
}

// GetCountries returns the original country codes to names map; callers must not modify it.
func GetCountries() map[string]string {
	return countriesMap
}

// CountryExists checks if a country code exists in the countries map
func CountryExists(code string) bool {
	_, exists := countriesMap[code]
	return exists
}
