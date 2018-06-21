package main

import (
	"fmt"
)

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
	fmt.Println()
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
			fmt.Print("  -> ")
		case false:
			fmt.Print("  <- ")
		}
	} else {
		fmt.Print(", ")
	}
	fmt.Print(s)
	return false
}

func prDefer(d bool) string {
	if d {
		return "1"
	}
	return ""
}

func makePattern(it, rt string, id, rd bool) {
	// have these DH taken place?
	var ee, es, se, ss bool
	// have initiator and responder sent e, s?
	var ie, is, re, rs bool

	fmt.Println(it + prDefer(id) + rt + prDefer(rd) + ":")
	// pre-message handling
	if it == "K" {
		is = true
		fmt.Println("  -> s")
	}
	if rt == "K" {
		rs = true
		fmt.Println("  <- s")
	}
	if it == "K" || rt == "K" {
		fmt.Println("  ...")
	}
	// direction, start with initiator writes
	initWrite := true
	line := 0
	var didLine bool
	var clearID, clearRD bool
	for {
		var didNothing bool
		var first = true
		for {
			switch initWrite {
			case true: // initiator writes
				switch {
				// send e if not sent
				case !ie:
					first = pr(first, initWrite, "e")
					ie = true
				// do ee as soon as possible
				case ie && re && !ee:
					first = pr(first, initWrite, "ee")
					ee = true
				// do se as soon as possible if not deferred
				case is && re && !se && !id:
					first = pr(first, initWrite, "se")
					se = true
				// do es as soon as possible if not deferred
				case ie && rs && !es && !rd:
					first = pr(first, initWrite, "es")
					es = true
				// do ss as soon as possible, if not done ee or se but have done se
				case is && rs && !ss && !se && !ee && es:
					first = pr(first, initWrite, "ss")
					ss = true
				// send s if I as soon as possible
				case it == "I" && !is:
					first = pr(first, initWrite, "s")
					is = true
				// send s if X, but not on first line
				case it == "X" && !is && line == 1:
					first = pr(first, initWrite, "s")
					is = true
				default:
					didNothing = true
				}
				// handle deferral
				if is && re && !se && id {
					clearID = true
				}
				if ie && rs && !es && rd {
					clearRD = true
				}
			case false: // recipient writes
				switch {
				// send e if not sent
				case !re:
					first = pr(first, initWrite, "e")
					re = true
				// do ee as soon as possible
				case ie && re && !ee:
					first = pr(first, initWrite, "ee")
					ee = true
				// do se as soon as possible if not deferred
				case re && is && !se && !id:
					first = pr(first, initWrite, "se")
					se = true
				// do es as soon as possible if not deferred
				case rs && ie && !es && !rd:
					first = pr(first, initWrite, "es")
					es = true
				// do ss as soon as possible if not done ee or es but have done se
				case is && rs && !ss && !es && !ee && se:
					first = pr(first, initWrite, "ss")
					ss = true
				// send s if X as soon as possible
				case rt == "X" && !rs:
					first = pr(first, initWrite, "s")
					rs = true
				default:
					didNothing = true
				}
				// handle deferral
				if rs && ie && !es && rd {
					clearRD = true
				}
				if re && is && !se && id {
					clearID = true
				}
			}
			if didNothing {
				fmt.Println()
				break
			}
			didLine = true
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
