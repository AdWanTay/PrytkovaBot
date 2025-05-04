package utils

import (
	"net/url"
	"strings"
)

func GetEncodedString(messengerUrl, text string) string {
	appUrl, err := url.Parse(messengerUrl)
	if err != nil {
		panic("boom")
	}
	params := url.Values{}
	params.Add("text", text)
	encodedString := strings.ReplaceAll(params.Encode(), "+", "%20")
	appUrl.RawQuery = encodedString
	return appUrl.String()
}

func PrepareForMarkdown(s string) string {
	chars := []string{".", "(", ")"}
	for _, char := range chars {
		s = strings.ReplaceAll(s, char, "\\"+char)
	}
	return s
}
