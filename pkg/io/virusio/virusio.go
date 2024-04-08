package virusio

import (
	"fmt"
	"index/suffixarray"
	"strings"

	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/internal/ent/virus"
	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/gnames/gnsys"
	"github.com/rs/zerolog/log"
)

type virusio struct {
	cfg           config.Config
	sufary        *suffixarray.Index
	mapMatchItems map[int]mlib.MatchItem
}

func New(cfg config.Config) virus.VirusMatcher {
	res := virusio{
		cfg:           cfg,
		mapMatchItems: make(map[int]mlib.MatchItem),
	}
	return &res
}

func (v *virusio) Init() {
	v.prepareDir()
	log.Info().Msg("Initializing viruses lookup data")
	v.prepareData()
}

// SetConfig updates configuration of the matcher.
func (v *virusio) SetConfig(cfg config.Config) {
	v.cfg = cfg
}

func (v *virusio) MatchVirus(s string) []mlib.MatchItem {
	bs := v.NameToBytes(s)
	idxs := v.sufary.Lookup(bs, 21)
	res := make([]mlib.MatchItem, len(idxs))
	for i := range idxs {
		if matchItem, ok := v.mapMatchItems[idxs[i]]; ok {
			res[i] = matchItem
		} else {
			err := fmt.Errorf("cannot find %d index", idxs[i])
			log.Warn().Err(err)
		}
	}
	return res
}

func (v *virusio) prepareDir() {
	log.Info().Msg("Preparing directory for viruses")
	bloomDir := v.cfg.VirusDir()
	err := gnsys.MakeDir(v.cfg.VirusDir())
	if err != nil {
		log.Fatal().Err(err).
			Msgf("Cannot create directory %s", bloomDir)
	}
}

var sep = "\x00"

func (v *virusio) NameToBytes(name string) []byte {
	name = strings.ToLower(name)
	words := strings.Fields(name)
	res := []byte(sep + strings.Join(words, " "))
	return res
}
