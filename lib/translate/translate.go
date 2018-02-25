package translate

import (
	"context"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// Translate words translate
func Translate(c context.Context, words []string, lang string) ([]string, error) {
	client, err := translate.NewClient(c)
	if err != nil {
		return nil, err
	}

	target, err := language.Parse(lang)
	if err != nil {
		return nil, err
	}

	w := []string{}
	for _, v := range words {
		translations, err := client.Translate(c, []string{v}, target, nil)
		if err != nil {
			return nil, err
		}

		w = append(w, translations[0].Text)
	}

	return w, nil
}
