package context_keys

const (
	DeviceInfoContextKey = "DeviceInfo"
	RequestIdContextKey  = "RequestID"
	InputDataContextKey  = "InputData"
)

type InputData map[string]interface{}
