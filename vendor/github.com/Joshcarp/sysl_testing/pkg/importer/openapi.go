package importer

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
)

func LoadOpenAPIText(args OutputData, text string, logger *logrus.Logger) (out string, err error) {
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData([]byte(text))
	if err != nil {
		return "", err
	}

	l := &loader{
		logger: logger,
		spec:   swagger,
		types:  TypeList{},
	}
	l.initTypes()
	endpoints := l.initEndpoints()

	result := &bytes.Buffer{}
	w := newWriter(result, logger)
	if err := w.Write(l.initInfo(args), l.types, "", endpoints...); err != nil {
		return "", err
	}
	return result.String(), nil
}

type loader struct {
	logger       *logrus.Logger
	spec         *openapi3.Swagger
	types        TypeList
	globalParams Parameters
}

func (l *loader) initInfo(args OutputData) SyslInfo {
	info := SyslInfo{
		OutputData:  args,
		Title:       l.spec.Info.Title,
		Description: l.spec.Info.Description,
		OtherFields: []string{},
	}
	values := []string{
		"version", l.spec.Info.Version,
		"license", l.spec.Info.Version,
		"termsOfService", l.spec.Info.TermsOfService,
	}
	for i := 0; i < len(values); i += 2 {
		key := values[i]
		val := values[i+1]
		if val != "" {
			info.OtherFields = append(info.OtherFields, key, val)
		}
	}
	return info
}

func (l *loader) initTypes() {
	// First init the swagger -> sysl mappings
	var swaggerToSyslMappings = map[string]string{
		"boolean": "bool",
		"date":    "date",
	}
	for swaggerName, syslName := range swaggerToSyslMappings {
		l.types.Add(&ImportedBuiltInAlias{
			name:   swaggerName,
			Target: &SyslBuiltIn{syslName},
		})
	}
	for name, schema := range l.spec.Components.Schemas {
		if _, has := l.types.Find(name); !has {
			_ = l.typeFromSchema(name, schema.Value)
		}
	}
	l.types.Sort()
}

func (l *loader) typeFromRef(path string) Type {
	path = strings.TrimPrefix(path, "#/components/schemas/")
	if t, has := checkBuiltInTypes(path); has {
		return t
	}
	if t, ok := l.types.Find(path); ok {
		return t
	}
	if schema, has := l.spec.Components.Schemas[path]; has {
		return l.typeFromSchema(path, schema.Value)
	}
	return nil
}

func (l *loader) typeFromSchemaRef(name string, ref *openapi3.SchemaRef) Type {
	if ref == nil {
		return nil
	}
	if t := l.typeFromRef(ref.Ref); t != nil {
		return t
	}
	return l.typeFromSchema(name, ref.Value)
}

func sortProperties(props FieldList) {
	sort.SliceStable(props, func(i, j int) bool {
		return strings.Compare(props[i].Name, props[j].Name) < 0
	})
}

func (l *loader) typeFromSchema(name string, schema *openapi3.Schema) Type {
	for _, check := range []string{name, mapSwaggerTypeAndFormatToType(schema.Type, schema.Format, l.logger)} {
		if t, found := l.types.Find(check); found {
			return t
		}
	}
	switch schema.Type {
	case ObjectTypeName, "":
		t := &StandardType{
			name:       name,
			Properties: FieldList{},
		}
		for pname, pschema := range schema.Properties {
			var fieldType Type
			if pschema.Value.Type == ArrayTypeName {
				if atype := l.typeFromRef(pschema.Value.Items.Ref); atype != nil {
					fieldType = &Array{Items: atype}
				} else if atype := l.typeFromRef(pschema.Value.Items.Value.Type); atype != nil {
					fieldType = &Array{Items: atype}
				}
			}
			if fieldType == nil {
				fieldType = l.typeFromSchemaRef(name+"_"+pname, pschema)
			}
			f := Field{
				Name: pname,
				Type: fieldType,
			}
			if !contains(pname, schema.Required) {
				f.Optional = true
			}
			t.Properties = append(t.Properties, f)
		}
		sortProperties(t.Properties)
		return l.types.AddAndRet(t)
	case ArrayTypeName:
		t := &Array{
			name:  name,
			Items: l.typeFromSchemaRef(name+"_obj", schema.Items),
		}
		if name != "" {
			return l.types.AddAndRet(t)
		}
		return t
	default:
		baseType := mapSwaggerTypeAndFormatToType(schema.Type, schema.Format, l.logger)
		if t, found := l.types.Find(baseType); found {
			return t
		}
		if s, has := l.spec.Components.Schemas[schema.Type]; has {
			return l.typeFromSchemaRef(schema.Type, s)
		}

		l.logger.Warnf("unknown schema.Type: %s", schema.Type)
		return l.types.AddAndRet(NewStringAlias(name))
	}
}

