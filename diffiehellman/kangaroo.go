package diffiehellman

import (
	"fmt"
	"log"
	"math/big"
)

var one *big.Int = big.NewInt(1)
var two *big.Int = big.NewInt(2)

//Kangaroo uses Pollard's Kangaroo algorithm to solve for x in y = g^x (mod p)
//given that x is in the range [a,b]. Kangaroo will return an error if it
//cannot find a valid index.
func Kangaroo(y, g, a, b, p *big.Int) (*big.Int, error) {

	//This is how sage generates N
	N := sqrt(new(big.Int).Sub(b, a))
	N = N.Add(N, one)

	k := big.NewInt(0)

	//This is how sage generates k
	for new(big.Int).Exp(two, k, nil).Cmp(N) < 0 {
		k = k.Add(k, one)
	}

	xT := big.NewInt(0)
	yT := new(big.Int).Exp(g, b, p)

	f := func(y *big.Int) *big.Int {
		ymk := new(big.Int).Mod(y, k)
		return new(big.Int).Exp(two, ymk, p)
	}

	log.Printf("[DEBUG] tame kangaroo...\n")
	for i := big.NewInt(1); i.Cmp(N) < 1; i = i.Add(i, one) {
		fyT := f(yT)
		gfyT := new(big.Int).Exp(g, fyT, p)
		xT = xT.Add(xT, fyT)
		yT = yT.Mul(yT, gfyT)
		yT = yT.Mod(yT, p)
	}

	xW := big.NewInt(0)
	yW := new(big.Int).SetBytes(y.Bytes())

	log.Printf("[DEBUG] wild kangaroo...\n")
	baxT := new(big.Int).Sub(b, a)
	baxT = baxT.Add(baxT, xT)
	gfyW := new(big.Int)
	for xW.Cmp(baxT) == -1 { // for xW < baxT
		fyW := f(yW)
		gfyW := gfyW.Exp(g, fyW, p)
		xW = xW.Add(xW, fyW)
		yW = yW.Mul(yW, gfyW)
		yW = yW.Mod(yW, p)

		if yW.Cmp(yT) == 0 {
			ret := new(big.Int).Add(b, xT)
			ret = ret.Sub(ret, xW)
			return ret, nil
		}
	}
	return big.NewInt(0), fmt.Errorf("failed to find index")
}

//integer sqrt root
func sqrt(number *big.Int) *big.Int {
	next := func(n, i *big.Int) *big.Int {
		t := new(big.Int).Div(i, n)
		t = t.Add(t, n)
		t = t.Rsh(t, 1)
		return t
	}

	n := big.NewInt(1)
	n1 := next(n, number)

	l1 := new(big.Int).Sub(n1, n)
	for new(big.Int).Abs(l1).Cmp(one) == 1 {
		n = n1
		n1 = next(n, number)
		l1 = l1.Sub(n1, n)
	}

	l2 := new(big.Int).Mul(n1, n1)
	for l2.Cmp(number) == 1 {
		n1 = n1.Sub(n1, one)
		l2 = l2.Mul(n1, n1)
	}
	return n1
}
