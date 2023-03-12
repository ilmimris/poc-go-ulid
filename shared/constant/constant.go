package constant

const (
	// AppName
	AppName        = "poc-go-ulid"
	HeaderXTraceID = "X-Trace-Id"

	// Default
	DefaultPort     int = 8080
	DefaultLogLevel     = "info"

	MaxBodyLimit    = 1024 * 1024 * 12 // 12MB
	MaxReadTimeOut  = 60               // 60 seconds
	MaxWriteTimeOut = 60               // 60 seconds
)
