package main

// Note: Be consistent with "Semantic Versioning 2.0.0" - see https://semver.org/
const version string = "1.2.0-2019.10.14"
const (
	// OutputText - output as normal text format
	OutputText uint = iota
	// OutputJSON - output as JSON, one item per line
	OutputJSON
)
