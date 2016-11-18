// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io/ioutil"
	"os/exec"
	"runtime"
)

const LICENSE = "" +
	"// Copyright 2016 Google Inc. All Rights Reserved.\n" +
	"//\n" +
	"// Licensed under the Apache License, Version 2.0 (the \"License\");\n" +
	"// you may not use this file except in compliance with the License.\n" +
	"// You may obtain a copy of the License at\n" +
	"//\n" +
	"//    http://www.apache.org/licenses/LICENSE-2.0\n" +
	"//\n" +
	"// Unless required by applicable law or agreed to in writing, software\n" +
	"// distributed under the License is distributed on an \"AS IS\" BASIS,\n" +
	"// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n" +
	"// See the License for the specific language governing permissions and\n" +
	"// limitations under the License.\n"

/// main program

func main() {
	base_schema := NewSchemaFromFile("schema.json")
	base_schema.resolveRefs(nil)
	base_schema.resolveAllOfs()

	openapi_schema := NewSchemaFromFile("openapi-2.0.json")
	// these non-object definitions are marked for handling as if they were objects
	// in the future, these could be automatically identified by their presence in a oneOf
	classNames := []string{
		"#/definitions/headerParameterSubSchema",
		"#/definitions/formDataParameterSubSchema",
		"#/definitions/queryParameterSubSchema",
		"#/definitions/pathParameterSubSchema"}
	openapi_schema.resolveRefs(classNames)
	openapi_schema.resolveAllOfs()
	openapi_schema.reduceOneOfs()

	// build a simplified model of the classes described by the schema
	cc := NewClassCollection(openapi_schema)
	// these pattern names are a bit of a hack until we find a more automated way to obtain them
	cc.PatternNames = map[string]string{
		"^x-": "vendorExtension",
		"^/":  "path",
		"^([0-9]{3})$|^(default)$": "responseCode",
	}
	cc.build()
	//log.Printf("%s\n", cc.display())
	if true {
		var err error

		// generate the protocol buffer description
		proto := cc.generateProto("OpenAPIv2")
		proto_filename := "openapi-v2.proto"
		err = ioutil.WriteFile(proto_filename, []byte(proto), 0644)
		if err != nil {
			panic(err)
		}

		// generate the compiler
		compiler := cc.generateCompiler("OpenAPIv2")
		go_filename := "openapi-v2.go"
		err = ioutil.WriteFile(go_filename, []byte(compiler), 0644)
		if err != nil {
			panic(err)
		}
		// autoformat the compiler
		err = exec.Command(runtime.GOROOT()+"/bin/gofmt", "-w", go_filename).Run()
	}
}