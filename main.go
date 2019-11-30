// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Program awsping pings Amazon Web Services regions and reports about the latency.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

// TODO Add more regions
var endpoints = map[string]string{
  "eu-west-1":            "awsping-eu-west-1.s3-website-eu-west-1.amazonaws.com",
  "eu-north-1":           "awsping-eu-north-1.s3-website.eu-north-1.amazonaws.com",
  "eu-west-3":            "awsping-eu-west-3.s3-website.eu-west-3.amazonaws.com",
  "eu-central-1":         "awsping-eu-central-1.s3-website.eu-central-1.amazonaws.com",
  "eu-west-2":            "awsping-eu-west-2.s3-website.eu-west-2.amazonaws.com",
  "ap-northeast-1":       "awsping-ap-northeast-1.s3-website-ap-northeast-1.amazonaws.com",
  "ap-northeast-2":       "awsping-ap-northeast-2.s3-website.ap-northeast-2.amazonaws.com",
  "ap-southeast-1":       "awsping-ap-southeast-1.s3-website-ap-southeast-1.amazonaws.com",
  "ap-southeast-2":       "awsping-ap-southeast-2.s3-website-ap-southeast-2.amazonaws.com",
  "ap-south-1":           "awsping-ap-south-1.s3-website.ap-south-1.amazonaws.com",
  "ca-central-1":         "awsping-ca-central-1.s3-website.ca-central-1.amazonaws.com",
  "sa-east-1":            "awsping-sa-east-1.s3-website-sa-east-1.amazonaws.com",
  "us-east-2":            "awsping-us-east-2.s3-website.us-east-2.amazonaws.com",
  "us-west-1":            "awsping-us-west-1.s3-website-us-west-1.amazonaws.com",
  "us-west-2":            "awsping-us-west-2.s3-website-us-west-2.amazonaws.com",
}

var (
	top         bool
	number      int // number of requests for each region
	concurrency int
	timeout     time.Duration
	csv         bool
	verbose     bool
	// TODO(jbd): Add payload options such as body size.

	client  *http.Client // TODO(jbd): One client per worker?
	inputs  chan input
	outputs chan output
)

func main() {
	flag.BoolVar(&top, "top", false, "")
	flag.IntVar(&number, "n", 10, "")
	flag.IntVar(&concurrency, "c", 10, "")
	flag.DurationVar(&timeout, "t", time.Duration(0), "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&csv, "csv", false, "")

	flag.Usage = usage
	flag.Parse()

	if number < 0 || concurrency <= 0 {
		usage()
	}
	if csv {
		verbose = false // if output is CSV, no need for verbose output
	}

	client = &http.Client{
		Timeout: timeout,
	}

	go start()
	inputs = make(chan input, concurrency)
	outputs = make(chan output, number*len(endpoints))
	for i := 0; i < number; i++ {
		for r, e := range endpoints {
			inputs <- input{region: r, endpoint: e}
		}
	}
	close(inputs)
	report()
}

func start() {
	for worker := 0; worker < concurrency; worker++ {
		go func() {
			for m := range inputs {
				m.HTTP()
			}
		}()
	}
}

func report() {
	m := make(map[string]output)
	for i := 0; i < number*len(endpoints); i++ {
		o := <-outputs

		a := m[o.region]

		a.region = o.region
		a.durations = append(a.durations, o.durations[0])
		a.errors += o.errors

		m[o.region] = a
	}
	all := make([]output, 0, len(m))
	for _, t := range m {
		all = append(all, t)
	}

	// sort all by median duration.
	sort.Slice(all, func(i, j int) bool {
		return all[i].median() < all[j].median()
	})

	if top {
		t := all[0].region
		if t == "global" {
			t = all[1].region
		}
		fmt.Print(t)
		return
	}

	tr := tabwriter.NewWriter(os.Stdout, 3, 2, 2, ' ', 0)
	for i, a := range all {
		fmt.Fprintf(tr, "%2d.\t[%v]\t%v", i+1, a.region, a.median())
		if a.errors > 0 {
			fmt.Fprintf(tr, "\t(%d errors)", a.errors)
		}
		fmt.Fprintln(tr)
	}
	tr.Flush()
}

func usage() {
	fmt.Println(usageText)
	os.Exit(0)
}

var usageText = `awsping [options...]

Options:
-n   Number of requests to be made to each region.
     By default 10; can't be negative.
-c   Max number of requests to be made at any time.
     By default 10; can't be negative or zero.
-t   Timeout. By default, no timeout.
     Examples: "500ms", "1s", "1s500ms".
-top If true, only the top (non-global) region is printed.

-csv CSV output; disables verbose output.
-v   Verbose output.
`
