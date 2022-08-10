package v3

import (
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v3"
)

type MediaType struct {
	Schema     low.NodeReference[*Schema]
	Example    low.NodeReference[any]
	Examples   low.NodeReference[map[low.KeyReference[string]]low.ValueReference[*Example]]
	Encoding   low.NodeReference[map[low.KeyReference[string]]low.ValueReference[*Encoding]]
	Extensions map[low.KeyReference[string]]low.ValueReference[any]
}

func (mt *MediaType) FindPropertyEncoding(eType string) *low.ValueReference[*Encoding] {
	return FindItemInMap[*Encoding](eType, mt.Encoding.Value)
}

func (mt *MediaType) FindExample(eType string) *low.ValueReference[*Example] {
	return FindItemInMap[*Example](eType, mt.Examples.Value)
}

func (mt *MediaType) GetAllExamples() map[low.KeyReference[string]]low.ValueReference[*Example] {
	return mt.Examples.Value
}

func (mt *MediaType) Build(root *yaml.Node) error {

	// extract extensions
	extensionMap, err := ExtractExtensions(root)
	if err != nil {
		return err
	}
	mt.Extensions = extensionMap

	// handle example if set.
	_, expLabel, expNode := utils.FindKeyNodeFull(ExampleLabel, root.Content)
	if expNode != nil {
		mt.Example = low.NodeReference[any]{Value: expNode.Value, KeyNode: expLabel, ValueNode: expNode}
	}

	// handle schema
	sch, sErr := ExtractSchema(root)
	if sErr != nil {
		return nil
	}
	if sch != nil {
		mt.Schema = *sch
	}

	// handle examples if set.
	exps, expsL, expsN, eErr := ExtractMapFlat[*Example](ExamplesLabel, root)
	if eErr != nil {
		return eErr
	}
	if exps != nil {
		mt.Examples = low.NodeReference[map[low.KeyReference[string]]low.ValueReference[*Example]]{
			Value:     exps,
			KeyNode:   expsL,
			ValueNode: expsN,
		}
	}

	// handle encoding
	encs, encsL, encsN, encErr := ExtractMapFlat[*Encoding](EncodingLabel, root)
	if encErr != nil {
		return err
	}
	if encs != nil {
		mt.Encoding = low.NodeReference[map[low.KeyReference[string]]low.ValueReference[*Encoding]]{
			Value:     encs,
			KeyNode:   encsL,
			ValueNode: encsN,
		}
	}
	return nil
}