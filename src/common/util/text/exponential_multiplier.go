// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package text

// NormalizeScientificMultiplier converts the exponential multiplier (10^exp)
// to a form with 'e'.
func NormalizeScientificMultiplier(s string) string {
	switch s {
	case "10^2", "10e2", "102":
		return "e2"
	case "10^3", "10e3", "103":
		return "e3"
	case "10^4", "10e4", "104":
		return "e4"
	case "10^5", "10e5", "105":
		return "e5"
	case "10^6", "10e6", "106":
		return "e6"
	case "10^7", "10e7", "107":
		return "e7"
	case "10^8", "10e8", "108":
		return "e8"
	case "10^9", "10e9", "109":
		return "e9"
	case "10^10", "10e10", "1010":
		return "e10"
	}
	return s
}
