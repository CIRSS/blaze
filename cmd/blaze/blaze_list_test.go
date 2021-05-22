package main

import (
	"strings"
	"testing"

	"github.com/cirss/go-cli/pkg/util"
)

func TestBlazegraphCmd_list_empty_store(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")

	Program.AssertExitCode(t, "blaze list", 0)

	util.LineContentsEqual(t, outputBuffer.String(), ``)
}

func TestBlazegraphCmd_list_default_dataset(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")
	Program.Invoke("blaze create --quiet")

	t.Run("default count", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze list", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`kb         0
			`)
	})

	t.Run("no count", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze list --count none", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`kb
			`)
	})

	t.Run("exact count", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze list --count exact", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`kb         0
			`)
	})

	t.Run("estimate count", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze list --count estimate", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`kb         0
			`)
	})
}

func TestBlazegraphCmd_list_custom_dataset(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")
	Program.Invoke("blaze create --quiet --dataset foo")

	Program.AssertExitCode(t, "blaze list", 0)

	util.LineContentsEqual(t, outputBuffer.String(),
		`foo        0
		`)
}

func TestBlazegraphCmd_list_custom_datasets(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")
	Program.Invoke("blaze create --quiet --dataset foo")
	Program.Invoke("blaze create --quiet --dataset bar")
	Program.Invoke("blaze create --quiet --dataset baz")

	t.Run("default count", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze list", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`bar        0
			baz        0
			foo        0
	   `)
	})

	t.Run("exact count", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze list --count exact", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`bar        0
			 baz        0
			 foo        0
			`)
	})

	t.Run("estimate count", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze list --count estimate", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`bar        0
			 baz        0
			 foo        0
			`)
	})
}

var expectedListHelpOutput = string(
	`blaze list: Lists the names of the RDF datasets in the Blazegraph instance.

	usage: blaze list [<flags>]

	flags:
		-count string
				Include count of triples in each dataset [none, estimate, exact] (default "exact")
		-instance URL
				URL of Blazegraph instance (default "http://127.0.0.1:9999/blazegraph")
		-quiet
				Discard normal command output
		-silent
				Discard normal and error command output
	`)

func TestBlazegraphCmd_list_help(t *testing.T) {
	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer
	Program.AssertExitCode(t, "blaze list help", 0)
	util.LineContentsEqual(t, outputBuffer.String(), expectedListHelpOutput)
}

func TestBlazegraphCmd_help_list(t *testing.T) {
	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer
	Program.AssertExitCode(t, "blaze help list", 0)
	util.LineContentsEqual(t, outputBuffer.String(), expectedListHelpOutput)
}

func TestBlazegraphCmd_list_bad_flag(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.AssertExitCode(t, "blaze list --not-a-flag", 1)

	util.LineContentsEqual(t, outputBuffer.String(),
		`blaze list: flag provided but not defined: -not-a-flag

		usage: blaze list [<flags>]

		flags:
			-count string
					Include count of triples in each dataset [none, estimate, exact] (default "exact")
			-instance URL
					URL of Blazegraph instance (default "http://127.0.0.1:9999/blazegraph")
			-quiet
					Discard normal command output
			-silent
					Discard normal and error command output
		`)
}
