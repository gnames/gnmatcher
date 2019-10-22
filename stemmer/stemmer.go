package stemmer

import (
	"strings"
)

// stemmer package is responsible for extracting a stem of a latinized word. It
// is used to create a stem for latinized specific epithets in scientific names.
// Specific epithets are always nouns, so we need to take this into account.

// http://snowballstem.org/otherapps/schinke/
// http://caio.ueberalles.net/a_stemming_algorithm_for_latin_text_databases-schinke_et_al.pdf
//
// The Schinke Latin stemming algorithm is described in,
// Schinke R, Greengrass M, Robertson AM and Willett P (1996)
// A stemming algorithm for Latin text databases. Journal of Documentation, 52: 172-187.
//
// It has the feature that it stems each word to two forms, noun and verb. For example,
//
//                NOUN        VERB
//                ----        ----
//    aquila      aquil       aquila
//    portat      portat      porta
//    portis      port        por
//
// Here (slightly reformatted) are the rules of the stemmer,
//
// 1. (start)
//
// 2.  Convert all occurrences of the letters 'j' or 'v' to 'i' or 'u',
//     respectively.
//
// 3.  If the word ends in '-que' then
//         if the word is on the list shown in Figure 4, then
//             write the original word to both the noun-based and verb-based
//             stem dictionaries and go to 8.
//         else remove '-que'
//
//     [Figure 4 was
//
//         atque quoque neque itaque absque apsque abusque adaeque adusque denique
//         deque susque oblique peraeque plenisque quandoque quisque quaeque
//         cuiusque cuique quemque quamque quaque quique quorumque quarumque
//         quibusque quosque quasque quotusquisque quousque ubique undique usque
//         uterque utique utroque utribique torque coque concoque contorque
//         detorque decoque excoque extorque obtorque optorque retorque recoque
//         attorque incoque intorque praetorque]
//
// 4.  Match the end of the word against the suffix list show in Figure 6(a),
//     removing the longest matching suffix, (if any).
//
//     [Figure 6(a) was
//
//         -ibus -ius  -ae   -am   -as   -em   -es   -ia
//         -is   -nt   -os   -ud   -um   -us   -a    -e
//         -i    -o    -u]
//
// 5.  If the resulting stem contains at least two characters then write this stem
//     to the noun-based stem dictionary.
//
// 6.  Match the end of the word against the suffix list show in Figure 6(b),
//     identifying the longest matching suffix, (if any).
//
//     [Figure 6(b) was
//
//     -iuntur-beris -erunt -untur -iunt  -mini  -ntur  -stis
//     -bor   -ero   -mur   -mus   -ris   -sti   -tis   -tur
//     -unt   -bo    -ns    -nt    -ri    -m     -r     -s
//     -t]
//
//     If any of the following suffixes are found then convert them as shown:
//
//         '-iuntur', '-erunt', '-untur', '-iunt', and '-unt', to '-i';
//         '-beris', '-bor', and '-bo' to '-bi';
//         '-ero' to '-eri'
//
//     else remove the suffix in the normal way.
//
// 7.  If the resulting stem contains at least two characters then write this stem
//     to the verb-based stem dictionary.
//
// 8.  (end)
//

// object LatinStemmer {
// 	case class Word(originalStem: String, mappedStem: String, suffix: String)

//     val sb = new StringBuilder(word)

//     var sbIdx = 0
//     while (sbIdx < sb.length) {
//       (sb.charAt(sbIdx): @switch) match {
//         case 'j' => sb(sbIdx) = 'i'
//         case 'v' => sb(sbIdx) = 'u'
//         case _ =>
//       }
//       sbIdx += 1
//     }
var empty = struct{}{}

