package admin_tabs

import (
	"github.com/aghape/admin"
	"github.com/aghape/core"
	"github.com/aghape/core/utils"
	"github.com/moisespsena-go/aorm"
	"github.com/moisespsena/go-path-helpers"
)

var (
	PKG      = path_helpers.GetCalledDir()
	KEY_TABS = PKG + ".tabs"
	KEY_TAB  = PKG + ".tab"
	THEME    = "tabbed"
)

func PrepareResource(res *admin.Resource, tabs Tabs, defaultTab *Tab) {
	index := res.Router.FindHandler("GET", admin.P_INDEX).(*admin.RouteHandler)
	indexHandler := index.Handle

	if defaultTab != nil {
		index.Handle = func(c *admin.Context) {
			c.Data().Set(KEY_TAB, defaultTab)
			indexHandler(c)
		}
	}

	scopesMap := &TabsData{tabs, make(map[string]*Tab)}

	for _, scope := range tabs {
		if scope.Path == "" {
			scope.Path = utils.ToParamString(scope.Title)
		}
		if scope.TitleKey == "" {
			scope.TitleKey = res.I18nPrefix + ".tabs." + scope.Path
		}
		scopesMap.ByPath[scope.Path] = scope
	}

	res.Data.Set(KEY_TABS, scopesMap)

	for _, tab := range tabs {
		res.Router.Get("/"+tab.Path, TabHandler(res, index.Config, indexHandler, tab))
	}

	res.DefaultFilter(func(context *core.Context, db *aorm.DB) *aorm.DB {
		scopePath := GetTabPath(context)
		if scope, ok := scopesMap.ByPath[scopePath]; ok {
			return scope.Handler(res, context, db)
		}
		return db
	})

	res.UseTheme(THEME)
}
