package nr

import (
	"log"

	"github.com/gophr-pm/gophr/common/config"

	newrelic "github.com/newrelic/go-agent"
)

// CreateNewRelicApp lol
func CreateNewRelicApp(conf *config.Config) (newrelic.Application, error) {
	var app newrelic.Application
	if !conf.IsDev {
		log.Println("Creating New Relic Application.")
		newRelicKey, err := generateKey(conf)
		if err != nil {
			log.Fatalln("Failed to read newrelic credentials secret: ", err)
			return app, err
		}
		config := newrelic.NewConfig("Gophr", newRelicKey)
		app, err = newrelic.NewApplication(config)
		if err != nil {
			log.Fatalln("Failed to create new relic monitoring application: ", err)
			return app, err
		}
	}

	return app, nil
}
