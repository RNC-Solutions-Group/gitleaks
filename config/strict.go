package config

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// ValidateStrictTOML rejects unknown TOML fields, incorrect table nesting and
// misplaced allowlists. It does not return field values, so configuration
// contents cannot leak through diagnostics.
func ValidateStrictTOML(raw []byte) error {
	var root map[string]any
	if err := toml.Unmarshal(raw, &root); err != nil {
		return fmt.Errorf("malformed TOML: %w", err)
	}
	return validateStrictRoot(root, "")
}

func ValidateStrictFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read configuration %q: %w", path, err)
	}
	if err := ValidateStrictTOML(raw); err != nil {
		return fmt.Errorf("strict configuration %q: %w", path, err)
	}
	return nil
}

var (
	strictRootKeys      = keySet("title", "description", "extend", "auto-overwrite", "aggregation", "informative-report", "console-only", "rules", "allowlist", "allowlists", "minVersion")
	strictExtendKeys    = keySet("path", "url", "useDefault", "disabledRules")
	strictRuleKeys      = keySet("id", "description", "path", "regex", "secretGroup", "entropy", "keywords", "tags", "allowlist", "allowlists", "required", "skipReport")
	strictAllowlistKeys = keySet("description", "condition", "commits", "paths", "regexTarget", "regexes", "stopWords", "targetRules")
	strictRequiredKeys  = keySet("id", "withinLines", "withinColumns")
)

func keySet(keys ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		m[key] = struct{}{}
	}
	return m
}

func validateStrictRoot(root map[string]any, prefix string) error {
	if err := rejectUnknown(root, strictRootKeys, "top-level"); err != nil {
		return err
	}
	if value, ok := root["extend"]; ok {
		table, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("extend must be a table")
		}
		if err := rejectUnknown(table, strictExtendKeys, "extend"); err != nil {
			return err
		}
	}
	if value, ok := root["allowlist"]; ok {
		table, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("allowlist must be a table")
		}
		if err := rejectUnknown(table, strictAllowlistKeys, "allowlist"); err != nil {
			return err
		}
	}
	if value, ok := root["allowlists"]; ok {
		tables, ok := value.([]any)
		if !ok {
			return fmt.Errorf("allowlists must be an array of tables")
		}
		for i, item := range tables {
			table, ok := item.(map[string]any)
			if !ok {
				return fmt.Errorf("allowlists[%d] must be a table", i)
			}
			if err := rejectUnknown(table, strictAllowlistKeys, fmt.Sprintf("allowlists[%d]", i)); err != nil {
				return err
			}
		}
	}
	if value, ok := root["rules"]; ok {
		rules, ok := value.([]any)
		if !ok {
			return fmt.Errorf("rules must be an array of tables")
		}
		for i, item := range rules {
			rule, ok := item.(map[string]any)
			if !ok {
				return fmt.Errorf("rules[%d] must be a table", i)
			}
			if err := validateStrictRule(rule, fmt.Sprintf("rules[%d]", i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateStrictRule(rule map[string]any, label string) error {
	if err := rejectUnknown(rule, strictRuleKeys, label); err != nil {
		return err
	}
	if value, ok := rule["allowlist"]; ok {
		table, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("%s.allowlist must be a table", label)
		}
		if err := rejectUnknown(table, strictAllowlistKeys, label+".allowlist"); err != nil {
			return err
		}
	}
	if value, ok := rule["allowlists"]; ok {
		tables, ok := value.([]any)
		if !ok {
			return fmt.Errorf("%s.allowlists must be an array of tables", label)
		}
		for i, item := range tables {
			table, ok := item.(map[string]any)
			if !ok {
				return fmt.Errorf("%s.allowlists[%d] must be a table", label, i)
			}
			if err := rejectUnknown(table, strictAllowlistKeys, fmt.Sprintf("%s.allowlists[%d]", label, i)); err != nil {
				return err
			}
		}
	}
	if value, ok := rule["required"]; ok {
		tables, ok := value.([]any)
		if !ok {
			return fmt.Errorf("%s.required must be an array of tables", label)
		}
		for i, item := range tables {
			table, ok := item.(map[string]any)
			if !ok {
				return fmt.Errorf("%s.required[%d] must be a table", label, i)
			}
			if err := rejectUnknown(table, strictRequiredKeys, fmt.Sprintf("%s.required[%d]", label, i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func rejectUnknown(table map[string]any, allowed map[string]struct{}, label string) error {
	unknown := make([]string, 0)
	for key := range table {
		if _, ok := allowed[key]; !ok {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) == 0 {
		return nil
	}
	sort.Strings(unknown)
	return fmt.Errorf("%s: unsupported field(s): %s", label, strings.Join(unknown, ", "))
}
