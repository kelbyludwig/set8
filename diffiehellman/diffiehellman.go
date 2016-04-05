package diffiehellman

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

var ZERO *big.Int = big.NewInt(0)
var ONE *big.Int = big.NewInt(1)
var TWO *big.Int = big.NewInt(2)

type Person struct {
	Generator *big.Int
	Modulus   *big.Int
	Secret    *big.Int
	Public    *big.Int
}

func NewPerson(g, m *big.Int) Person {

	person := Person{}

	buf := make([]byte, len(m.Bytes()))
	rand.Read(buf)
	randomInt := new(big.Int).SetBytes(buf)

	person.Secret = new(big.Int).Mod(randomInt, m)
	person.Generator = new(big.Int).SetBytes(g.Bytes())
	person.Modulus = new(big.Int).SetBytes(m.Bytes())
	person.Public = person.ComputePublic()
	return person
}

func (alice Person) SameGroup(bob Person) bool {
	if alice.Generator.Cmp(bob.Generator) != 0 {
		return false
	}

	if alice.Modulus.Cmp(bob.Modulus) != 0 {
		return false
	}

	return true
}

func (alice Person) ComputePublic() (public *big.Int) {
	m := alice.Modulus
	g := alice.Generator
	a := alice.Secret
	A := new(big.Int).Exp(g, a, m)
	return A
}

func (alice Person) KeyExchange(bob Person) (secret *big.Int) {

	if !alice.SameGroup(bob) {
		fmt.Printf("[!] Looks like Alice and Bob have two different groups.")
	}

	A := alice.Public
	b := bob.Secret
	m := bob.Modulus

	return new(big.Int).Exp(A, b, m)

}

//BruteForceDLP solves for x in e = g^x (mod p) using brute force
func BruteForceDLP(g, e, p *big.Int) *big.Int {
	i := big.NewInt(1)
	for i.Cmp(p) < 1 {
		exp := new(big.Int).Exp(g, i, p)
		if exp.Cmp(e) == 0 {
			return i
		}
		i = i.Add(i, ONE)
	}
	return ONE
}

//Kangaroo uses Pollard's Kangaroo algorithm to solve for x in y = g^x (mod p) given that
//x is in the range  [a,b].
func Kangaroo(a, b, g, y, p *big.Int) (index *big.Int, err error) {

	k := big.NewInt(150)
	c := big.NewInt(4)

	log.Printf("Kangaroo: k = %v\n", k)
	log.Printf("Kangaroo: c = %v\n", c)

	f := func(y *big.Int) *big.Int {
		two := big.NewInt(2)
		return two.Exp(two, y, k)
	}

	output := make([]*big.Int, int(k.Int64()))

	l := big.NewInt(int64(len(output)))
	N := big.NewInt(0)

	for i := big.NewInt(0); i.Cmp(k) != 1; i = i.Add(i, big.NewInt(1)) {
		x := f(i)
		if x.Cmp(ZERO) != 0 {
			output[int(x.Int64())] = big.NewInt(1)
		} else {
			l = i
			break
		}
	}

	for i, j := range output {
		bi := big.NewInt(int64(i))
		if j != nil {
			x := new(big.Int).Mul(bi, j)
			N = N.Add(x, N)
		}
	}
	N = N.Div(N, l)
	N = N.Mul(N, c)
	log.Printf("Kangaroo: N = %v\n", N)

	index, err = kangaroo(a, b, g, y, p, N, k, f)
	return

}

func kangaroo(a, b, g, y, p, N, k *big.Int, f func(*big.Int) *big.Int) (index *big.Int, err error) {

	xT := big.NewInt(0)
	yT := new(big.Int).Exp(g, b, p) //The end of our tame kangaroo's range.

	log.Printf("kangaroo: End of range yT = %v^%v (mod %v) =  %v\n", g, b, p, yT)

	log.Printf("kangaroo: Setting trap for wild kangaroo!\n")
	for i := big.NewInt(1); i.Cmp(N) != 1; i = i.Add(i, ONE) {
		fyt := f(yT)
		xT = xT.Add(xT, fyt) //xT keeps track of the sum of the exponents
		yT = yT.Mul(yT, new(big.Int).Exp(g, fyt, p))
	}

	xW := big.NewInt(0)
	yW := new(big.Int).SetBytes(y.Bytes())

	cond := new(big.Int).Sub(b, a)
	cond = cond.Add(cond, xT)

	log.Printf("kangaroo: Starting wild kangaroo!\n")
	for xW.Cmp(cond) == -1 {

		fyw := f(yW)
		xW = xW.Add(xW, fyw)
		yW = yW.Mul(yW, new(big.Int).Exp(g, fyw, p))

		if yW.Cmp(yT) == 0 {
			index := new(big.Int).Add(b, xT)
			index = index.Sub(index, xW)
			return index, nil
		}

	}

	return index, fmt.Errorf("unable to find the index")

}
