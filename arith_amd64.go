//go:build amd64 && !purego

package uint256

import "unsafe"

//go:noescape
func andAVX2(z, x, y *[32]byte)

//go:noescape
func orAVX2(z, x, y *[32]byte)

//go:noescape
func xorAVX2(z, x, y *[32]byte)

//go:noescape
func notAVX2(z, x *[32]byte)

func (z *Int) And(x, y *Int) *Int {
	if hasAVX2 {
		andAVX2((*[32]byte)(unsafe.Pointer(z)), // #nosec G103
			(*[32]byte)(unsafe.Pointer(x)), // #nosec G103
			(*[32]byte)(unsafe.Pointer(y))) // #nosec G103
		return z
	}
	return z.andScalar(x, y)
}

func (z *Int) Or(x, y *Int) *Int {
	if hasAVX2 {
		orAVX2((*[32]byte)(unsafe.Pointer(z)), // #nosec G103
			(*[32]byte)(unsafe.Pointer(x)), // #nosec G103
			(*[32]byte)(unsafe.Pointer(y))) // #nosec G103
		return z
	}
	return z.orScalar(x, y)
}

func (z *Int) Xor(x, y *Int) *Int {
	if hasAVX2 {
		xorAVX2((*[32]byte)(unsafe.Pointer(z)), // #nosec G103
			(*[32]byte)(unsafe.Pointer(x)), // #nosec G103
			(*[32]byte)(unsafe.Pointer(y))) // #nosec G103
		return z
	}
	return z.xorScalar(x, y)
}

func (z *Int) Not(x *Int) *Int {
	if hasAVX2 {
		notAVX2((*[32]byte)(unsafe.Pointer(z)), // #nosec G103
			(*[32]byte)(unsafe.Pointer(x))) // #nosec G103
		return z
	}
	return z.notScalar(x)
}

// 标量实现作为回退
func (z *Int) andScalar(x, y *Int) *Int {
	z[0] = x[0] & y[0]
	z[1] = x[1] & y[1]
	z[2] = x[2] & y[2]
	z[3] = x[3] & y[3]
	return z
}

func (z *Int) orScalar(x, y *Int) *Int {
	z[0] = x[0] | y[0]
	z[1] = x[1] | y[1]
	z[2] = x[2] | y[2]
	z[3] = x[3] | y[3]
	return z
}

func (z *Int) xorScalar(x, y *Int) *Int {
	z[0] = x[0] ^ y[0]
	z[1] = x[1] ^ y[1]
	z[2] = x[2] ^ y[2]
	z[3] = x[3] ^ y[3]
	return z
}

func (z *Int) notScalar(x *Int) *Int {
	z[0] = ^x[0]
	z[1] = ^x[1]
	z[2] = ^x[2]
	z[3] = ^x[3]
	return z
}
