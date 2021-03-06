package admin_tabs

import (
	"github.com/ecletus/admin"
	"github.com/ecletus/core"
	"github.com/ecletus/core/utils"
	"github.com/moisespsena-go/aorm"
	"context"
)

type Tab struct {
	Title    string
	Path     string
	TitleKey string
	Handler  func(t *Tabs, res *admin.Resource, context *core.Context, db *aorm.DB) *aorm.DB
	Default  bool
	Enabled  func(tabs *Tabs, ctx *admin.Context) bool
	scheme   *admin.Scheme
}

func (s *Tab) Scheme() *admin.Scheme {
	return s.scheme
}

func (s *Tab) URL(res *admin.Resource, context *core.Context) string {
	if s.Default {
		return res.GetContextIndexURI(context)
	}
	return res.GetContextIndexURI(context) + "/" + s.Path
}

type Tabs struct {
	Resource   *admin.Resource
	Tabs       []*Tab
	ByPath     map[string]*Tab
	defaultTab *Tab
}

func (t *Tabs) interseptor(chain *admin.Chain) {
	ctx := chain.Context
	var indexTabs []*Tab
	var currentTab *Tab
	for _, tab := range t.Tabs {
		if tab.Enabled == nil || tab.Enabled(t, ctx) {
			indexTabs = append(indexTabs, tab)
		}
	}
	if currentTab == nil {
		currentTab = t.defaultTab
	}
	ctx.SetValue(KEY_TABS, indexTabs)
	ctx.SetValue(KEY_TAB, currentTab)
	chain.Pass()
}
func (t *Tabs) Register(scheme *admin.Scheme) {
	var defaul bool
	for _, cat := range scheme.Categories {
		if cat == DEFAULT_SCHEME_CATEGORY {
			defaul = true
			break
		}
	}
	tab := &Tab{
		Title:    scheme.SchemeName,
		TitleKey: scheme.I18nKey(),
		scheme:   scheme,
		Path:     scheme.Path(),
		Default:  defaul,
	}
	if t.ByPath == nil {
		t.ByPath = map[string]*Tab{}
	}
	if tab.Path == "" {
		tab.Path = utils.ToParamString(tab.Title)
	}

	if tab.TitleKey == "" {
		tab.TitleKey = t.Resource.I18nPrefix + ".tabs." + tab.Path
	}

	if !tab.Default {
		t.ByPath[tab.Path] = tab
	}
	t.Tabs = append(t.Tabs, tab)
}

func GetTabPath(ctx context.Context) string {
	if tab := GetTab(ctx); tab != nil {
		return tab.Path
	}
	return ""
}

func GetTab(ctx context.Context) *Tab {
	if tab := ctx.Value(KEY_TAB); tab != nil {
		return tab.(*Tab)
	}
	return nil
}

func TabHandler(res *admin.Resource, config *admin.RouteConfig, indexHandler admin.Handler, scope *Tab) *admin.RouteHandler {
	return admin.NewHandler(func(c *admin.Context) {
		c.Breadcrumbs().Append(core.NewBreadcrumb(res.GetContextIndexURI(c.Context), res.GetLabelKey(true), ""))
		c.SetValue("page_title", c.T(scope.TitleKey, scope.Title))
		c.SetValue(KEY_TAB, scope)
		indexHandler(c)
	}, config)
}
