package otpetushi

import "testing"

const (
	Ko  = "ко"
	Koo = Ko + "о"
)

func TestWordToKoKo_3(t *testing.T) {
	const src3 string = "ori"

	src3_2 := wordToKoKo(src3)
	if src3_2 != Ko+Ko {
		t.Error("3 symb -> koko")
	}
}

func TestWordToKoKo_0(t *testing.T) {
	const src3 string = ""

	src3_2 := wordToKoKo(src3)
	if len(src3_2) != 0 {
		t.Error("no src -> no koko")
	}
}

func TestWordToKoKo_1(t *testing.T) {
	const src3 string = "h"

	src3_2 := wordToKoKo(src3)
	if src3_2 != Ko {
		t.Error("1 symb -> ko")
	}
}

func TestWordToKoKo_even(t *testing.T) {
	const src3 string = "source"

	src3_2 := wordToKoKo(src3)
	if src3_2 != Ko+Ko+Ko {
		t.Error("even src -> even ko's")
	}
}

func TestWordToKoKo_odd(t *testing.T) {
	const src3 string = "sourc"

	src3_2 := wordToKoKo(src3)
	if src3_2 != Ko+Koo {
		t.Error("odd src -> even ko's + tail 'o'")
	}
}

func TestWordToKoKo_pi(t *testing.T) {
	const src3 string = "3.14"

	src3_2 := wordToKoKo(src3)
	if src3_2 != "3.14" {
		t.Error("numbers src -> numbers out")
	}
}

// dont support numbers inside word. Just replace with koko as usual symbols
func TestWordToKoKo_number_inside(t *testing.T) {
	const src3 string = "he11ooo"

	src3_2 := wordToKoKo(src3)
	if src3_2 != Ko+Ko+Koo {
		t.Error("word_with_num_inside src -> ordinary koko's output")
	}
}

func TestMessageUnicode(t *testing.T) {
	const src3 string = "ляля"

	src3_2 := wordToKoKo(src3)
	if src3_2 != Ko+Ko {
		t.Error("any_unicode N chars src -> ~N koko's")
	}
}

/*func TestMessage0(t *testing.T) {
	const src3 string = "!отпетуши"

	src3_2 := messageToKoKo(src3)
	if src3_2 != "!"+Ko+Ko+Ko+Ko {
		t.Error("any_unicode N chars src -> ~N koko's")
	}
}*/
