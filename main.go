package main

import (
	"fmt"
)

type done struct {
	// what have we done
	e bool
	s bool
}

type half struct {
	// config
	symbol string
	def    bool
	// state
	d done
}

type pattern struct {
	i half
	r half
	// DH done
	ee bool
	es bool
	se bool
	ss bool
}

var symbols = []string{"N", "X", "K"}

func main() {
	// first the standard patterns
	for _, i := range symbols {
		for _, r := range symbols {
			makePattern(i, r, false, false)
			// also make the equivalent I patterns from X
			if i == "X" {
				makePattern("I", r, false, false)
			}
		}
	}
	fmt.Printf("\n")
	// then standard plus deferred
	for _, i := range symbols {
		for _, r := range symbols {
			if i == "N" && r == "N" {
				continue
			}
			makePattern(i, r, false, false)
			if i != "N" {
				makePattern(i, r, true, false)
			}
			if r != "N" {
				makePattern(i, r, false, true)
			}
			if i != "N" && r != "N" {
				makePattern(i, r, true, true)
			}
			if i == "X" {
				makePattern("I", r, false, false)
				makePattern("I", r, true, false)
				if r != "N" {
					makePattern("I", r, true, true)
				}
			}
		}
	}
}

func makeInitHalf(s string) half {
	return half{
		symbol: s,
		def:    false, // add later
	}
}

func makeInit(i, r string) pattern {
	p := pattern{
		i: makeInitHalf(i),
		r: makeInitHalf(r),
	}
	// pre-message
	if i == "K" {
		p.i.d.s = true
		fmt.Printf("  -> s\n")
	}
	if r == "K" {
		p.r.d.s = true
		fmt.Printf("  <- s\n")
	}
	if i == "K" || r == "K" {
		fmt.Printf("  ...\n")
	}
	return p
}

func pr(first, initWrite bool, s string) bool {
	if first {
		switch initWrite {
		case true:
			fmt.Printf("  -> ")
		case false:
			fmt.Printf("  <- ")
		}
	} else {
		fmt.Printf(", ")
	}
	fmt.Printf(s)
	return false
}

func makePattern(i, r string, id, rd bool) {
	var ids, rds string
	if id {
		ids = "1"
	}
	if rd {
		rds = "1"
	}
	fmt.Printf("%s%s%s%s:\n", i, ids, r, rds)
	p := makeInit(i, r)
	// direction
	initWrite := true
	line := 0
	var didLine bool
	var clearID, clearRD bool
	for {
		var didSomething bool
		var first = true
		for {
			switch initWrite {
			case true: // initiator writes
				switch {
				// send e if not sent
				case !p.i.d.e:
					first = pr(first, initWrite, "e")
					p.i.d.e = true
					didSomething = true
				// do ee as soon as possible
				case p.i.d.e && p.r.d.e && !p.ee:
					first = pr(first, initWrite, "ee")
					p.ee = true
					didSomething = true
				// do se as soon as possible if not deferred
				case p.i.d.s && p.r.d.e && !p.se && !id:
					first = pr(first, initWrite, "se")
					p.se = true
					didSomething = true
				// do es as soon as possible if not deferred
				case p.i.d.e && p.r.d.s && !p.es && !rd:
					first = pr(first, initWrite, "es")
					p.es = true
					didSomething = true
				// do ss as soon as possible
				case p.i.d.s && p.r.d.s && !p.ss:
					first = pr(first, initWrite, "ss")
					p.ss = true
					didSomething = true
				// send s if I as soon as possible
				case i == "I" && !p.i.d.s:
					first = pr(first, initWrite, "s")
					p.i.d.s = true
					didSomething = true
				case i == "X" && !p.i.d.s && line == 1:
					first = pr(first, initWrite, "s")
					p.i.d.s = true
					didSomething = true
				}
				// handle deferral
				if p.i.d.s && p.r.d.e && !p.se && id {
					clearID = true
				}
				if p.i.d.e && p.r.d.s && !p.es && rd {
					clearRD = true
				}
			case false: // recipient writes
				switch {
				// send e if not sent
				case !p.r.d.e:
					first = pr(first, initWrite, "e")
					p.r.d.e = true
					didSomething = true
				// do ee as soon as possible
				case p.i.d.e && p.r.d.e && !p.ee:
					first = pr(first, initWrite, "ee")
					p.ee = true
					didSomething = true
				// do se as soon as possible if not deferred
				case p.r.d.e && p.i.d.s && !p.se && !id:
					first = pr(first, initWrite, "se")
					p.se = true
					didSomething = true
				// do es as soon as possible if not deferred
				case p.r.d.s && p.i.d.e && !p.es && !rd:
					first = pr(first, initWrite, "es")
					p.es = true
					didSomething = true
				// do ss as soon as possible
				case p.i.d.s && p.r.d.s && !p.ss:
					first = pr(first, initWrite, "ss")
					p.ss = true
					didSomething = true
				case r == "X" && !p.r.d.s:
					first = pr(first, initWrite, "s")
					p.r.d.s = true
					didSomething = true
				}
				// handle deferral
				if p.r.d.s && p.i.d.e && !p.es && rd {
					clearRD = true
				}
				if p.r.d.e && p.i.d.s && !p.se && id {
					clearID = true
				}
			}
			if !didSomething {
				fmt.Printf("\n")
				break
			}
			didLine = true
			didSomething = false
		}
		initWrite = !initWrite
		// handle clearing deferral
		if clearID {
			id = false
		}
		if clearRD {
			rd = false
		}
		if initWrite {
			line++
		}
		if !didLine {
			break
		}
		didLine = false
	}
}
