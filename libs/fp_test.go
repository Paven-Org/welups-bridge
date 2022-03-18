package libs

import (
	"fmt"
	"testing"
)

type TestCase[A any, B any] struct {
	Case     string
	Given    A
	Expected B
}

type TestSuite[A any, B any] []TestCase[A, B]

func TestMapReduce(t *testing.T) {
	test1 := TestCase[[]int, []string]{
		"Test 1",
		[]int{1, 2, 3, 4, 5, 6},
		[]string{"1", "2", "3", "4", "5", "6"},
	}
	test2 := TestCase[[]int, string]{
		"Test 2",
		[]int{1, 2, 3, 4, 5, 6},
		"!1!2!3!4!5!6",
	}

	fmt.Println(Map(func(a int) string { return fmt.Sprintf("%d", a) }, test1.Given), test1.Expected)
	fmt.Println(Reduce(func(b string, a int) string { return fmt.Sprintf("%s!%d", b, a) }, "", test2.Given) == test2.Expected)

}
