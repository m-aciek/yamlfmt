// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package features

import (
	"bufio"
	"bytes"
	"context"
	"strings"

	"github.com/google/yamlfmt"
)

func MakeFeatureLineBreaksBetweenTopLevelBlocks(linebreakStr string) yamlfmt.Feature {
	return yamlfmt.Feature{
		Name:        "Line Breaks Between Top-Level Blocks",
		AfterAction: lineBreaksBetweenTopLevelBlocksFeature(linebreakStr),
	}
}

func lineBreaksBetweenTopLevelBlocksFeature(linebreakStr string) yamlfmt.FeatureFunc {
	return func(_ context.Context, content []byte) (context.Context, []byte, error) {
		var buf bytes.Buffer
		scanner := bufio.NewScanner(bytes.NewReader(content))
		var seenTopLevelBlock bool
		var prevLineBlank bool

		for scanner.Scan() {
			txt := scanner.Text()
			if isTopLevelBlockStart(txt) {
				if seenTopLevelBlock && !prevLineBlank {
					buf.WriteString(linebreakStr)
				}
				seenTopLevelBlock = true
			}

			if strings.TrimSpace(txt) == "---" || strings.TrimSpace(txt) == "..." {
				seenTopLevelBlock = false
			}

			buf.WriteString(txt)
			buf.WriteString(linebreakStr)
			prevLineBlank = strings.TrimSpace(txt) == ""
		}

		return nil, buf.Bytes(), scanner.Err()
	}
}

func isTopLevelBlockStart(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	if line[0] == ' ' || line[0] == '\t' {
		return false
	}
	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "%") {
		return false
	}
	if trimmed == "---" || trimmed == "..." {
		return false
	}
	return strings.Contains(trimmed, ":")
}
