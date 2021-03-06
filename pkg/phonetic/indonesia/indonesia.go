// Package indonesia implements indonesia-phonetic encoding.
// Implementation of encoding based on: https://github.com/lafzi/lafzi-web/blob/master/lib/fonetik_id.php.
package indonesia

import (
	"bytes"
	"regexp"

	"github.com/dlclark/regexp2"
)

// Encoder implements indonesia-phonetic encoding.
type Encoder struct {
	vowel bool
}

// SetVowel ...
func (enc *Encoder) SetVowel(vowel bool) {
	enc.vowel = vowel
}

// Encode returns encoded of src using encoding enc.
func (enc *Encoder) Encode(src []byte) []byte {
	b := praprocess(src)
	b = vowelSub(b)
	b = joinConsonant(b)
	b = joinVowel(b)
	b = diphthongSub(b)
	b = markHamzah(b)
	b = ikhfaSub(b)
	b = iqlabSub(b)
	b = idghamSub(b)
	b = encode1consonant(b)
	b = encode2consonant(b)
	b = removeSpace(b)
	if !enc.vowel {
		b = removeVowel(b)
	}

	return b
}

func praprocess(b []byte) []byte {
	// uppercase
	b = bytes.ToUpper(b)
	// change hyphen (-) into space
	b = bytes.Replace(b, []byte("-"), []byte(" "), -1)
	// single space
	b = regexp.MustCompile("\\s+").
		ReplaceAll(bytes.TrimSpace(b), []byte(" "))
	// remove all character except alphabet, grave (`), apostrophe ('), and space
	b = regexp.MustCompile("[^A-Z`'\\s]").
		ReplaceAll(b, []byte(""))

	return b
}

func vowelSub(b []byte) []byte {
	return bytes.Map(func(r rune) rune {
		switch r {
		case 'O':
			return 'A'
		case 'E':
			return 'I'
		default:
			return r
		}
	}, b)
}

func joinConsonant(b []byte) []byte {
	str := string(b)
	// single consonant
	str, _ = regexp2.MustCompile("(?<single>B|C|D|F|G|H|J|K|L|M|N|P|Q|R|S|T|V|W|X|Y|Z)\\s?\\1+", 0).
		Replace(str, "${single}", -1, -1)
	// double consonant
	str, _ = regexp2.MustCompile("(?<double>KH|CH|SH|TS|SY|DH|TH|ZH|DZ|GH)\\s?\\1+", 0).
		Replace(str, "${double}", -1, -1)

	return []byte(str)
}

func joinVowel(b []byte) []byte {
	str := string(b)
	// single vocal
	str, _ = regexp2.MustCompile("(?<single>A|I|U|E|O)\\1+", 0).
		Replace(str, "${single}", -1, -1)

	return []byte(str)
}

func diphthongSub(b []byte) []byte {
	b = regexp.MustCompile("AI").
		ReplaceAll(b, []byte("AY"))
	b = regexp.MustCompile("AU").
		ReplaceAll(b, []byte("AW"))

	return b
}

func markHamzah(b []byte) []byte {
	// beginning of the string
	b = regexp.MustCompile("^(?P<hamzah>A|I|U)").
		ReplaceAll(b, []byte("X${hamzah}"))
	// after space
	b = regexp.MustCompile("\\s(?P<hamzah>A|I|U)").
		ReplaceAll(b, []byte(" X${hamzah}"))
	// IA, IU => IXA, IXU
	b = regexp.MustCompile("I(?P<hamzah>A|U)").
		ReplaceAll(b, []byte("IX${hamzah}"))
	// UA, UI => UXA, UXI
	b = regexp.MustCompile("U(?P<hamzah>A|I)").
		ReplaceAll(b, []byte("UX${hamzah}"))

	return b
}

func ikhfaSub(b []byte) []byte {
	// [vowel][NG][ikhfa] => [vowel][N][ikhfa]
	return regexp.MustCompile("(?P<vowel>A|I|U)NG\\s?(?P<ikhfa>D|F|J|K|P|Q|S|T|V|Z)").
		ReplaceAll(b, []byte("${vowel}N${ikhfa}"))
}

func iqlabSub(b []byte) []byte {
	// NB => MB
	return regexp.MustCompile("N\\s?B").
		ReplaceAll(b, []byte("MB"))
}

func idghamSub(b []byte) []byte {
	// exception
	b = bytes.Replace(b, []byte("DUNYA"), []byte("DUN_YA"), -1)
	b = bytes.Replace(b, []byte("BUNYAN"), []byte("BUN_YAN"), -1)
	b = bytes.Replace(b, []byte("QINWAN"), []byte("KIN_WAN"), -1)
	b = bytes.Replace(b, []byte("KINWAN"), []byte("KIN_WAN"), -1)
	b = bytes.Replace(b, []byte("SINWAN"), []byte("SIN_WAN"), -1)
	b = bytes.Replace(b, []byte("SHINWAN"), []byte("SIN_WAN"), -1)

	// N,M,L,R,Y,W
	b = regexp.MustCompile("N\\s?(?P<idgham>N|M|L|R|Y|W)").
		ReplaceAll(b, []byte("${idgham}"))

	// reverse the exception
	b = bytes.Replace(b, []byte("DUN_YA"), []byte("DUNYA"), -1)
	b = bytes.Replace(b, []byte("BUN_YAN"), []byte("BUNYAN"), -1)
	b = bytes.Replace(b, []byte("KIN_WAN"), []byte("KINWAN"), -1)
	b = bytes.Replace(b, []byte("SIN_WAN"), []byte("SINWAN"), -1)

	return b
}

func encode2consonant(b []byte) []byte {
	b = regexp.MustCompile("KH|CH").
		ReplaceAll(b, []byte("H"))
	b = regexp.MustCompile("SH|TS|SY").
		ReplaceAll(b, []byte("S"))
	b = regexp.MustCompile("DH").
		ReplaceAll(b, []byte("D"))
	b = regexp.MustCompile("ZH|DZ").
		ReplaceAll(b, []byte("Z"))
	b = regexp.MustCompile("TH").
		ReplaceAll(b, []byte("T"))
	b = regexp.MustCompile("NG(?P<sub>A|I|U)").
		ReplaceAll(b, []byte("X${sub}"))
	b = regexp.MustCompile("GH").
		ReplaceAll(b, []byte("G"))

	return b
}

func encode1consonant(b []byte) []byte {
	b = regexp.MustCompile("'|`").
		ReplaceAll(b, []byte("X"))
	b = regexp.MustCompile("Q|K").
		ReplaceAll(b, []byte("K"))
	b = regexp.MustCompile("F|V|P").
		ReplaceAll(b, []byte("F"))
	b = regexp.MustCompile("J|Z").
		ReplaceAll(b, []byte("Z"))

	return b
}

func removeSpace(b []byte) []byte {
	return regexp.MustCompile("\\s").
		ReplaceAll(b, []byte(""))
}

func removeVowel(b []byte) []byte {
	return regexp.MustCompile("A|I|U").
		ReplaceAll(b, []byte(""))
}
