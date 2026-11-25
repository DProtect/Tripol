package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestApexRecurrence(t *testing.T) {
	cases := []string{
		`func f(n){if n<=1{return n};return f(n-1)+f(n-2)};show f(20),f(-9),f(0),f(1)`,
		`func g(x){if x<3{return 2};return g(x-2)-g(x-3)};show g(12),g(-1)`,
		`func h(n){if n<=1{return 2};return h(n-1)*h(n-2)};show h(5)`,
		`func h(n){if n<=1{return 9223372036854775807};return h(n-1)+h(n-2)};show h(2)`,
		`func h(n){if n<=1{return 9223372036854775807};return h(n-1)+h(n-2)};try{show h(2)}catch code{show code}`,
		`func f(n){if n<=1{return n};return f(n-1)+f(n-2)};func g(n){return n+f(n)};show g(12),f(f(5))`,
		`func f(n=8){if n<=1{return n};return f(n-1)+f(n-2)};show f(),f(9)`,
		`func f(n){if n<=1{return n};show n;return f(n-1)+f(n-2)};show f(5)`,
		`func f(n){if n<=1{return n};return f(n-1)+f(n-2)};show f(-9223372036854775807-1)`,
	}
	for i, src := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			w, b := hyperRun(t, src, "debug")
			for _, p := range []string{"nano", "balanced", "turbo"} {
				o, e := hyperRun(t, src, p)
				if o != w || e != b {
					t.Fatalf("%s: %q %v != %q %v", p, o, e, w, b)
				}
			}
		})
	}
}
func TestApexLoopProofs(t *testing.T) {
	cases := []string{
		`n=123;s=0;for i in range(n){x=i-n/2;s+=x/8;s+=x%8};show s`,
		`n=-91;s=0;for i in range(30){s+=n/4;s+=n%4};show s`,
		`n=123;s=0;for i in range(30){s+=n/4;n-=3;s+=n%4};show s,n`,
		`n=123;s=0;for i in range(0){s+=n/4;s+=n%4};show s`,
		`n=123;s=0;for i in range(31){if i%2==0{s+=n/4}else{s+=n%4};if i==20{break}};show s`,
		`s=0;for i in range(-15,16){if i%8==0{s+=i}else{s-=i};if i%4!=0{s+=1}};show s`,
		`s=0;for i in range(0,101){s+=i/4;s+=i%4};show s`,
		`s=0;for i in range(-100,-3){s+=i/4;s+=i%4};show s`,
		`s=9223372036854775807;for i in range(1,2){s+=i/2;s+=i%2};show s`,
		`s=9223372036854775807;for i in range(4,5){s+=i/2;s+=i%2};show s`,
		`s=100;for i in range(20){s+=s/8;s+=s%8};show s`,
		`s=0;for i in range(3){x=-9223372036854775807-1;s+=x/8;s+=x%8};show s`,
		`s=0;for i in range(3){x=9223372036854775807;s+=x/8;s+=x%8};show s`,
		`s=0;n=100;try{for i in range(3){s+=n/4;s+=n%4;show 1/(1-i)}}catch code{show code};show s`,
		`s=0;for i in range(10){s+=i/4;s+=i%4};for i in range(-10,0){s+=i/4;s+=i%4};show s`,
	}
	for i, src := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			w, b := hyperRun(t, src, "debug")
			for _, p := range []string{"balanced", "turbo"} {
				o, e := hyperRun(t, src, p)
				if o != w || e != b {
					t.Fatalf("%s: %q %v != %q %v", p, o, e, w, b)
				}
			}
		})
	}
}
func TestApexDivisionPairsRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(19))
	var s, w strings.Builder
	for i := 0; i < 80; i++ {
		n := rng.Int63n(1<<40) - (1 << 39)
		d := int64(1) << uint(1+rng.Intn(30))
		fmt.Fprintf(&s, "s%d=0;x%d=%d;for i in range(3){s%d+=x%d/%d;s%d+=x%d%%%d};show s%d\n", i, i, n, i, i, d, i, i, d, i)
		fmt.Fprintln(&w, 3*(n/d+n%d))
	}
	for _, p := range []string{"debug", "turbo"} {
		o, e := hyperRun(t, s.String(), p)
		if e || o != w.String() {
			t.Fatalf("%s mismatch", p)
		}
	}
}
