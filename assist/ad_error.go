package assist

const (
	AdErrorExist    = 1
	AdErrorNotExist = 2
)

var (
	ErrExist    = &AdError{Code: AdErrorExist, Message: "has been exist"}
	ErrNotExist = &AdError{Code: AdErrorNotExist, Message: "not exist"}
)

type AdError struct {
	Code    int
	Message string
}

func (s *AdError) Error() string {
	return s.Message
}
