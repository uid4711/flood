package main

//
//
//
type IntToByteArr struct {
	in  int64
	out []byte
}

type cntArr struct {
	a, b, c, d byte
}

//
//
//
func NewIntToByteArr() *IntToByteArr {
	return &IntToByteArr{
		out: make([]byte, 25, 50),
	}
}

//
//
//
func incCounter(ca *cntArr) {

	if ca.a == 10 {
		ca.b++
		ca.a = 0
	}

	if ca.b == 10 {
		ca.c++
		ca.b = 0
		ca.a = 0
	}

	if ca.c == 10 {
		ca.d++
		ca.c = 0
		ca.b = 0
		ca.a = 0
	}

	if ca.d == 10 {
		ca.d = 0
		ca.c = 0
		ca.b = 0
		ca.a = 0
	}

}

//
//
//
func getExponent(in int64) int {
	i := 0
	for {
		if in /= 10; in == 0 {
			i++
			break
		}
		i++
	}
	return i
}

//
//
//
func intToByteArr(t *IntToByteArr) []byte {

	i := getExponent(t.in)

	in := t.in

	out := t.out

	out = out[:i]

	for {

		i--
		t.out[i] = byte(in%10) + 48

		if in /= 10; in == 0 {
			break
		}

	}

	return out
}
