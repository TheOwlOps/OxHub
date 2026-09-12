package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type SkillInfo struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type Manifest struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Description string      `json:"description"`
	Skills      []SkillInfo `json:"skills"`
}

func main() {
	rootDir := "hub/skills"
	var skills []SkillInfo

	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), "skill.md") {
			relPath, _ := filepath.Rel(".", path)
			relPath = filepath.ToSlash(relPath)
			parts := strings.Split(relPath, "/")

			category := "general"
			if len(parts) >= 3 {
				category = parts[2]
			}

			// Read description snippet
			content, err := os.ReadFile(path)
			desc := "Skill definition"
			if err == nil {
				lines := strings.Split(string(content), "\n")
				for _, l := range lines {
					if strings.HasPrefix(strings.TrimSpace(l), "description:") {
						desc = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(l), "description:"))
						desc = strings.Trim(desc, `"'`)
						break
					}
				}
			}

			skillName := filepath.Base(filepath.Dir(path))
			skills = append(skills, SkillInfo{
				Name:        skillName,
				Category:    category,
				Path:        relPath,
				Description: desc,
			})
		}
		return nil
	})

	m := Manifest{
		Name:        "oxhub-all-skills",
		Version:     "1.1.0",
		Description: "The complete ecosystem of skills & plugins for Claude Code, Hermes Agent, and Codex",
		Skills:      skills,
	}

	data, _ := json.MarshalIndent(m, "", "  ")
	os.WriteFile("hub/manifest.json", data, 0644)
	fmt.Printf("Generated manifest with %d skills\n", len(skills))
}
