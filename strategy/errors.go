package strategy

type NoStrategyError struct {
	Err error
}

func (e NoStrategyError) Error() string {
	return e.Err.Error()
}
