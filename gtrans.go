package main

import (
	"fmt"
	"github.com/Conight/go-googletrans"
	"github.com/chzyer/readline"
	"os"
	"strings"
)

type TranslatorCLI struct {
	Env        map[string]string
	Reader     *readline.Instance
	Translator *translator.Translator
}

var LanguageCodes = map[string]string{
	"af":    "afrikaans",
	"sq":    "albanian",
	"am":    "amharic",
	"ar":    "arabic",
	"hy":    "armenian",
	"az":    "azerbaijani",
	"eu":    "basque",
	"be":    "belarusian",
	"bn":    "bengali",
	"bs":    "bosnian",
	"bg":    "bulgarian",
	"ca":    "catalan",
	"ceb":   "cebuano",
	"ny":    "chichewa",
	"zh-cn": "chinese (simplified)",
	"zh-tw": "chinese (traditional)",
	"co":    "corsican",
	"hr":    "croatian",
	"cs":    "czech",
	"da":    "danish",
	"nl":    "dutch",
	"en":    "english",
	"eo":    "esperanto",
	"et":    "estonian",
	"tl":    "filipino",
	"fi":    "finnish",
	"fr":    "french",
	"fy":    "frisian",
	"gl":    "galician",
	"ka":    "georgian",
	"de":    "german",
	"el":    "greek",
	"gu":    "gujarati",
	"ht":    "haitian creole",
	"ha":    "hausa",
	"haw":   "hawaiian",
	"iw":    "hebrew",
	"he":    "hebrew",
	"hi":    "hindi",
	"hmn":   "hmong",
	"hu":    "hungarian",
	"is":    "icelandic",
	"ig":    "igbo",
	"id":    "indonesian",
	"ga":    "irish",
	"it":    "italian",
	"ja":    "japanese",
	"jw":    "javanese",
	"kn":    "kannada",
	"kk":    "kazakh",
	"km":    "khmer",
	"ko":    "korean",
	"ku":    "kurdish (kurmanji)",
	"ky":    "kyrgyz",
	"lo":    "lao",
	"la":    "latin",
	"lv":    "latvian",
	"lt":    "lithuanian",
	"lb":    "luxembourgish",
	"mk":    "macedonian",
	"mg":    "malagasy",
	"ms":    "malay",
	"ml":    "malayalam",
	"mt":    "maltese",
	"mi":    "maori",
	"mr":    "marathi",
	"mn":    "mongolian",
	"my":    "myanmar (burmese)",
	"ne":    "nepali",
	"no":    "norwegian",
	"or":    "odia",
	"ps":    "pashto",
	"fa":    "persian",
	"pl":    "polish",
	"pt":    "portuguese",
	"pa":    "punjabi",
	"ro":    "romanian",
	"ru":    "russian",
	"sm":    "samoan",
	"gd":    "scots gaelic",
	"sr":    "serbian",
	"st":    "sesotho",
	"sn":    "shona",
	"sd":    "sindhi",
	"si":    "sinhala",
	"sk":    "slovak",
	"sl":    "slovenian",
	"so":    "somali",
	"es":    "spanish",
	"su":    "sundanese",
	"sw":    "swahili",
	"sv":    "swedish",
	"tg":    "tajik",
	"ta":    "tamil",
	"te":    "telugu",
	"th":    "thai",
	"tr":    "turkish",
	"uk":    "ukrainian",
	"ur":    "urdu",
	"ug":    "uyghur",
	"uz":    "uzbek",
	"vi":    "vietnamese",
	"cy":    "welsh",
	"xh":    "xhosa",
	"yi":    "yiddish",
	"yo":    "yoruba",
	"zu":    "zulu",
}

func cmdSet(tr *TranslatorCLI, pairs []string) {
	for _, pair := range pairs {
		if key, val, found := strings.Cut(pair, "="); found {
			tr.Env[key] = val
		} else {
			fmt.Printf("Invalid key=value pair: %s\n", pair)
		}
	}
}

func cmdGet(tr *TranslatorCLI, key string) {
	if val, ok := tr.Env[key]; ok {
		fmt.Printf("= %s\n", val)
	} else {
		fmt.Println("(nil)")
	}
}

func validForTranslation(tr *TranslatorCLI) bool {
	src := tr.Env["source"]
	tgt := tr.Env["target"]
	if src != "auto" {
		if _, isLangCode := LanguageCodes[src]; !isLangCode {
			fmt.Printf(
				"Expected valid language code for `source`, found '%s'\n",
				src)
			return false
		}
	}
	if _, isLangCode := LanguageCodes[tgt]; !isLangCode {
		fmt.Printf(
			"Expected valid language code for `target`, found '%s'\n",
			tgt)
		return false
	}
	return true
}

