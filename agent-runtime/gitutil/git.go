package gitutil

import (
	"os"

	"agent-runtime/llm"
)

func ApplyChanges(files []llm.FileChange) error {

	for _, f := range files {
		err := os.WriteFile(f.Path, []byte(f.Content), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
