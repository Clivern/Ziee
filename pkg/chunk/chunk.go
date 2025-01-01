// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package chunk

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/tmc/langchaingo/textsplitter"
)

// Splitter identifies the text splitting strategy.
type Splitter string

const (
	SplitterRecursive Splitter = "recursive"
	SplitterMarkdown  Splitter = "markdown"
)

// Config holds text chunking settings provided by the caller.
type Config struct {
	Splitter Splitter
	Size     int
	Overlap  int

	Recursive RecursiveConfig
	Markdown  MarkdownConfig
}

// RecursiveConfig holds recursive character splitter settings.
type RecursiveConfig struct {
	Separators    []string
	KeepSeparator bool
}

// MarkdownConfig holds markdown splitter settings.
type MarkdownConfig struct {
	CodeBlocks           bool
	ReferenceLinks       bool
	KeepHeadingHierarchy bool
	JoinTableRows        bool
}

// Chunker splits document text into chunks for RAG indexing.
type Chunker struct {
	splitter textsplitter.TextSplitter
	config   Config
}

// New returns a chunker for the given configuration.
func New(config Config) (*Chunker, error) {
	splitter, err := newSplitter(config)
	if err != nil {
		return nil, err
	}

	return &Chunker{
		splitter: splitter,
		config:   config,
	}, nil
}

// Config returns the chunking configuration used by this chunker.
func (c *Chunker) Config() Config {
	return c.config
}

// Split divides text into chunks using the configured splitter.
func (c *Chunker) Split(text string) ([]string, error) {
	chunks, err := c.splitter.SplitText(text)
	if err != nil {
		return nil, fmt.Errorf("chunk split text: %w", err)
	}

	return chunks, nil
}

// newSplitter creates a text splitter for the given strategy.
func newSplitter(config Config) (textsplitter.TextSplitter, error) {
	switch splitterType(config.Splitter) {
	case SplitterRecursive:
		return newRecursiveSplitter(config), nil
	case SplitterMarkdown:
		return newMarkdownSplitter(config), nil
	default:
		return nil, fmt.Errorf("chunk unsupported splitter %q", config.Splitter)
	}
}

// splitterType returns the splitter type for a document format.
func splitterType(splitter Splitter) Splitter {
	if lo.IsEmpty(splitter) {
		return SplitterRecursive
	}

	return splitter
}

// newRecursiveSplitter creates a recursive character text splitter.
func newRecursiveSplitter(config Config) textsplitter.TextSplitter {
	opts := []textsplitter.Option{
		textsplitter.WithChunkSize(config.Size),
		textsplitter.WithChunkOverlap(config.Overlap),
	}

	if len(config.Recursive.Separators) > 0 {
		opts = append(opts, textsplitter.WithSeparators(config.Recursive.Separators))
	}

	if config.Recursive.KeepSeparator {
		opts = append(opts, textsplitter.WithKeepSeparator(true))
	}

	return textsplitter.NewRecursiveCharacter(opts...)
}

// newMarkdownSplitter creates a markdown-aware text splitter.
func newMarkdownSplitter(config Config) textsplitter.TextSplitter {
	opts := []textsplitter.Option{
		textsplitter.WithChunkSize(config.Size),
		textsplitter.WithChunkOverlap(config.Overlap),
		textsplitter.WithSecondSplitter(newRecursiveSplitter(config)),
	}

	if config.Markdown.CodeBlocks {
		opts = append(opts, textsplitter.WithCodeBlocks(true))
	}

	if config.Markdown.ReferenceLinks {
		opts = append(opts, textsplitter.WithReferenceLinks(true))
	}

	if config.Markdown.KeepHeadingHierarchy {
		opts = append(opts, textsplitter.WithHeadingHierarchy(true))
	}

	if config.Markdown.JoinTableRows {
		opts = append(opts, textsplitter.WithJoinTableRows(true))
	}

	return textsplitter.NewMarkdownTextSplitter(opts...)
}
