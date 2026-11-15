package report

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

const groupingDomain = "rnc-mod/report-grouping/v1\x00"

// GroupedReporter writes a separate, report-only view of repeated secret
// values. It never modifies or suppresses primary findings.
type GroupedReporter struct {
	Key []byte
}

type GroupedDocument struct {
	SchemaVersion int     `json:"schema_version"`
	Algorithm     string  `json:"algorithm"`
	Domain        string  `json:"domain"`
	Groups        []Group `json:"groups"`
}

type Group struct {
	ID          string            `json:"id"`
	Count       int               `json:"count"`
	Occurrences []GroupOccurrence `json:"occurrences"`
}

type GroupOccurrence struct {
	RuleID      string `json:"rule_id"`
	File        string `json:"file"`
	StartLine   int    `json:"start_line"`
	EndLine     int    `json:"end_line"`
	Commit      string `json:"commit,omitempty"`
	Fingerprint string `json:"fingerprint"`
}

var _ Reporter = (*GroupedReporter)(nil)

func (r *GroupedReporter) Write(w io.WriteCloser, findings []Finding) error {
	if len(r.Key) < 32 {
		return fmt.Errorf("grouping key must contain at least 32 bytes")
	}
	groups := make(map[string][]GroupOccurrence)
	for _, finding := range findings {
		secret := finding.groupingSecret
		if secret == "" {
			secret = finding.Secret
		}
		if secret == "" {
			continue
		}
		secret = strings.ReplaceAll(secret, "\r\n", "\n")
		mac := hmac.New(sha256.New, r.Key)
		_, _ = mac.Write([]byte(groupingDomain))
		_, _ = mac.Write([]byte(secret))
		id := "hmac-sha256:" + hex.EncodeToString(mac.Sum(nil))
		groups[id] = append(groups[id], GroupOccurrence{
			RuleID:      finding.RuleID,
			File:        finding.File,
			StartLine:   finding.StartLine,
			EndLine:     finding.EndLine,
			Commit:      finding.Commit,
			Fingerprint: finding.Fingerprint,
		})
	}

	document := GroupedDocument{
		SchemaVersion: 1,
		Algorithm:     "hmac-sha256",
		Domain:        "rnc-mod/report-grouping/v1",
		Groups:        make([]Group, 0, len(groups)),
	}
	for id, occurrences := range groups {
		sort.Slice(occurrences, func(i, j int) bool {
			a, b := occurrences[i], occurrences[j]
			if a.File != b.File {
				return a.File < b.File
			}
			if a.Commit != b.Commit {
				return a.Commit < b.Commit
			}
			if a.StartLine != b.StartLine {
				return a.StartLine < b.StartLine
			}
			return a.RuleID < b.RuleID
		})
		document.Groups = append(document.Groups, Group{
			ID:          id,
			Count:       len(occurrences),
			Occurrences: occurrences,
		})
	}
	sort.Slice(document.Groups, func(i, j int) bool {
		if document.Groups[i].Count != document.Groups[j].Count {
			return document.Groups[i].Count > document.Groups[j].Count
		}
		return document.Groups[i].ID < document.Groups[j].ID
	})

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	return encoder.Encode(document)
}
