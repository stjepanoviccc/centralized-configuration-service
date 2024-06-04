package model

// ConfigurationGroup TODO implement version as struct

// swagger:model ConfigurationGroup
type ConfigurationGroup struct {
	Name           string          `json:"name"`
	Id             int64           `json:"id"`
	Version        Version         `json:"version"`
	Configurations []Configuration `json:"configurations"`
}

func (cg *ConfigurationGroup) SetName(name string) {
	cg.Name = name
}

func (cg *ConfigurationGroup) SetId(id int64) {
	cg.Id = id
}

func (cg *ConfigurationGroup) SetVersion(version Version) {
	cg.Version = version
}

func (cg *ConfigurationGroup) SetConfigurations(configs []Configuration) {
	cg.Configurations = configs
}

// TODO add methods for struct

/*
	Ne znam sta nam jos fali za config group, s obzirom da odvajamo na dva repoa predpostavljam da

ce config group imati neke dodatne nacine za pretragu. To cemo prodiskutovati na discordu
*/
type ConfigurationGroupRepository interface {
	Add(configGroup *ConfigurationGroup) error
	Get(name string, version Version) (ConfigurationGroup, error)
	Delete(configGroup ConfigurationGroup) error
	Save(configGroup *ConfigurationGroup) error
}
