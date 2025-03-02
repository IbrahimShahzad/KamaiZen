package document_manager

import (
	// "KamaiZen/logger"
	"errors"
	"fmt"
)

// Holds a map of function names to their corresponding documentation.
// It is used to store and retrieve documentation for multiple functions within a module.
type FunctionDocumentationMap struct {
	Functions map[string]FunctionDocumentation
}

// AddFunctionDoc adds a function documentation entry to the FunctionDocumentationMap.
// If the function documentation already exists and overwrite is set to false, it returns an error.
// If overwrite is set to true, it overwrites the existing function documentation.
//
// functionDoc: The FunctionDocumentation struct containing the function documentation to add.
// overwrite: A boolean indicating whether to overwrite existing documentation if it exists.
// return: An error if the function documentation already exists and overwrite is false.
func (f *FunctionDocumentationMap) AddFunctionDoc(functionDoc FunctionDocumentation, overwrite bool) error {
	if _, exists := f.Functions[functionDoc.Name]; exists {
		if !overwrite {
			return errors.New("Function Documentation already exists")
		} else {
			// logger.Debug("Overwrting function documentation: ", functionDoc.Name)
		}
	}
	f.Functions[functionDoc.Name] = functionDoc
	return nil
}

// Retrieves the documentation for a specific function as a formatted string.
// It looks up the function in the function documentation map and returns its string representation.
//
// functionName: The name of the function to retrieve documentation for.
//
//	A string containing the documentation for the specified function.
func (f *FunctionDocumentationMap) GetFunctionDocAsString(functionName string) string {
	return f.Functions[functionName].String()
}

// Holds the documentation details for a specific function.
// It includes the function's name, parameters, description, and example.
type FunctionDocumentation struct {
	Name        string // the name of the function.
	Parameters  string // the parameters of the function.
	Description string // a description of what the function does.
	Example     string // an example usage of the function.
	Parent      string // the parent module of the function.
}

// Returns a formatted string representation of the function documentation.
// It includes the function name, parameters, description, and example in a structured format.
//
// A string containing the formatted function documentation.
func (f FunctionDocumentation) String() string {
	return fmt.Sprintf("## Function:\n\t%s\n\n## Parameters:\n\t%s\n\n## Description:\n%s\n\n## Example:\n```\n%s\n```", f.Name, f.Parameters, f.Description, f.Example)
}
