package document_manager

import (
	"KamaiZen/settings"
	"fmt"
	"github.com/rs/zerolog/log"
	"iter"
	"maps"
	"os"
	"regexp"
	"strings"
)

var moduleDocumentationMapInstance = &moduleDocumentationMap{
	ModuleDocs: make(map[string]ModuleDocs),
}

const (
	_EXAMPLE_START           = "Example"
	_EXAMPLE_BLOCK_SPECIFIER = "..."
)

const (
	_MODULES_PATH = "/src/modules"
	_READEME_FILE = "/README"
)

const (
	_EXAMPLE_BLOCK_SPECIFIER_COUNT = 2
)

var (
	_TOC_REGX_PATTERN  *regexp.Regexp = regexp.MustCompile(`^\s*\d+\.\s*(\w+)\s*$`)
	_FUNC_REGX_PATTERN *regexp.Regexp = regexp.MustCompile(`^\s*\d+\.\d+\.\s*(\w+)\((.*)\)\s*$`)
)

// brief Parses the provided README content and prints the table of contents.
// It identifies section headers from lines that match a specific pattern.
//
// The function expects the README content to follow a specific format:
// - Section headers start with a line matching the pattern `TOC_REGX_PATTERN`.
//
// readme: A string containing the content of the README file.
func listTableOfContents(readme string) {
	lines := strings.Split(readme, "\n")
	for _, line := range lines {
		if match := _TOC_REGX_PATTERN.FindStringSubmatch(line); match != nil {
			fmt.Println(match[1])
		}
	}
}

// Parses the provided README content and extracts a list of unique function names.
// It identifies function names from lines that match a specific pattern.
//
// The function expects the README content to follow a specific format:
// - Function documentation starts with a line matching the pattern `FUNC_REGX_PATTERN`.
//
// readme: A string containing the content of the README file.
// return: A slice of strings containing the unique function names found in the README content.
func listFunctions(readme string) []string {
	lines := strings.Split(readme, "\n")
	listofFunctions := []string{}
	for _, line := range lines {
		if match := _FUNC_REGX_PATTERN.FindStringSubmatch(line); match != nil {
			// don't add the function if it is a duplicate
			if !strings.Contains(strings.Join(listofFunctions, ","), match[1]) {
				listofFunctions = append(listofFunctions, match[1])
			}
		}
	}
	return listofFunctions
}

// Parses the provided README content and extracts the documentation
// for a specific function by its name. It identifies the function's name and parameters
// from the README content.
//
// The function expects the README content to follow a specific format:
// - Function documentation starts with a line matching the pattern `FUNC_REGX_PATTERN`.
//
// readme: A string containing the content of the README file.
// functionName: The name of the function to extract documentation for.
// return: A FunctionDocumentation struct containing the name and parameters of the specified function.
func getFunction(readme string, functionName string) FunctionDocumentation {
	lines := strings.Split(readme, "\n")
	var functionDoc FunctionDocumentation
	for _, line := range lines {
		if match := _FUNC_REGX_PATTERN.FindStringSubmatch(line); match != nil {
			if match[1] == functionName {
				functionDoc.Name = match[1]
				functionDoc.Parameters = match[2]
			}
		}
	}

	return functionDoc
}

// Parses a slice of strings representing lines of documentation
// and extracts function documentation details. It identifies function names, parameters,
// descriptions, and examples from the provided lines.
//
// The function expects the documentation to follow a specific format:
// - Function documentation starts with a line matching the pattern `FUNC_REGEX_PATTERN`.
// - Descriptions are lines following the function declaration until an "Example" line is encountered.
// - Examples are lines following the "Example" line. They are contained within ... lines.
//
// lines: A slice of strings where each string is a line of documentation.
// return: A slice of FunctionDocumentation structs containing the parsed documentation details.
func extractFunctionDoc(lines []string) []FunctionDocumentation {
	var functionDocs []FunctionDocumentation
	var functionDoc FunctionDocumentation
	var inExample bool
	var example string
	var exampleLineCount int
	for _, line := range lines {
		if match := _FUNC_REGX_PATTERN.FindStringSubmatch(line); match != nil {
			if functionDoc.Name != "" {
				functionDocs = append(functionDocs, functionDoc)
			}
			functionDoc = FunctionDocumentation{Name: match[1], Parameters: match[2]}
			inExample = false
			example = ""
		} else if functionDoc.Name != "" {
			if !inExample {
				// example has indentation of some spaces
				if strings.Contains(line, _EXAMPLE_START) {
					inExample = true
				} else {
					functionDoc.Description += line + "\n"
				}
			} else {
				if strings.Contains(line, _EXAMPLE_BLOCK_SPECIFIER) {
					exampleLineCount++
					if exampleLineCount == _EXAMPLE_BLOCK_SPECIFIER_COUNT {
						// marks end of example block
						// end of example if there are 2 lines
						inExample = false
					}
					functionDoc.Example = example
				} else {
					example += line + "\n"
				}
			}
		}
	}
	if functionDoc.Name != "" {
		functionDocs = append(functionDocs, functionDoc)
	}
	return functionDocs
}

