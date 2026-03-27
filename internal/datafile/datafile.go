package datafile

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bfv/xref/internal/models"
)

const DefaultDataFile = "xref.json"

// Load reads and parses a JSON file containing xref data.
func Load(path string) ([]*models.XrefFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read data file '%s': %w", path, err)
	}

	var xreffiles []*models.XrefFile
	if err := json.Unmarshal(data, &xreffiles); err != nil {
		return nil, fmt.Errorf("cannot parse data file '%s': %w", path, err)
	}
	return xreffiles, nil
}
