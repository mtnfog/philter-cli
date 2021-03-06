/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mtnfog/philter-sdk-golang/philter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	hostnamePtr := flag.String("h", "https://localhost:8080", "Philter endpoint, e.g. https://localhost:8080")
	filePtr := flag.String("f", "", "The file to process (optional).")
	filterProfilePtr := flag.String("p", "default", "The filter profile (optional).")
	contextPtr := flag.String("c", "default", "The context (optional).")
	documentId := flag.String("d", "default", "The document ID (optional).")
	explain := flag.Bool("e", false, "Explain (optional).")
	ignoreSsl := flag.Bool("i", false, "Ignore certificate errors.")
	version := flag.Bool("v", false, "Show version.")
	token := flag.String("t", "", "The API authentication token (optional).")

	flag.Parse()

	if *version == true {
		fmt.Println("Philter CLI 1.0.1")
		os.Exit(0)
	}

	if *hostnamePtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var text = ""

	if *filePtr == "" {

		bytes, _ := ioutil.ReadAll(os.Stdin)
		text = string(bytes)

	} else {

		if _, err := os.Stat(*filePtr); err != nil {
			fmt.Printf("Input file was not found.\n");
			flag.PrintDefaults()
			os.Exit(1)
		}

		content, err := ioutil.ReadFile(*filePtr)
		if err != nil {
			log.Fatal(err)
		}

		text = string(content)

	}

	if *ignoreSsl == true {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if *explain == true {

		var filterResponse = philter.Filter(*hostnamePtr, text, *contextPtr, *documentId, *filterProfilePtr, *token)

		json, err := json.Marshal(filterResponse)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print(string(json))

	} else {

		var filterResponse = philter.Filter(*hostnamePtr, text, *contextPtr, *documentId, *filterProfilePtr, *token)
		fmt.Print(filterResponse.FilteredText)

	}

}
