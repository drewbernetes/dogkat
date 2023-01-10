package coreworkloads

type Resource interface {
	Generate(data map[string]string)
	Create() error
	Validate() error
	Delete() error
	IsReady() bool
}
