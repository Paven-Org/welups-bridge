package libs

func Map[A any, B any](f func(a A) B, amap []A) []B {
	bmap := make([]B, len(amap))
	for k, v := range amap {
		bmap[k] = f(v)
	}
	return bmap
}

func Reduce[A any, B any](f func(b B, a A) B, seed B, amap []A) B {
	b := seed
	for _, v := range amap {
		b = f(b, v)
	}
	return b
}

func Filter[A any](f func(a A) bool, amap []A) []A {
	var res []A
	for _, v := range amap {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func TakeWhile[A any](f func(a A) bool, amap []A) []A {
	var res []A
	for _, v := range amap {
		if f(v) {
			res = append(res, v)
		} else {
			break
		}
	}
	return res
}

func DropWhile[A any](f func(a A) bool, amap []A) []A {
	var res []A
	for i, v := range amap {
		if !f(v) {
			copy(res, amap[i:])
			break
		}
	}
	return res
}
