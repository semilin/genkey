package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"flag"
)

var Data TextData
var KPS []float64

//var SameKeyKPS []float64

func init() {
}

func main() {
	flag.StringVar(&ImproveFlag, "improve", "", "if set, decides which layout to improve")
	flag.BoolVar(&StaggerFlag, "stagger", false, "if true, calculates distance for row-stagger form factor")
	flag.BoolVar(&SlideFlag, "slide", false, "if true, ignores slideable sfbs")
	flag.BoolVar(&DynamicFlag, "dynamic", false, "")
	flag.Parse()
	origargs := os.Args[1:]
	var args []string
	for _, v := range origargs {
		if string(v[0]) != "-" {
			args = append(args, v)
		}
	}
	GeneratePositions()
	KPS = []float64{1.0, 4.0, 4.8, 5.7, 5.7, 4.8, 4.0, 1.0}
	//KPS = []float64{6, 16, 26.5, 40.36, 40.36, 26.5, 16, 6}
	//KPS = []float64{1, 3, 6, 8, 8, 6, 3, 1}
	//KPS = []float64{1, 1, 1, 1, 1, 1, 1, 1}

	Layouts = make(map[string]Layout)

	Layouts["qwerty"] = NewLayout("QWERTY", "qwertyuiopasdfghjkl;zxcvbnm,./")
	//Layouts["azerty"] = "azertyuiopqsdfghjklmwxcvbn',./"
	Layouts["dvorak"] = NewLayout("Dvorak", "',.pyfgcrlaoeuidhtns;qjkxbmwvz")
	Layouts["colemak"] = NewLayout("Colemak", "qwfpgjluy;arstdhneiozxcvbkm,./")
	Layouts["colemak dh"] = NewLayout("Colemak DH", "qwfpbjluy;arstgmneiozxcdvkh,./")
	// Layouts["funny colemak dh"] = "qwfpbjkuy;arstgmneiozxcdvlh,./"

	Layouts["colemaq"] = NewLayout("ColemaQ", ";wfpbjluyqarstgmneiozxcdkvh/.,")
	Layouts["colemaq-f"] = NewLayout("ColemaQ-f", ";wgpbjluyqarstfmneiozxcdkvh/.,")
	Layouts["colemak f"] = NewLayout("Colemak F", "qwgpbjluy;arstfmneiozxcdvkh,./")
	Layouts["colemak qi"] = NewLayout("Colemak Qi", "qlwmkjfuy'arstgpneiozxcdvbh,./")
	Layouts["colemak qix"] = NewLayout("Colemak Qi;x", ";lcmkjfuyqarstgpneiozxwdvbh/.,")
	// Layouts["NESO"] = "qylmkjfuc;airtgpnesoz.wdvbh/x,"
	// Layouts["NESO 2"] = "qylwvjfuc;airtgpneso.zkdmbh,x/"
	// "qulmkzbocyairtgpnesh.,wdjvf;x/"
	Layouts["isrt"] = NewLayout("ISRT", "yclmkzfu,'isrtgpneaoqvwdjbh/.x")
	// Layouts["hands down"] = "qchpvkyoj/rsntgwueiaxmldbzf',."
	// Layouts["norman"] = "qwdfkjurl;asetgyniohzxcvbpm,./"
	Layouts["mtgap"] = NewLayout("MTGAP", "ypoujkdlcwinea,mhtsrqz/.;bfgvx")
	Layouts["mtgap 2.0"] = NewLayout("MTGAP 2.0", ",fhdkjcul.oantgmseriqxbpzyw'v;")
	// Layouts["sind"] = "y,hwfqkouxsindcvtaerj.lpbgm;/z"
	// Layouts["rtna"] = "xdh.qbfoujrtna;gweislkm,/pczyv"
	// //Layouts["funny colemaq"] = "'wgdbmhuyqarstplneiozxcfkjv/.,"
	Layouts["workman"] = NewLayout("Workman", "qdrwbjfup;ashtgyneoizxmcvkl,./")
	// Layouts["workman ct"] = "wcldkjyru/ashtmpneoiqvgfbzx',."
	//Layouts["Colby's Funny"] = "/wgdbmho,qarstflneuizxcpkjv'.y"
	//Layouts["ISRT-AI"] = ",lcmkzfuy.arstgpneio;wvdjbh'qx"
	// Layouts["halmak"] = "wlrbz;qudjshnt,.aeoifmvc/gpxky"
	//Layouts["Balance Twelve but Funny"] = "pclmb'uoyknsrtg,aeihzfwdj/.'-x"
	//Layouts["Dynamica 0.1"] = "lfawqzghu,rnoibysetdjp/m'xckv."
	// Layouts["abc"] = "abcdefghijklmnopqrstuvwxyz,./'"
	//Layouts["TypeHack"] = "jghpfqvou;rsntkyiaelzwmdbc,'.x"
	// Layouts["qgmlwy"] = "qgmlwyfub;dstnriaeohzxcvjkp,./"
	//Layouts["TNWMLC"] = "tnwmlcbprhsgxjfkqzv;eadioyu,./"
	//Layouts["0.1"] = NewLayout("0.1", "vlafqzgu,ytronbmdeiskj/hpcw'.x")
	//Layouts["0.2"] = NewLayout("0.2", "ydlwkzfuo,strmcbneaiqj'gvph/x.")
	//Layouts["0.2mb"] = NewLayout("0.2mb", "kdl.gxfuoystrm,pneaivz'cwbh/qj")
	//Layouts["0.3"] = NewLayout("0.3", "kfawxqbulytsodchnerizv'gmp.,j/")
	//Layouts["0.4"] = NewLayout("0.4", "ymlkjqfau,scrtdbnoeixw'gvph/z.")
	//Layouts["0.5"] = NewLayout("0.5", "yluwqkfha.sredcmtnoixj'gpzvb/,")
	//Layouts["0.6"] = NewLayout("0.6", ".yuwfqzalvisedcmnort/x,gpbh'jk") // -rolling, +index balance
	//Layouts["0.7"] = NewLayout("0.7", "yhavzqwulfinotkcders/b.mjpg,'x")
	//Layouts["0.7a"] = NewLayout("0.7a", "yauvzqwhlfioetkcdnrs/.,mjpgb'x")
	//Layouts["0.7a3"] = NewLayout("0.7a3", "ylhvzqwuofirntkcdeas/'bmjpg,.x")
	Layouts["a"] = NewLayout("a", "yauwbxkclvioenpdhsrtj/,.qfmg'z")
	Layouts["1.0"] = NewLayout("1.0", "flhvzqwuoysrntkcdeaix'bmjpg,./")
	//Layouts["1.0x"] = NewLayout("1.0x", "flhvzqwuoysrntkcdeai'xbmjpg,./")

	// Layouts["0.7mv"] = NewLayout("0.7mv", "yhamzqwulfinotkcders/b.vjpg,'x")
	// Layouts["0.7idk"] = NewLayout("0.7idk", "yhamkqwulfinotvcders/b.jzpg,'x")

	Layouts["whorf"] = NewLayout("Whorf", "flhdmvwou,srntkgyaeixjbzqpc';.")
	Layouts["strtyp"] = NewLayout("strtyp", "jyuozkdlcwhiea,gtnsr'x/.qpbmfv")

	// Layouts["flaw"] = "flawpzkur/hsoycmtenibj'gvqd.x,"
	// Layouts["beakl"] = "qyouxgcrfzkhea.dstnbj/,i'wmlpv"
	// Layouts["owomak"] = "qwfpbjluy;arstdhneioxvcbzkm,./"
	Layouts["boo"] = NewLayout("Boo", ",.ucvzfmlyaoesgpntri;x'djbhkwq")
	Layouts["boo2"] = NewLayout("Boo 2", ",.ucvzfdlyaoesgpntri;x'wjbhmkq")

	Layouts["x1"] = NewLayout("X1", "kyo,'fclpvhieaudstnrzqj.;wgmbx")
	// Layouts["colemake"] = ";lgwvqpdu.arstkfnhio,jcmzb'y/x"
	// //Layouts["ctgap"] = "qwgdbmhuy'orstplneiazxcfkjv/.,"
	// Layouts["ctgap"] = "wcldkjyou/rsthmpneiazvgfbqx',."
	// Layouts["aptap"] = "wcdl'/youqrsthmpneiavbgk,.fjxz"
	// Layouts["rsthd"] = "jcyfkzl,uqrsthdmnaio/vgpbxw.;-"
	// Layouts["notgate"] = "youwg.vdlpiaescmhtrn'q;xzf,kjb"
	// Layouts["slider"] = "qwfpbjvuyzarscgmneio'ktdxlh/.,"
	// Layouts["paper 200"] = " wldk mic asthy nero bgf vpuj "

	//trigrams := Trigrams(Layouts["mtgap"].Keys)
	//ft100 := FastTrigrams(Layouts["mtgap"], 100)
	// ft250 := FastTrigrams(Layouts["qwerty"], 250)
	// ft500 := FastTrigrams(Layouts["qwerty"], 500)
	// ft1000 := FastTrigrams(Layouts["qwerty"], 1000)
	// ftfull := FastTrigrams(Layouts["qwerty"], 26000)
	//fmt.Println(ft100[0], ft100[1], ft100[2], ft100[3])
	// fmt.Println(ft250[0], ft250[1], ft250[2], ft250[3])
	// fmt.Println(ft500[0], ft500[1], ft500[2], ft500[3])
	// fmt.Println(ft1000[0], ft1000[1], ft1000[2], ft1000[3])
	// fmt.Println(ftfull[0], ftfull[1], ftfull[2], ftfull[3])
	//fmt.Println(trigrams[0], trigrams[1], trigrams[2], trigrams[3])

	if len(args) > 0 {
		if args[0] == "a" || args[0] == "analyze" {
			if len(args) == 1 {
				fmt.Println("You must provide the name of a layout to analyze")
				os.Exit(1)
			}
			Data = LoadData()

			input := strings.ToLower(args[1])
			PrintAnalysis(Layouts[input])
		} else if args[0] == "r" {
			Data = LoadData()

			type x struct {
				name string
				score float64
			}

			var sorted []x

			for _, v := range Layouts {
				sorted = append(sorted, x{v.Name, Score(v)})
			}

			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].score < sorted[j].score
			})

			for _, l := range sorted {
				spaces := strings.Repeat(".", 20-len(l.name))
				fmt.Printf("%s.%s%.2f\n", l.name, spaces, l.score)
			}
		} else if args[0] == "g" {
			Data = LoadData()
			start := time.Now()
			best := Populate(1000)
			end := time.Now()
			fmt.Println(end.Sub(start))

			optimal := Score(best)

			type x struct {
				name string
				score float64
			}

			var sorted []x

			for k, v := range Layouts {
				sorted = append(sorted, x{k, Score(v)})
			}

			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].score < sorted[j].score
			})
			
			for _, l := range sorted {
				spaces := strings.Repeat(".", 25-len(l.name))
				fmt.Printf("%s.%s%d%%\n", l.name, spaces, int(100*optimal/(Score(Layouts[l.name]))))
			}

		} else if args[0] == "sfbs" {
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := Layouts[input]
			total := 100*float64(SFBs(l.Keys))/float64(Data.TotalBigrams)
			sfbs := ListSFBs(l.Keys)
			SortFreqList(sfbs)
			fmt.Printf("%.2f%%\n", total)
			PrintFreqList(sfbs, 16)
		} else if args[0] == "dsfbs" {
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := Layouts[input]
			total := 100*float64(DSFBs(l.Keys))/float64(Data.TotalBigrams)
			dsfbs := ListDSFBs(l.Keys)
			SortFreqList(dsfbs)
			fmt.Printf("%.2f%%\n", total)
			PrintFreqList(dsfbs, 16)
		}else if args[0] == "bigrams" {
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := Layouts[input]
			sf := ListWeightedSameFinger(l.Keys)
			SortFreqList(sf)
			PrintFreqList(sf, 16)
		} else if args[0] == "dynamic" {
			Data = LoadData()
			if len(args) == 1 {
				fmt.Println("You must specify a layout")
				os.Exit(1)
			}
			input := strings.ToLower(args[1])
			l := Layouts[input]
			truecount, usage := SFBsMinusTop(l.Keys)
			total := 100*float64(usage)/float64(Data.TotalBigrams)
			dynamic, truesfbs := ListRepeats(l.Keys)
			SortFreqList(dynamic)
			SortFreqList(truesfbs)
			fmt.Printf("Dynamic Usage: %.2f%%\n", total)
			PrintFreqList(dynamic, 30)
			fmt.Printf("True SFBs: %.2f%%\n", 100*float64(truecount)/float64(Data.TotalBigrams))
			PrintFreqList(truesfbs, 8)
		} else if args[0] == "speed" {
			Data = LoadData()
			input := strings.ToLower(args[1])
			l := Layouts[input]
			speeds := FingerSpeed(l.Keys)
			fmt.Println("Unweighted Speed")
			for i, v := range speeds {
				fmt.Printf("\t%s: %.2f\n", FingerNames[i], v)
			}
			
		} else if args[0] == "h" {			
			Data = LoadData()
			Heatmap(Layouts[args[1]].Keys)
		} else if args[0] == "ngram" {
			Data = LoadData()
			total := float64(Data.Total)
			ngram := args[1]
			if len(ngram) == 1 {
				fmt.Printf("unigram: %.3f%%\n", 100*float64(Data.Letters[ngram]) / total)
			} else if len(ngram) == 2 {
				fmt.Printf("bigram: %.3f%%\n", 100*float64(Data.Bigrams[ngram]) / total)
				fmt.Printf("skipgram: %.3f%%\n", 100*Data.Skipgrams[ngram] / total)
			} else if len(ngram) == 3 {
				fmt.Printf("trigram: %.3f%%\n", 100*float64(Data.Trigrams[ngram]) / total)
			}
			// } else if args[0] == "i" {
			// 	LoadData()
			// 	input := strings.ToLower(args[1])
			// 	InteractiveAnalysis(Layouts[input])
		} else if args[0] == "load" {
			Data = GetTextData()
			WriteData(Data)
		}
	}
}
