package textenc

var iso6937 = []rune("" +
	"\u00a0¡¢£ ¥ §¤‘“«←↑→↓" +
	"°±²³×µ¶·÷’”»¼½¾¿" +
	"                " +
	"―¹®©™♪¬¦    ⅛⅜⅝⅞" +
	"ΩÆĐªĦ ĲĿŁØŒºÞŦŊŉ" +
	"ĸæđðħıĳŀłøœßþŧŋ\u00AD")

var iso8859_1 = []rune("" +
	"\u00a0¡¢£¤¥¦§¨©ª«¬\u00AD®¯" +
	"°±²³´µ¶·¸¹º»¼½¾¿" +
	"ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ" +
	"ÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞß" +
	"àáâãäåæçèéêëìíîï" +
	"ðñòóôõö÷øùúûüýþÿ")

var iso8859_2 = []rune("" +
	"\u00a0Ą˘Ł¤ĽŚ§¨ŠŞŤŹ\u00ADŽŻ" +
	"°ą˛ł´ľśˇ¸šşťź˝žż" +
	"ŔÁÂĂÄĹĆÇČÉĘËĚÍÎĎ" +
	"ĐŃŇÓÔŐÖ×ŘŮÚŰÜÝŢß" +
	"ŕáâăäĺćçčéęëěíîď" +
	"đńňóôőö÷řůúűüýţ˙")

var iso8859_3 = []rune("" +
	"\u00a0Ħ˘£¤ Ĥ§¨İŞĞĴ\u00AD Ż" +
	"°ħ²³´µĥ·¸ışğĵ½ ż" +
	"ÀÁÂ ÄĊĈÇÈÉÊËÌÍÎÏ" +
	" ÑÒÓÔĠÖ×ĜÙÚÛÜŬŜß" +
	"àáâ äċĉçèéêëìíîï" +
	" ñòóôġö÷ĝùúûüŭŝ˙")

var iso8859_4 = []rune("" +
	"\u00a0ĄĸŖ¤ĨĻ§¨ŠĒĢŦ\u00ADŽ¯" +
	"°ą˛ŗ´ĩļˇ¸šēģŧŊžŋ" +
	"ĀÁÂÃÄÅÆĮČÉĘËĖÍÎĪ" +
	"ĐŅŌĶÔÕÖ×ØŲÚÛÜŨŪß" +
	"āáâãäåæįčéęëėíîī" +
	"đņōķôõö÷øųúûüũū˙")

var iso8859_5 = []rune("" +
	"\u00a0ЁЂЃЄЅІЇЈЉЊЋЌ\u00ADЎЏ" +
	"АБВГДЕЖЗИЙКЛМНОП" +
	"РСТУФХЦЧШЩЪЫЬЭЮЯ" +
	"абвгдежзийклмноп" +
	"рстуфхцчшщъыьэюя" +
	"№ёђѓєѕіїјљњћќ§ўџ")

var iso8859_6 = []rune("" +
	"\u00a0   ¤       ،\u00AD  " +
	"           ؛   ؟" +
	" ءآأؤإئابةتثجحخد" +
	"ذرزسشصضطظعغ     " +
	"ـفقكلمنهوىي     " +
	"                ")

var iso8859_7 = []rune("" +
	"\u00a0ʽʼ£  ¦§¨© «¬\u00AD ―" +
	"°±²³΄΅Ά·ΈΉΊ»Ό½ΎΏ" +
	"ΐΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟ" +
	"ΠΡ ΣΤΥΦΧΨΩΪΫάέήί" +
	"ΰαβγδεζηθικλμνξο" +
	"πρςστυφχψωϊϋόύώ ")

var iso8859_8 = []rune("" +
	"\u00a0 ¢£¤¥¦§¨©×«¬\u00AD®‾" +
	"°±²³´µ¶·¸¹÷»¼½¾ " +
	"                " +
	"               ‗" +
	"אבגדהוזחטיךכלםמן" +
	"נסעףפץצקרשת     ")

