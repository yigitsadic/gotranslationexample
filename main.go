package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func loadLocales(b *i18n.Bundle, lang string) *i18n.Localizer {
	b.MustLoadMessageFile(fmt.Sprintf("./locales/%s.json", lang))

	return i18n.NewLocalizer(b, lang)
}

func runMessageTicker(message *string) {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Println(*message)
		}
	}()
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln("Unable to run!")
	}
	defer watcher.Close()

	done := make(chan bool)

	time.AfterFunc(time.Minute, func() {
		done <- true
	})

	var lang string

	flag.StringVar(&lang, "lang", "en", "language to run")
	flag.Parse()

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	localizer := loadLocales(bundle, lang)

	greetingMessage := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "WelcomeMessage"})

	go func() {
		for {
			select {
			case <-watcher.Events:
				localizer = loadLocales(bundle, lang)
				greetingMessage = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "WelcomeMessage"})
			case watcherErr := <-watcher.Errors:
				fmt.Println("Error ocurred: ", watcherErr)
			}
		}
	}()

	if err = watcher.Add(fmt.Sprintf("locales/%s.json", lang)); err != nil {
		fmt.Printf("Error watching %s locale: %s\n", lang, err)
	}

	runMessageTicker(&greetingMessage)

	<-done
}
