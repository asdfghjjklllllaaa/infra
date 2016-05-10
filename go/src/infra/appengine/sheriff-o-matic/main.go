// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package som implements HTTP server that handles requests to default module.
package som

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	"github.com/luci/gae/service/datastore"
	"github.com/luci/luci-go/appengine/gaeauth/server"
	"github.com/luci/luci-go/appengine/gaemiddleware"
	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/server/auth"
	"github.com/luci/luci-go/server/auth/identity"
	"github.com/luci/luci-go/server/middleware"
	"github.com/luci/luci-go/server/settings"
)

const authGroup = "sheriff-o-matic-access"

var (
	mainPage         = template.Must(template.ParseFiles("./index.html"))
	accessDeniedPage = template.Must(template.ParseFiles("./access-denied.html"))
)

var errStatus = func(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

var requireGoogler = func(w http.ResponseWriter, c context.Context) bool {
	isMember, err := auth.IsMember(c, authGroup)
	if !isMember || err != nil {
		msg := ""
		if !isMember {
			msg = "Access Denied"
		} else {
			msg = err.Error()
		}

		errStatus(w, http.StatusForbidden, msg)
		return false
	}
	return true
}

const settingsKey = "tree"

type settingsUIPage struct {
	settings.BaseUIPage
}

func (settingsUIPage) Title(c context.Context) (string, error) {
	return "Admin SOM settings", nil
}

func (settingsUIPage) Fields(c context.Context) ([]settings.UIField, error) {
	return []settings.UIField{
		{
			ID:    "Trees",
			Title: "Trees in SOM",
			Type:  settings.UIFieldText,
			Help:  `Trees listed in SOM. Comma separated values.`,
		},
	}, nil
}

func (settingsUIPage) ReadSettings(c context.Context) (map[string]string, error) {
	q := datastore.NewQuery("Tree")
	results := []*Tree{}
	datastore.Get(c).GetAll(q, &results)
	stringed := make([]string, len(results))
	for i, tree := range results {
		stringed[i] = tree.Name
	}

	return map[string]string{
		"Trees": strings.Join(stringed, ","),
	}, nil
}

func (settingsUIPage) WriteSettings(c context.Context, values map[string]string, who, why string) error {
	ds := datastore.Get(c)

	q := datastore.NewQuery("Tree")
	trees := []*Tree{}
	datastore.Get(c).GetAll(q, &trees)

	treeStr, ok := values["Trees"]

	if ok {
		toMake := strings.Split(treeStr, ",")
		for _, tree := range trees {
			for i, it := range toMake {
				if it == tree.Name {
					toMake[i] = ""
				}
			}
		}

		for _, it := range toMake {
			it = strings.TrimSpace(it)
			if it != "" {
				err := ds.Put(&Tree{
					Name:        it,
					DisplayName: strings.Replace(strings.Title(it), "_", " ", -1),
				})

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//// Handlers.
func indexPage(c context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if p.ByName("path") == "" {
		http.Redirect(w, r, "/chromium", http.StatusFound)
		return
	}

	user := auth.CurrentIdentity(c)

	if user.Kind() == identity.Anonymous {
		url, err := auth.LoginURL(c, "/")
		if err != nil {
			errStatus(w, http.StatusInternalServerError, fmt.Sprintf(
				"You must login. Additionally, an error was encountered while serving this request: %s", err.Error()))
		} else {
			http.Redirect(w, r, url, http.StatusFound)
		}

		return
	}

	isGoogler, err := auth.IsMember(c, authGroup)

	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	logoutURL, err := auth.LogoutURL(c, "/")

	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !isGoogler {
		err = accessDeniedPage.Execute(w, map[string]interface{}{
			"Group":     authGroup,
			"LogoutURL": logoutURL,
		})
		if err != nil {
			logging.Errorf(c, "while rendering index: %s", err)
		}
		return
	}

	data := map[string]interface{}{
		"User":      user.Email(),
		"LogoutUrl": logoutURL,
	}

	err = mainPage.Execute(w, data)
	if err != nil {
		logging.Errorf(c, "while rendering index: %s", err)
	}
}

func getTreesHandler(c context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !requireGoogler(w, c) {
		return
	}

	q := datastore.NewQuery("Tree")
	results := []*Tree{}
	err := datastore.Get(c).GetAll(q, &results)
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	txt, err := json.Marshal(results)
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(txt)
}

func getAlertsHandler(c context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !requireGoogler(w, c) {
		return
	}

	ds := datastore.Get(c)

	tree := p.ByName("tree")
	q := datastore.NewQuery("AlertsJSON")
	q = q.Ancestor(ds.MakeKey("Tree", tree))
	q = q.Order("-Date")
	q = q.Limit(1)

	results := []*AlertsJSON{}
	err := ds.GetAll(q, &results)
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(results) == 0 {
		logging.Warningf(c, "No alerts found for tree %s", tree)
		errStatus(w, http.StatusNotFound, fmt.Sprintf("Tree \"%s\" not found", tree))
		return
	}

	alertsJSON := results[0]
	w.Header().Set("Content-Type", "application/json")
	w.Write(alertsJSON.Contents)
}

func postAlertsHandler(c context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !requireGoogler(w, c) {
		return
	}

	tree := p.ByName("tree")
	ds := datastore.Get(c)

	alerts := AlertsJSON{
		Tree: ds.MakeKey("Tree", tree),
		Date: clock.Now(c),
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.Body.Close(); err != nil {
		errStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	out := make(map[string]interface{})
	err = json.Unmarshal(data, &out)

	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	out["date"] = alerts.Date.String()
	data, err = json.Marshal(out)

	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	alerts.Contents = data
	err = datastore.Get(c).Put(&alerts)
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func getAnnotationsHandler(c context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !requireGoogler(w, c) {
		return
	}

	q := datastore.NewQuery("Annotation")
	results := []*Annotation{}
	datastore.Get(c).GetAll(q, &results)

	data, err := json.Marshal(results)
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func postAnnotationsHandler(c context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !requireGoogler(w, c) {
		return
	}

	annKey := p.ByName("annKey")
	action := p.ByName("action")
	ds := datastore.Get(c)

	if !(action == "add" || action == "remove") {
		errStatus(w, http.StatusNotFound, "Invalid action")
		return
	}

	annotation := &Annotation{
		KeyDigest: fmt.Sprintf("%x", sha1.Sum([]byte(annKey))),
		Key:       annKey,
	}

	err := ds.Get(annotation)
	if action == "remove" && err != nil {
		errStatus(w, http.StatusNotFound, fmt.Sprintf("Annotation %s not found", annKey))
		return
	}
	// The annotation probably doesn't exist if we're adding something

	if action == "add" {
		err = annotation.add(c, r.Body)
	} else if action == "remove" {
		err = annotation.remove(c, r.Body)
	}

	if err != nil {
		errStatus(w, http.StatusBadRequest, err.Error())
		return
	}

	err = r.Body.Close()
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = ds.Put(annotation)
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	data, err := json.Marshal(annotation)
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// base is the root of the middleware chain.
func base(h middleware.Handler) httprouter.Handle {
	methods := auth.Authenticator{
		&server.OAuth2Method{Scopes: []string{server.EmailScope}},
		server.CookieAuth,
		&server.InboundAppIDAuthMethod{},
	}
	h = auth.Use(h, methods)
	if !appengine.IsDevAppServer() {
		h = middleware.WithPanicCatcher(h)
	}
	return gaemiddleware.BaseProd(h)
}

//// Routes.
func init() {
	settings.RegisterUIPage(settingsKey, settingsUIPage{})

	router := httprouter.New()
	gaemiddleware.InstallHandlers(router, base)
	router.GET("/api/v1/trees/", base(auth.Authenticate(getTreesHandler)))
	router.GET("/api/v1/alerts/:tree", base(auth.Authenticate(getAlertsHandler)))
	router.POST("/api/v1/alerts/:tree", base(auth.Authenticate(postAlertsHandler)))
	router.GET("/api/v1/annotations/", base(auth.Authenticate(getAnnotationsHandler)))
	router.POST("/api/v1/annotations/:annKey/:action", base(auth.Authenticate(postAnnotationsHandler)))

	rootRouter := httprouter.New()
	rootRouter.GET("/*path", base(auth.Authenticate(indexPage)))

	http.DefaultServeMux.Handle("/api/", router)
	http.DefaultServeMux.Handle("/admin/", router)
	http.DefaultServeMux.Handle("/auth/", router)
	http.DefaultServeMux.Handle("/_ah/", router)
	http.DefaultServeMux.Handle("/", rootRouter)
}
