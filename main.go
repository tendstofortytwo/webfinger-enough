package main

import (
	"net/http"
	"flag"
	"log"
	"encoding/json"
	"regexp"
	"strings"
	"fmt"
)

var (
	idp    = flag.String("idp", "", "OIDC IDP to point to")
	domain = flag.String("domain", "", "domain names for which to answer") 
)

type Response struct {
	Subject string `json:"subject,omitempty"`
	Links   []Link `json:"links,omitempty"`
}

type Link struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}

func die(w http.ResponseWriter) {
	w.WriteHeader(404)
	json.NewEncoder(w).Encode(Response {})
}

func main() {
	flag.Parse()
	if *idp == "" || *domain == "" {
		log.Fatalf("must provide idp and domain!")
	}
	atDomain := "@" + *domain
	http.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		resource := r.URL.Query().Get("resource")
		email := strings.TrimPrefix(resource, "acct:")
		if !strings.HasSuffix(email, atDomain) {
			die(w)
			return
		}
		username := strings.TrimSuffix(email, atDomain)
		// TODO: this is probably overly restrictive, but
		// good enough for me for now.
		if ok, err := regexp.Match("^[A-Za-z0-9_.]+$", []byte(username)); !ok || err != nil {
			die(w)
			return
		}
		json.NewEncoder(w).Encode(Response {
			Subject: fmt.Sprintf("acct:%s%s", username, atDomain),
			Links: []Link{
				{
					Rel: "http://openid.net/specs/connect/1.0/issuer",
					Href: *idp,
				},
			},
		})
	})
	log.Fatal(http.ListenAndServe(":9090", nil))
}