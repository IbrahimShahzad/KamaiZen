package document_manager

import (
	"errors"
)

// Holds the documentation for all modules.
// It maps module names to their corresponding ModuleDocs structs.
type moduleDocumentationMap struct {
	ModuleDocs map[string]ModuleDocs
}

// GetModuleDocs retrieves the documentation for a specific module from the module documentation map.
//
// moduleName: The name of the module to retrieve documentation for.
// return: A ModuleDocs struct containing the documentation for the specified module and a boolean indicating
//
//	whether the module was found. If the module is not found, it returns an empty ModuleDocs struct and false.
func (m *moduleDocumentationMap) GetModuleDocs(moduleName string) (ModuleDocs, bool) {
	for key, value := range m.ModuleDocs {
		if key == moduleName {
			return value, true
		}
	}
	return ModuleDocs{}, false
}

// AddModuleDocs adds module documentation to the module documentation map.
// If the module documentation already exists and overwrite is set to false, it returns an error.
// If overwrite is set to true, it overwrites the existing module documentation.
//
// moduleName: The name of the module to add documentation for.
// moduleDocs: The ModuleDocs struct containing the documentation to add.
// overwrite: A boolean indicating whether to overwrite existing documentation if it exists.
// return: An error if the module documentation already exists and overwrite is false.
func (m *moduleDocumentationMap) AddModuleDocs(moduleName string, moduleDocs ModuleDocs, overwrite bool) error {
	if _, exists := m.ModuleDocs[moduleName]; exists && !overwrite {
		return errors.New("Module already exists")
	}
	m.ModuleDocs[moduleName] = moduleDocs
	return nil
}

type ModuleDocs struct {
	Functions map[string]FunctionDocumentationMap
}

// AddFunctionDoc adds function documentation to the specified module in the ModuleDocs.
// If the function documentation already exists and overwrite is set to false, it returns an error.
// If overwrite is set to true, it overwrites the existing function documentation.
//
// moduleName: The name of the module to add function documentation for.
// functionDocs: The FunctionDocumentationMap containing the function documentation to add.
// overwrite: A boolean indicating whether to overwrite existing documentation if it exists.
// return: An error if the function documentation already exists and overwrite is false.
func (m *ModuleDocs) AddFunctionDoc(moduleName string, functionDocs FunctionDocumentationMap, overwrite bool) error {
	if _, exists := m.Functions[moduleName]; exists {
		if !overwrite {
			return errors.New("Function Documentation already exists")
		} else {
			// logger.Debug("Overwriting function documentation for module: ", moduleName)
		}
	}
	m.Functions[moduleName] = functionDocs
	return nil
}

// GetFunctionDocAsString retrieves the documentation for a specific function within a specified module
// and returns it as a string.
//
// moduleName: The name of the module containing the function.
// functionName: The name of the function to retrieve documentation for.
// return: A string containing the documentation for the specified function.
func (m *ModuleDocs) GetFunctionDocAsString(moduleName string, functionName string) string {
	return m.Functions[moduleName].Functions[functionName].String()
}

// newModuleDocs initializes and returns a new ModuleDocs instance with an empty Functions map.
//
// return: A new ModuleDocs instance with an initialized Functions map.
func newModuleDocs() ModuleDocs {
	return ModuleDocs{Functions: make(map[string]FunctionDocumentationMap)}
}
