package main

import (
	"fmt"
)

var Data TextData

func main() {
	GeneratePositions()
	
	Data = GetTextData()
	fmt.Println(Data)
	isrt := "qwertyuiopasdfghjkl;zxcvbnm,./"
	fmt.Println(CalcFingerSpeed(isrt))
}
