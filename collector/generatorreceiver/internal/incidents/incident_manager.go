package incidents

import (
	"go.uber.org/zap"
	"sync"
)

type IncidentManager struct {
	incidents map[string]*Incident

	sync.Mutex
}

var Manager *IncidentManager

func init() {
	Manager = &IncidentManager{incidents: make(map[string]*Incident)}
}

func (im *IncidentManager) LoadIncidents(incidents []Incident, logger *zap.Logger) {
	im.Lock()
	defer im.Unlock()
	for _, v := range incidents {
		inc := v
		inc.Setup(logger)
		im.incidents[inc.Name] = &inc
	}
}

func (im *IncidentManager) GetIncidents() []*Incident {
	im.Lock()
	defer im.Unlock()
	out := []*Incident{}
	for _, i := range Manager.incidents {
		out = append(out, i)
	}
	return out
}

func (im *IncidentManager) GetIncident(name string) *Incident {
	im.Lock()
	defer im.Unlock()
	incident, exists := im.incidents[name]
	if !exists {
		incident = &Incident{Name: name}
		im.incidents[name] = incident
	}
	return incident
}
