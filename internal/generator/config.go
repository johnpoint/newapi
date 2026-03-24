package generator

import "gopkg.in/yaml.v2"

// ApiDesc is the top-level structure of the YAML config file used by `newapi gen`.
type ApiDesc struct {
	Group []ApiGroup `yaml:"group"`
}

// ApiGroup groups a set of APIs under a common Swagger tag.
type ApiGroup struct {
	Name string `yaml:"name"`
	Apis []Api  `yaml:"apis"`
}

// Api describes a single API endpoint in the YAML config file.
type Api struct {
	ApiName string `yaml:"name"`
	Path    string `yaml:"path"`
	Desc    string `yaml:"desc"`
	Method  string `yaml:"method"`
}

// ParseConfig parses a YAML config file into an ApiDesc.
func ParseConfig(data []byte) (ApiDesc, error) {
	var d ApiDesc
	if err := yaml.Unmarshal(data, &d); err != nil {
		return ApiDesc{}, err
	}
	return d, nil
}
