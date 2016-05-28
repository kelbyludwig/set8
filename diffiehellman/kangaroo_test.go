package diffiehellman

import (
	"log"
	"math/big"
	"testing"
)

func TestKangarooSmall(t *testing.T) {
	gs := "2"
	ps := "23"
	//xs := "5"
	ys := "9"

	p, _ := new(big.Int).SetString(ps, 10)
	g, _ := new(big.Int).SetString(gs, 10)
	y, _ := new(big.Int).SetString(ys, 10)
	//x, _ := new(big.Int).SetString(xs, 10)

	a := big.NewInt(0)
	b := big.NewInt(8)

	log.Printf("[DEBUG] Starting kangaroo...\n")
	index, err := Kangaroo(y, g, a, b, p)
	if err != nil {
		t.Errorf("KangarooSmall failed to find an index.\n")
		return
	}
	log.Printf("[DEBUG] index %v...\n", index)
	res := new(big.Int).Exp(g, index, p)
	if res.Cmp(y) != 0 {
		t.Errorf("KangarooSmall failed to find the right index.\n")
		return
	}

}

func TestKangarooMedium(t *testing.T) {

	ps := "11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623"
	//qs := "335062023296420808191071248367701059461"
	//js := "34233586850807404623475048381328686211071196701374230492615844865929237417097514638999377942356150481334217896204702"
	gs := "622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357"
	ys := "7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119"
	p, _ := new(big.Int).SetString(ps, 10)
	//q, _ := new(big.Int).SetString(ps, 10)
	//j, _ := new(big.Int).SetString(ps, 10)
	g, _ := new(big.Int).SetString(gs, 10)
	y, _ := new(big.Int).SetString(ys, 10)

	a := big.NewInt(0)
	b := new(big.Int).Exp(big.NewInt(2), big.NewInt(20), nil)
	index, err := Kangaroo(y, g, a, b, p)
	if err != nil {
		t.Errorf("KangarooMedium failed to find an index.\n")
		return
	}
	res := new(big.Int).Exp(g, index, p)
	if res.Cmp(y) != 0 {
		t.Errorf("KangarooMedium failed to find the right index.\n")
		return
	}
	log.Printf("[DEBUG] index %v...\n", index)
}

func TestKangarooBig(t *testing.T) {
	ps := "11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623"
	//xs := "359579674340"
	//qs := "335062023296420808191071248367701059461"
	//js := "34233586850807404623475048381328686211071196701374230492615844865929237417097514638999377942356150481334217896204702"
	gs := "622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357"
	ys := "9388897478013399550694114614498790691034187453089355259602614074132918843899833277397448144245883225611726912025846772975325932794909655215329941809013733"
	p, _ := new(big.Int).SetString(ps, 10)
	//q, _ := new(big.Int).SetString(ps, 10)
	//j, _ := new(big.Int).SetString(ps, 10)
	g, _ := new(big.Int).SetString(gs, 10)
	y, _ := new(big.Int).SetString(ys, 10)

	a := big.NewInt(0)
	b := new(big.Int).Exp(big.NewInt(2), big.NewInt(40), nil)
	index, err := Kangaroo(y, g, a, b, p)
	if err != nil {
		t.Errorf("KangarooBig failed to find an index.\n")
		return
	}
	res := new(big.Int).Exp(g, index, p)

	if res.Cmp(y) != 0 {
		t.Errorf("KangarooBig failed to find the right index.\n")
		return
	}
	log.Printf("[DEBUG] index %v...\n", index)
}