var iso8859_9 = []rune("" +
	"\u00a0¡¢£¤¥¦§¨©ª«¬\u00AD®¯" +
	"°±²³´µ¶·¸¹º»¼½¾¿" +
	"ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ" +
	"ĞÑÒÓÔÕÖ×ØÙÚÛÜİŞß" +
	"àáâãäåæçèéêëìíîï" +
	"ğñòóôõö÷øùúûüışÿ")

var iso6937twoBytes = []map[byte]rune{
	{ // Grave 0xC1
		'A': 'À', 'E': 'È', 'I': 'Ì', 'O': 'Ò', 'U': 'Ù',
		'a': 'à', 'e': 'è', 'i': 'ì', 'o': 'ò', 'u': 'ù',
	}, { // Acute 0xC2
		'A': 'Á', 'C': 'Ć', 'E': 'É', 'I': 'Í', 'L': 'Ĺ', 'N': 'Ń',
		'O': 'Ó', 'R': 'Ŕ', 'S': 'Ś', 'U': 'Ú', 'Y': 'Ý', 'Z': 'Ź',
		'a': 'á', 'c': 'ć', 'e': 'é', 'g': 'ģ', 'i': 'í', 'l': 'ĺ',
		'n': 'ń', 'o': 'ó', 'r': 'ŕ', 's': 'ś', 'u': 'ú', 'y': 'ý', 'z': 'ź',
	}, { // Circumflex 0xC3
		'A': 'Â', 'C': 'Ĉ', 'E': 'Ê', 'G': 'Ĝ', 'H': 'Ĥ', 'I': 'Î',
		'J': 'Ĵ', 'O': 'Ô', 'S': 'Ŝ', 'U': 'Û', 'W': 'Ŵ', 'Y': 'Ŷ',
		'a': 'â', 'c': 'ĉ', 'e': 'ê', 'g': 'ĝ', 'h': 'ĥ', 'i': 'î',
		'j': 'ĵ', 'o': 'ô', 's': 'ŝ', 'u': 'û', 'w': 'ŵ', 'y': 'ŷ',
	}, { // Tilde 0xC4
		'A': 'Ã', 'I': 'Ĩ', 'N': 'Ñ', 'O': 'Õ', 'U': 'Ũ',
		'a': 'ã', 'i': 'ĩ', 'n': 'ñ', 'o': 'õ', 'u': 'ũ',
	}, { // Macron 0xC5
		'A': 'Ā', 'E': 'Ē', 'I': 'Ī', 'O': 'Ō', 'U': 'Ū',
		'a': 'ā', 'e': 'ē', 'i': 'ī', 'o': 'ō', 'u': 'ū',
	}, { //Breve 0xC6
		'A': 'Ă', 'G': 'Ğ', 'U': 'Ŭ', 'a': 'ă', 'g': 'ğ', 'u': 'ŭ',
	}, { // Dot 0xC7
		'C': 'Ċ', 'E': 'Ė', 'G': 'Ġ', 'I': 'İ', 'Z': 'Ż',
		'c': 'ċ', 'e': 'ė', 'g': 'ġ', 'z': 'ż',
	}, { // Umlaut 0xC8
		'A': 'Ä', 'E': 'Ë', 'I': 'Ï', 'O': 'Ö', 'U': 'Ü', 'Y': 'Ÿ',
		'a': 'ä', 'e': 'ë', 'i': 'ï', 'o': 'ö', 'u': 'ü', 'y': 'ÿ',
	}, { // unused 0xC9
	}, { // Ring 0xCA
		'A': 'Å', 'U': 'Ů', 'a': 'å', 'u': 'ů',
	}, { // Cedilla 0xCB
		'C': 'Ç', 'G': 'Ģ', 'K': 'Ķ', 'L': 'Ļ', 'N': 'Ņ', 'R': 'Ŗ', 'S': 'Ş',
		'T': 'Ţ',
		'c': 'ç', 'k': 'ķ', 'l': 'ļ', 'n': 'ņ', 'r': 'ŗ', 's': 'ş', 't': 'ţ',
	}, { //  unused 0xCC
	}, { //  DoubleAcute 0xCD
		'O': 'Ő', 'U': 'Ű', 'o': 'ő', 'u': 'ű',
	}, { // Ogonek 0xCE
		'A': 'Ą', 'E': 'Ę', 'I': 'Į', 'U': 'Ų',
		'a': 'ą', 'e': 'ę', 'i': 'į', 'u': 'ų',
	}, { // Caron 0xCF
		'C': 'Č', 'D': 'Ď', 'E': 'Ě', 'L': 'Ľ', 'N': 'Ň', 'R': 'Ř', 'S': 'Š',
		'T': 'Ť', 'Z': 'Ž',
		'c': 'č', 'd': 'ď', 'e': 'ě', 'l': 'ľ', 'n': 'ň', 'r': 'ř', 's': 'š',
		't': 'ť', 'z': 'ž',
	},
}

