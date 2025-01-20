package setup

type Setup struct {
	Data map[string]interface{}
}

type SetupService interface {
	Setup(setup Setup) error
}