// Initializes the document manager by reading the README files from the specified
// Kamailio source path and extracting function documentation from them. It then adds the
// extracted documentation to the module documentation map.
//
// The function expects the settings to provide a valid Kamailio source path.
//
// s: An instance of settings.LSPSettings containing the configuration settings.
//
// The function performs the following steps:
// 1. Reads the directory specified by the Kamailio source path.
// 2. Iterates over each module directory found in the path.
// 3. Reads the README file from each module directory.
// 4. Extracts function documentation from the README file.
// 5. Adds the extracted function documentation to the function documentation map.
// 6. Adds the function documentation map to the module documentation map.
//
// return: An error if there was an issue reading the directory or file.
func Initialise(s settings.LSPSettings) error {
	path := s.KamailioSourcePath + _MODULES_PATH
	listOfModules, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	// Get All Modules
	for _, module := range listOfModules {
		readme, err := os.ReadFile(path + "/" + module.Name() + _READEME_FILE)
		if err != nil {
			log.Error().Err(err)
			continue
		}
		functionDocs := extractFunctionDoc(strings.Split(string(readme), "\n"))
		functionDocsMap := FunctionDocumentationMap{Functions: make(map[string]FunctionDocumentation)}
		for _, functionDoc := range functionDocs {
			// we are overwriting the function documentation if it already exists
			err = functionDocsMap.AddFunctionDoc(functionDoc, true)
			if err != nil {
				log.Error().Str("function", functionDoc.Name).Msg("Error Adding function documentation...skipping")
			}
		}
		moduleDocs := newModuleDocs()
		err = moduleDocs.AddFunctionDoc(module.Name(), functionDocsMap, true)
		if err != nil {
			log.Error().Str("module", module.Name()).Msg("Error Adding function documentation...skipping")
		}
		moduleDocumentationMapInstance.AddModuleDocs(module.Name(), moduleDocs, true)
	}
	return nil
}

// Retrieves the documentation for a specific function within a specified module.
// It looks up the module documentation map to find the module and then retrieves the function
// documentation as a string.
//
// moduleName: The name of the module containing the function.
// functionName: The name of the function to retrieve documentation for.
// return: A string containing the documentation for the specified function. If the module is not found,
//
//	it returns "Module not found".
func GetFunctionDoc(moduleName string, functionName string) string {
	moduleDocs, exists := moduleDocumentationMapInstance.GetModuleDocs(moduleName)
	if !exists {
		return "Module not found"
	}
	return moduleDocs.GetFunctionDocAsString(moduleName, functionName)
}

// Searches for a specific function across all modules and retrieves its documentation.
// It iterates through the module documentation map to find the module containing the function and then
// returns the function documentation as a formatted string.
//
// functionName: The name of the function to search for.
// return: A string containing the documentation for the specified function, including the module name.
//
//	If the function is not found in any module, it returns "Function not found".
func FindFunctionInAllModules(functionName string) string {
	for moduleName, moduleDocs := range moduleDocumentationMapInstance.ModuleDocs {
		if _, exists := moduleDocs.Functions[moduleName].Functions[functionName]; exists {
			return "# Module: " + moduleName + "\n\n" + moduleDocs.GetFunctionDocAsString(moduleName, functionName)
		}
	}
	return "Function not found"
}

// GetAllAvailableModules retrieves the names of all available modules
// from the module documentation map.
//
// return: A slice of strings containing the names of all available modules.
func GetAllAvailableModules() iter.Seq[string] {
	return maps.Keys(moduleDocumentationMapInstance.ModuleDocs)
}

// GetAllFunctionsInModule retrieves all function documentation for a specific module.
//
// moduleName: The name of the module to retrieve function documentation for.
// return: A FunctionDocumentationMap containing the documentation for all functions
//
//	in the specified module. If the module is not found, it returns an empty FunctionDocumentationMap.
func GetAllFunctionsInModule(moduleName string) FunctionDocumentationMap {
	moduleDocs, exists := moduleDocumentationMapInstance.GetModuleDocs(moduleName)
	if !exists {
		return FunctionDocumentationMap{}
	}
	return moduleDocs.Functions[moduleName]
}

// GetAllAvailableFunctionDocs retrieves all available function documentation
// from the module documentation map.
//
// return: A slice of FunctionDocumentation structs containing the documentation
//
//	for all functions across all modules.
func GetAllAvailableFunctionDocs() []FunctionDocumentation {
	var functionDocs []FunctionDocumentation
	for _, moduleDocs := range moduleDocumentationMapInstance.ModuleDocs {
		for _, functionDoc := range moduleDocs.Functions {
			for _, doc := range functionDoc.Functions {
				functionDocs = append(functionDocs, doc)
			}
		}
	}
	return functionDocs
}
