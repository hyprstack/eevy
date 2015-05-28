package listener

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/hevnly/eevy/event"
)

type OAuth2Config interface {
	ListenerConfig

	GetClientId() string
	GetClientSecret() string
	GetTokenUrl() string
	GetScope() []string
	GetEndPoint() string
	GetVerb() string
	GetBody() string
}

// Make an http call authnticating via an OAuth2 url
type OAuth2 struct {
	ListenerBase

	Config OAuth2Config
}

// Satifies the Listener interface and makes the http call after authenticating
func (this *OAuth2) Exec(evt event.Event) {

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
		gLog.Error("Unsupported verb: %s", verb)
		return
	}
	if err != nil {
		gLog.Error("ep: %s", ep)
		gLog.Error(err.Error())
		return
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		gLog.Error(err.Error())
	}
	gLog.Debug("OAuth: %s", robots)
}

// Gets the end point for this listener
func (this *OAuth2) getEndPoint(evt event.Event) string {
	return magicString(this.Config.GetEndPoint(), evt)
}

// Gets the verb to be used in the http call
func (this *OAuth2) getVerb(evt event.Event) string {
	return magicString(this.Config.GetVerb(), evt)
}
