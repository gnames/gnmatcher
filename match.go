package gnmatcher

import (
	"fmt"

	"github.com/gnames/gnmatcher/bloom"
)

func (gnm GNmatcher) MatchAndFormat(str string) (string, error) {
	filters, err := bloom.GetFilters(gnm.FiltersDir(), gnm.Dbase)
	_ = filters
	if err != nil {
		return "", err
	}
	fmt.Println("Got here")
	return "", nil
}
