package webizen

import (
	"strings"

	"github.com/kierdavis/argo"
	"github.com/linkeddata/gold"
)

var (
	foaf = argo.NewNamespace("http://xmlns.com/foaf/0.1/")
)

func assertURI(uri string) (uris map[string]int64) {
	uris = make(map[string]int64)
	names := make(map[string]string)
	mboxes := make(map[string]string)
	images := make(map[string]string)

	g := gold.NewGraph(uri)
	err := g.LoadURI(uri)
	if err != nil {
		return
	}

	for elt := range g.IterTriples() {
		var k, v string

		switch s := elt.Subject.(type) {
		case *argo.Resource:
			k = s.URI
		}
		switch o := elt.Object.(type) {
		case *argo.Resource:
			v = o.URI
		case *argo.Literal:
			v = o.Value
		}

		if len(k) < len(uri) || k[:len(uri)] != uri {
			continue
		}

		if elt.Predicate.Equal(foaf.Get("name")) {
			uris[k] = 0
			names[k] = v
		}
		if elt.Predicate.Equal(foaf.Get("img")) ||
			elt.Predicate.Equal(foaf.Get("depiction")) {
			uris[k] = 0
			images[k] = v
		}
		if elt.Predicate.Equal(foaf.Get("mbox")) {
			uris[k] = 0
			mboxes[k] = v
		}
	}

	for k, _ := range uris {
		user := &User{Uri: k}
		db.InsertOne(user)
		db.Get(user)
		uris[k] = user.Id
	}
	for k, v := range names {
		db.Delete(&UserName{User: uris[k]})
		db.InsertOne(&UserName{uris[k], v})
	}
	for k, v := range images {
		db.Delete(&UserImage{User: uris[k]})
		db.InsertOne(&UserImage{uris[k], v})
	}
	for k, v := range mboxes {
		db.Delete(&UserMbox{User: uris[k]})
		lst := strings.SplitN(v, "@", 2)
		m, n := lst[0], ""
		if len(lst) > 1 {
			n = lst[1]
		}
		if len(m) > 7 && m[:7] == "mailto:" {
			m = m[7:]
		}
		db.InsertOne(&UserMbox{uris[k], m, n})
	}

	return uris
}

type result struct {
	Image []string `json:"img,omitempty"`
	Mbox  []string `json:"mbox,omitempty"`
	Name  []string `json:"name,omitempty"`
}

func search(query string) (r map[string]result) {
	r = map[string]result{}
	id := map[string]int64{}
	uri := map[int64]string{}
	urik := map[int64]bool{}

	for _, elt := range strings.Split(query, " ") {
		if len(elt) < 6 {
			continue
		}
		if elt[:6] == "https:" || elt[:5] == "http:" {
			for k, v := range assertURI(elt) {
				id[k] = v
				uri[v] = k
				urik[v] = true
			}
		}
	}

	lookup := func(i int64) string {
		if len(uri[i]) == 0 {
			user := new(User)
			db.Id(i).Get(user)
			id[user.Uri] = user.Id
			uri[user.Id] = user.Uri
		}
		return uri[i]
	}

	db.Where("name LIKE ?", `%`+query+`%`).Iterate(new(UserName), func(i int, bean interface{}) error {
		elt := bean.(*UserName)
		v := r[lookup(elt.User)]
		v.Name = append(v.Name, elt.Name)
		r[lookup(elt.User)] = v
		return nil
	})

	db.Where("local LIKE ?", `%`+query+`%`).Iterate(new(UserMbox), func(i int, bean interface{}) error {
		elt := bean.(*UserMbox)
		v := r[lookup(elt.User)]
		w := elt.Local
		if len(elt.Domain) > 0 {
			w += "@" + elt.Domain
		}
		v.Mbox = append(v.Mbox, w)
		r[lookup(elt.User)] = v
		return nil
	})

	for k := range r {
		v := r[k]

		var images []UserImage
		db.Find(&images, &UserImage{User: id[k]})
		for _, elt := range images {
			v.Image = append(v.Image, elt.Image)
		}

		r[k] = v
	}

	for k := range urik {
		if k < 1 {
			continue
		}
		v := r[uri[k]]
		if len(v.Name) == 0 {
			db.Where("user = ?", k).Iterate(new(UserName), func(i int, bean interface{}) error {
				elt := bean.(*UserName)
				v.Name = append(v.Name, elt.Name)
				return nil
			})
		}
		if len(v.Image) == 0 {
			db.Where("user = ?", k).Iterate(new(UserImage), func(i int, bean interface{}) error {
				elt := bean.(*UserImage)
				v.Image = append(v.Image, elt.Image)
				return nil
			})
		}
		r[uri[k]] = v
	}
	return
}
