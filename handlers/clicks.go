package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"zgo.at/goatcounter/v2"
	"zgo.at/zhttp"
	"zgo.at/zstd/ztime"
)

type clicks struct{}

type clickBreadcrumb struct {
	Name string
	URL  string
}

type clickViewRow struct {
	goatcounter.ClickStat
	URL string
}

func (clicks) mount(a chi.Router) {
	a.Get("/clicks", zhttp.Wrap(clicksPage))
	a.Get("/clicks/{project}", zhttp.Wrap(clicksPage))
	a.Get("/clicks/{project}/{page}", zhttp.Wrap(clicksPage))
	a.Get("/clicks/{project}/{page}/{type}", zhttp.Wrap(clicksPage))
}

func clicksPage(w http.ResponseWriter, r *http.Request) error {
	user := User(r.Context())
	selection := goatcounter.ClickSelection{
		Project: chi.URLParam(r, "project"),
		Page:    chi.URLParam(r, "page"),
		Type:    chi.URLParam(r, "type"),
	}

	var (
		rng ztime.Range
		err error
	)
	if r.URL.Query().Get("period-start") == "" && r.URL.Query().Get("period-end") == "" {
		rng = timeRange(r.Context(), "month", user.Settings.Timezone.Loc(), user.Settings.SundayStartsWeek)
	} else {
		rng, err = getPeriod(w, r, Site(r.Context()), user)
		if err != nil {
			return err
		}
	}

	stats, err := goatcounter.ListClickStats(r.Context(), rng, selection)
	if err != nil {
		return err
	}

	query := r.URL.RawQuery
	withQuery := func(path string) string {
		if query != "" {
			return path + "?" + query
		}
		return path
	}
	base := newGlobals(w, r).Base + "/clicks"
	breadcrumbs := []clickBreadcrumb{{Name: "全部项目", URL: withQuery(base)}}
	path := base
	for _, part := range []string{selection.Project, selection.Page, selection.Type} {
		if part == "" {
			break
		}
		path += "/" + url.PathEscape(part)
		breadcrumbs = append(breadcrumbs, clickBreadcrumb{Name: part, URL: withQuery(path)})
	}

	level := selection.NextLevel()
	rows := make([]clickViewRow, 0, len(stats))
	var total, unique int
	for _, stat := range stats {
		rowURL := ""
		if level != goatcounter.ClickLevelElement {
			rowURL = withQuery(path + "/" + url.PathEscape(stat.Name))
		}
		rows = append(rows, clickViewRow{ClickStat: stat, URL: rowURL})
		total += stat.TotalClicks
		unique += stat.UniqueClicks
	}

	labels := map[goatcounter.ClickLevel]string{
		goatcounter.ClickLevelProject: "项目",
		goatcounter.ClickLevelPage:    "页面",
		goatcounter.ClickLevelType:    "点击类型",
		goatcounter.ClickLevelElement: "按钮 / 元素",
	}
	location := user.Settings.Timezone.Loc()
	return zhttp.Template(w, "clicks.gohtml", struct {
		Globals
		Rows        []clickViewRow
		Breadcrumbs []clickBreadcrumb
		Level       goatcounter.ClickLevel
		LevelLabel  string
		Total       int
		Unique      int
		StartDate   string
		EndDate     string
		Selection   goatcounter.ClickSelection
	}{
		Globals:     newGlobals(w, r),
		Rows:        rows,
		Breadcrumbs: breadcrumbs,
		Level:       level,
		LevelLabel:  labels[level],
		Total:       total,
		Unique:      unique,
		StartDate:   rng.Start.In(location).Format("2006-01-02"),
		EndDate:     rng.End.In(location).Format("2006-01-02"),
		Selection:   selection,
	})
}

func clickPath(parts ...string) string {
	return strings.Join(parts, "/")
}
