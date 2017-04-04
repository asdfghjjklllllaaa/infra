// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dashboard

import (
	"net/http"
	"time"

	"github.com/luci/gae/service/info"
	"github.com/luci/luci-go/appengine/gaemiddleware"
	"github.com/luci/luci-go/server/router"
	"github.com/luci/luci-go/server/templates"

	"infra/appengine/dashboard/backend"
)

var templateBundle = &templates.Bundle{
	Loader:    templates.FileSystemLoader("templates"),
	DebugMode: info.IsDevAppServer,
}

func pageBase() router.MiddlewareChain {
	return gaemiddleware.BaseProd().Extend(
		templates.WithTemplates(templateBundle),
	)
}

func init() {
	r := router.New()
	gaemiddleware.InstallHandlers(r, pageBase())
	r.GET("/", pageBase(), dashboard)
	http.DefaultServeMux.Handle("/", r)
}

func dashboard(ctx *router.Context) {
	c, w := ctx.Context, ctx.Writer

	dates := []string{}
	for i := 0; i < 7; i++ {
		dates = append(dates, time.Now().AddDate(0, 0, -i).Format("01-02-2006"))
	}

	// TODO(jojwang): not using returned Service entity because
	// the Start/EndTime fields must be converted to string
	// before passing info to template.
	_, err := backend.GetService(c, "monorail")
	if err != nil {
		http.Error(w, "Failed to query datastore, see logs", http.StatusInternalServerError)
		return
	}

	services := []backend.Service{}
	nonSLAServices := []backend.Service{}

	templates.MustRender(c, w, "pages/dash.tmpl", templates.Args{
		"ChopsServices":  services,
		"NonSLAServices": nonSLAServices,
		"Dates":          dates,
	})
}
