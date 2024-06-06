package test

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestString(t *testing.T) {
	name := "Lee"
	newName := (*[]byte)(unsafe.Pointer(&name))
	fmt.Println(string(*newName))

}
