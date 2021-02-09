package parser

import "taxes-be/internal/core"

type Parser interface {
	Parse(lines []string) (*core.Report, error)
}
