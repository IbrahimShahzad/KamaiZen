package document_manager

import (
	"KamaiZen/logger"
	_ "embed"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
)

//go:embed cookbooks_5.2.x.json
var cookbooksFile []byte

type DocEntry struct {
	Name          string `json:"name"`
	Documentation string `json:"documentation"`
}

type Docs struct {
	Docs []DocEntry `json:"docs"`
}

var CookBookDocs map[string]string

func readJSONFromFile() error {
	if len(cookbooksFile) == 0 {
		logger.Error("Cookbooks file is empty")
		return fmt.Errorf("cookbooks file is empty")
	}
	var docs Docs
	err := json.Unmarshal(cookbooksFile, &docs)
	if err != nil {
		logger.Error("Error reading JSON", err)
		return err
	}
	for _, doc := range docs.Docs {
		CookBookDocs[doc.Name] = doc.Documentation
	}
	return nil

}

func init() {
	logger.Debug("Initializing CookBookDocs")
	CookBookDocs = make(map[string]string)
	readJSONFromFile()
}

func GetCookBookDocs(name string) string {
	return CookBookDocs[name]
}

func GetAllCookBookKeys() iter.Seq[string] {
	return maps.Keys(CookBookDocs)
}
