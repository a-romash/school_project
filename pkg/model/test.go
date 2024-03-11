package model

import "encoding/json"

type Test struct {
	Uri         string
	Author      string
	AuthorId    string
	SolutionsId []int
	Question    []map[string]string
	Id          int
}

func (t *Test) GetJson() (data []byte, err error) {
	data, err = json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return data, nil
}
