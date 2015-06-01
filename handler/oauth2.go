package handler

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/hevnly/eevy/event"
	"github.com/hevnly/eevy/handler/config"
)

// Make an http call authnticating via an OAuth2 url
type OAuth2 struct {
	HandlerBase

	Config config.OAuth2
}

// Satifies the Listener interface and makes the http call after authenticating
func (this *OAuth2) Exec(evt event.Event) {

	this.Log.Handler(this, &evt)

	ep := this.getEndPoint(evt)
	verb := this.getVerb(evt)
	cid := magicString(this.Config.GetClientId(), evt)
	sec := magicString(this.Config.GetClientSecret(), evt)
	sco := this.Config.GetScope()
	tul := magicString(this.Config.GetTokenUrl(), evt)

	conf := &clientcredentials.Config{
		ClientID:     cid,
		ClientSecret: sec,
		Scopes:       sco,
		TokenURL:     tul,
	}
	client := conf.Client(oauth2.NoContext)

	var res *http.Response
	var err error
	switch verb {
	case "get":
		res, err = client.Get(ep)
	case "post":
		res, err = client.Post(ep, "", nil)
	default:
		this.Log.Error("Unsupported verb: %s", verb)
		return
	}
	if err != nil {
		this.Log.Error("ep: %s", ep)
		this.Log.Error(err.Error())
		return
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		this.Log.Error(err.Error())
	}
	this.Log.Debug("OAuth: %s", robots)
}

// Gets the end point for this listener
func (this *OAuth2) getEndPoint(evt event.Event) string {
	return magicString(this.Config.GetEndPoint(), evt)
}

// Gets the verb to be used in the http call
func (this *OAuth2) getVerb(evt event.Event) string {
	return magicString(this.Config.GetVerb(), evt)
}

func (this *OAuth2) GetType() string {

	return this.GetConfig().GetType()
}

func (this *OAuth2) GetConfig() config.Handler {
	return this.Config
}
