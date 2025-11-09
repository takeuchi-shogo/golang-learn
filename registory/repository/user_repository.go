package repository

import (
	"database/sql"

	"github.com/takeuchi-shogo/golang-learn/registory/cache"
	"github.com/takeuchi-shogo/golang-learn/registory/datastore"
	"github.com/takeuchi-shogo/golang-learn/registory/search"
)

type userRepository struct {
	db *sql.DB
	// cache  *cache.UserCache
	// search *search.UserSearch
}

type UserRepository interface {
	Store() datastore.UserStore
	Cache() cache.UserCache
	Search() search.UserSearch
}

func NewUserRepository(
	db *sql.DB,
	// redis *cache.UserCache,
	// search *search.UserSearch,
) UserRepository {
	return &userRepository{
		db: db,
		// cache:  cache,
		// search: search,
	}
}

func (r *userRepository) Store() datastore.UserStore {
	return datastore.NewUserStore(r.db)
}

func (r *userRepository) Cache() cache.UserCache {
	return cache.NewUserCache()
}

func (r *userRepository) Search() search.UserSearch {
	return search.NewUserSearch()
}
