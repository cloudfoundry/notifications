package strategies

type StrategyInterface interface {
	Dispatch(dispatch Dispatch) ([]Response, error)
}
