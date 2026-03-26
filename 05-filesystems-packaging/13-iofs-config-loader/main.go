package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

func LoadConfigs(fsys fs.FS, root string) (map[string][]byte, error) {
	if root == "" {
		return nil, fmt.Errorf("root cannot be empty")
	}

	cleanRoot := path.Clean(root)
	if cleanRoot == "/" || !fs.ValidPath(cleanRoot) {
		return nil, fmt.Errorf("invalid root path: %q", root)
	}

	configs := make(map[string][]byte)

	err := fs.WalkDir(fsys, cleanRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(p, ".md") {
			return nil
		}

		content, err := fs.ReadFile(fsys, p)
		if err != nil {
			return fmt.Errorf("read %q: %w", p, err)
		}

		configs[p] = content
		return nil
	})
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func main() {
	configs, err := LoadConfigs(os.DirFS("C:/Users/abdel/OneDrive/Desktop/Projects/go-kata"), ".")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for name, content := range configs {
		fmt.Printf("%s:\n%s\n\n", name, content)
	}
}
