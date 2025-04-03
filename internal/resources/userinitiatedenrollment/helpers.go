package userinitiatedenrollment

import (
	"strings"
)

// GetISO639Code returns the ISO 639-1 two-letter code for a given language name
// If the language name is not found, it returns an empty string
func GetISO639Code(languageName string) string {
	// Normalize the input
	normalizedName := strings.TrimSpace(strings.ToLower(languageName))

	// Map of language names to ISO 639-1 codes
	languageMap := map[string]string{
		"abkhazian":           "ab",
		"afar":                "aa",
		"afrikaans":           "af",
		"akan":                "ak",
		"albanian":            "sq",
		"amharic":             "am",
		"arabic":              "ar",
		"aragonese":           "an",
		"armenian":            "hy",
		"assamese":            "as",
		"avaric":              "av",
		"avestan":             "ae",
		"aymara":              "ay",
		"azerbaijani":         "az",
		"azeri":               "az", // Alternative name for Azerbaijani
		"bambara":             "bm",
		"bashkir":             "ba",
		"basque":              "eu",
		"belarusian":          "be",
		"bengali":             "bn",
		"bangla":              "bn", // Alternative name for Bengali
		"bislama":             "bi",
		"bosnian":             "bs",
		"breton":              "br",
		"bulgarian":           "bg",
		"burmese":             "my",
		"myanmar":             "my", // Alternative name for Burmese
		"catalan":             "ca",
		"valencian":           "ca", // Alternative name for Catalan
		"chamorro":            "ch",
		"chechen":             "ce",
		"chichewa":            "ny",
		"chewa":               "ny", // Alternative name for Chichewa
		"nyanja":              "ny", // Alternative name for Chichewa
		"chinese":             "zh",
		"mandarin":            "zh", // Common variant of Chinese
		"church slavonic":     "cu",
		"old slavonic":        "cu", // Alternative name for Church Slavonic
		"old church slavonic": "cu", // Alternative name for Church Slavonic
		"chuvash":             "cv",
		"cornish":             "kw",
		"corsican":            "co",
		"cree":                "cr",
		"croatian":            "hr",
		"czech":               "cs",
		"danish":              "da",
		"divehi":              "dv",
		"dhivehi":             "dv", // Alternative name for Divehi
		"maldivian":           "dv", // Alternative name for Divehi
		"dutch":               "nl",
		"flemish":             "nl", // Alternative name for Dutch
		"dzongkha":            "dz",
		"bhutanese":           "dz", // Alternative name for Dzongkha
		"english":             "en",
		"esperanto":           "eo",
		"estonian":            "et",
		"ewe":                 "ee",
		"faroese":             "fo",
		"faeroese":            "fo", // Alternative spelling
		"fijian":              "fj",
		"finnish":             "fi",
		"french":              "fr",
		"western frisian":     "fy",
		"frisian":             "fy", // Simplified name
		"west frisian":        "fy", // Alternative name
		"fries":               "fy", // Alternative name
		"fulah":               "ff",
		"fula":                "ff", // Alternative name
		"fulani":              "ff", // Alternative name
		"gaelic":              "gd",
		"scottish gaelic":     "gd", // Full name
		"scots gaelic":        "gd", // Alternative name
		"galician":            "gl",
		"galego":              "gl", // Alternative name
		"ganda":               "lg",
		"luganda":             "lg", // Alternative name
		"georgian":            "ka",
		"german":              "de",
		"greek":               "el",
		"modern greek":        "el", // Clarification
		"kalaallisut":         "kl",
		"greenlandic":         "kl", // Alternative name
		"guarani":             "gn",
		"gujarati":            "gu",
		"haitian":             "ht",
		"haitian creole":      "ht", // Alternative name
		"hausa":               "ha",
		"hausan":              "ha", // Alternative name
		"hebrew":              "he",
		"herero":              "hz",
		"otjiherero":          "hz", // Alternative name
		"hindi":               "hi",
		"hiri motu":           "ho",
		"hungarian":           "hu",
		"magyar":              "hu", // Alternative name
		"icelandic":           "is",
		"ido":                 "io",
		"igbo":                "ig",
		"indonesian":          "id",
		"interlingua":         "ia",
		"interlingue":         "ie",
		"occidental":          "ie", // Alternative name
		"inuktitut":           "iu",
		"inupiaq":             "ik",
		"inupiat":             "ik", // Alternative name
		"inupiatun":           "ik", // Alternative name
		"irish":               "ga",
		"irish gaelic":        "ga", // Alternative name
		"italian":             "it",
		"japanese":            "ja",
		"javanese":            "jv",
		"kannada":             "kn",
		"kannadan":            "kn", // Alternative name
		"canarese":            "kn", // Alternative name
		"kanuri":              "kr",
		"kashmiri":            "ks",
		"koshur":              "ks", // Alternative name
		"kazakh":              "kk",
		"qazaq":               "kk", // Alternative name
		"central khmer":       "km",
		"khmer":               "km", // Simplified name
		"cambodian":           "km", // Alternative name
		"kikuyu":              "ki",
		"gikuyu":              "ki", // Alternative name
		"kinyarwanda":         "rw",
		"rwandan":             "rw", // Alternative name
		"rwanda":              "rw", // Alternative name
		"kirghiz":             "ky",
		"kyrgyz":              "ky", // Alternative spelling
		"komi":                "kv",
		"zyran":               "kv", // Alternative name
		"zyrian":              "kv", // Alternative name
		"komi-zyryan":         "kv", // Alternative name
		"kongo":               "kg",
		"kikongo":             "kg", // Alternative name
		"korean":              "ko",
		"kuanyama":            "kj",
		"kwanyama":            "kj", // Alternative name
		"cuanhama":            "kj", // Alternative name
		"oshikwanyama":        "kj", // Alternative name
		"kurdish":             "ku",
		"lao":                 "lo",
		"laotian":             "lo", // Alternative name
		"latin":               "la",
		"latvian":             "lv",
		"lettish":             "lv", // Alternative name
		"limburgan":           "li",
		"limburger":           "li", // Alternative name
		"limburgish":          "li", // Alternative name
		"lingala":             "ln",
		"lithuanian":          "lt",
		"luba-katanga":        "lu",
		"luba-shaba":          "lu", // Alternative name
		"luxembourgish":       "lb",
		"letzeburgesch":       "lb", // Alternative name
		"luxembourgian":       "lb", // Alternative name
		"macedonian":          "mk",
		"malagasy":            "mg",
		"malay":               "ms",
		"malayalam":           "ml",
		"maltese":             "mt",
		"manx":                "gv",
		"manx gaelic":         "gv", // Alternative name
		"maori":               "mi",
		"marathi":             "mr",
		"maharashtran":        "mr", // Alternative name
		"marshallese":         "mh",
		"ebon":                "mh", // Alternative name
		"mongolian":           "mn",
		"mongol":              "mn", // Alternative name
		"nauru":               "na",
		"nauruan":             "na", // Alternative name
		"navajo":              "nv",
		"navaho":              "nv", // Alternative name
		"north ndebele":       "nd",
		"northern ndebele":    "nd", // Alternative name
		"south ndebele":       "nr",
		"southern ndebele":    "nr", // Alternative name
		"ndonga":              "ng",
		"oshindonga":          "ng", // Alternative name
		"nepali":              "ne",
		"nepalese":            "ne", // Alternative name
		"gorkhali":            "ne", // Alternative name
		"norwegian":           "no",
		"norwegian bokmål":    "nb",
		"bokmål":              "nb", // Simplified name
		"norwegian nynorsk":   "nn",
		"nynorsk":             "nn", // Simplified name
		"occitan":             "oc",
		"provençal":           "oc", // Alternative name
		"provential":          "oc", // Alternative name
		"provencal":           "oc", // Alternative name
		"ojibwa":              "oj",
		"ojibwe":              "oj", // Alternative name
		"ojibway":             "oj", // Alternative name
		"otchipwe":            "oj", // Alternative name
		"ojibwemowin":         "oj", // Alternative name
		"oriya":               "or",
		"odia":                "or", // Alternative name
		"odian":               "or", // Alternative name
		"odishan":             "or", // Alternative name
		"orissan":             "or", // Alternative name
		"oromo":               "om",
		"oromoo":              "om", // Alternative name
		"ossetian":            "os",
		"ossetic":             "os", // Alternative name
		"ossete":              "os", // Alternative name
		"pali":                "pi",
		"pali-magadhi":        "pi", // Alternative name
		"pashto":              "ps",
		"pushto":              "ps", // Alternative name
		"persian":             "fa",
		"farsi":               "fa", // Alternative name
		"polish":              "pl",
		"portuguese":          "pt",
		"punjabi":             "pa",
		"panjabi":             "pa", // Alternative name
		"quechua":             "qu",
		"quechuan":            "qu", // Alternative name
		"romanian":            "ro",
		"moldavian":           "ro", // Alternative name
		"moldovan":            "ro", // Alternative name
		"romansh":             "rm",
		"romansch":            "rm", // Alternative name
		"rundi":               "rn",
		"kirundi":             "rn", // Alternative name
		"russian":             "ru",
		"northern sami":       "se",
		"north sami":          "se", // Alternative name
		"samoan":              "sm",
		"sango":               "sg",
		"sangoic":             "sg", // Alternative name
		"sanskrit":            "sa",
		"sardinian":           "sc",
		"sard":                "sc", // Alternative name
		"serbian":             "sr",
		"shona":               "sn",
		"sindhi":              "sd",
		"sinhala":             "si",
		"sinhalese":           "si", // Alternative name
		"slovak":              "sk",
		"slovakian":           "sk", // Alternative name
		"slovenian":           "sl",
		"slovene":             "sl", // Alternative name
		"somali":              "so",
		"somalian":            "so", // Alternative name
		"southern sotho":      "st",
		"sesotho":             "st", // Alternative name
		"sotho":               "st", // Alternative name
		"spanish":             "es",
		"castilian":           "es", // Alternative name
		"sundanese":           "su",
		"swahili":             "sw",
		"kiswahili":           "sw", // Alternative name
		"swati":               "ss",
		"swazi":               "ss", // Alternative name
		"swedish":             "sv",
		"tagalog":             "tl",
		"tahitian":            "ty",
		"tajik":               "tg",
		"tajiki":              "tg", // Alternative name
		"tamil":               "ta",
		"thamizh":             "ta", // Alternative name
		"tatar":               "tt",
		"telugu":              "te",
		"thai":                "th",
		"siamese":             "th", // Alternative name
		"central thai":        "th", // Alternative name
		"tibetan":             "bo",
		"standard tibetan":    "bo", // Alternative name
		"lhasa tibetan":       "bo", // Alternative name
		"tigrinya":            "ti",
		"tigrigna":            "ti", // Alternative name
		"tonga":               "to",
		"tongan":              "to", // Alternative name
		"tsonga":              "ts",
		"xitsonga":            "ts", // Alternative name
		"tswana":              "tn",
		"setswana":            "tn", // Alternative name
		"sechuana":            "tn", // Alternative name
		"turkish":             "tr",
		"turkmen":             "tk",
		"twi":                 "tw",
		"uighur":              "ug",
		"uyghur":              "ug", // Alternative name
		"ukrainian":           "uk",
		"urdu":                "ur",
		"uzbek":               "uz",
		"venda":               "ve",
		"tshivenda":           "ve", // Alternative name
		"vietnamese":          "vi",
		"volapük":             "vo",
		"walloon":             "wa",
		"welsh":               "cy",
		"wolof":               "wo",
		"xhosa":               "xh",
		"xosa":                "xh", // Alternative name
		"sichuan yi":          "ii",
		"nuosu":               "ii", // Alternative name
		"northern yi":         "ii", // Alternative name
		"liangshan yi":        "ii", // Alternative name
		"nosu":                "ii", // Alternative name
		"yiddish":             "yi",
		"judeo-german":        "yi", // Alternative name
		"yoruba":              "yo",
		"zhuang":              "za",
		"chuang":              "za", // Alternative name
		"zulu":                "zu",
		"isizulu":             "zu", // Alternative name
	}

	// Try to find an exact match
	if code, exists := languageMap[normalizedName]; exists {
		return code
	}

	// If no exact match, try to find a partial match
	// This is helpful for cases like "English (United States)" or "Brazilian Portuguese"
	for name, code := range languageMap {
		if strings.Contains(normalizedName, name) {
			return code
		}
	}

	// Handle special cases
	if strings.Contains(normalizedName, "filipino") || strings.Contains(normalizedName, "pilipino") {
		return "fil" // Filipino has a special ISO 639-2 code different from Tagalog
	}

	// No match found
	return ""
}

