package repo

import (
	"github.com/AnNosov/communications_info/internal/entity"
)

func checkCountry(country string) bool {

	for key := range entity.Countries {
		if key == country {
			return true
		}
	}
	return false
}
