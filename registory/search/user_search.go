package search

type userSearch struct{}

type UserSearch interface {
	Search(query string) (any, error)
}

func NewUserSearch() UserSearch {
	return &userSearch{}
}

func (s *userSearch) Search(query string) (any, error) {
	return nil, nil
}
