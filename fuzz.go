// +build gofuzz

package uint256

const (
	opUdivrem = 0
	opMul     = 1
)

func Fuzz(data []byte) int {
	if len(data) != 65 {
		return 0
	}

	op := data[0]

	var x, y Int
	x.SetBytes(data[1:33])
	y.SetBytes(data[33:])

	bx := x.ToBig()
	by := y.ToBig()

	switch op {
	case opUdivrem:
		if y.IsZero() {
			return 0
		}

		var q, r Int
		q.Div(&x, &y)
		r.Mod(&x, &y)
		bx.QuoRem(bx, by, by)
		eq, _ := FromBig(bx)
		er, _ := FromBig(by)

		if !q.Eq(eq) {
			panic("invalid quotient")
		}

		if !r.Eq(er) {
			panic("invalid remainder")
		}

	case opMul:
		var p Int
		p.Mul(&x, &y)

		bx.Mul(bx, by)
		ep, _ := FromBig(bx)

		if !p.Eq(ep) {
			panic("invalid multiplication result")
		}
	}

	return 0
}
