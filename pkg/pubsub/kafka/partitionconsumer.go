package kafka

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/twmb/franz-go/pkg/kgo"

	sdkkafka "github.com/scribd/go-sdk/pkg/instrumentation/kafka"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

type pconsumer struct {
	pool *pool

	quit chan struct{}
	done chan struct{}
	recs chan *sdkkafka.FetchPartition
}

func (pc *pconsumer) consume(cl *kgo.Client, logger sdklogger.Logger, shouldCommit bool, handler func(*kgo.Record)) {
	defer func() {
		if rec := recover(); rec != nil {
			sentry.CurrentHub().Recover(rec)
			sentry.Flush(time.Second * 5)
			logger.Fatalf("kafka consumer: panic error: %v", rec)
		}
	}()

	defer close(pc.done)

	for {
		select {
		case <-pc.quit:
			return
		case p := <-pc.recs:
			p.EachRecord(func(rec *kgo.Record) {
				pc.pool.Schedule(func() {
					defer p.ConsumeRecord(rec)

					handler(rec)
				})
			})
			if shouldCommit {
				if err := cl.CommitRecords(context.Background(), p.Records...); err != nil {
					logger.WithError(err).Errorf("Partition consumer failed to commit records")
				}
			}
		}
	}
}
