package ogen

import (
	"errors"
	"iter"
	"maps"
	"slices"
)

type (
	ProcessSchema = func(schema *Schema) error
	FixSchema     = func(err error, name string, schema *Schema) error
)

var (
	ErrMissingItemsInArray = errors.New("missing items in array")
)

// WalkAllSchemas walks through all [Schema] and call the given [ProcessSchema]
// function on each primitive [Schema], that is, all schemas except "object" and
// "array". If one of the following errors is occurred during the walking:
//   - [ErrMissingItemsInArray];
//
// then the given [FixSchema] function is called that accepts an error,
// an optional name of the schema taken from [Property.Name], and the current
// [Schema]. If no [FixSchema] is given then no-op function is using.
func WalkAllSchemas(spec *Spec, do ProcessSchema, fix FixSchema) error {
	if do == nil {
		return errors.New("do function is nil")
	}
	if fix == nil {
		fix = func(err error, name string, schema *Schema) error { return err }
	}

	for _, item := range spec.Paths {
		if err := walkAllSchemasInParameters(slices.Values(item.Parameters), do, fix); err != nil {
			return err
		}

		for _, op := range [8]*Operation{
			item.Get,
			item.Put,
			item.Post,
			item.Delete,
			item.Options,
			item.Head,
			item.Patch,
			item.Trace,
		} {
			if err := walkAllSchemasInOperation(op, do, fix); err != nil {
				return err
			}
		}
	}

	if err := walkAllSchemasInParameters(maps.Values(spec.Components.Parameters), do, fix); err != nil {
		return err
	}

	for _, s := range spec.Components.Schemas {
		if err := walkSchema("", s, do, fix); err != nil {
			return err
		}
	}

	return nil
}

func walkAllSchemasInParameters(params iter.Seq[*Parameter], do ProcessSchema, fix FixSchema) error {
	for p := range params {
		if err := walkSchema("", p.Schema, do, fix); err != nil {
			return err
		}

		for k := range p.Content {
			if err := walkSchema("", p.Content[k].Schema, do, fix); err != nil {
				return err
			}
		}
	}

	return nil
}

func walkAllSchemasInOperation(operation *Operation, do ProcessSchema, fix FixSchema) error {
	if operation == nil {
		return nil
	}

	if err := walkAllSchemasInParameters(slices.Values(operation.Parameters), do, fix); err != nil {
		return err
	}

	if operation.RequestBody != nil {
		for _, m := range operation.RequestBody.Content {
			if err := walkSchema("", m.Schema, do, fix); err != nil {
				return err
			}
		}
	}

	for _, r := range operation.Responses {
		for _, m := range r.Content {
			if err := walkSchema("", m.Schema, do, fix); err != nil {
				return err
			}
		}
	}

	return nil
}

func walkSchema(name string, schema *Schema, do ProcessSchema, fix FixSchema) error {
	if schema == nil || schema.Ref != "" {
		return nil
	}

	for _, s := range schema.AllOf {
		if err := walkSchema("", s, do, fix); err != nil {
			return err
		}
	}

	for _, s := range schema.OneOf {
		if err := walkSchema("", s, do, fix); err != nil {
			return err
		}
	}

	switch schema.Type {
	default:
		return do(schema)

	case "object":
		for _, prop := range schema.Properties {
			if err := walkSchema(prop.Name, prop.Schema, do, fix); err != nil {
				return err
			}
		}

	case "array":
		if schema.Items == nil {
			if err := fix(ErrMissingItemsInArray, name, schema); err != nil {
				return err
			}
		}

		if err := walkSchema("", schema.Items.Item, do, fix); err != nil {
			return err
		}

		for _, s := range schema.Items.Items {
			if err := walkSchema("", s, do, fix); err != nil {
				return err
			}
		}
	}

	return nil
}