var queExceptions = map[string]struct{}{
	"atque": empty, "quoque": empty, "neque": empty, "itaque": empty,
	"absque": empty, "apsque": empty, "abusque": empty, "adaeque": empty,
	"adusque": empty, "denique": empty, "deque": empty, "susque": empty,
	"oblique": empty, "peraeque": empty, "plenisque": empty, "quandoque": empty,
	"quisque": empty, "quaeque": empty, "cuiusque": empty, "cuique": empty,
	"quemque": empty, "quamque": empty, "quaque": empty, "quique": empty,
	"quorumque": empty, "quarumque": empty, "quibusque": empty,
	"quosque": empty, "quasque": empty, "quotusquisque": empty,
	"quousque": empty, "ubique": empty, "undique": empty, "usque": empty,
	"uterque": empty, "utique": empty, "utroque": empty, "utribique": empty,
	"torque": empty, "coque": empty, "concoque": empty, "contorque": empty,
	"detorque": empty, "decoque": empty, "excoque": empty, "extorque": empty,
	"obtorque": empty, "optorque": empty, "retorque": empty, "recoque": empty,
	"attorque": empty, "incoque": empty, "intorque": empty, "praetorque": empty,
}

var nounSuffixes = []string{
	"ibus", "ius", "ae", "am", "as",
	"em", "es", "ia", "is",
	"nt", "os", "ud", "um", "us",
	"a", "e", "i", "o", "u",
}

type StemmedWord struct {
	Orig   string
	Stem   string
	Suffix string
}

func Stem(wrd string) StemmedWord {
	wrdR := []rune(wrd)
	for i, v := range wrdR {
		switch v {
		case 'j':
			wrdR[i] = 'i'
		case 'v':
			wrdR[i] = 'u'
		}
	}
	var sw StemmedWord
	var isException bool
	if sw, isException = processEndsWithQue(wrd, wrdR); isException {
		return sw
	}
	return checkNounSuffix(sw)
}

func processEndsWithQue(wrd string, wrdR []rune) (StemmedWord, bool) {
	sw := StemmedWord{Orig: wrd, Stem: string(wrdR)}

	if len(wrdR) < 3 {
		return sw, false
	}
	suffix := string(wrdR[len(wrdR)-3:])
	endsWithQue := suffix == "que"
	if endsWithQue {
		if _, ok := queExceptions[sw.Stem]; ok {
			return sw, true
		} else {
			sw.Stem = string(wrdR[:len(wrdR)-3])
		}
	}
	return sw, false
}

func checkNounSuffix(sw StemmedWord) StemmedWord {
	var found bool
	for _, v := range nounSuffixes {
		if strings.HasSuffix(sw.Stem, v) {
			if found {
				break
			}
			found = true
			wrdR := []rune(sw.Stem)
			stem := string(wrdR[:len(wrdR)-len(v)])
			if len(stem) >= 2 {
				sw.Stem = stem
				sw.Suffix = string(wrdR[len(v):])
			}
		}
	}
	return sw
}

//     val wordMapped = sb.toString
//     val wordEndsWithQue = word.endsWith("que")
//     if (wordEndsWithQue && queExceptions.contains(wordMapped)) {
//       Word(originalStem = word, mappedStem = wordMapped, suffix = "")
//     } else {
//       if (wordEndsWithQue) {
//         sb.delete(sb.length - 3, sb.length)
//       }
//       var found = false
//       val nounSuffixesIterator = nounSuffixes.iterator
//       while (!found && nounSuffixesIterator.hasNext) {
//         val suffix = nounSuffixesIterator.next()
//         val lastIndex = sb.lastIndexOf(suffix)
//         if (lastIndex >= 0 && sb.length - suffix.length == lastIndex) {
//           found = true
//           if (sb.length - suffix.length >= 2) {
//             sb.delete(sb.length - suffix.length, sb.length)
//           }
//         }
//       }
//       val stem = sb.toString()
//       Word(originalStem = word.substring(0, stem.length), mappedStem = stem,
//            suffix = word.substring(stem.length))
//     }
//   }
// }
