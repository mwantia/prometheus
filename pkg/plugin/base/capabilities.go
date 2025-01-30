package base

const (
	None PluginCapabilityType = 0
)

type PluginCapabilityType int32

type PluginCapabilities struct {
	Types []PluginCapabilityType `json:"types"`
}

var pluginCapabilityType_Names = map[int32]string{
	0: "none",
}

var pluginCapabilityType_Values = map[string]int32{
	"none": 0,
}

func (t PluginCapabilityType) String() string {
	name, exist := pluginCapabilityType_Names[int32(t)]
	if !exist {
		return ""
	}
	return name
}
