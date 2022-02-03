package virusio

import (
	"index/suffixarray"
	"strings"

	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/config"
	"github.com/gnames/gnmatcher/ent/virus"
	"github.com/gnames/gnsys"
	log "github.com/sirupsen/logrus"
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
	log.Print("Initializing viruses lookup data.")
	v.prepareData()
}

func (v *virusio) MatchVirus(s string) []mlib.MatchItem {
	bs := v.NameToBytes(s)
	idxs := v.sufary.Lookup(bs, 21)
	res := make([]mlib.MatchItem, len(idxs))
	for i := range idxs {
		if matchItem, ok := v.mapMatchItems[idxs[i]]; ok {
			res[i] = matchItem
		} else {
			log.Errorf("cannot find %d index", idxs[i])
		}
	}
	return res
}

func (v *virusio) prepareDir() {
	log.Print("Preparing directory for viruses.")
	bloomDir := v.cfg.VirusDir()
	err := gnsys.MakeDir(v.cfg.VirusDir())
	if err != nil {
		log.Fatalf("Cannot create directory %s: %s.", bloomDir, err)
	}
}

var sep = "\x00"

func (v *virusio) NameToBytes(name string) []byte {
	name = strings.ToLower(name)
	words := strings.Fields(name)
	res := []byte(sep + strings.Join(words, " "))
	return res
}
