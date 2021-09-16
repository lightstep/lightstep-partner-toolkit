package webhookprocessor

import (
	"fmt"
	"go.opentelemetry.io/collector/model/pdata"
)

type ActionKeyValue struct {
	Key string `mapstructure:"key"`
	Value interface{} `mapstructure:"value"`
	Action Action `mapstructure:"action"`
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

func (ap *AttrProc) Process(attrs pdata.AttributeMap) {
	for _, action := range ap.actions {
		switch action.Action {
		case DELETE:
			attrs.Delete(action.Key)
		case UPSERT:
			attrs.Upsert(action.Key, pdata.NewAttributeValueString(fmt.Sprintf("%v", action.Value)))
		}
	}
}