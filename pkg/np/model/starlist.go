package model

type StarList map[int]Star

func (sl StarList) Count() int {
	return len(sl)
}

func (sl StarList) Get(uid int) (Star, bool) {
	s, ok := sl[uid]
	return s, ok
}
