package model

type StarList map[int]Star

func (sl StarList) Count() int {
	return len(sl)
}
