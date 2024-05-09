package validator

import (
	"regexp"
	"strings"
	"unicode"
)

var countryCodes = []string{
	"93", "355", "213", "1-684", "376", "244", "1-264", "672", "1-268",
	"54", "374", "297", "43", "994", "1-242", "973", "880", "1-246",
	"375", "32", "501", "229", "1-441", "975", "591", "387", "267", "55",
	"246", "1-284", "673", "359", "226", "257", "855", "237", "1", "238",
	"1-345", "236", "235", "56", "86", "61", "57", "269", "682", "506",
	"385", "53", "599", "357", "420", "243", "45", "253", "1-767", "1-809",
	"1-829", "1-849", "670", "593", "20", "503", "240", "291", "372", "251",
	"500", "298", "679", "358", "33", "689", "241", "220", "995", "49", "233", "350", "30", "299", "1-473", "1-671", "502", "44-1481", "224", "245", "592", "509", "504", "852", "36", "354", "91", "62", "98", "964", "353", "44-1624", "972", "39", "225", "1-876", "81", "44-1534", "962", "7", "254", "686", "383", "965", "996", "856", "371", "961", "266", "231", "218", "423", "370", "352", "853", "389", "261", "265", "60", "960", "223", "356", "692", "222", "230", "262", "52", "691", "373", "377", "976", "382", "1-664", "212", "258", "95", "264", "674", "977", "31", "599", "687", "64", "505", "227", "234", "683", "850", "1-670", "47", "968", "92", "680", "970", "507", "675", "595", "51", "63", "64", "48", "351", "1-787", "1-939", "974", "242", "262", "40", "7", "250", "590", "290", "1-869", "1-758", "590", "508", "1-784", "685", "378", "239", "966", "221", "381", "248", "232", "65", "1-721", "421", "386", "677", "252", "27", "82", "211", "34", "94", "249", "597", "47", "268", "46", "41", "963", "886", "992", "255", "66", "228", "690", "676", "1-868", "216", "90", "993", "1-649", "688", "1-340", "256", "380", "971", "44", "1", "598", "998", "678", "379", "58", "84", "681", "212", "967", "260", "263",
}

const minPasswordLen = 5
const maxPasswordLen = 15
const phoneNumberLen = 10

func IsValidPhoneNumber(rawPhoneNumber string) bool {
	if !strings.HasPrefix(rawPhoneNumber, "+") {
		return false
	}

	offset := len(rawPhoneNumber) - phoneNumberLen
	countryCode := rawPhoneNumber[1:offset]
	foundCountryCode := false
	for _, cc := range countryCodes {
		if strings.Compare(countryCode, cc) == 0 {
			foundCountryCode = true
		}
	}

	if !foundCountryCode {
		return false
	}

	phoneNumber := rawPhoneNumber[offset:]
	return len(phoneNumber) == phoneNumberLen
}

func IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func IsValidFullName(fullname string) bool {
	// contains only letters and spaces, and must be between 3 and 40 characters long
	re := regexp.MustCompile(`^[a-zA-Z\s]{5,15}$`)
	return re.MatchString(fullname)
}

func IsSolidPassword(s string) bool {
	var (
		hasMinMaxLen = false
		hasNumber    = false
		hasLetter    = false
	)

	if len(s) >= minPasswordLen && len(s) <= maxPasswordLen {
		hasMinMaxLen = true
	}

	for _, char := range s {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	return hasMinMaxLen && (hasLetter || hasNumber)
}
