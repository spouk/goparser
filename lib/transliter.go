package lib

import (
	"strings"
	"fmt"
	d "github.com/fiam/gounidecode/unidecode"
)
//---------------------------------------------------------------------------
//  TRANSLITE: транслитерация
//---------------------------------------------------------------------------
type Transliter struct {
	Validate []int
	Replacer []int
	InValid  []int
}

func NewTransliter() *Transliter {
	//	hex 65-122 A-z (допустимые )
	//	hex 48-57 0-9 ( допустимые )
	// 	hex 20 (зменяемые)
	//	hex 123-126, 91-96, 58-64, 33-47 punctuations ( запретные )
	n := new(Transliter)
	n.Validate = n.convert(65, 122)
	n.InValid = n.convert(123, 126)

	n.InValid = append(n.InValid, n.convert(91, 96)...)
	n.InValid = append(n.InValid, n.convert(58, 64)...)
	n.InValid = append(n.InValid, n.convert(33, 47)...)
	n.InValid = append(n.InValid, n.convert(33, 47)...)
	n.Replacer = n.convert(32, 32)
	//n.Replacer = append(n.Replacer, 20)
	return n
}
func (s *Transliter) convert(start, end int) []int {
	stack := []int{}
	for ; start <= end; start++ {
		stack = append(stack, start)
	}
	return stack
}
func (s *Transliter) correct(str string) string {

	var result []string
	for _, x := range strings.Split(strings.TrimSpace(str), " ") {
		if x != "" {
			result = append(result, x)
		}
	}
	return strings.Join(result, " ")
}
func (s *Transliter) preCorrect(str string) string {
	str = s.correct(str)
	var tmp []string
	for _, sym := range str {
		switch {
		case s.InSlice(s.InValid, int(sym)):
			continue
		case s.InSlice(s.Validate, int(sym)):
			tmp = append(tmp, string(sym))
		case s.InSlice(s.Replacer, int(sym)):
			tmp = append(tmp, " ")
		default:
			tmp = append(tmp, string(sym))
		}
	}
	return s.correct(strings.Join(tmp, ""))
}

func (s *Transliter) TransliterCyr(str string) string {
	str = s.preCorrect(str)
	var result []string
	for _, sym := range d.Unidecode(str) {
		switch {
		case s.InSlice(s.InValid, int(sym)):
			continue
		case s.InSlice(s.Validate, int(sym)):
			result = append(result, string(sym))
		case s.InSlice(s.Replacer, int(sym)):
			result = append(result, "-")
		default:
			result = append(result, string(sym))
		}
	}
	return strings.Join(result, "")
}
func (s *Transliter) InSlice(str []int, target int) bool {
	for x := 0; x < len(str); x++ {
		if str[x] == target {
			return true
		}
	}
	return false
}
func (s *Transliter) ShowAscii() {
	var i int
	for i = 0; i < 255; i++ {
		fmt.Printf("Dec: %3d Sym: %3c Hex: %3x\n", i, i, i)
	}
}

