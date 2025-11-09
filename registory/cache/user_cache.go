package cache

type userCache struct{}

type UserCache interface {
	Get(key string) (any, error)
}

func NewUserCache() UserCache {
	return &userCache{}
}

func (c *userCache) Get(key string) (any, error) {
	return nil, nil
}
