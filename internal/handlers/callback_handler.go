package handlers

import (
	"github.com/antonioo83/license-server/config"
	"github.com/antonioo83/license-server/internal/services/interfaces"
	"github.com/robfig/cron/v3"
)

func InitCallbackCronJob(config config.Callback, service interfaces.LicenseCallbackService) error {
	cronHandler := cron.New(cron.WithSeconds())
	_, err := cronHandler.AddFunc(config.CronSpec, func() {
		service.SendCallbacks()
	})
	if err != nil {

		return err
	}

	cronHandler.Start()

	return nil
}
