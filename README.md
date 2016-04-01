Welcome back, Cryptopals! It's the long-awaited eighth series!

This set focuses on abstract algebra, including DH, GCM, and (most
importantly) elliptic curve cryprography.

Fair warning - it's really tough! There's a ton of content here, and
it's more demanding than anything we've released so far. By the time
you're done, you will have written an ad hoc, informally-specified,
bug-ridden, slow implementation of one percent of SageMath.

BONUS: There is a secret ninth problem in this set! To unlock it,
submit your completed solutions to the first eight problem to
set8.cryptopals@gmail.com with subject "Shackling the Masses with
Drastic Rap Tactics".

Special thanks for this section go to Dan Bernstein, Tanja Lange, and
Watson Ladd, all of whom gave us great advice on problem selection.

// ------------------------------------------------------------

57. Diffie-Hellman Revisited: Small Subgroup Confinement

This set is going to focus on elliptic curves. But before we get to
that, we're going to kick things off with some classic Diffie-Hellman.

Trust me, it's gonna make sense later.

Let's get right into it. First, build your typical Diffie-Hellman key
agreement: Alice and Bob exchange public keys and derive the same
shared secret. Then Bob sends Alice some message with a MAC over
it. Easy as pie.

Use these parameters:

    p =
7199773997391911030609999317773941274322764333428698921736339643928346453700085358802973900485592910475480089726140708102474957429903531369589969318716771
    g =
4565356397095740655436854503483826832136106141639563487732438195343690437606117828318042418238184896212352329118608100083187535033402010599512641674644143

The generator g has order q:

    q = 236234353446506858198510045061214171961

"Order" is a new word, but it just means g^q = 1 mod p. You might
notice that q is a prime, just like p. This isn't mere chance: in
fact, we chose q and p together such that q divides p-1 (the order or
size of the group itself) evenly. This guarantees that an element g of
order q will exist. (In fact, there will be q-1 such elements.)

Back to the protocol. Alice and Bob should choose their secret keys as
random integers mod q. There's no point in choosing them mod p; since
g has order q, the numbers will just start repeating after that. You
can prove this to yourself by verifying g^x mod p = g^(x + k*q) mod p
for any x and k.

The rest is the same as before.

How can we attack this protocol? Remember what we said before about
order: the fact that q divides p-1 guarantees the existence of
elements of order q. What if there are smaller divisors of p-1?

Spoiler alert: there are. I chose j = (p-1) / q to have many small
factors because I want you to be happy. Find them by factoring j,
which is:

    j =
30477252323177606811760882179058908038824640750610513771646768011063128035873508507547741559514324673960576895059570

You don't need to factor it all the way. Just find a bunch of factors
smaller than, say, 2^16. There should be plenty. (Friendly tip: maybe
avoid any repeated factors. They only complicate things.)

Got 'em? Good. Now, we can use these to recover Bob's secret key using
the Pohlig-Hellman algorithm for discrete logarithms. Here's how:

1. Take one of the small factors j. Call it r. We want to find an
   element h of order r. To find it, do:

       h := rand(1, p)^((p-1)/r) mod p

   If h = 1, try again.

2. You're Eve. Send Bob h as your public key. Note that h is not a
   valid public key! There is no x such that h = g^x mod p. But Bob
   doesn't know that.

3. Bob will compute:

       K := h^x mod p

   Where x is his secret key and K is the output shared secret. Bob
   then sends back (m, t), with:

       m := "crazy flamboyant for the rap enjoyment"
       t := MAC(K, m)

4. We (Eve) can't compute K, because h isn't actually a valid public
   key. But we're not licked yet.

   Remember how we saw that g^x starts repeating when x > q? h has the
   same property with r. This means there are only r possible values
   of K that Bob could have generated. We can recover K by doing a
   brute-force search over these values until t = MAC(K, m).

   Now we know Bob's secret key x mod r.

5. Repeat steps 1 through 4 many times. Eventually you will know:

       x = b1 mod r1
       x = b2 mod r2
       x = b3 mod r3
       ...

   Once (r1*r2*...*rn) > q, you'll have enough information to
   reassemble Bob's secret key using the Chinese Remainder Theorem.

// ------------------------------------------------------------

58. Pollard's Method for Catching Kangaroos

The last problem was a little contrived. It only worked because I
helpfully foisted those broken group parameters on Alice and
Bob. While real-world groups may include some small subgroups, it's
improbable to find this many in a randomly generated group.

So what if we can only recover some fraction of the Bob's secret key?
It feels like there should be some way to use that knowledge to
recover the rest. And there is: Pollard's kangaroo algorithm.

This is a generic attack for computing a discrete logarithm (or
"index") known to lie within a certain contiguous range [a, b]. It has
a work factor approximately the square root of the size of the range.

The basic strategy is to try to find a collision between two
pseudorandom sequences of elements. One will start from an element of
known index, and the other will start from the element y whose index
we want to find.

It's important to understand how these sequences are
generated. Basically, we just define some function f mapping group
elements (like the generator g, or a public key y) to scalars (a
secret exponent, like x), i.e.:

    f(y) = <some x>

Don't worry about how f is implemented for now. Just know that it's a
function mapping where we are (some y) to the next jump we're going to
take (some x). And it's deterministic: for a given y, it should always
return the same x.

Then we do a loop like this:

    y := y * g^f(y)

The key thing here is that the next step we take is a function whose
sole input is the current element. This means that if our two
sequences ever happen to visit the same element y, they'll proceed in
lockstep from there.

Okay, let's get a bit more specific. I mentioned we're going to
generate two sequences this way. The first is our control
sequence. This is the tame kangaroo in Pollard's example. We do
something like this:

    xT := 0
    yT := g^b

    for i in 1..N:
        xT := xT + f(yT)
        yT := yT * g^f(yT)

Recall that b is the upper bound on the index of y. So we're starting
the tame kangaroo's run at the very end of that range. Then we just
take N leaps and accumulate our total distance traveled in xT. At the
end of the loop, yT = g^(b + xT). This will be important later.

Note that this algorithm doesn't require us to build a big look-up
table a la Shanks' baby-step giant-step, so its space complexity is
constant. That's kinda neat.

Now: let's catch that wild kangaroo. We'll do a similar loop, this
time starting from y. Our hope is that at some point we'll collide
with the tame kangaroo's path. If we do, we'll eventually end up at
the same place. So on each iteration, we'll check if we're there.

    xW := 0
    yW := y

    while xW < b - a + xT:
        xW := xW + f(yW)
        yW := yW * g^f(yW)

        if yW = yT:
            return b + xT - xW

Take a moment to puzzle out the loop condition. What that relation is
checking is whether we've gone past yT and missed it. In other words,
that we didn't collide. This is a probabilistic algorithm, so it's not
guaranteed to work.

Make sure also that you understand the return statement. If you think
through how we came to the final values for yW and yT, it should be
clear that this value is the index of the input y.

There are a couple implementation details we've glossed over -
specifically the function f and the iteration count N. I do something
like this:

    f(y) = 2^(y mod k)

For some k, which you can play around with. Making k bigger will allow
you to take bigger leaps in each loop iteration.

N is then derived from f - take the mean of all possible outputs of f
and multiply it by a small constant, e.g. 4. You can make the constant
bigger to better your chances of finding a collision at the (obvious)
cost of extra computation. The reason N needs to depend on f is that f
governs the size of the jumps we can make. If the jumps are bigger, we
need a bigger runway to land on, or else we risk leaping past it.

Implement Pollard's kangaroo algorithm. Here are some (less
accommodating) group parameters:

    p =
11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623
    q = 335062023296420808191071248367701059461
    j =
34233586850807404623475048381328686211071196701374230492615844865929237417097514638999377942356150481334217896204702
    g =
622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357

And here's a sample y:

    y =
7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119

The index of y is in the range [0, 2^20]. Find it with the kangaroo
algorithm.

Wait, that's small enough to brute force. Here's one whose index is in
[0, 2^40]:

    y =
9388897478013399550694114614498790691034187453089355259602614074132918843899833277397448144245883225611726912025846772975325932794909655215329941809013733

Find that one, too. It might take a couple minutes.

    ~~ later ~~

Enough about kangaroos, let's get back to Bob. Suppose we know Bob's
secret key x = n mod r for some r < q. It's actually not totally
obvious how to apply this algorithm to get the rest! Because we only
have:

    x = n mod r

Which means:

    x = n + m*r

For some unknown m. This relation defines a set of values that are
spread out at intervals of r, but Pollard's kangaroo requires a
continuous range!

