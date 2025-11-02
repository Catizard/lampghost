package random

const (
	multiplier uint64 = 0x5DEECE66D
	mask       uint64 = (1 << 48) - 1
	addend     uint64 = 0xB
)

type JavaRandom struct {
	seed uint64
}

func NewJavaRandom(seed int64) JavaRandom {
	r := JavaRandom{
		seed: (uint64(seed) ^ multiplier) & mask,
	}
	return r
}

func (r *JavaRandom) Next(bit int32) int32 {
	r.seed = (r.seed*multiplier + addend) & mask
	// TODO: This might cause a diverage
	return int32((r.seed >> (48 - bit)) & ((1 << bit) - 1))
}

func (r *JavaRandom) NextInt(bound int32) int32 {
	m := bound - 1
	ret := r.Next(31)
	if (bound & m) == 0 {
		ret = int32((int64(bound) * int64(ret)) >> 31)
	} else {
		u := ret
		for {
			ret = u % bound
			if u-ret+m >= 0 {
				return ret
			}
			u = r.Next(31)
		}
	}
	return ret
}

func MakeRandom(seed int64, keys []int) []int {
	random := NewJavaRandom(seed)
	// TODO: I know I shouldn't do this, I will fix this later
	l := make([]int, len(keys))
	result := make([]int, len(keys))
	for i := range l {
		result[i] = i
		l[i] = i
	}
	for lane := range keys {
		r := random.NextInt(int32(len(l)))
		result[keys[lane]] = l[r]
		l = append(l[:r], l[r+1:]...)
	}
	return result
}
