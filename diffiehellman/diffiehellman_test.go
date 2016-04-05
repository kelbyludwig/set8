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

func TestKangarooBig(t *testing.T) {

	ps := "11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623"
	qs := "335062023296420808191071248367701059461"
	//js := "34233586850807404623475048381328686211071196701374230492615844865929237417097514638999377942356150481334217896204702"
	gs := "622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357"
	ys := "7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119"

	p, _ := new(big.Int).SetString(ps, 10)
	q, _ := new(big.Int).SetString(qs, 10)
	g, _ := new(big.Int).SetString(gs, 10)
	y, _ := new(big.Int).SetString(ys, 10)

	a := big.NewInt(0)
	b := big.NewInt(2)
	b = b.Exp(b, big.NewInt(20), q)

	x, err := Kangaroo(a, b, g, y, p)

	if err != nil {
		t.Errorf("TestKangaroo failed")
	}

	log.Printf("index %v\n", x)

}

func TestKangarooSmall(t *testing.T) {

	gs := "2"
	ps := "23"
	//xs := "5"
	ys := "9"

	p, _ := new(big.Int).SetString(ps, 10)
	g, _ := new(big.Int).SetString(gs, 10)
	y, _ := new(big.Int).SetString(ys, 10)

	a := big.NewInt(0)
	b := big.NewInt(8)

	x, err := Kangaroo(a, b, g, y, p)

	if err != nil {
		t.Errorf("TestKangarooSmall failed")
	}

	log.Printf("index %v\n", x)

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
