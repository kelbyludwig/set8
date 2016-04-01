package diffiehellman

import (
	"log"
	"math/big"
	"testing"
)

func TestKeyExchange(t *testing.T) {
	ps := "7199773997391911030609999317773941274322764333428698921736339643928346453700085358802973900485592910475480089726140708102474957429903531369589969318716771"
	gs := "4565356397095740655436854503483826832136106141639563487732438195343690437606117828318042418238184896212352329118608100083187535033402010599512641674644143"
	p, _ := new(big.Int).SetString(ps, 10)
	g, _ := new(big.Int).SetString(gs, 10)

	alice := NewPerson(g, p)
	bob := NewPerson(g, p)

	log.Printf("TestKeyExchange: Alice Secret %v\n", alice.Secret)
	log.Printf("TestKeyExchange: Bob Secret %v\n", bob.Secret)

	sharedAlice := alice.KeyExchange(bob)
	sharedBob := bob.KeyExchange(alice)

	log.Printf("TestKeyExchange: Alice Shared %v\n", sharedAlice)
	log.Printf("TestKeyExchange: Bob Shared %v\n", sharedBob)

	if sharedAlice.Cmp(sharedBob) != 0 {
		t.Errorf("TestKeyExchange: shared keys were not equal\n")
	}
}