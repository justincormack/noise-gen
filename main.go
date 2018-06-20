package main

import (
	"fmt"
)

type half struct {
	// public keys sent
	e bool
	s bool
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
	// then deferred patterns
	for _, i := range symbols {
		for _, r := range symbols {
			if i == "N" && r == "N" {
				continue
			}
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
				makePattern("I", r, true, false)
				if r != "N" {
					makePattern("I", r, false, true)
					makePattern("I", r, true, true)
				}
			}
		}
	}
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
	p := pattern{i: half{}, r: half{}}
	// pre-message handling
	if i == "K" {
		p.i.s = true
		fmt.Printf("  -> s\n")
	}
	if r == "K" {
		p.r.s = true
		fmt.Printf("  <- s\n")
	}
	if i == "K" || r == "K" {
		fmt.Printf("  ...\n")
	}
	// direction, start with initiator writes
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
				case !p.i.e:
					first = pr(first, initWrite, "e")
					p.i.e = true
					didSomething = true
				// do ee as soon as possible
				case p.i.e && p.r.e && !p.ee:
					first = pr(first, initWrite, "ee")
					p.ee = true
					didSomething = true
				// do se as soon as possible if not deferred
				case p.i.s && p.r.e && !p.se && !id:
					first = pr(first, initWrite, "se")
					p.se = true
					didSomething = true
				// do es as soon as possible if not deferred
				case p.i.e && p.r.s && !p.es && !rd:
					first = pr(first, initWrite, "es")
					p.es = true
					didSomething = true
				// do ss as soon as possible
				case p.i.s && p.r.s && !p.ss:
					first = pr(first, initWrite, "ss")
					p.ss = true
					didSomething = true
				// send s if I as soon as possible
				case i == "I" && !p.i.s:
					first = pr(first, initWrite, "s")
					p.i.s = true
					didSomething = true
				// send s if X, but not on first line
				case i == "X" && !p.i.s && line == 1:
					first = pr(first, initWrite, "s")
					p.i.s = true
					didSomething = true
				}
				// handle deferral
				if p.i.s && p.r.e && !p.se && id {
					clearID = true
				}
				if p.i.e && p.r.s && !p.es && rd {
					clearRD = true
				}
			case false: // recipient writes
				switch {
				// send e if not sent
				case !p.r.e:
					first = pr(first, initWrite, "e")
					p.r.e = true
					didSomething = true
				// do ee as soon as possible
				case p.i.e && p.r.e && !p.ee:
					first = pr(first, initWrite, "ee")
					p.ee = true
					didSomething = true
				// do se as soon as possible if not deferred
				case p.r.e && p.i.s && !p.se && !id:
					first = pr(first, initWrite, "se")
					p.se = true
					didSomething = true
				// do es as soon as possible if not deferred
				case p.r.s && p.i.e && !p.es && !rd:
					first = pr(first, initWrite, "es")
					p.es = true
					didSomething = true
				// do ss as soon as possible
				case p.i.s && p.r.s && !p.ss:
					first = pr(first, initWrite, "ss")
					p.ss = true
					didSomething = true
				// send s if X as soon as possible
				case r == "X" && !p.r.s:
					first = pr(first, initWrite, "s")
					p.r.s = true
					didSomething = true
				}
				// handle deferral
				if p.r.s && p.i.e && !p.es && rd {
					clearRD = true
				}
				if p.r.e && p.i.s && !p.se && id {
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
