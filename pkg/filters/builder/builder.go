/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package builder

import (
	"bytes"
	"fmt"
	"text/template"

	sprig "github.com/go-task/slim-sprig"
	"github.com/megaease/easegress/pkg/context"
	"gopkg.in/yaml.v3"
)

const (
	resultBuildErr = "buildErr"
)

type (
	// Builder is the base HTTP builder.
	Builder struct {
		template *template.Template
	}

	// Spec is the spec of Builder.
	Spec struct {
		LeftDelim       string `yaml:"leftDelim" jsonschema:"omitempty"`
		RightDelim      string `yaml:"rightDelim" jsonschema:"omitempty"`
		SourceNamespace string `yaml:"sourceNamespace" jsonschema:"omitempty"`
		Template        string `yaml:"template" jsonschema:"omitempty"`
	}
)

// Validate validates the Builder Spec.
func (spec *Spec) Validate() error {
	if spec.SourceNamespace == "" && spec.Template == "" {
		return fmt.Errorf("sourceNamespace or template must be specified")
	}

	if spec.SourceNamespace != "" && spec.Template != "" {
		return fmt.Errorf("sourceNamespace and template cannot be specified at the same time")
	}

	return nil
}

func (b *Builder) reload(spec *Spec) {
	if spec.SourceNamespace != "" {
		return
	}

	t := template.New("").Delims(spec.LeftDelim, spec.RightDelim)
	t.Funcs(sprig.TxtFuncMap()).Funcs(extraFuncs)
	b.template = template.Must(t.Parse(spec.Template))
}

func (b *Builder) build(data map[string]interface{}, v interface{}) error {
	var result bytes.Buffer

	if err := b.template.Execute(&result, data); err != nil {
		return err
	}

	if err := yaml.NewDecoder(&result).Decode(v); err != nil {
		return err
	}

	return nil
}

// Status returns status.
func (b *Builder) Status() interface{} {
	return nil
}

// Close closes Builder.
func (b *Builder) Close() {
}

func prepareBuilderData(ctx *context.Context) (map[string]interface{}, error) {
	requests := make(map[string]interface{})
	responses := make(map[string]interface{})

	for k, v := range ctx.Requests() {
		requests[k] = v.ToBuilderRequest(k)
	}

	for k, v := range ctx.Responses() {
		responses[k] = v.ToBuilderResponse(k)
	}

	return map[string]interface{}{
		"requests":  requests,
		"responses": responses,
		"data":      ctx.Data(),
	}, nil
}
