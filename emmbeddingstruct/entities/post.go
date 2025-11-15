package entities

import "github.com/takeuchi-shogo/golang-learn/emmbeddingstruct/base"

type Post struct {
	base.Base
	title   string
	content string
}

func NewPost(id int, title, content string) *Post {
	if id == 0 {
		return nil
	}
	postId := base.NewBase(id)
	return &Post{Base: *postId, title: title, content: content}
}

func (p *Post) GetTitle() string {
	return p.title
}

func (p *Post) GetContent() string {
	return p.content
}
