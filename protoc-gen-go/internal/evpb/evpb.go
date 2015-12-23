// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2015 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package evpb

import (
	"path"
	"strconv"

	"github.com/golang/protobuf/protoc-gen-go/generator"
)

func init() {
	generator.RegisterPlugin(new(evpb))
}

type evpb struct {
	gen *generator.Generator
}

var (
	evpbPkg     string
	evpbPkgPath = "github.com/crackcomm/evpb"
)

func (g *evpb) P(args ...interface{}) { g.gen.P(args...) }
func (g *evpb) Name() string          { return "evpb" }

func (g *evpb) Init(gen *generator.Generator) {
	g.gen = gen
	evpbPkg = generator.RegisterUniquePackageName("evpb", nil)
}

func (g *evpb) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.MessageType) == 0 {
		return
	}
	g.P("import (")
	g.P(evpbPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, evpbPkgPath)))
	g.P(")")
	g.P()
}

func (g *evpb) Generate(file *generator.FileDescriptor) {
	fd := file.FileDescriptorProto
	if len(fd.MessageType) == 0 {
		return
	}

	for _, msg := range fd.MessageType {
		g.P()
		g.P("var topicName", msg.GetName(), " = ", `"`, fd.GetPackage(), ".", msg.GetName(), `"`)
		g.P()
		g.P("func Consume", msg.GetName(), "(q ", evpbPkg, ".Interface, h func(*", msg.GetName(), ") error) error {")
		g.P("  return q.Consume(topicName", msg.GetName(), ", func(body []byte) error {")
		g.P("    msg := new(", msg.GetName(), ")")
		g.P("    if err := proto.Unmarshal(body, msg); err != nil {")
		g.P("        return err")
		g.P("    }")
		g.P("    return h(msg)")
		g.P("  })")
		g.P("}")
		g.P()
		g.P("func Send", msg.GetName(), "(q ", evpbPkg, ".Interface, msg *", msg.GetName(), ") error {")
		g.P("  body, err := proto.Marshal(msg)")
		g.P("  if err != nil {")
		g.P("    return err")
		g.P("  }")
		g.P("  return q.Send(topicName", msg.GetName(), ", body)")
		g.P("}")
		g.P()
	}
}
