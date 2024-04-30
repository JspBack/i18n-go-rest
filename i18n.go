package main

import "github.com/nicksnyder/go-i18n/v2/i18n"

func TranslateMessage(T *i18n.Localizer, messageID string) (string, error) {
	message, err := T.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	return message, err
}
