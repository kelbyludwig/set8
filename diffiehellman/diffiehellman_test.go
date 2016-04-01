package diffiehellman

import (
	"github.com/kelbyludwig/cryptopalscont/sieve"
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

func TestSmallFactors(t *testing.T) {

	//ps := "7199773997391911030609999317773941274322764333428698921736339643928346453700085358802973900485592910475480089726140708102474957429903531369589969318716771"
	//gs := "4565356397095740655436854503483826832136106141639563487732438195343690437606117828318042418238184896212352329118608100083187535033402010599512641674644143"
	js := "30477252323177606811760882179058908038824640750610513771646768011063128035873508507547741559514324673960576895059570"
	//p, _ := new(big.Int).SetString(ps, 10)
	//g, _ := new(big.Int).SetString(gs, 10)
	j, _ := new(big.Int).SetString(js, 10)

	limit := 65536
	uniqueSmallFactors := func(x *big.Int) []*big.Int {
		primes := sieve.GeneratePrimes(limit)
		zero := big.NewInt(0)
		factors := make([]*big.Int, 0)
		for _, p := range primes {
			bigP := big.NewInt(int64(p))
			rem := new(big.Int).Mod(x, bigP)
			if rem.Cmp(zero) == 0 {
				factors = append(factors, bigP)
			}
		}
		return factors
	}

	log.Printf("Small factors of j: %v\n", uniqueSmallFactors(j))

	//alice := NewPerson(g, p)
	//bob := NewPerson(g, p)

}
