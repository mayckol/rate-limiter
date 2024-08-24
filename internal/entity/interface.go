package entity

type RequestRepositoryInterface interface {
	CheckRateLimit(key string, limit int) (bool, error)
}
