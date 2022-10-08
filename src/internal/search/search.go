package search

import (
	"fmt"

	"github.com/nullify005/exporter-weather/internal/bom"
)

func Search(name string) (err error) {
	locations, err := bom.Location(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(locations) <= 0 {
		err = fmt.Errorf("no locations found for: %s", name)
		fmt.Println(err.Error())
		return
	}
	for _, v := range locations {
		fmt.Println(v.String())
	}
	return
}