Actually, this isn't a big deal. Because check it out - we can just
apply the following transformations:

    x = n + m*r
    y = g^x = g^(n + m*r)
    y = g^n * g^(m*r)
    y' = y * g^-n = g^(m*r)
    g' = g^r
    y' = (g')^m

Now simply search for the index m of y' to the base element g'. Notice
that we have a rough bound for m: [0, (q-1)/r]. After you find m, you
can plug it into your existing knowledge of x to recover the rest of
the secret.

Take the above group parameters and generate a key pair for Bob. Use
your subgroup-confinement attack from the last problem to recover as
much of Bob's secret as you can. You'll be able to get a good chunk of
it, but not the whole thing. Then use the kangaroo algorithm to run
down the remaining bits.

// ------------------------------------------------------------

59. Elliptic Curve Diffie-Hellman and Invalid-Curve Attacks

I'm not going to show you any graphs - if you want to see one, you can
find them in, like, every other elliptic curve tutorial on the
internet. Personally, I've never been able to gain much insight from
them.

They're also really hard to draw in ASCII.

The key thing to understand about elliptic curves is that they're a
setting analogous in many ways to one we're more familiar with, the
multiplicative integers mod p. So if we learn how certain primitive
operations are defined, we can reason about them using a lot of tools
we already have in our utility belts.

Let's dig in. An elliptic curve E is just an equation like this:

    y^2 = x^3 + a*x + b

The choice of the a and b coefficients defines the curve.

The elements in our group are going to be (x, y) coordinates
satisfying the curve equation. Now, there are infinitely many pairs
like that on the curve, but we only want to think about some of
them. We'll trim our set of points down by considering the curve in
the context of a finite field.

For the moment, it's not too important to know what a finite field
is. You can basically just think of it as "integers mod p" with all
the usual operations you expect: multiplication, division (via modular
inversion), addition, and subtraction.

We'll use the notation GF(p) to talk about a finite field of size
p. (The "GF" is for "Galois field", another name for a finite field.)
When we take a curve E over field GF(p) (written E(GF(p))), what we're
saying is that only points with both x and y in GF(p) are valid.

For example, (3, 6) might be a valid point in E(GF(7)), but it
wouldn't be a valid point in E(GF(5)); 6 is not a member of GF(5).

(3, 4.7) wouldn't be a valid point on either curve, since 4.7 is not
an integer and thus not a member of either field.

What about (3, -1)? This one is on the curve, but remember we're in
some GF(p). So in GF(7), -1 is actually 6. That means (3, -1) and (3,
6) are the same point. In GF(5), -1 is 4, so blah blah blah you get
what I'm saying.

Okay: if these points are going to form a group analogous to the
multiplicative integers mod p, we need to have an analogous set of
primitive functions to work with them.

1. In the multiplicative integers mod p, we combined two elements by
   multiplying them together and taking the remainder modulo p.

   We combine elliptic curve points by adding them. We'll talk about
   what that means in a hot second.

2. We used 1 as a multiplicative identity: y * 1 = y for all y.

   On an elliptic curve, we define the identity O as an abstract
   "point at infinity" that doesn't map to any actual (x, y)
   pair. This might feel like a bit of a hack, but it works.

   On the curve, we have the straightforward rule that P + O = P for
   all P.

   In your code, you can just write something like O := object(),
   since it only ever gets used in pointer comparisons. Or you can use
   some sentinel coordinate that doesn't satisfy the curve equation;
   (0, 1) is popular.

3. We had a modinv function to invert an integer mod p. This acted as
   a stand-in for division. Given y, it finds x such that y * x = 1.

   Inversion is way easier in elliptic curves. Just flip the sign on
   y, and remember that we're in GF(p):

       invert((x, y)) = (x, -y) = (x, p-y)

   Just like with multiplicative inverses, we have this rule on
   elliptic curves:

       P + (-P) = P + invert(P) = O

Incidentally, these primitives, along with a finite set of elements,
are all we need to define a finite cyclic group, which is all we need
to define the Diffie-Hellman function. Not important to understand the
abstract jargon, just FYI.

Let's talk about addition. Here it is:

    function add(P1, P2):
        if P1 = O:
            return P2

        if P2 = O:
            return P1

        if P1 = invert(P2):
            return O

        x1, y1 := P1
        x2, y2 := P2

        if P1 = P2:
            m := (3*x1^2 + a) / 2*y1
        else:
            m := (y2 - y1) / (x2 - x1)

        x3 := m^2 - x1 - x2
        y3 := m*(x1 - x3) - y1

        return (x3, y3)

The first three checks are simple - they pretty much just implement
the rules we have for the identity and inversion.

After that we, uh, use math. You can read more about that part
elsewhere, if you're interested. It's not too important to us, but it
(sort of) makes sense in the context of those graphs I'm not showing
you.

There's one more thing we need. In the multiplicative integers, we
expressed repeated multiplication as exponentiation, e.g.:

    y * y * y * y * y = y^5

We implemented this using a modexp function that walked the bits of
the exponent with a square-and-multiply inner loop.

On elliptic curves, we'll use scalar multiplication to express
repeated addition, e.g.:

    P + P + P + P + P = 5*P

Don't be confused by the shared notation: scalar multiplication is not
analogous to multiplication in the integers. It's analogous to
exponentiation.

Your scalarmult function will look pretty much exactly the same as
your modexp function, except with the primitives swapped out.

Actually, you wanna hear something great? You could define a generic
scale function parameterized over a group that works as a drop-in
implementation for both. Like this:

    function scale(x, k):
        result := identity
        while k > 0:
            if odd(k):
                result := combine(result, x)
            x := combine(x, x)
            k := k >> 1
        return result

The combine function would delegate to modular multiplication or
elliptic curve point depending on the group. It's kind of like the
definition of a group constitutes a kind of interface, and we have
these two different implementations we can swap out freely.

To extend this metaphor, here's a generic Diffie-Hellman:

    function generate_keypair():
        secret := random(1, baseorder)
        public := scale(base, secret)
        return (secret, public)

    function compute_secret(peer_public, self_secret):
        return scale(peer_public, self_secret)

Simplicity itself! The base and baseorder attributes map to g and q in
the multiplicative integer setting. It's pretty much the same on a
curve: we'll have a base point G and its order n such that:

    n*G = O

The fact that these two settings share so many similarities (and can
even share a naive implementation) is great news. It means we already
have a lot of the tools we need to reason about (and attack) elliptic
curves!

Let's put this newfound knowledge into action. Implement a set of
functions up to and including elliptic curve scalar
multiplication. (Remember that all computations are in GF(p), i.e. mod
p.) You can use this curve:

    y^2 = x^3 - 95051*x + 11279326

Over GF(233970423115425145524320034830162017933). Use this base point:

    (182, 85518893674295321206118380980485522083)

It has order 29246302889428143187362802287225875743.

Oh yeah, order. Finding the order of an elliptic curve group turns out
to be a bit tricky, so just trust me when I tell you this one has
order 233970423115425145498902418297807005944. That factors to 2^3 *
29246302889428143187362802287225875743.

Note: it's totally possible to pick an elliptic curve group whose
order is just a straight-up prime number. This would mean that every
point on the curve (except the identity) would have the same order,
since the group order would have no other divisors. The NIST P-curves
are like this.

Our curve has almost-prime order. There's just that small cofactor of
2^3, which is beneficial for reasons we'll cover later. Don't worry
about it for now.

If your implementation works correctly, it should be easy to verify:
remember that multiplying the base point by its order should yield the
group identity.

Implement ECDH and verify that you can do a handshake correctly. In
this case, Alice and Bob's secrets will be scalars modulo the base
point order and their public elements will be points. If you
implemented the primitives correctly, everything should "just work".

Next, reconfigure your protocol from #57 to use it.

Can we apply the subgroup-confinement attacks from #57 in this
setting? At first blush, it seems like it will be pretty difficult,
since the cofactor is so small. We can recover, like, three bits by
sending a point with order 8, but that's about it. There just aren't
enough small-order points on the curve.

How about not on the curve?

Wait, what? Yeah, points *not* on the curve. Look closer at our
combine function. Notice anything missing? The b parameter of the
curve is not accounted for anywhere. This is because we have four
inputs to the calculation: the curve parameters (a, b) and the point
coordinates (x, y). Given any three, you can calculate the fourth. In
other words, we don't need b because b is already baked into every
valid (x, y) pair.

There's a dangerous assumption there: namely, that the peer will
submit a valid (x, y) pair. If Eve can submit an invalid pair, that
really opens up her play: now she can pick points from any curve that
differs only in its b parameter. All she has to do is find some curves
with small subgroups and cherry-pick a few points of small
order. Alice will unwittingly compute the shared secret on the wrong
curve and leak a few bits of her private key in the process.

