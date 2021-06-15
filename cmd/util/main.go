package main

import (
	"fmt"
)

func getA() string {
	fmt.Println(">> getA() called!")
	return "AAA"
}

func getAWrapper() string {
	a := getA()
	fmt.Println(">> a mem address", &a)
	return a
}

type StructWrapper struct {
	Title string
	Desc  string
	Child StructChild
}

type StructChild struct {
	Subtitle string
	Kind     string
}

// This is a temp placeholder stub. Fill out later. 6/11/21
func main() {
	// ev := config.GetEnvVars()
	// fmt.Println(">> Run main util. Env vars:", ev)

	aa := getAWrapper()
	fmt.Println(">> getAWrapper():", aa, &aa)
	fmt.Println(">> getAWrapper():", aa, &aa)
	fmt.Println()

	sw := new(StructWrapper)
	sw.Title = "this is the title"
	sw.Child.Subtitle = "this is a subtitle"
	fmt.Println(">> sw:", sw)
	fmt.Println(">> &sw:", &sw)
	fmt.Println(">> *sw:", *sw)
	fmt.Println()

	sw2 := sw
	sw2.Title = "this is the title2"
	sw2.Child.Subtitle = "this is a subtitle2"
	fmt.Println(">> sw2:", sw2)
	fmt.Println(">> &sw2:", &sw2)
	fmt.Println(">> *sw2:", *sw2)
	fmt.Println()

	fmt.Println(">> sw:", sw)
}
