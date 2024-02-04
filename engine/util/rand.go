package util

var randCtx struct {
	a uint64
	b uint64
	c uint64
	d uint64
}

func rotate(x uint64, k int) uint64 {
	return (((x) << (k)) | ((x) >> (64 - (k))))
}

func RandU64() uint64 {
	e := randCtx.a - rotate(randCtx.b, 7)
	randCtx.a = randCtx.b ^ rotate(randCtx.c, 13)
	randCtx.b = randCtx.c + rotate(randCtx.d, 37)
	randCtx.c = randCtx.d + e
	randCtx.d = e + randCtx.a
	return randCtx.d
}

func RandInit(seed uint64) {
	randCtx.a = 0xf1ea5eed
	randCtx.b, randCtx.c, randCtx.d = seed, seed, seed
	for i := 0; i < 20; i++ {
		RandU64()
	}
}