How do we find suitable curves? Well, remember that I mentioned
counting points on elliptic curves is tricky. If you're very brave,
you can implement Schoof-Elkies-Atkins. Or you can use a computer
algebra system like SageMath. Or you can just use these curves I
generated for you:

    y^2 = x^3 - 95051*x + 210
    y^2 = x^3 - 95051*x + 504
    y^2 = x^3 - 95051*x + 727

They have orders:

    233970423115425145550826547352470124412
    233970423115425145544350131142039591210
    233970423115425145545378039958152057148

They should have a fair few small factors between them. So: find some
points of small order and send them to Alice. You can use the same
trick from before to find points of some prime order r. Suppose the
group has order q. Pick some random point and multiply by q/r. If you
land on the identity, start over.

It might not be immediately obvious how to choose random points, but
you can just pick an x and calculate y. This will require you to
implement a modular square root algorithm; use Tonelli-Shanks, it's
pretty straightforward.

Implement the key-recovery attack from #57 using small-order points
from invalid curves.

// ------------------------------------------------------------

60. Single-Coordinate Ladders and Insecure Twists

All our hard work is about to pay some dividends. Here's a list of
cool-kids jargon you'll be able to deploy after completing this
challenge:

* Montgomery curve
* single-coordinate ladder
* isomorphism
* birational equivalence
* quadratic twist
* trace of Frobenius

Not that you'll understand it all; you won't. But you'll at least be
able to silence crypto-dilettantes on Twitter.

Now, to the task at hand. In the last problem, we implemented ECDH
using a short Weierstrass curve form, like this:

    y^2 = x^3 + a*x + b

For a long time, this has been the most popular curve form. The NIST
P-curves standardized in the 90s look like this. It's what you'll see
first in most elliptic curve tutorials (including this one).

We can do a lot better. Meet the Montgomery curve:

    B*v^2 = u^3 + A*u^2 + u

Although it's almost as old as the Weierstrass form, it's been buried
in the literature until somewhat recently. The Montgomery curve has a
killer feature in the form of a simple and efficient algorithm to
compute scalar multiplication: the Montgomery ladder.

Here's the ladder:

    function ladder(u, k):
        u2, w2 := (1, 0)
        u3, w3 := (u, 1)
        for i in reverse(range(bitlen(p))):
            b := 1 & (k >> 1)
            u2, u3 := cswap(u2, u3, b)
            w2, w3 := cswap(w2, w3, b)
            u3, w3 := ((u2*u3 - w2*w3)^2,
                       u1 * (u2*w3 - w2*u3)^2)
            u2, w2 := ((u2^2 - w2^2)^2,
                       4*u2*w2 * (u2^2 + A*u2*w2 + w2^2))
            u2, u3 := cswap(u2, u3, b)
            w2, w3 := cswap(w2, w3, b)
        return u2 * w2^(p-2)

You are not expected to understand this.

