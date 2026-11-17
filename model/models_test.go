package model

import (
	"slices"
	"testing"
)

// Actuator state properties are driven by actions, never by a property write:
// the TD must advertise them as read-only.
func TestCapabilityStatePropertiesReadOnly(t *testing.T) {
	models := []string{"LightBasic", "LightDimmer", "ShutterBasic", "ShutterPosition", "GarageDoor", "Trigger"}
	th, err := NewThingFromModels("urn:test", "1.0.0", "test", "test", models)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"on", "brightness", "open", "position", "moving", "closed", "triggered"} {
		p, ok := th.Properties[name]
		if !ok {
			t.Fatalf("property %q missing", name)
		}
		if !p.ReadOnly {
			t.Errorf("property %q should be read-only", name)
		}
	}
}

// The TD names its models in @type, once each: auth reads them.
func TestModelsInType(t *testing.T) {
	th, err := NewThingFromModels("urn:test", "1.0.0", "test", "test", []string{"LightBasic", "Trigger"})
	if err != nil {
		t.Fatal(err)
	}
	if err := AddModel(th, "LightBasic", ""); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(th.AtTypes, []string{"LightBasic", "Trigger"}) {
		t.Errorf("@type = %v", th.AtTypes)
	}
}
