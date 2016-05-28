package diffiehellman

import (
	"fmt"
	"github.com/kelbyludwig/cryptopalscont/sieve"
	"log"
	"math/big"
	"math/rand"
	"testing"
	"time"
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

	js := "30477252323177606811760882179058908038824640750610513771646768011063128035873508507547741559514324673960576895059570"
	j, _ := new(big.Int).SetString(js, 10)

	log.Printf("Small factors of j: %v\n", sieve.UniqueSmallFactors(j))

}

func TestBruteForce(t *testing.T) {
	g := big.NewInt(2)
	p := big.NewInt(17)
	x := big.NewInt(7)
	e := new(big.Int).Exp(g, x, p)
	x2 := BruteForceDLP(g, e, p)
	if x2.Cmp(x) != 0 {
		t.Errorf("TestBruteForce failed\n")
		t.Errorf("expected %v got %v\n", x, x2)
	}
}

func TestPohligHellman(t *testing.T) {
	ps := "7199773997391911030609999317773941274322764333428698921736339643928346453700085358802973900485592910475480089726140708102474957429903531369589969318716771"
	gs := "4565356397095740655436854503483826832136106141639563487732438195343690437606117828318042418238184896212352329118608100083187535033402010599512641674644143"
	js := "30477252323177606811760882179058908038824640750610513771646768011063128035873508507547741559514324673960576895059570"
	qs := "236234353446506858198510045061214171961"

	p, _ := new(big.Int).SetString(ps, 10)
	g, _ := new(big.Int).SetString(gs, 10)
	j, _ := new(big.Int).SetString(js, 10)
	q, _ := new(big.Int).SetString(qs, 10)

	bob := NewPerson(g, p)
	eve := NewPerson(g, p)

	factors := sieve.UniqueSmallFactors(j)

	bobSecret := bob.Secret.Mod(bob.Secret, q)
	log.Printf("Bob's secret %v\n", bobSecret)
	a := make([]*big.Int, 0)
	n := make([]*big.Int, 0)
	enough := new(big.Int)
	for _, r := range factors {

		h := new(big.Int)
		for {
			//Generate a group of small order r
			src := rand.NewSource(time.Now().UnixNano())
			randsrc := rand.New(src)
			h = h.Rand(randsrc, p)
			exp := new(big.Int).Sub(p, ONE)
			exp = exp.Div(exp, r)
			h = h.Exp(h, exp, p)

			if h.Cmp(ONE) != 0 {
				break
			}
		}

		eve.Public = h

		//Send Bob my malcious public key
		sharedSecret := eve.KeyExchange(bob)
		b := BruteForceDLP(h, sharedSecret, p)

		ai := new(big.Int).SetBytes(b.Bytes())
		ni := new(big.Int).SetBytes(r.Bytes())
		a = append(a, ai)
		n = append(n, ni)

		enough = enough.Mul(enough, r)
		if enough.Cmp(q) == 1 {
			break
		}
		log.Printf("x = %v mod %v\n", b, r)
	}
	discoveredSecret, err := CRT(a, n)
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Printf("Eve discovered %v\n", discoveredSecret)

	if discoveredSecret.Cmp(bobSecret) != 0 {
		t.Errorf("TestPohligHellman failed to recover secret")
	}

}

//No shame: https://rosettacode.org/wiki/Chinese_remainder_theorem#Go
func CRT(a, n []*big.Int) (*big.Int, error) {
	one := big.NewInt(1)
	p := new(big.Int).Set(n[0])
	for _, n1 := range n[1:] {
		p.Mul(p, n1)
	}
	var x, q, s, z big.Int
	for i, n1 := range n {
		q.Div(p, n1)
		z.GCD(nil, &s, n1, &q)
		if z.Cmp(one) != 0 {
			return nil, fmt.Errorf("%d not coprime", n1)
		}
		x.Add(&x, s.Mul(a[i], s.Mul(&s, &q)))
	}
	return x.Mod(&x, p), nil
}