func (l *loader) initEndpoints() []MethodEndpoints {
	epMap := map[string][]Endpoint{}

	l.initGlobalParams()

	for path, item := range l.spec.Paths {
		ops := map[string]*openapi3.Operation{
			"GET":    item.Get,
			"PUT":    item.Put,
			"POST":   item.Post,
			"DELETE": item.Delete,
			"PATCH":  item.Patch,
		}

		params := l.buildParams(item.Parameters)

		for method, op := range ops {
			if op != nil {
				epMap[method] = append(epMap[method], l.initEndpoint(path, op, params))
			}
		}
	}

	for key := range epMap {
		key := key
		sort.SliceStable(epMap[key], func(i, j int) bool {
			return strings.Compare(epMap[key][i].Path, epMap[key][j].Path) < 0
		})
	}

	var result []MethodEndpoints
	for _, method := range methodDisplayOrder {
		if eps, ok := epMap[method]; ok {
			result = append(result, MethodEndpoints{
				Method:    method,
				Endpoints: eps,
			})
		}
	}
	return result
}

func (l *loader) initEndpoint(path string, op *openapi3.Operation, params Parameters) Endpoint {
	var responses []Response
	typePrefix := strings.NewReplacer(
		"/", "_",
		"{", "_",
		"}", "_",
		"-", "_").Replace(path) + "_"
	for statusCode, resp := range op.Responses {
		text := "error"
		if statusCode[0] == '2' {
			text = "ok"
		}
		respType := &StandardType{
			name:       typePrefix + text,
			Properties: FieldList{},
		}
		for mediaType, val := range resp.Value.Content {
			t := l.typeFromSchemaRef("", val.Schema)
			f := Field{
				Name:       t.Name(),
				Attributes: []string{fmt.Sprintf("mediatype=\"%s\"", mediaType)},
				Type:       t,
			}
			respType.Properties = append(respType.Properties, f)
		}
		for name, header := range resp.Value.Headers {
			f := Field{
				Name:       name,
				Attributes: []string{"~header"},
				Type:       l.typeFromSchemaRef("", header.Value.Schema),
			}
			if f.Type == nil {
				f.Type = StringAlias
			}
			respType.Properties = append(respType.Properties, f)
		}

		r := Response{Text: text}
		if len(respType.Properties) > 0 {
			if len(respType.Properties) == 1 && respType.Properties[0].Attributes[0] != "~header" {
				r.Type = respType.Properties[0].Type
			} else {
				sortProperties(respType.Properties)
				l.types.Add(respType)
				r.Type = respType
			}
		}
		responses = append(responses, r)
	}

	res := Endpoint{
		Path:        path,
		Description: op.Description,
		Responses:   responses,
		Params:      params.Extend(l.buildParams(op.Parameters)),
	}

	if op.RequestBody != nil {
		for mediaType, content := range op.RequestBody.Value.Content {
			t := l.typeFromSchemaRef("", content.Schema)
			p := Param{
				Field: Field{
					Name:       t.Name() + "Request",
					Type:       t,
					Optional:   !op.RequestBody.Value.Required,
					Attributes: []string{fmt.Sprintf("mediatype=\"%s\"", mediaType)},
					SizeSpec:   nil,
				},
				In: "body",
			}
			res.Params.Add(p)
		}
	}
	return res
}

func (l *loader) initGlobalParams() {
	l.globalParams = Parameters{}
	for _, param := range l.spec.Components.Parameters {
		l.globalParams.Add(l.buildParam(param.Value))
	}
}

func (l *loader) buildParams(params openapi3.Parameters) Parameters {
	out := Parameters{}
	for _, param := range params {
		out.Add(l.buildParam(param.Value))
	}
	return out
}

func (l *loader) buildParam(p *openapi3.Parameter) Param {
	t := l.typeFromSchemaRef("", p.Schema)
	return Param{
		Field: Field{
			Name:       p.Name,
			Type:       t,
			Optional:   !p.Required,
			Attributes: nil,
			SizeSpec:   nil,
		},
		In: p.In,
	}
}
