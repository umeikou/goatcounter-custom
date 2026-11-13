package cron

import (
	"context"

	"zgo.at/errors"
	"zgo.at/goatcounter/v2"
	"zgo.at/zdb"
)

func updateClickCounts(ctx context.Context, hits []goatcounter.Hit) error {
	type groupKey struct {
		hour   string
		pathID goatcounter.PathID
	}
	type counts struct {
		total  int
		unique int
	}

	grouped := make(map[groupKey]counts)
	for _, h := range hits {
		if h.Bot > 0 || !h.Event {
			continue
		}
		if _, ok := goatcounter.ParseClickPath(h.Path); !ok {
			continue
		}

		key := groupKey{
			hour:   h.CreatedAt.UTC().Format("2006-01-02 15:00:00"),
			pathID: h.PathID,
		}
		value := grouped[key]
		value.total++
		if h.FirstVisit {
			value.unique++
		}
		grouped[key] = value
	}

	err := zdb.TX(ctx, func(ctx context.Context) error {
		ins, err := goatcounter.Tables.ClickCounts.Bulk(ctx)
		if err != nil {
			return err
		}

		siteID := goatcounter.MustGetSite(ctx).ID
		for key, value := range grouped {
			ins.Values(siteID, key.pathID, key.hour, value.total, value.unique)
		}
		return ins.Finish()
	})
	return errors.Wrap(err, "cron.updateClickCounts")
}
