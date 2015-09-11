package zetton

type InvalidJsonError struct {
	Err error
}

func (self InvalidJsonError) Error() string {
	return "invalid json: " + self.Err.Error()
}
