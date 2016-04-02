package diffiehellman

import (
	"crypto/rand"
	"fmt"
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
