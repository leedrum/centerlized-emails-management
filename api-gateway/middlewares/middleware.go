package middlewares

type Middleware interface {
	JWTAuth()
	Logger()
	RateLimit()
}
