package translate

import (
	"cloud.google.com/go/translate"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
)

// Translate words translate
func Translate(words []string, lang string) ([]string, error) {
	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	target, err := language.Parse(lang)
	if err != nil {
		return nil, err
	}

	w := []string{}
	for _, v := range words {
		translations, err := client.Translate(ctx, []string{v}, target, nil)
		if err != nil {
			return nil, err
		}

		w = append(w, translations[0].Text)
	}

	return w, nil
}
