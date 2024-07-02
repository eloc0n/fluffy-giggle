package dynamicvariables

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Blacklisted tokens
var blacklisted = []string{
	"__locals__",
	"__globals__",
	"()",
	"eval(",
	"exec(",
	"__import__",
	"__call__",
	"#!/bin/bash",
	"#!/bin/sh",
}

// GenericDynamicVariables implements the common notation for requested dynamic variables
type GenericDynamicVariables struct {
	TaskToVariableDelimiter string
	DependencyPattern       *regexp.Regexp
	VariablePattern         *regexp.Regexp
}

// NewGenericDynamicVariables initializes a new instance of GenericDynamicVariables
func NewGenericDynamicVariables() *GenericDynamicVariables {
	dependencyPattern := regexp.MustCompile(`\$\[\[(.*?)\]\]\$`)
	variablePattern := regexp.MustCompile(`\<\<.*?\>\>`)
	return &GenericDynamicVariables{
		TaskToVariableDelimiter: ":!:",
		DependencyPattern:       dependencyPattern,
		VariablePattern:         variablePattern,
	}
}

// SplitNotation splits the dependency into 'node-label' and 'variables'
func (gdv *GenericDynamicVariables) SplitNotation(dependency string) (string, string, error) {
	parts := strings.SplitN(dependency, gdv.TaskToVariableDelimiter, 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid dependency format")
	}
	return parts[0], parts[1], nil
}

// Retrieve retrieves the dependency from the dependency pattern
func (gdv *GenericDynamicVariables) Retrieve(target string) ([]string, error) {
	_target := replaceTargetEditorTags(target)
	matches := gdv.DependencyPattern.FindAllStringSubmatch(_target, -1)

	var dependencies []string
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		dependencies = append(dependencies, match[1])
	}

	for _, dependency := range dependencies {
		nodeName, vars, err := gdv.SplitNotation(dependency)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(nodeName) == "" {
			return nil, errors.New("node-name cannot be a whitespace string")
		}
		for _, v := range gdv.VariablesInDependency(vars) {
			if strings.TrimSpace(v) == "" {
				return nil, errors.New("all variables cannot be whitespace only")
			}
		}
	}
	return removeDuplicates(dependencies), nil
}

// VariablesInDependency retrieves the variables within a dependency
func (gdv *GenericDynamicVariables) VariablesInDependency(dependency string) []string {
	matches := gdv.VariablePattern.FindAllStringSubmatch(dependency, -1)
	var variables []string
	for _, match := range matches {
		if len(match) < 1 {
			continue
		}
		variables = append(variables, strings.Trim(match[0], "<>"))
	}
	return variables
}

// AssertAllValid asserts that a list of variables does not contain blacklisted tokens
func (gdv *GenericDynamicVariables) AssertAllValid(variables []string) error {
	for _, v := range variables {
		if err := assertIsVariableValid(v); err != nil {
			return err
		}
	}
	return nil
}

// VariableToNotation transforms a dependency to initial notation
func (gdv *GenericDynamicVariables) VariableToNotation(variable string) string {
	return fmt.Sprintf("$[[%s]]$", variable)
}

// IsPureVariable checks if a string is an exact match of a dynamic variable
func (gdv *GenericDynamicVariables) IsPureVariable(str string, variables []string) bool {
	if len(variables) != 1 {
		return false
	}
	return gdv.VariableToNotation(variables[0]) == str
}

// Inject injects dependencies' values into a target
func (gdv *GenericDynamicVariables) Inject(target interface{}, dependencies map[string]interface{}) interface{} {
	switch t := target.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, v := range t {
			result[k] = gdv.Inject(v, dependencies)
		}
		return result
	case string:
		_target := replaceTargetEditorTags(t)
		variables, err := gdv.Retrieve(_target)
		if err != nil {
			return t
		}
		if gdv.IsPureVariable(_target, variables) {
			return dependencies[variables[0]]
		}
		for _, v := range variables {
			_target = strings.ReplaceAll(_target, gdv.VariableToNotation(v), fmt.Sprint(dependencies[v]))
		}
		return _target
	case []interface{}:
		return gdv.InjectIntoList(t, dependencies)
	default:
		return target
	}
}

// InjectIntoList injects dependencies' values into a list of values
func (gdv *GenericDynamicVariables) InjectIntoList(targetList []interface{}, dependencies map[string]interface{}) []interface{} {
	var result []interface{}
	for _, item := range targetList {
		injected := gdv.Inject(item, dependencies)
		if items, ok := injected.([]interface{}); ok {
			result = append(result, items...)
		} else {
			result = append(result, injected)
		}
	}
	return result
}

// assertIsVariableValid checks if a variable is valid (not blacklisted)
func assertIsVariableValid(variable string) error {
	for _, b := range blacklisted {
		if strings.Contains(variable, b) {
			return fmt.Errorf("invalid dynamic variable: %s", variable)
		}
	}
	return nil
}

// replaceTargetEditorTags replaces target editor tags in a string
func replaceTargetEditorTags(target string) string {
	return strings.ReplaceAll(strings.ReplaceAll(target, "&lt;", "<"), "&gt;", ">")
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	var result []string
	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

// func main() {
// 	// Example usage
// 	gdv := NewGenericDynamicVariables()
// 	// target := "$[[Task-Title:!:<<Target-Label>>]]$"
// 	target := []string{
// 		"Task-Title:!:<<Target-Label>>",
// 		"Task-Title:!:<<Target-Label1>>",
// 		"Task-Title:!:<<Target-Label2>>",
// 		// "Task-Title:!:<<Target-Label3>>",
// 	}
// 	dependencies := map[string]interface{}{
// 		"Task-Title:!:<<Target-Label>>":  "aaaaaaaaaaa",
// 		"Task-Title:!:<<Target-Label1>>": "skata",
// 		"Task-Title:!:<<Target-Label2>>": "poutsa",
// 		"Task-Title:!:<<Target-Label3>>": "psoli",
// 	}
// 	fmt.Println(gdv.Inject(target, dependencies))
// }
