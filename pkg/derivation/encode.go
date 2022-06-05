package derivation

import (
	"bytes"
	"io"
	"sort"
)

// Adds quotation marks around a string.
// This is primarily meant for non-user provided strings.
func quoteString(s string) []byte {
	buf := make([]byte, len(s)+2)

	buf[0] = '"'

	for i := 0; i < len(s); i++ {
		buf[i+1] = s[i]
	}

	buf[len(s)+1] = '"'

	return buf
}

// Convert a slice of strings to a slice of byte slices.
func stringsToBytes(elems []string) [][]byte {
	b := make([][]byte, len(elems))

	for i, s := range elems {
		b[i] = []byte(s)
	}

	return b
}

// Encode a list of elements staring with `opening` character and ending with a `closing` character.
func encodeArray(opening byte, closing byte, quote bool, elems ...[]byte) []byte {
	if len(elems) == 0 {
		return []byte{opening, closing}
	}

	n := 3 * (len(elems) - 1)
	if quote {
		n += 2
	}

	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var buf bytes.Buffer

	buf.Grow(n)
	buf.WriteByte(opening)

	writeElem := func(b []byte) {
		if quote {
			buf.WriteByte('"')
		}

		buf.Write(b)

		if quote {
			buf.WriteByte('"')
		}
	}

	writeElem(elems[0])

	for _, s := range elems[1:] {
		buf.WriteByte(',')
		writeElem(s)
	}

	buf.WriteByte(closing)

	return buf.Bytes()
}

// WriteDerivation writes the textual representation of the derivation to the passed writer.
func (d *Derivation) WriteDerivation(writer io.Writer) error {
	// we need to sort outputs by their name, which is the key of the map.
	// get the list of keys, sort them, then add each one by one.
	// Due to the "sorted paths" requirement, we know there's no two
	// outputs with the same path.
	outputNames := make([]string, len(d.Outputs))
	{
		i := 0
		for k := range d.Outputs {
			outputNames[i] = k
			i++
		}
		sort.Strings(outputNames)
	}

	encOutputs := make([][]byte, len(d.Outputs))
	{
		for i, outputName := range outputNames {
			o := d.Outputs[outputName]

			encOutputs[i] = encodeArray(
				'(', ')',
				true,
				[]byte(outputName),
				[]byte(o.Path),
				[]byte(o.HashAlgorithm),
				[]byte(o.Hash),
			)
		}
	}

	// input derivations are sorted by their path, which is the key of the map.
	// get the list of keys, sort them, then add each one by one.
	inputDerivationPaths := make([]string, len(d.InputDerivations))
	{
		i := 0
		for inputDerivationPath := range d.InputDerivations {
			inputDerivationPaths[i] = inputDerivationPath
			i++
		}
		sort.Strings(inputDerivationPaths)
	}

	encInputDerivations := make([][]byte, len(d.InputDerivations))
	{
		for i, inputDerivationPath := range inputDerivationPaths {
			names := encodeArray('[', ']', true, stringsToBytes(d.InputDerivations[inputDerivationPath])...)
			encInputDerivations[i] = encodeArray('(', ')', false, quoteString(inputDerivationPath), names)
		}
	}

	// environment variables need to be sorted by their key.
	// extract the list of keys, sort them, then add each one by one
	envKeys := make([]string, len(d.Env))
	{
		i := 0
		for k := range d.Env {
			envKeys[i] = k
			i++
		}
		sort.Strings(envKeys)
	}

	encEnv := make([][]byte, len(d.Env))
	{
		for i, k := range envKeys {
			encEnv[i] = encodeArray('(', ')', false, quoteString(k), quoteString(d.Env[k]))
		}
	}

	_, err := writer.Write([]byte("Derive"))
	if err != nil {
		return err
	}

	_, err = writer.Write(
		encodeArray('(', ')', false,
			encodeArray('[', ']', false, encOutputs...),
			encodeArray('[', ']', false, encInputDerivations...),
			encodeArray('[', ']', true, stringsToBytes(d.InputSources)...),
			quoteString(d.Platform),
			quoteString(d.Builder),
			encodeArray('[', ']', true, stringsToBytes(d.Arguments)...),
			encodeArray('[', ']', false, encEnv...),
		),
	)

	return err
}
