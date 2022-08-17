package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
	This tool convert a tab-delimited Mass Spectrometry SNP genotyping results to a software output, and output to standard out.
	Author: xuweixw27@gmail.com
	Date: Aug 16, 2022
*/

/*
	Example:
		Primers	r1311895	r1406821	r145504	r17661	r2493899
		MB128	A	C	G	CA	C
		MB140	A	C	G	C	C
		MB155	A	C	G	C	C
		MB164	AG	C	G		C
	Output:
		Primers	r1311895-A	r1311895-B	r1406821-A	r1406821-B	r145504-A	r145504-B	r17661-A	r17661-B	r2493899-A	r2493899-B
		MB128	A	A	C	C	G	G	C	A	C	C
		MB140	A	A	C	C	G	G	C	C	C	C
		MB155	A	A	C	C	G	G	C	C	C	C
		MB164	A	G	C	C	G	G	0	0	C	C
	Description:
		Each SNP site will contain two columns, and genotype missing will fill up by 0 value.
*/
const Example string = `Primers	r1311895	r1406821	r145504	r17661	r2493899
MB128	A	C	G	CA	C
MB140	A	C	G	C	C
MB155	A	C	G	C	C
MB164	AG	C	G		C`

type Header struct {
	SNPs []string
	Len  uint8
}

func (h *Header) String() string {
	var s strings.Builder
	s.WriteString("SNPMarkers")
	for _, marker := range h.SNPs {
		for _, chromatid := range [2]string{"A", "B"} {
			s.WriteString(fmt.Sprintf("\t%s-%s", marker, chromatid))
		}
	}
	return s.String()
}

type GenoTypeSet struct {
	header     *Header
	SampleName string
	Genotype   [][2]rune
}

func (gts *GenoTypeSet) String() string {
	var s strings.Builder
	s.WriteString(gts.SampleName)
	for _, gt := range gts.Genotype {
		s.WriteString(fmt.Sprintf("\t%s\t%s", string(gt[0]), string(gt[1])))
	}
	return s.String()
}

func Read(file string) *[]GenoTypeSet {
	lines := strings.Split(file, "\r\n")

	// first line for Header information
	var header Header
	head := strings.Split(lines[0], "\t")
	header.Len = uint8(len(head) - 1)
	header.SNPs = head[1:]

	// second and other lines for each sample genotypes
	var samples []GenoTypeSet
	for i := 1; i < len(lines); i++ { // each line

		genotypes := strings.Split(lines[i], "\t")
		sample := GenoTypeSet{header: &header,
			SampleName: genotypes[0],
			Genotype:   make([][2]rune, header.Len, header.Len),
		}
		sample.SampleName = genotypes[0]
		for j := 1; j < len(genotypes); j++ { // each genotype
			genotype := genotypes[j]
			switch len(genotype) {
			case 0: // blank genotype
				sample.Genotype[j-1][0], sample.Genotype[j-1][1] = '0', '0'
			case 1: // Homozygosity
				sample.Genotype[j-1][0], sample.Genotype[j-1][1] = rune(genotype[0]), rune(genotype[0])
			case 2: // Heterozygosity
				sample.Genotype[j-1][0], sample.Genotype[j-1][1] = rune(genotype[0]), rune(genotype[1])
			default:
				fmt.Printf("%b\n", genotype[:])
				log.Panicf("Unexpect line: %s %d unluckly\n", genotypes[0], j)
			}
		}
		samples = append(samples, sample)
	}
	return &samples
}

var (
	path = flag.String("in", "", "specify a tab-delimited file")
)

func main() {
	//a := Read(Example)
	//const path string = "/Users/apple/"

	flag.Parse()
	if *path == "" {
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(*path)
	if err != nil {
		log.Fatalln(err)
	}
	read, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	a := Read(string(read))
	fmt.Println((*a)[0].header.String())
	for _, b := range *a {
		fmt.Println(b.String())
	}
}
