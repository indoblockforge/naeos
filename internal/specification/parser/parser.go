package parser

import (
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Parser interface {
	Parse(input string) (*SpecDocument, error)
}

type ParserFunc func(input string) (*SpecDocument, error)

func (f ParserFunc) Parse(input string) (*SpecDocument, error) {
	return f(input)
}

type SpecDocument struct {
	Raw  string
	Data any
}

func NewParser() Parser {
	return ParserFunc(func(input string) (*SpecDocument, error) {
		if input == "" {
			return nil, fmt.Errorf("input cannot be empty")
		}

		var root yaml.Node
		if err := yaml.Unmarshal([]byte(input), &root); err != nil {
			return nil, fmt.Errorf("parse spec: %w", err)
		}

		if len(root.Content) == 0 {
			return nil, fmt.Errorf("empty specification document")
		}

		value, err := parseYAMLNode(root.Content[0])
		if err != nil {
			return nil, err
		}

		return &SpecDocument{Raw: input, Data: value}, nil
	})
}

func parseYAMLNode(node *yaml.Node) (any, error) {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			return nil, fmt.Errorf("empty document")
		}
		return parseYAMLNode(node.Content[0])
	case yaml.MappingNode:
		result := map[string]any{}
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			if keyNode.Kind != yaml.ScalarNode {
				return nil, fmt.Errorf("map keys must be scalar")
			}
			value, err := parseYAMLNode(valueNode)
			if err != nil {
				return nil, err
			}
			result[keyNode.Value] = value
		}
		return result, nil
	case yaml.SequenceNode:
		result := make([]any, len(node.Content))
		for i, child := range node.Content {
			value, err := parseYAMLNode(child)
			if err != nil {
				return nil, err
			}
			result[i] = value
		}
		return result, nil
	case yaml.ScalarNode:
		return parseYAMLScalar(node)
	case yaml.AliasNode:
		if node.Alias == nil {
			return nil, fmt.Errorf("invalid alias node")
		}
		return parseYAMLNode(node.Alias)
	default:
		return nil, fmt.Errorf("unsupported YAML node kind %d", node.Kind)
	}
}

func parseYAMLScalar(node *yaml.Node) (any, error) {
	if node.Tag == "!!null" {
		return nil, nil
	}

	switch node.Tag {
	case "!!bool":
		return strconv.ParseBool(node.Value)
	case "!!int":
		return strconv.ParseInt(node.Value, 10, 64)
	case "!!float":
		return strconv.ParseFloat(node.Value, 64)
	case "!!str":
		return node.Value, nil
	default:
		if node.Value == "true" || node.Value == "false" {
			return strconv.ParseBool(node.Value)
		}
		if node.Value == "null" || node.Value == "~" {
			return nil, nil
		}
		if i, err := strconv.ParseInt(node.Value, 10, 64); err == nil {
			return i, nil
		}
		if f, err := strconv.ParseFloat(node.Value, 64); err == nil {
			return f, nil
		}
		return node.Value, nil
	}
}
