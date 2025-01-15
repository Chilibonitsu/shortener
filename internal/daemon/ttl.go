package daemon

import (
	"context"
	"short-url-api/storage"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type ExpDelete interface {
	GetAllExpired() ([]storage.Url, error)
	DeleteUrl(id int64) error
}

const (
	ttl      = 12
	timeUnit = time.Second
)

func TtlDelete(ctx context.Context, log *logrus.Logger, db ExpDelete) {
	ticker := time.NewTicker(ttl * timeUnit)
	for {
		select {
		case <-ticker.C:
			go DeleteExpired(ctx, log, db)
		case <-ctx.Done():
			return
		}

	}

}

type resultWithError struct {
	Url storage.Url
	Err error
}

func deleteExpired(wg *sync.WaitGroup, inputCh <-chan storage.Url, outputCh chan<- resultWithError, db ExpDelete) {

	defer wg.Done()

	for url := range inputCh {
		err := db.DeleteUrl(url.ID)
		outputCh <- resultWithError{
			Url: url,
			Err: err,
		}
	}

}
func DeleteExpired(ctx context.Context, log *logrus.Logger, db ExpDelete) {
	const op = "internal.api.daemon.DeleteExpired"
	inputCh := make(chan storage.Url)
	outputCh := make(chan resultWithError)

	urls, err := db.GetAllExpired()
	urlsLen := len(urls)

	wg := &sync.WaitGroup{}

	if err != nil {
		log.WithField("Cant get expired urls", err).Warn()
	}

	log.WithField(op, "Deleting expired urls...")

	go func() {
		defer close(inputCh)

		for i := range urls {
			inputCh <- urls[i]
		}
	}()

	go func() {
		for i := 0; i < urlsLen; i++ {
			wg.Add(1)

			go deleteExpired(wg, inputCh, outputCh, db)
		}
		wg.Wait()
		close(outputCh)
	}()

	output := make([]storage.Url, 0, len(urls))

	for res := range outputCh {
		if res.Err != nil {
			log.WithFields(logrus.Fields{
				"url alias": res.Url.Alias,
				"err":       res.Err,
			}).Warn("error deleting url")
		} else {
			output = append(output, res.Url)
		}
	}

	for url := range output {
		log.WithFields(logrus.Fields{
			"url alias": output[url].Alias,
			"url":       output[url].Url,
		}).Info("url deleted")
	}

}

//https://habr.com/ru/companies/tuturu/articles/755072/
//Нужно было сделать воркер пул буквально как тут?
//По одной ссылке параллельно удалять из базы как будто проблемы будут

func Del(ctx context.Context, log *logrus.Logger, db ExpDelete) {
	const op = "internal.api.daemon.Del"

	urls, err := db.GetAllExpired()
	urlsLen := len(urls)

	if err != nil {
		log.WithField("Cant get expired urls", err).Warn()
	}

	log.WithField(op, "Deleting expired urls...")

	for i := 0; i < urlsLen; i++ {
		err := db.DeleteUrl(urls[i].ID)
		if err != nil {
			log.WithFields(logrus.Fields{
				"Cant delete url": urls[i].Url,
				"err":             err,
			})
		}
	}
}
