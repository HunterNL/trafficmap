package description

import (
	"strconv"
	"strings"
	"unicode"
)

type DescriptionDerivatives struct {
	Organization string
	RoadId       string
	RoadOffset   int
	RoadSide     string
	Name         string
}

var orgMap = map[string]string{
	"PZH": "Provincie Zuid-Holland",
	"GDH": "Gemeente Den Haag",
}

var blackList = map[string]bool{
	"Links":  true,
	"Rechts": true,
}

var sideLookup = map[string]string{
	"links":  "L",
	"li":     "L",
	"l":      "L",
	"rechts": "R",
	"re":     "R",
	"r":      "R",
}

var roadPrefixes = map[rune]bool{
	'A': true,
	'N': true,
	'S': true,
}

func sideSuffix(str string) (string, bool) {
	if strings.HasSuffix(str, "L") {
		return "L", true
	}
	if strings.HasSuffix(str, "R") {
		return "R", true
	}
	return "", false
}

func isRoadId(str string) bool {
	for i, r := range str {
		// First rune must be a road type prefix
		if i == 0 {
			if _, found := roadPrefixes[r]; !found {
				return false
			}
		} else {
			// Allow final rune to be a roadSide suffix
			if i == len(str)-1 {
				if _, hasSuffix := sideSuffix(str); hasSuffix {
					return true
				}
			}
			if !unicode.IsDigit(r) {
				return false
			}
		}
	}
	return true
}

func parseRoadOffset(str string) (int, bool) {
	str = strings.ReplaceAll(strings.TrimRight(str, "km"), ",", ".")
	num, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return -1, false
	}

	return int(num * 1000), true
}

// Takes a string, returning bits of road data it can find
// Reads the string from left to right, outputing remaining "uninteresting" string as `description`
func parseRoadData(in string) (roadId string, roadOffset int, roadSide string, description string) {
	roadOffset = -1
	description = in
	remainingBits := strings.FieldsFunc(in, func(r rune) bool {
		return unicode.IsSpace(r) || r == '-' || r == '_'
	})
	bitCount := len(remainingBits)

	// FieldFunc failed, return input
	if bitCount == 0 {
		description = in
		return
	}

	for i, field := range remainingBits {
		if roadId == "" && isRoadId(field) {
			roadId = field

			if side, hasSuffix := sideSuffix(roadId); hasSuffix {
				roadSide = side
				roadId = roadId[:len(roadId)-1]
			}
			continue
		}

		if roadSide == "" {
			if side, isSide := sideLookup[strings.ToLower(field)]; isSide {
				roadSide = side
				continue
			}
		}

		if roadOffset == -1 {
			offset, isOffset := parseRoadOffset(field)
			if isOffset {
				roadOffset = offset
				continue
			}
		}

		if field == "km" {
			continue
		}

		// Did nothing, return the rest joined together
		description = strings.Join(remainingBits[i:], " ")
		return
	}
	description = ""

	// if isRoadId(remainingBits[0]) {
	// 	roadId = remainingBits[0]
	// 	description = strings.Join(remainingBits[1:], " ")

	// 	if sideSuffix(roadId) {
	// 		roadId = roadId[:len(roadId-2)]
	// 	}

	// 	if bitCount == 1 {
	// 		return
	// 	}

	// 	var isSide bool
	// 	if roadSide, isSide = sideLookup[strings.ToLower(remainingBits[1])]; isSide {
	// 		description = strings.Join(remainingBits[2:], " ")

	// 		if bitCount == 2 {
	// 			return
	// 		}

	// 		var isOffset bool
	// 		if roadOffset, isOffset = parseRoadOffset(remainingBits[2]); isOffset {
	// 			description = strings.Join(remainingBits[3:], " ")
	// 			return
	// 		}

	// 		if bitCount >= 3 && remainingBits[2] == "km" {
	// 			if offset, isOffset := parseRoadOffset(remainingBits[3]); isOffset {
	// 				roadOffset = offset
	// 				description = strings.Join(remainingBits[4:], " ")
	// 			}
	// 		}
	// 	}
	// }

	return
}

// // Takes a string, returning bits of road data it can find
// // Given string should start with the road id and be seperated by whitespace or '-'
// func parseRoadData(in string) (roadId string, roadOffset int, roadSide string, description string) {
// 	roadOffset = -1
// 	description = in
// 	remainingBits := strings.FieldsFunc(in, func(r rune) bool {
// 		return unicode.IsSpace(r) || r == '-'
// 	})
// 	bitCount := len(remainingBits)

// 	// FieldFunc failed, return input
// 	if bitCount == 0 {
// 		description = in
// 		return
// 	}

// 	for i, field := range remainingBits {
// 		if isRoadId(field) {
// 			roadId = field

// 			if side, hasSuffis := sideSuffix(roadId); hasSuffis {
// 				roadSide = side
// 			}
// 			continue
// 		}
// 	}

// 	if isRoadId(remainingBits[0]) {
// 		roadId = remainingBits[0]
// 		description = strings.Join(remainingBits[1:], " ")

// 		if sideSuffix(roadId) {
// 			roadId = roadId[:len(roadId-2)]
// 		}

// 		if bitCount == 1 {
// 			return
// 		}

// 		var isSide bool
// 		if roadSide, isSide = sideLookup[strings.ToLower(remainingBits[1])]; isSide {
// 			description = strings.Join(remainingBits[2:], " ")

// 			if bitCount == 2 {
// 				return
// 			}

// 			var isOffset bool
// 			if roadOffset, isOffset = parseRoadOffset(remainingBits[2]); isOffset {
// 				description = strings.Join(remainingBits[3:], " ")
// 				return
// 			}

// 			if bitCount >= 3 && remainingBits[2] == "km" {
// 				if offset, isOffset := parseRoadOffset(remainingBits[3]); isOffset {
// 					roadOffset = offset
// 					description = strings.Join(remainingBits[4:], " ")
// 				}
// 			}
// 		}
// 	}

// 	return
// }

func Parse(description string) DescriptionDerivatives {
	description = strings.TrimSpace(description)

	// Remove parentheses and everything contained, starting from the right
	var inParens = false
	description = strings.TrimRightFunc(description, func(r rune) bool {
		if r == ')' {
			inParens = true
			return true
		}

		if r == '(' {
			inParens = false
			return true
		}

		if inParens {
			return true
		}

		return false
	})

	out := DescriptionDerivatives{}

	// Remove and save organization prefix
	description = strings.TrimSpace(description)
	if strings.Contains(description, " - ") {
		if left, right, found := strings.Cut(description, " - "); found {
			description = right

			for prefix, fullName := range orgMap {
				if (left[:3]) == prefix {
					out.Organization = fullName
				}
			}
		}
	}

	out.RoadId, out.RoadOffset, out.RoadSide, description = parseRoadData(description)

	// Remove any leftover meaningless words
	if blackList[description] {
		description = ""
	}

	out.Name = description

	return out
}
