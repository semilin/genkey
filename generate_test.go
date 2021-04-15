package main

import (
	"testing"
	"fmt"
)

func TestCycleRandKeys(t *testing.T) {
	l := "qwertyuiopasdfghjkl;zxcvbnm,./"
	fmt.Println(l)
	fmt.Println(cycleRandKeys(l,1))
	fmt.Println(cycleRandKeys(l,2))
	fmt.Println(cycleRandKeys(l,3))
}
