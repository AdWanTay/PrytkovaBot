package utils

import "net/url"

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
