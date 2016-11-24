package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/brentp/irelate"
	"github.com/brentp/irelate/interfaces"
	"github.com/brentp/irelate/parsers"
	"github.com/xuyangy/vcfanno/api"
	"github.com/xuyangy/vcfanno/shared"
	"github.com/brentp/xopen"
)

func benchmarkAnno(b *testing.B) {
	var configs shared.Config
	if _, err := toml.DecodeFile("example/conf.toml", &configs); err != nil {
		panic(err)
	}

	out := bufio.NewWriter(ioutil.Discard)
	Lua, _ := xopen.Ropen("example/custom.lua")
	lbytes, _ := ioutil.ReadAll(Lua)
	l_string := string(lbytes)

	srcs, err := configs.Sources()
	if err != nil {
		log.Fatal(err)
	}
	empty := make([]api.PostAnnotation, 0)
	for n := 0; n < b.N; n++ {
		if err != nil {
			log.Fatal(err)
		}
		a := api.NewAnnotator(srcs, l_string, false, true, empty)
		qrdr, err := xopen.Ropen("example/query.vcf.gz")
		if err != nil {
			log.Fatal(err)
		}
		qstream, query, err := parsers.VCFIterator(qrdr)
		queryables, err := a.Setup(query)
		if err != nil {
			log.Fatal(err)
		}

		//files := []string{"example/cadd.sub.txt.gz", "example/query.vcf.gz", "example/fitcons.bed.gz"}
		//files := []string{"example/query.vcf.gz", "example/fitcons.bed.gz"}

		fn := func(v interfaces.Relatable) {
			a.AnnotateEnds(v, api.INTERVAL)
		}
		stream := irelate.PIRelate(6000, 20000, qstream, false, fn, queryables...)

		for interval := range stream {
			fmt.Fprintf(out, "%s\n", interval)
		}
		out.Flush()
	}
}

func BenchmarkNormal(b *testing.B) { benchmarkAnno(b) }