func DecodeISO6937(in []byte) string {
	out := make([]rune, 0, len(in))
	for i := 0; i < len(in); i++ {
		c := in[i]
		if c < 0xA0 {
			out = append(out, rune(c))
			continue
		}
		if c >= 0xC1 && c <= 0xCF && c != 0xC9 && c != 0xCC && i+1 < len(in) {
			r, ok := iso6937twoBytes[c-0xC1][in[i+1]]
			if ok {
				out = append(out, r)
				i++
				continue
			}
		}
		out = append(out, iso6937[c-0xA0])
	}
	return string(out)
}

func decodeISO8859(table []rune, in []byte) string {
	out := make([]rune, len(in))
	for i, c := range in {
		if c < 0xa0 {
			out[i] = rune(c)
			continue
		}
		out[i] = table[c-0xa0]
	}
	return string(out)
}
func DecodeISO8859_1(in []byte) string {
	return decodeISO8859(iso8859_1, in)
}

func DecodeISO8859_2(in []byte) string {
	return decodeISO8859(iso8859_2, in)
}

func DecodeISO8859_3(in []byte) string {
	return decodeISO8859(iso8859_3, in)
}

func DecodeISO8859_4(in []byte) string {
	return decodeISO8859(iso8859_4, in)
}

func DecodeISO8859_5(in []byte) string {
	return decodeISO8859(iso8859_5, in)
}

func DecodeISO8859_6(in []byte) string {
	// BUG: two bytes characters not handled
	return decodeISO8859(iso8859_6, in)
}

func DecodeISO8859_7(in []byte) string {
	return decodeISO8859(iso8859_7, in)
}

func DecodeISO8859_8(in []byte) string {
	return decodeISO8859(iso8859_8, in)
}

func DecodeISO8859_9(in []byte) string {
	return decodeISO8859(iso8859_9, in)
}

var iso8859tab = []func([]byte) string{
	DecodeISO8859_1,
	DecodeISO8859_2,
	DecodeISO8859_3,
	DecodeISO8859_4,
	DecodeISO8859_5,
	DecodeISO8859_6,
	DecodeISO8859_7,
	DecodeISO8859_8,
	DecodeISO8859_9,
}

// DecodeISO8859
// TODO: It works only for n > 0 && n <= 9. Update for n <= 16
func DecodeISO8859(n int, in []byte) string {
	if n == 0 || n > 9 {
		panic("not supported table number")
	}
	return iso8859tab[n-1](in)
}

// Decode treats in as text encoded according to EN 300 468 Annex A. It uses
// appropriate conversion according to selection byte
func Decode(in []byte) string {
	if len(in) == 0 {
		return ""
	}
	sel := in[0]
	if sel >= 0x20 {
		return DecodeISO6937(in)
	}
	in = in[1:]
	switch sel {
	case 1:
		return DecodeISO8859_5(in)
	case 2:
		return DecodeISO8859_6(in)
	case 3:
		return DecodeISO8859_7(in)
	case 4:
		return DecodeISO8859_8(in)
	case 5:
		return DecodeISO8859_9(in)
	case 0x10:
		if len(in) < 2 {
			break
		}
		n := int(uint16(in[0])<<8 | uint16(in[1]))
		if n > 0 && n <= 9 { // TODO: support n <= 16
			return DecodeISO8859(n, in)
		}
	}
	// Assume UTF8
	return string(in)
}
