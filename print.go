package main

import (
	"fmt"
)

type printer interface {
	PrintI(s string)
	PrintR(s string)
	Println()
	EndPremessage()
	PrintHeader(it, rt string, id, rd bool)
}

type pretty struct {
	pos int
}

func (pr *pretty) PrintI(s string) {
	if pr.pos == 0 {
		fmt.Print("  -> ")
	} else {
		fmt.Print(", ")
	}
	fmt.Print(s)
	pr.pos++
}

func (pr *pretty) PrintR(s string) {
	if pr.pos == 0 {
		fmt.Print("  <- ")
	} else {
		fmt.Print(", ")
	}
	fmt.Print(s)
	pr.pos++
}

func (pr *pretty) Println() {
	fmt.Println()
	pr.pos = 0
}

func (pr *pretty) EndPremessage() {
	fmt.Println("  ...")
	pr.pos = 0
}

func (pr *pretty) PrintHeader(it, rt string, id, rd bool) {
	fmt.Println(it + prDefer(id) + rt + prDefer(rd) + ":")
	pr.pos = 0
}

func prDefer(d bool) string {
	if d {
		return "1"
	}
	return ""
}
