package cnb

import (
	"encoding/json"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

const buildMetadataLabel = "io.buildpacks.build.metadata"

type BuildMetadata struct {
	Processes  []Process   `toml:"processes" json:"processes"`
	Buildpacks []Buildpack `toml:"buildpacks" json:"buildpacks"`
	BOM        []BOMEntry  `toml:"bom" json:"bom"`
}

type Buildpack struct {
	ID      string `toml:"id" json:"id"`
	Version string `toml:"version" json:"version"`
}

type BOMEntry struct {
	Name      string                 `toml:"name" json:"name"`
	Version   string                 `toml:"version" json:"version"`
	Metadata  map[string]interface{} `toml:"metadata" json:"metadata"`
	Buildpack Buildpack
}

type Process struct {
	Type    string   `toml:"type" json:"type"`
	Command string   `toml:"command" json:"command"`
	Args    []string `toml:"args" json:"args"`
	Direct  bool     `toml:"direct" json:"direct"`
}

func ParseBuildMetadata(img v1.Image) (BuildMetadata, error) {
	cfg, err := img.ConfigFile()
	if err != nil {
		return BuildMetadata{}, err
	}
	label, ok := cfg.Config.Labels[buildMetadataLabel]
	if !ok {
		return BuildMetadata{}, nil
	}
	var md BuildMetadata
	if err := json.Unmarshal([]byte(label), &md); err != nil {
		return BuildMetadata{}, err
	}
	return md, nil
}
