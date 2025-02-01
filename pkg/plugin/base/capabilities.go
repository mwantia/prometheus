package base

const (
	None     PluginCapabilityType = 0
	Generate PluginCapabilityType = 1
	Embed    PluginCapabilityType = 2
)

type PluginCapabilityType int32

type PluginCapabilities struct {
	Types []PluginCapabilityType `json:"types"`
}

var pluginCapabilityType_Names = map[PluginCapabilityType]string{
	None:     "none",
	Generate: "generate",
	Embed:    "embed",
}

var pluginCapabilityType_Values = map[string]PluginCapabilityType{
	"none":     None,
	"generate": Generate,
	"embed":    Embed,
}

func (t PluginCapabilityType) String() string {
	name, exist := pluginCapabilityType_Names[t]
	if !exist {
		return ""
	}
	return name
}
