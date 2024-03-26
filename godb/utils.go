package godb
type Pair[T, U any] struct {
    First  T
    Second U
}
func Zip[T, U any](ts []T, us []U) []Pair[T,U] {
    if len(ts) != len(us) {
        panic("slices have different length")
    }
    pairs := make([]Pair[T,U], len(ts))
    for i := 0; i < len(ts); i++ {
        pairs[i] = Pair[T,U]{ts[i], us[i]}
    }
    return pairs
}

func bool2int(b bool) uint8 {
	if b {
	   return 1
	}
	return 0
 } 