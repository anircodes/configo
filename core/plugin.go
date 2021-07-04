package core

type ExecOutput struct {
	ErrorCode    int
	ErrorMessage string
}

type Result struct {
	Status int
	Error  ExecOutput
	Output string
}

type Plugin interface {
	Validate() ExecOutput
	Run() Result
}