No, really! Most people don't understand it. Instead, they visit the
Explicit-Formulas Database (https://www.hyperelliptic.org/EFD/), the
one-stop shop for state-of-the-art ECC implementation techniques. It's
like cheat codes for elliptic curves. Worth visiting for the
bibliography alone.

With that said, we should try to demystify this a little bit. Here's
the CliffsNotes:

1. Points on a Montgomery curve are (u, v) pairs, but this function
   only takes u as an input. Given *just* the u coordinate of a point
   P, this function computes *just* the u coordinate of k*P. Since we
   only care about u, this is a single-coordinate ladder.

2. So what the heck is w? It's part of an alternate point
   representation. Instead of a coordinate u, we have a coordinate
   u/w. Think of it as a way to defer expensive division (read:
   inversion) operations until the very end.

3. cswap is a function that swaps its first two arguments (or not)
   depending on whether its third argument is one or zero. Choosy
   implementers choose arithmetic implementations of cswap, not
   branching ones.

4. The core of the inner loop is a differential addition followed by a
   doubling operation. Differential addition means we can add two
   points P and Q only if we already know P - Q. We'll take this
   difference to be the input u and maintain it as an invariant
   throughout the ladder. Indeed, our two initial points are:

       u2, w2 := (1, 0)
       u3, w3 := (u, 1)

   Representing the identity and the input u.

5. The return statement performs the modular inversion using a trick
   due to Fermat's Little Theorem:

       a^p     = a    mod p
       a^(p-1) = 1    mod p
       a^(p-2) = a^-1 mod p

6. A consequence of the Montgomery ladder is that we conflate (u, v)
   and (u, -v). But this encoding also conflates zero and
   infinity. Both are represented as zero. Note that the usual
   exceptional case where w = 0 is handled gracefully: our trick for
   doing the inversion with exponentiation outputs zero as expected.

   This is fine: we're still working in a subgroup of prime order.

Go ahead and implement the ladder. Remember that all computations are
in GF(233970423115425145524320034830162017933).

Oh yeah, the curve parameters. You might be thinking that since we're
switching to a new curve format, we also need to pick out a whole new
curve. But you'd be totally wrong! It turns out that some short
Weierstrass curves can be converted into Montgomery curves.

This is because all finite cyclic groups with an equal number of
elements share a kind of equivalence we call "isomorphism". It makes
sense, if you think about it - if the order is the same, all the same
subgroups will be present, and in the same proportions.

So all we need to do is:

1. Find a Montgomery curve with an equal order to our curve.

2. Figure out how to map points back and forth between curves.

You can perform this conversion algebraically. But it's kind of a
pain, so here you go:

    v^2 = u^3 + 534*u^2 + u

Through cunning and foresight, I have chosen this curve specifically
to have a really simple map between Weierstrass and Montgomery
forms. Here it is:

    u = x - 178
    v = y

Which makes our base point:

    (4, 85518893674295321206118380980485522083)

Or, you know. Just 4.

Anyway, implement the ladder. Verify ladder(4, n) = 0. Map some points
back and forth between your Weierstrass and Montgomery representations
and verify them.

One nice thing about the Montgomery ladder is its lack of special
cases. Specifically, no special handling of: P1 = O; P2 = O; P1 = P2;
or P1 = -P2. Contrast that with our Weierstrass addition function and
its battalion of ifs.

And there's a security benefit, too: by ignoring the v coordinate, we
take away a lot of leeway from the attacker. Recall that the ability
to choose arbitrary (x, y) pairs let them cherry-pick points from any
curve they can think of. The single-coordinate ladder robs the
attacker of that freedom.

But hang on a tick! Give this a whirl:

    ladder(76600469441198017145391791613091732004, 11)

What the heck? What's going on here?

Let's do a quick sanity check. Here's the curve equation again:

    v^2 = u^3 + 534*u^2 + u

Plug in u and take the square root to recover v.

You should detect that something is quite wrong. This u does not
represent a point on our curve! Not every u does.

This means that even though we can only submit one coordinate, we
still have a little bit of leeway to find invalid
points. Specifically, an input u such that u^3 + 534*u^2 + u is not a
quadratic residue can never represent a point on our curve. So where
the heck are we?

The other curve we're on is a sister curve called a "quadratic twist",
or simply "the twist". There is actually a whole family of quadratic
twists to our curve, but they're all isomorphic to each
other. Remember that that means they have the same number of points,
the same subgroups, etc. So it doesn't really matter which particular
twist we use; in fact, we don't even need to pick one.

We're mostly interested in the subgroups present on the twist, which
means we need to know how many points it contains. Fortunately, it
turns out to be easier to count the combined set of points on the
curve and its twist at the same time. Let's do it:

1. For every nonzero u up to the modulus p, if u^3 + A*u^2 + u is a
   square in GF(p), there are two points on the original curve.

2. If the above sum is a nonsquare in GF(p), there are two points on
   the twisted curve.

It should be clear that these add up to 2*(p-1) points in total, since
there are p-1 nonzero integers in GF(p) and two points for each. Let's
continue:

3. Both the original curve and its twist have a point (0, 0). This is
   just a regular point, not the group identity.

4. Both the original curve and its twist have an abstract point at
   infinity which serves as the group identity.

So we have 2*p + 2 points across both curves. Since we already know
how many points are on the original curve, we can easily calculate the
order of the twist.

If Alice chose a curve with an insecure twist, i.e. one with a
partially smooth order, then some doors open back up for Eve. She can
choose low-order points on the twisted curve, send them to Alice, and
perform the invalid-curve attack as before.

The only caveat is that she won't be able to recover the full secret
using off-curve points, only a fraction of it. But we know how to
handle that.

So:

1. Calculate the order of the twist and find its small factors. This
   one should have a bunch under 2^24.

2. Find points with those orders. This is simple:

   a. Choose a random u mod p and verify that u^3 + A*u^2 + u is a
      nonsquare in GF(p).

   b. Call the order of the twist n. To find an element of order q,
      calculate Mladder(u, n/q).

3. Send these points to Alice to recover portions of her secret.

4. When you've exhausted all the small subgroups in the twist, recover
   the remainder of Alice's secret with the kangaroo attack.

// ------------------------------------------------------------

61. Duplicate-Signature Key Selection in ECDSA (and RSA)

Suppose you have a message-signature pair. If I give you a public key
that verifies the signature, can you trust that I'm the author?

You shouldn't. It turns out to be pretty easy to solve this problem
across a variety of digital signature schemes. If you have a little
flexibility in choosing your public key, that is.

Let's consider the case of ECDSA.

First, implement ECDSA. If you still have your old DSA implementation
lying around, this should be straightforward. All the same, here's a
refresher if you need it:

    function sign(m, d):
        k := random_scalar(1, n)
        r := (k * G).x
        s := (H(m) + d*r) * k^-1
        return (r, s)

    function verify(m, (r, s), Q):
        u1 := H(m) * s^-1
        u2 := r * s^-1
        R := u1*G + u2*Q
        return r = R.x

Remember that all the scalar operations are mod n, the order of the
base point G. (d, Q) is the signer's key pair. H(m) is a hash of the
message.

Note that the verification function requires arbitrary point
addition. This means your Montgomery ladder (which only performs
scalar multiplication) won't work here. This is no big deal; just fall
back to your old Weierstrass imlpementation.

Once you've got this implemented, generate a key pair for Alice and
use it to sign some message m.

It would be tough for Eve to find a Q' to verify this signature if all
the domain parameters are fixed. But the domain parameters might not
be fixed - some protocols let the user specify them as part of their
public key.

Let's rearrange some terms. Consider this equality:

    R = u1*G + u2*Q

Let's do some regrouping:

    R = u1*G + u2*(d*G)
    R = (u1 + u2*d)*G

Consider R, u1, and u2 to be fixed. That leaves Alice's secret d and
the base point G. Since we don't know d, we'll need to choose a new
pair of values for which the equality holds. We can do it by starting
from the secret and working backwards.

1. Choose a random d' mod n.

2. Calculate t := u1 + u2*d'.

3. Calculate G' := t^-1 * R.

4. Calculate Q' := d' * G'.

5. Eve's public key is Q' with domain parameters (E(GF(p)), n, G').
   E(GF(p)) is the elliptic curve Alice originally chose.

Note that Eve's public key is totally valid: both the base point and
her public point are members of the subgroup of prime order n. Since
E(GF(p)) and n are unchanged from Alice's public key, they should pass
the same validation rules.

Assuming the role of Eve, derive a public key and domain parameters to
verify Alice's signature over the message.

Let's do the same thing with RSA. Same setup: we have some message and
a signature over it. How do we craft a public key to verify the
signature?

Well, first let's refresh ourselves on RSA. Signature verification
looks like this:

    s^e = pad(m) mod N

Where (m, s) is the message-signature pair and (e, N) is Alice's
public key.

So what we're really looking for is the pair (e', N') to make that
equality hold up. If this is starting to look a little familiar, it
should: what we're doing here is looking for the discrete logarithm of
pad(m) with base s.

We know discrete logarithms are easy to solve with Pohlig-Hellman in
groups with many small subgroups. And the choice of group is up to us,
so we can't fail!

But we should exercise some care. If we choose our primes incorrectly,
the discrete logarithm won't exist.

Okay, check the method:

1. Pick a prime p. Here are some conditions for p:

   a. p-1 should be smooth. How smooth is up to you, but you will need
      to find discrete logarithms in each of these subgroups. You can
      use something like Shanks or Pollard's rho to compute these in
      square-root time.

   b. s shouldn't be in any subgroup that pad(m) is not in. If it is,
      the discrete logarithm won't exist. The simplest thing to do is
      make sure they're both primitive roots. To check if an element g
      is a primitive root mod p, check that:

          g^((p-1)/q) != 1 mod p

      For every factor q of p-1.

2. Now pick a prime q. Ensure the same conditions as before, but add these:

   a. Don't reuse any factors of p-1 other than 2. It's possible to
      make this work with repeated factors, but it's a huge
      headache. Better just to avoid it.

   b. Make sure p*q is greater than Alice's modulus N. This is just to
      make sure the signature and padded message will fit under your
      new modulus.

3. Use Pohlig-Hellman to derive ep = e' mod p and eq = e' mod q.

4. Use the Chinese Remainder Theorem to put ep and eq together:

       e' = crt([ep, eq], [p-1, q-1])

5. Your public modulus is N' = p * q.

6. You can derive d' in the normal fashion.

Easy as pie. e' will be a lot larger than the typical public exponent,
but that's still legal.

Since RSA signing and decryption are equivalent operations, you can
use this same technique for other surprising results. Try generating a
random (or chosen) ciphertext and creating a key to decrypt it to a
plaintext of your choice!

// ------------------------------------------------------------

62. Key-Recovery Attacks on ECDSA with Biased Nonces

Back in set 6 we saw how "nonce" is kind of a misnomer for the k value
in DSA. It's really more like an ephemeral key. And distressingly, the
security of your long-term private key hinges on it.

Nonce disclosure? Congrats, you just coughed up your secret key.

Predictable nonce? Ditto.

Even by repeating a nonce you lose everything.

How far can we take this? Turns out, pretty far: even a slight bias in
nonce generation is enough for an attacker to recover your private
key. Let's see how.

First, let's clarify what we mean by a "biased" nonce. All we really
need for this attack is knowledge of a few bits of the nonce. For
simplicity, let's say the low byte of each nonce is zero. So take
whatever code you were using for nonce generation and just mask off
the eight least significant bits.

How does this help us? Let's review the signing algorithm:

    function sign(m, d):
        k := random_scalar(1, q)
        r := (k * G).x
        s := (H(m) + d*r) * k^-1
        return (r, s)

(Quick note: before we used "n" to mean the order of the base
point. In this problem I'm going to use "q" to avoid naming
collisions. Deal with it.)

Focus on the s calculation. Observe that if the low l bits of k are
biased to some constant c, we can rewrite k as b*2^l + c. In our case,
c = 0, so we'll instead rewrite k as b*2^l. This means we can relate
the public r and s values like this:

    s = (H(m) + d*r) * (b*2^l)^-1

Some straightforward algebra gets us from there to here:

    d*r / (s*2^l) = (H(m) / (-s*2^l)) + b

Remember that these calculations are all modulo q, the order of the
base point. Now, let's define some stand-ins:

    t =    r / ( s*2^l)
    u = H(m) / (-s*2^l)

Now our equation above can be written like this:

    d*t = u + b

Remember that b is small. Whereas t, u, and the secret key d are all
roughly the size of q, b is roughly q/2^l. It's a rounding
error. Since b is so small, we can basically just ignore it and say:

    d*t ~ u

In other words, u is an approximation for d*t mod q. Let's massage the
numbers some more. Since this is mod q, we can instead say this:

    d*t ~ u + m*q
      0 ~ u + m*q - d*t

That sum won't really be zero - it's just an approximation. But it
will be less than some bound, say q/2^l. The point is that it will be
very small relative to the other terms in play.

We can use this property to recover d if we have enough (u, t)
pairs. But to do that, we need to know a little bit of linear
algebra. Not too much, I promise.

Linear algebra is about vectors. A vector could be almost anything,
but for simplicity we'll say a vector is a fixed-length sequence of
numbers. There are two main things we can do with vectors: we can add
them and we can multiply them by scalars. To add two vectors, simply
sum their pairwise components. To multiply a vector by a scalar k,
simply add it to itself k times. (Equivalently, multiply each of its
elements by the scalar.) Together, these operations are called linear
combinations.

If we have a set of vectors, we say they span a vector space. The
vector space is simply the full range of possible vectors we can
generate by adding and scaling the vectors in our set. We call a
minimal spanning set a basis for the vector space. "Minimal" means
that dropping any of our vectors from the set would result in a
smaller vector space. Added vectors would either be redundant
(i.e. they could be defined as sums of existing vectors) or they would
give us a larger vector space. So you can think of a basis as "just
right" for the vector space it spans.

We're only going to use integers as our scalars. A vector space
generated using only integral scalars is called a lattice. It's best
to picture this in the two-dimensional plane. Suppose our set of
vectors is {(3, 4), (2, 1)}. The lattice includes all integer
combinations of these two pairs. You can graph this out on paper to
get the idea; you should end up with a polka dot pattern of sorts.

We said that a basis is just the right size for the vector space it
spans, but that shouldn't be taken to imply uniqueness. Indeed, any of
the lattices we will care about have infinite possible bases. The only
requirements are that the basis spans the space and the basis is
minimal in size. In that sense, all bases of a given lattice are
equal.

But some bases are more equal than others. In practice, people like to
use bases comprising shorter vectors. Here "shorter" means, roughly,
"containing smaller components on average". A handy measuring stick
here is the Euclidean norm: simply take the dot product of a vector
with itself and take the square root. Or don't take the square root, I
don't care. It won't affect the ordering.

Why do people like these smaller bases? Mostly because they're more
efficient for computation. Honestly, it doesn't matter too much why
people like them. The important thing is that we have relatively
efficient methods for "reducing" a basis. Given an input basis, we can
produce an equivalent-but-with-much-shorter-vectors basis. How much
shorter? Well, maybe not the very shortest possible, but pretty darn
short.

This implies a really neat approach to problem-solving:

1. Encode your problem space as a set of vectors forming the basis for
   a lattice. The lattice you choose should contain the solution
   you're looking for as a short vector. You don't need to know the
   vector (obviously, since you're looking for it), you just need to
   know that it exists as some integral combination of your basis
   vectors.

2. Derive a reduced basis for the lattice. We'll come back to this.

3. Fish your solution vector out of the reduced basis.

4. That's it.

Wait, that's it? Yeah, you heard me - lattice basis reduction is an
incredibly powerful technique. It single-handedly shattered knapsack
cryptosystems back in the '80s, and it's racked up a ton of trophies
since then. As long as you can define a lattice containing a short
vector that encodes the solution to your problem, you can put it to
work for you.

Obviously, defining the lattice is the tricky bit. How do we encode
ECDSA? Well, when we left off, we had the following approximation:

    0 ~ u + m*q - d*t

Suppose we collect a bunch of signatures. Then that one approximation
becomes many:

    0 ~ u1 + m1*q - d*t1
    0 ~ u2 + m2*q - d*t2
    0 ~ u3 + m3*q - d*t3
    0 ~ u4 + m4*q - d*t4
    0 ~ u5 + m5*q - d*t5
    0 ~ u6 + m6*q - d*t6
    ...
    0 ~ un + mn*q - d*tn

The coefficient for each u is always 1, and the coefficient for t is
always the secret key d. So it seems natural that we should line those
up in two vectors:

    bt = [ t1 t2 t3 t4 t5 t6 ... tn ]

    bu = [ u1 u2 u3 u4 u5 u6 ... un ]

Each approximation also contains some factor of q. But the coefficient
m is different each time. That means we'll need a separate vector for
each one:

    b1 = [  q  0  0  0  0  0 ...  0 ]

    b2 = [  0  q  0  0  0  0 ...  0 ]

    b3 = [  0  0  q  0  0  0 ...  0 ]

    b4 = [  0  0  0  q  0  0 ...  0 ]

    b5 = [  0  0  0  0  q  0 ...  0 ]

    b6 = [  0  0  0  0  0  q ...  0 ]

            ...              ...

    bn = [  0  0  0  0  0  0 ...  q ]

    bt = [ t0 t1 t2 t3 t4 t5 ... tn ]

    bu = [ u0 u1 u2 u3 u4 u5 ... un ]

See how the columns cutting across our row vectors match up with the
approximations we collected above? Notice also that the lattice
defined by this basis contains at least one reasonably short vector
we're interested in:

    bu - d*bt + m0*b1 + m1*b2 + m2*b3 ... + mn*bn

But we have a problem: even if this vector is included in our reduced
basis, we have no way to identify it. We can solve this by adding a
couple new columns.

    b1 = [  q  0  0  0  0  0 ...  0  0  0 ]

    b2 = [  0  q  0  0  0  0 ...  0  0  0 ]

    b3 = [  0  0  q  0  0  0 ...  0  0  0 ]

    b4 = [  0  0  0  q  0  0 ...  0  0  0 ]

    b5 = [  0  0  0  0  q  0 ...  0  0  0 ]

    b6 = [  0  0  0  0  0  q ...  0  0  0 ]

            ...              ...

    bn = [  0  0  0  0  0  0 ...  q  0  0 ]

    bt = [ t0 t1 t2 t3 t4 t5 ... tn ct  0 ]

    bu = [ u0 u1 u2 u3 u4 u5 ... un  0 cu ]

We've added two new columns with sentinel values in bt and bu. This
will allow us to determine whether these two vectors are included in
any of the output vectors and in what proportions. (That's not the
only problem this solves. Our last set of vectors wasn't really a
basis, because we had n+2 vectors of degree n, so there were clearly
some redundancies in there.)

We can identify the vector we're looking for by looking for cu in the
last slot of each vector in our reduced basis. Our hunch is that the
adjacent slot will contain -d*ct, and we can divide through by -ct to
recover d.

Okay. To go any further, we need to dig into the nuts and bolts of
basis reduction. There are different strategies for finding a reduced
basis for a lattice, but we're going to focus on a simple and
efficient polynomial-time algorithm: Lenstra-Lenstra-Lovasz (LLL).

Most people don't implement LLL. They use a library, of which there
are several excellent ones. NTL is a popular choice.

For instructional purposes only, we're going to write our own.

Here's some pseudocode:

    function LLL(B, delta):
        B := copy(B)
        Q := gramschmidt(B)

        function mu(i, j):
            v := B[i]
            u := Q[j]
            return (v*u) / (u*u)

        n := len(B)
        k := 1

        while k < n:
            for j in reverse(range(k)):
                if abs(mu(k, j)) > 1/2:
                    B[k] := B[k] - round(mu(k, j))*B[j]
                    Q := gramschmidt(B)

            if (Q[k]*Q[k]) >= (delta - mu(k, k-1)^2) * (Q[k-1]*Q[k-1]):
                k := k + 1
            else:
                B[k], B[k-1] := B[k-1], B[k]
                Q := gramschmidt(B)
                k := max(k-1, 1)

        return B

B is our input basis. Delta is a parameter such that 0.25 < delta <=
1. You can just set it to 0.99 and forget about it.

Gram-Schmidt is an algorithm to convert a basis into an equivalent
basis of mutually orthogonal (a fancy word for "perpendicular")
vectors. It's dead simple:

    function proj(u, v):
        if u = 0:
            return 0
        return ((v*u) / (u*u)) * u

    function gramschmidt(B):
        Q := []
        for i, v in enumerate(B):
            Q[i] := v - sum(proj(u, v) for u in Q[:i])
        return Q

Proj finds the projection of v onto u. This is basically the part of v
going in the same "direction" as u. If u and v are orthogonal, this is
the zero vector. Gram-Schmidt orthogonalizes a basis by iterating over
the original and shaving off these projections.

Back to LLL. The best way to get a sense for how and why it works is
to implement it and test it on some small examples with lots of debug
output. But basically: we walk up and down the basis B, comparing each
vector b against the orthogonalized basis Q. Whenever we find a vector
q in Q that mostly aligns with b, we shave off an integral
approximation of q's projection onto b. Remember that the lattice
deals in integral coefficients, and so must we. After each iteration,
we use some heuristics to decide whether we should move forward or
backward in B, whether we should swap some rows, etc.

One more thing: the above description of LLL is very naive and
inefficient. It probably won't be fast enough for our purposes, so you
may need to optimize it a little. A good place to start would be not
recalculating the entire Q matrix on every update.

Here's a test basis:

    b1 = [  -2    0    2    0]
    b2 = [ 1/2   -1    0    0]
    b3 = [  -1    0   -2  1/2]
    b4 = [  -1    1    1    2]

It reduces to this (with delta = 0.99):

    b1 = [ 1/2   -1    0    0]
    b2 = [  -1    0   -2  1/2]
    b3 = [-1/2    0    1    2]
    b4 = [-3/2   -1    2    0]

I forgot to mention: you'll want to write your implementation to work
on vectors of rationals. If you have infinite-precision floats,
those'll work too.

All that's left is to tie up a few loose ends. First, how do we choose
our sentinel values ct and cu? This is kind of an implementation
detail, but we want to "balance" the size of the entries in our target
vector. And since we expect all of the other entries to be roughly
size q/2^l:

    ct = 1/2^l
    cu = q/2^l

Remember that ct will be multiplied by -d, and d is roughly the size
of q.

Okay, you're finally ready to run the attack:

1. Generate your ECDSA secret key d.

2. Sign a bunch of messages using d and your biased nonce generator.

3. As the attacker, collect your (u, t) pairs. You can experiment with
   the amount. With an eight-bit nonce bias, I get good results with
   as few as 20 signatures. YMMV.

4. Stuff your values into a matrix and reduce it with LLL. Consider
   playing with some smaller matrices to get a sense for how long this
   will take to run.

5. In the reduced basis, find a vector with q/2^l as the final
   entry. There's a good chance it will have -d/2^l as the
   second-to-last entry. Extract d.

// ------------------------------------------------------------

63. Key-Recovery Attacks on GCM with Repeated Nonces

GCM is the most widely deployed block cipher mode for authenticated
encryption with associated data (AEAD). It's basically just CTR mode
with a weird MAC function wrapped around it. The MAC function works by
evaluating a polynomial over GF(2^128).

Remember how much trouble a repeated nonce causes for CTR mode
encryption? The same thing is true here: an attacker can XOR
ciphertexts together and recover plaintext using statistical methods.

But there's an even more devastating consequence for GCM: it leaks the
authentication key immediately!

Here's the high-level view:

1. The GCM MAC function (GMAC) works by building up a polynomial whose
   coefficients are the blocks of associated data (AD), the blocks of
   ciphertext (C), a block encoding the length of AD and C, and a
   block used to "mask" the output. Sort of like this:

       AD*y^3 + C*y^2 + L*y + S

   To calculate the MAC, we plug in the authentication key for y and
   evaluate.

2. AD, C, and their respective lengths are known. For a given message,
   the attacker knows everything about the MAC polynomial except the
   masking block

3. The masking block is generated using only the key and the nonce. If
   the nonce is repeated, the mask is the same. If we can collect two
   messages encrypted under the same nonce, they'll have used the same
   mask.

4. In this field, addition is XOR. We can XOR our two messages
   together to "wash out" the mask and recover a known polynomial with
   the authentication key as a root. We can factor the polynomial to
   recover the authentication key immediately.

The last step probably feels a little magical, but you don't actually
need to understand it to implement the attack: you can literally just
plug the right values into a computer algebra system like SageMath and
hit "factor".

But that's not satisfying. You didn't come this far to beat the game
on Easy Mode, did you?

I didn't think so. Now, let's dig into that MAC function. Like I said,
a polynomial over GF(2^128).

So far, all the fields we've worked with have been of prime size p,
i.e. GF(p). It turns out we can construct GF(q) for any q = p^k for
any positive integer k. GF(p) = GF(p^1) is just one form that's common
in cryptography. Another is GF(2^k), and in this case we have
GF(2^128).

For GF(2^k), we'll represent its elements as polynomials with
coefficients in GF(2). That just means each coefficient will be 0 or
1. Before we start talking about a particular field (i.e. a particular
choice of k), let's just talk about these polynomials. Here are some
of them:

                    0
                    1
                x
                x + 1
          x^2
          x^2     + 1
          x^2 + x
          x^2 + x + 1
    x^3
    x^3           + 1
    x^3       + x
    x^3       + x + 1
    x^3 + x^2
    x^3 + x^2     + 1
    x^3 + x^2 + x
    x^3 + x^2 + x + 1
    ...

And so forth.

If you squint a little they look like the binary expansions of the
integers counting up from zero. This is convenient because it gives us
an obvious choice of representation in unsigned integers.

Now we need some primitive functions for operating on these
polynomials. Let's tool up:

1. Addition and subtraction between GF(2) polynomials are really
   simple: they're both just the XOR function.

2. Multiplication and division are a little trickier, but they both
   just approximate the algorithms you learned in grade school. Here
   they are:

       function mul(a, b):
            p := 0

            while a > 0:
                if a & 1:
                    p := p ^ b

                a := a >> 1
                b := b << 1

            return p

        function divmod(a, b):
            q, r := 0, a

            while deg(r) >= deg(b):
                d := deg(r) - deg(b)
                q := q ^ (1 << d)
                r := r ^ (b << d)

            return q, r

   deg(a) is a function returning the degree of a polynomial. For the
   polynomial x^4 + x + 1, it should return 4. For 1, it should return
   0. For 0, it should return some negative value.

Now that we have a small nucleus of functions to work on polynomials
in GF(2), let's see how we can use them to represent elements in
GF(2^k). To be concrete, let's say k = 4.

Our set of elements is the same one enumerated above:

                    0
                    1
                x
                x + 1
          x^2
          x^2     + 1
          x^2 + x
          x^2 + x + 1
    x^3
    x^3           + 1
    x^3       + x
    x^3       + x + 1
    x^3 + x^2
    x^3 + x^2     + 1
    x^3 + x^2 + x
    x^3 + x^2 + x + 1

Addition and subtraction are unchanged. We still just XOR elements
together.

Multiplication is different. As with fields of size p, we need to
perform modular reductions after each multiplication to keep our
elements in range. Our modulus will be x^4 + x + 1. If that seems a
little arbitrary, it is - we could use any fourth-degree monic
polynomial that's irreducible over GF(2). An irreducible polynomial is
sort of analogous to a prime number in this setting.

Here's a naive modmul:

    function modmul(a, b, m):
        p := mul(a, b)
        q, r := divmod(p, m)
        return r

In practice, we'll want to be more efficient. So we'll interleave the
steps of the multiplication with steps of the reduction:

    function modmul(a, b, m):
        p := 0

        while a > 0:
            if a & 1:
                p := p ^ b

            a := a >> 1
            b := b << 1

            if deg(b) = deg(m):
                b := b ^ m

        return p

You can implement both versions to prove to yourself that the output
is the same.

Division is also different. Remember that in fields of size p we
defined it as multiplication by the inverse. So you'll need to write a
modinv function. It should be pretty easy to translate your existing
integer modinv function. I'll leave that to you.

You may find yourself in want of other functions you take for granted
in the integer setting, e.g. modexp. Most of these should have
straightforward equivalents in our polynomial setting. Do what you
need to.

Okay, now that you are the master of GF(2^k), we can finally talk
about GCM. Like I said (many words ago): CTR mode for encryption,
weird MAC in GF(2^128).

Here's the modulus for that field:

    x^128 + x^7 + x^2 + x + 1

The size of this field was chosen very specifically to match up with
the width of a 128-bit block cipher. We can convert a block into a
field element trivially; the leftmost bit is the coefficient of x^0,
and so on.

I described the MAC at a very high level. Here's a more detailed view:

1. Take your AES key and use it to encrypt a block of zeros:

       h := E(K, 0)

   h is your authentication key. Convert it into a field element.

1. Zero-pad the bytes of associated data (AD) to be divisible by the
   block length. If it's already aligned on a block, leave it
   alone. Do the same with the ciphertext. Chain them together so you
   have something like:

       a0 || a1 || c0 || c1 || c2

2. Add one last block describing the length of the AD and the length
   of the ciphertext. Original lengths, not padded lengths; bit
   lengths, not byte lengths. Like this:

       len(AD) || len(C)

3. Take h and your string of blocks and do this:

       g := 0
       for b in bs:
           g := g + b
           g := g * h

   Convert the blocks into field elements first, of course. The
   resulting value of g is a keyed hash of the input blocks.

4. GCM takes a 96-bit nonce. Do this with it:

       s := E(K, nonce || 1)
       t := g + s

   Conceptually, we're masking the hash with a nonce-derived secret
   value. More on that later.

   t is your tag. Convert it back to a block and ship it.

Implement GCM. Use AES-128 as your block cipher. You can probably
reuse whatever you had before for CTR mode. The important new thing to
implement here is the MAC. The above description is brief and
informal; check out the spec for the finer points. Since you've
already got the tools for working in GF(2^k), this shouldn't take too
long.

Okay. Let's rethink our view of the MAC. We'll use our example payload
from above. Here it is:

    t = (((((((((((h * a0) + a1) * h) + c0) * h) + c1) * h) + c2) * h) +
len) * h) + s

Kind of a mouthful. Let's rewrite it:

    t = a0*h^6 + a1*h^5 + c0*h^4 + c1*h^3 + c2*h^2 + len*h + s

In other words, we calculate the MAC by constructing this polynomial:

    f(y) = a0*y^6 + a1*y^5 + c0*y^4 + c1*y^3 + c2*y^2 + len*y + s

And computing t = f(h).

Remember: as the attacker, we don't know that whole polynomial. We
know all the AD and ciphertext coefficients, and we know t = f(h), but
we don't know the mask s.

What happens if we repeat a nonce? Let's posit this additional payload
encrypted under the same nonce:

    b0 || d0 || d1

That's one block of AD and two blocks of ciphertext. The MAC will look
like this:

    t = b0*h^4 + d0*h^3 + d1*h^2 + len*h + s

Let's put them side by side (and rewrite them a little):

    t0 = a0*h^6 + a1*h^5 + c0*h^4 + c1*h^3 + c2*h^2 + l0*h + s
    t1 =                   b0*h^4 + d0*h^3 + d1*h^2 + l1*h + s

See how the s masks are identical? They depend only on the nonce and
the encryption key. Since addition is XOR in our field, we can add
these two equations together and that mask will wash right out:

    t0 + t1 = a0*h^6 + a1*h^5 + (c0 + b0)*h^4 + (c1 + d0)*h^3 +
              (c2 + d1)*h^2 + (l0 + l1)*h

Finally, we'll collect all the terms on one side:

    0 = a0*h^6 + a1*h^5 + (c0 + b0)*h^4 + (c1 + d0)*h^3 +
        (c2 + d1)*h^2 + (l0 + l1)*h + (t0 + t1)

And rewrite it as a polynomial in y:

    f(y) = a0*y^6 + a1*y^5 + (c0 + b0)*y^4 + (c1 + d0)*y^3 +
           (c2 + d1)*y^2 + (l0 + l1)*y + (t0 + t1)

Now we know a polynomial f(y), and we know that f(h) = 0. In other
words, the authentication key is a root. That means all we have to do
is factor the equation to get an extremely short list of candidates
for the authentication key. This turns out not to be so hard, but we
will need some more tools.

First, we need to be able to operate on polynomials with coefficients
in GF(2^128).

(Don't get confused: before we used polynomials with coefficients in
GF(2) to represent elements in GF(2^128). Now we're building on top of
that to work with polynomials with coefficients in GF(2^128).)

The simplest representation is probably just an array of field
elements. The algorithms are all going to be basically the same as
above, so I'm not going to reiterate them here. The only difference is
that you will need to call your primitive functions for GF(2^k)
polynomials in place of your language's built-in arithmetic operators.

With that out of the way, let's get factoring. Factoring a polynomial
over a finite field means separating it out into smaller polynomials
that are irreducible over the field. Remember that irreducible
polynomials are sort of like prime numbers.

To factor a polynomial, we proceed in three (well, four) phases:

0. As a preliminary step, we need to convert our polynomial to a monic
   polynomial. That just means the leading coefficient is 1. So take
   your polynomial and divide it by the coefficient of the leading
   term. You can save the coefficient as a degree zero factor if you
   want, but it's not really important for our purposes.

1. This is the real first step: we perform a square-free
   factorization. We find any doubled factors (i.e. "squares") and
   split them out.

2. Next, We take each of our square-free polynomials and find its
   distinct-degree factorization. This separates out polynomials that
   are products of smaller polynomials of equal degree. So if our
   polynomial has three irreducible factors of degree four, this will
   separate out a polynomial of degree twelve that is the product of
   all of them.

3. Finally, we take each output from the last step and perform an
   equal-degree factorization. This is pretty much like it sounds. In
   that last example, we'd take that twelfth-degree polynomial and
   factor it into its fourth-degree components.

Square-free factorization and distinct-degree factorization are both
easy to implement. Just find them on Wikipedia and go to town.

I want to focus on equal-degree factorization. Meet Cantor-Zassenhaus:

    function edf(f, d):
        n := deg(f)
        r := n / d
        S := {f}

        while len(S) < r:
            h := random_polynomial(1, f)
            g := gcd(h, f)

            if g = 1:
                g := h^((q^d - 1)/3) - 1 mod f

            for u in S:
                if deg(u) = d:
                    continue

                if gcd(g, u) =/= 1 and gcd(g, u) =/= u:
                    S := union(S - {u}, {gcd(g, u), u / gcd(g, u)})

        return S

It's kind of brilliant.

Remember earlier that we said a finite field of size p^k can be
represented as a polynomial in GF(p) modulo any monic, irreducible
degree-k polynomial? Take a moment to convince yourself that a field
of size q^d can be represented by polynomials in GF(q) modulo a
polynomial of degree d.

f is the product of r polynomials of degree d. Each of them is a valid
modulus for a finite field of size q^d. And each field of that size
contains a multiplicative group of size q^d - 1. And since q^d - 1 is
always divisible by 3 (in our case), each group of that size has a
subgroup of size 3. It contains the multiplicative identity (1) and
two other elements.

We have a simple trick for forcing elements into that subgroup: simply
raise them to the exponent (q^d - 1)/3. When we force our random
element into that subgroup, there's a 1/3 chance we'll land on 1. This
means that when we subtract 1, there's a 1/3 chance we'll be sitting
on 0.

The only hitch is that we don't know what these moduli are. But we
don't need to! Since f is their product, we can perform these
operations mod f and implicitly apply the Chinese Remainder Theorem.

So we do it: generate some polynomial, raise it to (q^d - 1)/3, and
subtract 1. (All modulo f, of course.) Compute the GCD of this
polynomial and each of our remaining composites. For any given
remaining factor, there's a 1/3 chance our polynomial is a
multiple. In other words, our factors should reveal themselves pretty
quickly.

Just keep doing this until the whole thing is factored into
irreducible parts.

Once the polynomial is factored, we can pick out the authentication
key at our leisure. There will be at least one first-degree
factor. Like this:

    y + c

Where c is a constant element of GF(2^128). It's also a candidate for
the key. It's possible you'll end up with a few first-degree factors
like this. The key will be one of them.

If you do have more than one candidate, there are two ways to narrow
the list:

1. Recover another pair of messages encrypted under the same
   nonce. Perform the factorization again and identify the common
   factors. The key will probably be the only one.

2. Just attempt a forgery with each candidate. This is probably
   easier.

// ------------------------------------------------------------

64. Key-Recovery Attacks on GCM with a Truncated MAC

This one is my favorite.

It's somewhat common to use a truncated MAC tag. For instance, you
might be authenticating with HMAC-SHA256 and shorten the tag to 128
bits. The idea is that you can save some bandwidth or storage and
still have an acceptable level of security.

This is a totally reasonable thing to want.

In some protocols, you might take this to the extreme. If two parties
are exchanging lots of small packets, and the value of forging any one
packet is pretty low, they might use 16-bit tag and expect 16 bits of
security.

In GCM, this is a disaster.

To see how, we'll first review the GCM MAC function. We make a
calculation like this:

    t = s + c1*h + c2*h^2 + c3*h^3 + ... + cn*h^n

We're making some notational changes here for convenience; most
notably, we using one-based indexing. The c1 block encodes the length,
and [c2, cn] are blocks of ciphertext.

We'll also ignore the possibility of AD blocks here, since they don't
matter too much for our purposes.

Recall that our coefficients (and our authentication key h) are
elements in GF(2^128). We've seen a few different representations,
namely polynomials in GF(2) and unsigned integers. To this we'll add
one more: a 128-degree vector over GF(2). It's basically just a bit
vector; not so different from either of our previous representations,
but we want to get our heads in linear algebra mode.

The key concept we want to explore here is that certain operations in
GF(2^128) are linear.

One of them is multiplication by a constant. Suppose we have this
function:

    f(y) = c*y

With c and y both elements of GF(2^128). This function is linear in
the bits of y. That means that this function is equivalent:

    g(y) = c*y[0] + c*y[1] + c*y[2] + ... + c*y[127]

Any linear function can be represented by matrix multiplication. So if
we think of y as a vector, we can construct a matrix Mc such that:

    Mc*y = f(y) = c*y

To construct Mc, just calculate c*1, c*x, c*x^2, ..., c*x^127, convert
each product to a vector, and let the vectors be the columns of your
matrix. You can verify by performing the matrix multiplication against
y and checking the result. I'm going to assume you either know how
matrix multiplication works or have access to Wikipedia to look it up.

Squaring is also linear. This is because (a + b)^2 = a^2 + b^2 in
GF(2^128). Again, this means we can replace this function:

    f(y) = y^2

With a matrix multiplication. Again, compute 1^2, x^2, (x^2)^2, ...,
(x^127)^2, convert the results to vectors, and make them your
columns. Then verify that:

    Ms*y = f(y) = y^2

Okay, let's put these matrices on the back burner for now. To forge a
ciphertext c', we'll start with a valid ciphertext c and flip some
bits, hoping that:

    sum(ci * h^i) = sum(ci' * h^i)

Another way to write this is:

    sum((ci - ci') * h^i) = 0

If we let ei = ci - ci', we can simplify this:

    sum(ei * h^i) = 0

Note that if we leave a block unmolested, then ei = ci - ci' =
0. We're going to leave most ei = 0. In fact, we're only going to flip
bits in blocks di = e(2^i). These are blocks d0, d1, d2, ..., dn (note
that we're back to zero-based indexing) such that:

    sum(di * h^(2^i)) = 0

We hope it equals zero, anyway. Maybe it's better to say:

    sum(di * h^(2^i)) = e

Where e is some error polynomial. In other words, the difference in
the MAC tag induced by our bit-flipping.

At this point, we'll recall the matrix view. Recall that
multiplications by a constant and squaring operations are both
linear. That means we can rewrite the above equation as a linear
operation on h:

    sum(Mdi * Ms^i * h) = e
    sum(Mdi * Ms^i) * h = e

    Ad = sum(Mdi * Ms^i)

    Ad * h = e

We want to find an Ad such that Ad * h = 0.

Let's think about how the bits in the vector e are calculated. This
just falls out of the basic rules of matrix multiplication:

    e[0] = Ad[0] * h
    e[1] = Ad[1] * h
    e[2] = Ad[2] * h
    e[3] = Ad[3] * h
    ...

In other words, e[i] is the inner product of row i of Ad with h. If we
can force rows of Ad to zero, we can force terms of the error
polynomial to zero. Every row we force to zero will basically double
our chances of a forgery.

Suppose the MAC is 16 bits. If we can flip bits and force eight rows
of Ad to zero, that's eight bits of the MAC we know are right. We can
flip whatever bits are left over with a 2^-8 chance of a forgery, way
better than the expected 2^-16!

It turns out to be really easy to force rows of Ad to zero. Ad is the
sum of a bunch of linear operations. That means we can simply
determine which bits of d0, d1, ..., dn affect which bits of Ad and
flip them accordingly.

Actually, let's leave d0 alone. That's the block that encodes the
ciphertext length. Things could get tricky pretty quickly if we start
messing with it.

We still have d1, ..., dn to play with. That means n*128 bits we can
flip. Since the rows of Ad are each 128 bits, we'll have to settle for
forcing n-1 of them to zero. We need some bits left over to play with.

Check the strategy: we'll build a dependency matrix T with n*128
columns and (n-1)*128 rows. Each column represents a bit we can flip,
and each row represents a cell of Ad (reading left-to-right,
top-to-bottom). The cells where they intersect record whether a
particular free bit affects a particular bit of Ad.

Iterate over the columns. Build the hypothetical Ad you'd get by
flipping only the corresponding bit. Iterate over the first (n-1)*128
cells of Ad and set the corresponding cells in this column of T.

After doing this for each column, T will be full of ones and
zeros. We're looking for sets of bit flips that will zero out those
first n-1 rows. In other words, we're looking for solutions to this
equation:

    T * d = 0

Where d is a vector representing all n*128 bits you have to play with.

If you know a little bit of linear algebra, you'll know that what we
really want to find is a basis for N(T), the null space of T. The null
space is exactly that set of vectors that solve the equation
above. Just what we're looking for. Recall that a basis is a minimal
set of vectors whose linear combinations span the whole space. So if
we find a basis for N(T), we can just take random combinations of its
vectors to get viable candidates for d.

Finding a basis for the null space is not too hard. What you want to
do is transpose T (i.e. flip it across its diagonal) and find the
reduced row echelon form using Gaussian elimination. Now perform the
same operations on an identity matrix of size n*128. The rows that
correspond to the zero rows in the reduced row echelon form of T
transpose form a basis for N(T).

Gaussian elimination is pretty simple; you can more or less figure it
out yourself once you know what it's supposed to do.

Now that we have a basis for N(T), we're ready to start forging
messages. Take a random vector from N(T) and decode it to a bunch of
bit flips in your known good ciphertext C. (Remember that you'll be
flipping bits only in the blocks that are multiplied by h^(2*i) for
some i.) Send the adjusted message C' to the oracle and see if it
passes authentication. If it fails, generate a new vector and try
again.

If it succeeds, we've gained more than just an easy forgery. Examine
your matrix Ad. It should be a bunch of zero rows followed by a bunch
of nonzero rows. We care about the nonzero rows corresponding to the
bits of the tag. So if your tag is 16 bits, and you forced eight bits
to zero, you should have eight nonzero rows of interest.

Pick those rows out and stuff them in a matrix of their own. Call it,
I don't know, K. Here's something neat we know about K:

    K * h = 0

In other words, h is in the null space of K! In our example, K is an
8x128 matrix. Assuming all its rows are independent (none of them is a
combination of any of the others), N(K) is a 120-dimensional subspace
of the larger 128-dimensional space. Since we know h is in there, the
range of possible values for h went from 2^128 to 2^120.

2^120 is still a lot of values, but hey - it's a start.

If we can produce more forgeries, we can find more vectors to add to
K, reducing the range of values further and further. And check this
out: our newfound knowledge of h is going to make the next forgery
even easier. Find a basis for N(K) and put the vectors in the columns
of a matrix. Call it X. Now we can rewrite h like this:

    h = X * h'

Where h' is some unknown 120-bit vector. Now instead of saying:

    Ad * h = e

We can say:

    Ad * X * h' = e

Whereas Ad is a 128x128 matrix, X is 128x120. And Ad * X is also
128x120. Instead of zeroing out 128-degree row vectors, now we can
zero out 120-degree vectors. Since we still have the same number of
bits to play with, we can (maybe) zero out more rows than before. The
general picture is that if we have n*128 bits to play with, we can
zero out (n*128) / (ncols(X)) rows. Just remember to leave at least
one nonzero row in each attempt; otherwise you won't learn anything
new.

So: start over and build a new T matrix, but this time to nullify rows
of Ad * X. Forge another message and harvest some new linear equations
on h. Stuff them in K and recalculate X.

Lather, rinse, repeat.

The endgame comes when K has 127 linearly independent rows. N(K) will
be a 1-dimensional subspace containing exactly one nonzero vector, and
that vector will be h.

Let's try it out:

1. Build a toy system with a 32-bit MAC. This is the smallest tag
   length NIST defines in the GCM specification.

2. Generate valid messages of 2^17 blocks for the attacker to play
   with.

3. As the attacker, build your matrices, forge messages, and recover
   the key. You should be able to zero out 16 bits of each tag to
   start, and you'll only gain leverage from there.

// ------------------------------------------------------------

Here's a few interesting papers:

1. "A Key Recovery Attack on Discrete Log-based Schemes Using a Prime
Order Subgroup" by Chae Hoon Lim and Pil Joong Lee.

2. "Monte Carlo Methods for Index Computation (mod p)" by John Pollard
(http://www.ams.org/journals/mcom/1978-32-143/S0025-5718-1978-0491431-9/S0025-5718-1978-0491431-9.pdf).

3. "Differential Fault Attacks on Elliptic Curve Cryptosystems (Extended
Abstract" by Ingrid Biehl, Bernd Meyer, and Volker Muller
(https://www.iacr.org/archive/crypto2000/18800131/18800131.pdf).

4. "Speeding the Pollard and Elliptic Curve Methods of Factorization" by
Peter L. Montgomery
(http://www.ams.org/journals/mcom/1987-48-177/S0025-5718-1987-0866113-7/S0025-5718-1987-0866113-7.pdf).

5. "Unknown Key-Share Attacks on the Station-to-Station (STS) Protocol"
by Simon Blake-Wilson and Alfred Menezes.

6. "The Insecurity of the Digital Signature Algorithm with Partially
Known Nonces" by Phong Q. Nguyen and Igor E. Shparlinski.

7. "Authentication Failures in the NIST version of GCM" by Antoine Joux
(http://csrc.nist.gov/groups/ST/toolkit/BCM/documents/comments/800-38_Series-Drafts/GCM/Joux_comments.pdf).

8. "Authentication weaknesses in GCM" by Niels Ferguson
(http://csrc.nist.gov/groups/ST/toolkit/BCM/documents/comments/CWC-GCM/Ferguson2.pdf).

And here are some useful links:

1. http://safecurves.cr.yp.to/ - A great starting point to learn about
   attacks on elliptic curves and the design choices in modern curves
   that prevent them.

2. https://www.hyperelliptic.org/EFD/ - Cheat codes for ECC
   implementation.

3. http://cacr.uwaterloo.ca/hac/ - This book is free online and it has
   everything you need to do the math in this problem set.
