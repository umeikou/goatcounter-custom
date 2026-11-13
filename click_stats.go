package goatcounter

import (
	"cmp"
	"context"
	"slices"

	"zgo.at/errors"
	"zgo.at/zdb"
	"zgo.at/zstd/ztime"
)

type ClickLevel string

const (
	ClickLevelProject ClickLevel = "project"
	ClickLevelPage    ClickLevel = "page"
	ClickLevelType    ClickLevel = "type"
	ClickLevelElement ClickLevel = "element"
)

type ClickSelection struct {
	Project string
	Page    string
	Type    string
}

func (s ClickSelection) NextLevel() ClickLevel {
	switch {
	case s.Project == "":
		return ClickLevelProject
	case s.Page == "":
		return ClickLevelPage
	case s.Type == "":
		return ClickLevelType
	default:
		return ClickLevelElement
	}
}

type ClickStat struct {
	Name         string
	TotalClicks  int
	UniqueClicks int
}

type clickPathStat struct {
	Path         string `db:"path"`
	TotalClicks  int    `db:"total_clicks"`
	UniqueClicks int    `db:"unique_clicks"`
}

// ListClickStats returns the next level in the project → page → type → element
// hierarchy for a site and time range.
func ListClickStats(ctx context.Context, rng ztime.Range, selection ClickSelection) ([]ClickStat, error) {
	var paths []clickPathStat
	err := zdb.Select(ctx, &paths, `/* ListClickStats */
		select p.path,
			sum(c.total_clicks) as total_clicks,
			sum(c.unique_clicks) as unique_clicks
		from click_counts c
		join paths p on p.site_id = c.site_id and p.path_id = c.path_id
		where c.site_id = :site and c.hour >= :start and c.hour <= :end and p.event = true
		group by p.path`, map[string]any{
		"site":  MustGetSite(ctx).ID,
		"start": rng.Start,
		"end":   rng.End,
	})
	if err != nil {
		return nil, errors.Wrap(err, "ListClickStats")
	}

	level := selection.NextLevel()
	grouped := make(map[string]ClickStat)
	for _, path := range paths {
		click, ok := ParseClickPath(path.Path)
		if !ok ||
			(selection.Project != "" && click.Project != selection.Project) ||
			(selection.Page != "" && click.Page != selection.Page) ||
			(selection.Type != "" && click.Type != selection.Type) {
			continue
		}

		var name string
		switch level {
		case ClickLevelProject:
			name = click.Project
		case ClickLevelPage:
			name = click.Page
		case ClickLevelType:
			name = click.Type
		case ClickLevelElement:
			name = click.Element
		}
		stat := grouped[name]
		stat.Name = name
		stat.TotalClicks += path.TotalClicks
		stat.UniqueClicks += path.UniqueClicks
		grouped[name] = stat
	}

	stats := make([]ClickStat, 0, len(grouped))
	for _, stat := range grouped {
		stats = append(stats, stat)
	}
	slices.SortFunc(stats, func(a, b ClickStat) int {
		if n := cmp.Compare(b.TotalClicks, a.TotalClicks); n != 0 {
			return n
		}
		return cmp.Compare(a.Name, b.Name)
	})
	return stats, nil
}
