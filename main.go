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
	b.MustLoadMessageFile("locales/en.json")
	b.MustLoadMessageFile("locales/tr.json")
	b.MustLoadMessageFile("locales/fr.json")

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

	if err = watcher.Add("./locales/tr.json"); err != nil {
		fmt.Println("Error watching Turkish locale", err)
	}

	if err = watcher.Add("./locales/en.json"); err != nil {
		fmt.Println("Error watching English locale", err)
	}

	if err = watcher.Add("./locales/fr.json"); err != nil {
		fmt.Println("Error watching French locale", err)
	}

	runMessageTicker(&greetingMessage)

	<-done
}
