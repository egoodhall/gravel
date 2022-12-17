package cli

import (
	"io"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

// Allows parsing of kong-defined flags from a
// YAML configuration file
func YAML(r io.Reader) (kong.Resolver, error) {
	values := map[string]interface{}{}
	err := yaml.NewDecoder(r).Decode(&values)
	if err != nil {
		return nil, err
	}

	var f kong.ResolverFunc = func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
		name := strings.ReplaceAll(flag.Name, "-", "_")
		raw, ok := values[name]
		if ok {
			return raw, nil
		}
		raw = values
		for _, part := range strings.Split(name, ".") {
			if values, ok := raw.(map[string]any); ok { // Standard JSON object
				raw, ok = values[part]
				if !ok {
					return nil, nil
				}
			} else if values, ok := raw.(map[any]any); ok { // YAML can also have map[any]any
				raw, ok = values[part]
				if !ok {
					return nil, nil
				}
			} else { // No matching key
				return nil, nil
			}
		}
		return raw, nil
	}

	return f, nil
}
