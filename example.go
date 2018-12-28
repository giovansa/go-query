package main

import "fmt"

func main(){
	test := NewQuery(SetParam("select", "id = 1"))
	fmt.Print(test.GetQuery("anjay"))
}