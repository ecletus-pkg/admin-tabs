package admin_tabs

import (
	"github.com/ecletus/admin"
	"github.com/ecletus/plug"
	"github.com/moisespsena/go-path-helpers"
)

var (
	PKG      = path_helpers.GetCalledDir()
	KEY_TABS = PKG + ".tabs"
	KEY_TAB  = PKG + ".tab"
	THEME    = "tabbed"
	SCHEME_CATEGORY = PKG
	DEFAULT_SCHEME_CATEGORY = PKG + ".default"
)

func PrepareResource(res *admin.Resource) *Tabs {
	tabs := &Tabs{Resource: res}
	res.On(admin.E_SCHEME_ADDED, func(e plug.EventInterface) {
		s := e.(*admin.SchemeEvent).Scheme
		scat := SCHEME_CATEGORY
		for _, cat := range s.Categories {
			if cat == scat {
				tabs.Register(s)
				return
			}
		}
	})
	res.Data.Set(KEY_TABS, tabs)
	res.UseTheme(THEME)
	return tabs
}
