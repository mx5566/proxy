package base

import (
	"testing"

	"fmt"
)

func TestNew(t *testing.T) {
	set := New()

	set.Add(1)
	set.Add(1)
	set.Add(2)

	fmt.Println(set.m)

	fmt.Println(set.Test(1))
}
