package interfaces

type UseCaser[TRequest any, TResponse any] interface {
	Execute(TRequest) (TResponse, error)
}
