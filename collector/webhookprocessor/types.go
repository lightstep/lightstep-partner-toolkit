package webhookprocessor

import (
	"fmt"
	"go.opentelemetry.io/collector/model/pdata"
)

type ActionKeyValue struct {
	ServiceName string
	Key string
	Value interface{}
	Action Action
}

type Action string

const (
	// UPSERT performs the INSERT or UPDATE action. The key/value is
	// inserted to attributes that did not originally have the key. The key/value is
	// updated for attributes where the key already existed.
	UPSERT Action = "upsert"

	// DELETE deletes the attribute. If the key doesn't exist, no action is performed.
	DELETE Action = "delete"
)

type AttrProc struct {
	actions []ActionKeyValue
}

func (ap *AttrProc) Process(attrs pdata.AttributeMap, serviceName string) {
	for _, action := range ap.actions {

		// skip if there is a service name set that does not match
		if action.ServiceName != "" && serviceName != action.ServiceName {
			continue
		}

		switch action.Action {
		case DELETE:
			attrs.Delete(action.Key)
		case UPSERT:
			attrs.Upsert(action.Key, pdata.NewAttributeValueString(fmt.Sprintf("%v", action.Value)))
		}
	}
}