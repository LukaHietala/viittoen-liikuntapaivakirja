package api

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "api context value " + k.name
}

var (
	collectionContextKey  = &contextKey{"collection"}
	challengeContextKey   = &contextKey{"challenge"}
	performanceContextKey = &contextKey{"performance"}
	userContextKey        = &contextKey{"user"}
)