// GetLanguageNameFromCode returns the standard language name for a given ISO 639-1 code
// If the code is not found, it returns an empty string
func GetLanguageNameFromCode(code string) string {
	codeMap := map[string]string{
		"ab":  "Abkhazian",
		"aa":  "Afar",
		"af":  "Afrikaans",
		"ak":  "Akan",
		"sq":  "Albanian",
		"am":  "Amharic",
		"ar":  "Arabic",
		"an":  "Aragonese",
		"hy":  "Armenian",
		"as":  "Assamese",
		"av":  "Avaric",
		"ae":  "Avestan",
		"ay":  "Aymara",
		"az":  "Azerbaijani",
		"bm":  "Bambara",
		"ba":  "Bashkir",
		"eu":  "Basque",
		"be":  "Belarusian",
		"bn":  "Bengali",
		"bi":  "Bislama",
		"bs":  "Bosnian",
		"br":  "Breton",
		"bg":  "Bulgarian",
		"my":  "Burmese",
		"ca":  "Catalan",
		"ch":  "Chamorro",
		"ce":  "Chechen",
		"ny":  "Chichewa",
		"zh":  "Chinese",
		"cu":  "Church Slavonic",
		"cv":  "Chuvash",
		"kw":  "Cornish",
		"co":  "Corsican",
		"cr":  "Cree",
		"hr":  "Croatian",
		"cs":  "Czech",
		"da":  "Danish",
		"dv":  "Divehi",
		"nl":  "Dutch",
		"dz":  "Dzongkha",
		"en":  "English",
		"eo":  "Esperanto",
		"et":  "Estonian",
		"ee":  "Ewe",
		"fo":  "Faroese",
		"fj":  "Fijian",
		"fi":  "Finnish",
		"fr":  "French",
		"fy":  "Western Frisian",
		"ff":  "Fulah",
		"gd":  "Scottish Gaelic",
		"gl":  "Galician",
		"lg":  "Ganda",
		"ka":  "Georgian",
		"de":  "German",
		"el":  "Greek",
		"kl":  "Kalaallisut",
		"gn":  "Guarani",
		"gu":  "Gujarati",
		"ht":  "Haitian Creole",
		"ha":  "Hausa",
		"he":  "Hebrew",
		"hz":  "Herero",
		"hi":  "Hindi",
		"ho":  "Hiri Motu",
		"hu":  "Hungarian",
		"is":  "Icelandic",
		"io":  "Ido",
		"ig":  "Igbo",
		"id":  "Indonesian",
		"ia":  "Interlingua",
		"ie":  "Interlingue",
		"iu":  "Inuktitut",
		"ik":  "Inupiaq",
		"ga":  "Irish",
		"it":  "Italian",
		"ja":  "Japanese",
		"jv":  "Javanese",
		"kn":  "Kannada",
		"kr":  "Kanuri",
		"ks":  "Kashmiri",
		"kk":  "Kazakh",
		"km":  "Central Khmer",
		"ki":  "Kikuyu",
		"rw":  "Kinyarwanda",
		"ky":  "Kyrgyz",
		"kv":  "Komi",
		"kg":  "Kongo",
		"ko":  "Korean",
		"kj":  "Kuanyama",
		"ku":  "Kurdish",
		"lo":  "Lao",
		"la":  "Latin",
		"lv":  "Latvian",
		"li":  "Limburgan",
		"ln":  "Lingala",
		"lt":  "Lithuanian",
		"lu":  "Luba-Katanga",
		"lb":  "Luxembourgish",
		"mk":  "Macedonian",
		"mg":  "Malagasy",
		"ms":  "Malay",
		"ml":  "Malayalam",
		"mt":  "Maltese",
		"gv":  "Manx",
		"mi":  "Maori",
		"mr":  "Marathi",
		"mh":  "Marshallese",
		"mn":  "Mongolian",
		"na":  "Nauru",
		"nv":  "Navajo",
		"nd":  "North Ndebele",
		"nr":  "South Ndebele",
		"ng":  "Ndonga",
		"ne":  "Nepali",
		"no":  "Norwegian",
		"nb":  "Norwegian Bokmål",
		"nn":  "Norwegian Nynorsk",
		"oc":  "Occitan",
		"oj":  "Ojibwa",
		"or":  "Oriya",
		"om":  "Oromo",
		"os":  "Ossetian",
		"pi":  "Pali",
		"ps":  "Pashto",
		"fa":  "Persian",
		"pl":  "Polish",
		"pt":  "Portuguese",
		"pa":  "Punjabi",
		"qu":  "Quechua",
		"ro":  "Romanian",
		"rm":  "Romansh",
		"rn":  "Rundi",
		"ru":  "Russian",
		"se":  "Northern Sami",
		"sm":  "Samoan",
		"sg":  "Sango",
		"sa":  "Sanskrit",
		"sc":  "Sardinian",
		"sr":  "Serbian",
		"sn":  "Shona",
		"sd":  "Sindhi",
		"si":  "Sinhala",
		"sk":  "Slovak",
		"sl":  "Slovenian",
		"so":  "Somali",
		"st":  "Southern Sotho",
		"es":  "Spanish",
		"su":  "Sundanese",
		"sw":  "Swahili",
		"ss":  "Swati",
		"sv":  "Swedish",
		"tl":  "Tagalog",
		"ty":  "Tahitian",
		"tg":  "Tajik",
		"ta":  "Tamil",
		"tt":  "Tatar",
		"te":  "Telugu",
		"th":  "Thai",
		"bo":  "Tibetan",
		"ti":  "Tigrinya",
		"to":  "Tonga",
		"ts":  "Tsonga",
		"tn":  "Tswana",
		"tr":  "Turkish",
		"tk":  "Turkmen",
		"tw":  "Twi",
		"ug":  "Uighur",
		"uk":  "Ukrainian",
		"ur":  "Urdu",
		"uz":  "Uzbek",
		"ve":  "Venda",
		"vi":  "Vietnamese",
		"vo":  "Volapük",
		"wa":  "Walloon",
		"cy":  "Welsh",
		"wo":  "Wolof",
		"xh":  "Xhosa",
		"ii":  "Sichuan Yi",
		"yi":  "Yiddish",
		"yo":  "Yoruba",
		"za":  "Zhuang",
		"zu":  "Zulu",
		"fil": "Filipino",
	}

	if name, exists := codeMap[strings.ToLower(code)]; exists {
		return name
	}
	return ""
}
