package main

import (
	"flag"
)

var symbols = []string{"K", "X"}

func main() {
	var oneway, standard, deferred, noss, ss bool
	flag.BoolVar(&oneway, "oneway", false, "Print the one way patterns")
	flag.BoolVar(&standard, "standard", false, "Print the standard two way patterns")
	flag.BoolVar(&deferred, "deferred", false, "Print the deferred two way patterns")
	flag.BoolVar(&noss, "noss", false, "Print the noss patterns")
	flag.BoolVar(&ss, "ss", false, "Print the ss patterns")
	flag.Parse()
	if !oneway && !standard && !deferred {
		oneway, standard, deferred = true, true, true
	}

	var pr printer
	pr = new(pretty)

	// the one way patterns
	if oneway {
		for _, i := range symbols {
			makePattern(pr, i, "", false, false, false, false)
		}
		if standard || deferred {
			pr.Println()
		}
	}
	if standard {
		// the standard patterns
		for _, i := range symbols {
			for _, r := range symbols {
				makePattern(pr, i, r, false, false, false, false)
				// also make the equivalent I patterns from X
				if i == "X" {
					makePattern(pr, "I", r, false, false, false, false)
				}
			}
		}
		if deferred {
			pr.Println()
		}
	}
	if deferred {
		// the deferred patterns
		for _, i := range symbols {
			for _, r := range symbols {
				if i == "N" && r == "N" {
					continue
				}
				if i != "N" {
					makePattern(pr, i, r, true, false, false, false)
				}
				if r != "N" {
					makePattern(pr, i, r, false, true, false, false)
				}
				if i != "N" && r != "N" {
					makePattern(pr, i, r, true, true, false, false)
				}
				if i == "X" {
					makePattern(pr, "I", r, true, false, false, false)
					if r != "N" {
						makePattern(pr, "I", r, false, true, false, false)
						makePattern(pr, "I", r, true, true, false, false)
					}
				}
			}
		}
	}
	if noss && oneway {
		makePattern(pr, "K", "", false, false, true, false)
		makePattern(pr, "X", "", false, false, true, false)
	}
	if noss && standard {
		makePattern(pr, "K", "K", false, false, true, false)
		makePattern(pr, "I", "K", false, false, true, false)
	}
	if ss && standard {
		for _, i := range symbols {
			for _, r := range symbols {
				if i == "N" || r == "N" {
					continue
				}
				makePattern(pr, i, r, false, false, false, true)
				if i == "X" {
					makePattern(pr, "I", r, false, false, false, true)
				}
			}
		}
	}
	if ss && deferred {
		for _, i := range symbols {
			for _, r := range symbols {
				if i == "N" || r == "N" {
					continue
				}
				makePattern(pr, i, r, true, false, false, true)
				makePattern(pr, i, r, false, true, false, true)
				makePattern(pr, i, r, true, true, false, true)
				if i == "X" {
					makePattern(pr, "I", r, true, false, false, true)
					makePattern(pr, "I", r, false, true, false, true)
					makePattern(pr, "I", r, true, true, false, true)
				}
			}
		}
	}
}

// makePattern outputs a single pattern based on the two tokens and two booleans for deferral
func makePattern(pr printer, it, rt string, id, rd bool, noss bool, ssmod bool) {
	// have these DH taken place?
	var ee, es, se, ss bool
	// have initiator and responder sent e, s?
	var ie, is, re, rs bool

	pr.PrintHeader(it, rt, id, rd, noss, ssmod)
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
	first := true
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
				// do ss if we cannot send se on first line and not deferred
				case !noss && !ssmod && is && rs && !ss && es && !se && !id && first:
					pr.PrintI("ss")
					ss = true
				// for ss modifer, do ss last after se, es
				case ssmod && !ss && is && rs && !ss && es && se:
					pr.PrintI("ss")
					ss = true
				// send s if I or one way X as soon as possible
				case (it == "I" || (it == "X" && rt == "")) && !is:
					pr.PrintI("s")
					is = true
				// send s if X, but not on first line
				case it == "X" && !is && !first:
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
				// send ss if have ss modifier
				case ssmod && is && rs && !ss && se && es:
					pr.PrintI("ss")
					ss = true
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
		first = false
		if !didLine {
			break
		}
		didLine = false
	}
}
