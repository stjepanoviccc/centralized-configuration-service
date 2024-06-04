package model

import (
	"fmt"
	"sort"
)

// TODO add version as struct, add labels as field (used for filtering)

// swagger:model Configuration
type Configuration struct {
	Name       string            `json:"name"`
	Id         int64             `json:"id"`
	Version    Version           `json:"version"`
	Parameters map[string]string `json:"parameters"`
	Labels     map[string]string `json:"labels"`
}

/*
	mozda bi bilo dobro dodati dve odvojene strukture ako treba da postoji

neka apstrakcija izmedju klijenta i beka U tom slucaju bi imali dve strukture,
kao da pravimo dto sloj, mada s obzirom kakva je primena projekta msm da nema
potrebe Gettera nema zato sto su exportovana polja pa nema potrebe, vec ima
gettere sam po sebi
*/
func (c *Configuration) SetName(name string) {
	c.Name = name
}

func (c *Configuration) SetId(id int64) {
	c.Id = id
}

func (c *Configuration) SetVersion(version Version) {
	c.Version = version
}

func (c *Configuration) SetParameters(params map[string]string) {
	c.Parameters = params
}

func (c *Configuration) SetLabels(labels map[string]string) {
	c.Labels = labels
}

func SortLabels(labels map[string]string) string {
	keys := make([]string, 0, len(labels))

	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	reorderedMap := make(map[string]string)

	for _, k := range keys {
		reorderedMap[k] = labels[k]
	}

	var slice []string
	var counter int
	for k, v := range reorderedMap {
		if counter != len(reorderedMap)-1 {
			counter++
			slice = append(slice, fmt.Sprintf("%s:%s", k, v))
			slice = append(slice, ";")
		} else {
			slice = append(slice, fmt.Sprintf("%s:%s", k, v))
		}
	}

	var result string
	for _, v := range slice {
		result += v
	}

	return result
}

/* Ovo nisam hteo vise nista dodavati, msm da je dovoljno za pocetak, samo osnovan CRUD
mislim da nam nece biti potreban FindAll zbog toga sto moze samo po IDu da se povuce
*/

type ConfigurationRepository interface {
	Add(config *Configuration) error
	Get(name string, version Version) (Configuration, error)
	Delete(config Configuration) error
}
