package services

type StrategyInterface interface {
	Dispatch(dispatch Dispatch) ([]Response, error)
}
