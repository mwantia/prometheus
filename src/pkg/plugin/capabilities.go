package plugin

const (
	None PluginCapabilityType = 0
)

type PluginCapabilityType int32

type PluginCapabilities struct {
	Types []PluginCapabilityType `json:"types"`
}

var PluginCapabilityType_Name = map[int32]string{
	0: "None",
}

var PluginCapabilityType_Value = map[string]int32{
	"None": 0,
}

func (t PluginCapabilityType) String() string {
	value := int32(t)

	name, exists := PluginCapabilityType_Name[value]
	if !exists {
		return ""
	}

	return name
}
