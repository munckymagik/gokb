package btree

import "golang.org/x/exp/slices"

type items[E any] []E

func (self items[E]) first() *E {
	if len(self) == 0 {
		return nil
	}
	return &self[0]
}
func (self items[E]) last() *E {
	if len(self) == 0 {
		return nil
	}
	return &self[len(self)-1]
}

func (self *items[E]) add(elem E) {
	*self = append(*self, elem)
}

func (self *items[E]) pop() E {
	elem := *self.last()
	*self = (*self)[:len(*self)-1]
	return elem
}

func (self items[E]) sort(less func(a E, b E) bool) {
	slices.SortFunc(self, less)
}

func (src *items[E]) splitAt(fromIndex int) items[E] {
	ret := (*src)[fromIndex:]
	*src = (*src)[:fromIndex:fromIndex]
	return ret
}
