package model

type Test struct {
	Uri       string
	AuthorId  string
	Solutions []*Solution
	QA        map[string]string // question/answer; map, where key - question, value - answer
}
