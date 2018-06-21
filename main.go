package main

import (
	"fmt"
)

var symbols = []string{"N", "K", "X"}

func main() {
	// the one way patterns
	for _, i := range symbols {
		makePattern(i, "", false, false)
	}
	fmt.Println()
	// the standard patterns
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

type printer struct {
	pos int
}

func (pr *printer) PrintI(s string) {
	if pr.pos == 0 {
		fmt.Print("  -> ")
	} else {
		fmt.Print(", ")
	}
	fmt.Print(s)
	pr.pos++
}

func (pr *printer) PrintR(s string) {
	if pr.pos == 0 {
		fmt.Print("  <- ")
	} else {
		fmt.Print(", ")
	}
	fmt.Print(s)
	pr.pos++
}

func (pr *printer) Println() {
	fmt.Println()
	pr.pos = 0
}

func (pr *printer) EndPremessage() {
	fmt.Println("  ...")
	pr.pos = 0
}

func prDefer(d bool) string {
	if d {
		return "1"
	}
	return ""
}

func printHeader(it, rt string, id, rd bool) {
	fmt.Println(it + prDefer(id) + rt + prDefer(rd) + ":")
}

func makePattern(it, rt string, id, rd bool) {
	// have these DH taken place?
	var ee, es, se, ss bool
	// have initiator and responder sent e, s?
	var ie, is, re, rs bool

	pr := new(printer)

	printHeader(it, rt, id, rd)
	// pre-message handling
	if it == "K" {
		is = true
		pr.PrintI("s")
		pr.Println()
	}
	if rt == "K" || rt == "" {
		rs = true
		pr.PrintR("s")
		pr.Println()
	}
	if it == "K" || rt == "K" || rt == "" {
		pr.EndPremessage()
	}
	// direction, start with initiator writes
	initWrite := true
	line := 0
	var didLine bool
	var clearID, clearRD bool
	for {
		var didNothing bool
		for {
			switch initWrite {
			case true: // initiator writes
				switch {
				// send e if not sent
				case !ie:
					pr.PrintI("e")
					ie = true
				// do ee as soon as possible
				case ie && re && !ee:
					pr.PrintI("ee")
					ee = true
				// do se as soon as possible if not deferred
				case is && re && !se && !id:
					pr.PrintI("se")
					se = true
				// do es as soon as possible if not deferred
				case ie && rs && !es && !rd:
					pr.PrintI("es")
					es = true
				// do ss if we cannot send se on first line
				case is && rs && !ss && es && !se && line == 0:
					pr.PrintI("ss")
					ss = true
				// send s if I or one way X as soon as possible
				case (it == "I" || (it == "X" && rt == "")) && !is:
					pr.PrintI("s")
					is = true
				// send s if X, but not on first line
				case it == "X" && !is && line != 0:
					pr.PrintI("s")
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
					pr.PrintR("e")
					re = true
				// do ee as soon as possible
				case ie && re && !ee:
					pr.PrintR("ee")
					ee = true
				// do se as soon as possible if not deferred
				case is && re && !se && !id:
					pr.PrintR("se")
					se = true
				// do es as soon as possible if not deferred
				case ie && rs && !es && !rd:
					pr.PrintR("es")
					es = true
				// send s if X as soon as possible
				case rt == "X" && !rs:
					pr.PrintR("s")
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
				pr.Println()
				break
			}
			didLine = true
		}
		// finish if one way pattern
		if rt == "" {
			pr.Println()
			break
		}
		initWrite = !initWrite
		// handle clearing deferral
		if clearID {
			id = false
		}
		if clearRD {
			rd = false
		}
			line++
		if !didLine {
			break
		}
		didLine = false
	}
}
