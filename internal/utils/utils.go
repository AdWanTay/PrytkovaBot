package utils

import (
	"net/url"
	"strings"
)

func GetWhatsAppString() string {
	whatsAppUrl, err := url.Parse("https://wa.me/79659413788")
	if err != nil {
		panic("boom")
	}
	params := url.Values{}
	params.Add("text", "Здравствуйте, хочу вам на программу.")
	whatsAppUrl.RawQuery = params.Encode()
	return whatsAppUrl.String()
}

func PrepareForMarkdown(s string) string {
	chars := []string{".", "(", ")"}
	for _, char := range chars {
		s = strings.ReplaceAll(s, char, "\\"+char)
	}
	return s
}