func cmdTranslateText(tr *TranslatorCLI, text string) {
	if len(text) > 5000 {
		fmt.Printf(
			"Maximum text length is 5000, received text has length %d\n",
			len(text))
	} else if !validForTranslation(tr) {
		return
	}
	src := tr.Env["source"]
	tgt := tr.Env["target"]
	response, err := tr.Translator.Translate(text, src, tgt)
	if err != nil {
		fmt.Printf("Encountered error during translation: %v\n", err)
	}
	fmt.Printf(
		"%s [%s] -> %s [%s]\n",
		response.Src,
		LanguageCodes[src],
		response.Dest,
		LanguageCodes[tgt])
	fmt.Println(response.Text)
}

func cmdTranslateFile(tr *TranslatorCLI, text string) {
	if !validForTranslation(tr) {
		return
	}
	srcFilePath, destFilePath, ok := strings.Cut(text, " ")
	if !ok {
		var dir, fileName string
		if sepIndex := strings.LastIndex(srcFilePath, "/"); sepIndex != -1 {
			dir = srcFilePath[0 : sepIndex+1]
			fileName = "translated-" + srcFilePath[sepIndex+1:]
		} else {
			dir = ""
			fileName = srcFilePath
		}
		destFilePath = dir + fileName
	}
	inputFile, err := os.Open(srcFilePath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer inputFile.Close()
	var inputBuilder strings.Builder
	buffer := make([]byte, 8192)
	for {
		count, err := inputFile.Read(buffer)
		if err != nil {
			panic(err)
		}
		if count < len(buffer) {
			inputBuilder.Write(buffer[0:count])
			break
		} else {
			inputBuilder.Write(buffer)
		}
	}
	src := tr.Env["source"]
	tgt := tr.Env["target"]
	response, err := tr.Translator.Translate(
		inputBuilder.String(),
		src,
		tgt)
	if err != nil {
		fmt.Printf("Encountered error during translation: %v\n", err)
		return
	}
	fmt.Printf(
		"%s [%s] -> %s [%s]\n",
		response.Src,
		LanguageCodes[src],
		response.Dest,
		LanguageCodes[tgt])
	if destFilePath == "\\stdout" {
		fmt.Println(response.Dest)
	} else {
		err := os.WriteFile(destFilePath, []byte(response.Dest), 0644)
		if err != nil {
			fmt.Printf("Error when writing file: %v\n", err)
		}
	}
}

func printHelp() {
	fmt.Fprintln(os.Stderr, "| Get internal environment variable")
	fmt.Fprintln(os.Stderr, "| > get envVar")
	fmt.Fprintln(os.Stderr, "| Set internal environment variables in the form key=value")
	fmt.Fprintln(os.Stderr, "| > set source=es target=en")
	fmt.Fprintln(os.Stderr, "| Translate text after command from Env[\"source\"] to Env[\"target\"]")
	fmt.Fprintln(os.Stderr, "| By default Env[\"source\"] = \"auto\"")
	fmt.Fprintln(os.Stderr, "| > Krankenhaus")
	fmt.Fprintln(os.Stderr, "| Translate text in file in path {source} and output to file in path {dest}")
	fmt.Fprintln(os.Stderr, "| If {dest} is not an absolute path, it will be interpreted relative to the directory of {source}")
	fmt.Fprintln(os.Stderr, "| > trf {source} {dest}")
	fmt.Fprintln(os.Stderr, "| If {dest} is not provided translated text will be outputed to Stdout instead.")
	fmt.Fprintln(os.Stderr, "| > trf {source}")
	fmt.Fprintln(os.Stderr, "| Print this help message")
	fmt.Fprintln(os.Stderr, "| > help")
}

func (tr *TranslatorCLI) Start() {
	for {
		line, err := tr.Reader.Readline()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("Error occurred when reading input: %v\n", err)
		}
		switch {
		case strings.HasPrefix(line, "set "):
			cmdSet(tr, strings.Split(line[len("set "):], " "))
		case strings.HasPrefix(line, "get "):
			cmdGet(tr, line[len("get "):])
		case strings.HasPrefix(line, "trt "):
			cmdTranslateText(tr, line[len("trt "):])
		case strings.HasPrefix(line, "trf "):
			cmdTranslateFile(tr, line[len("trf "):])
		case strings.HasPrefix(line, "help"):
			printHelp()
		case line == "":
			// nop
		default:
			fmt.Println("Invalid input")
			printHelp()
		}
	}
}

func main() {
	config := translator.Config{
		UserAgent: []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
	}
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	tr := TranslatorCLI{
		Env: map[string]string{
			"source": "auto",
			"target": "sq",
		},
		Reader:     rl,
		Translator: translator.New(config),
	}
	tr.Start()
}
