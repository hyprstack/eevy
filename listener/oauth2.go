package listener

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/event"
)

type OAuth2 struct {
	ListenerBase

	ClientId     string
	ClientSecret string
	Scope        []string
	TokenUrl     string
	EndPoint     string
	Verb         string
}

func (this *OAuth2) Exec(evt event.Event) {

	ep := this.createEndPoint(evt)
	verb := this.getVerb(evt)

	conf := &clientcredentials.Config{
		ClientID:     this.ClientId,
		ClientSecret: this.ClientSecret,
		Scopes:       this.Scope,
		TokenURL:     this.TokenUrl,
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
		gLog.Error("Unsupported verb: %s", this.Verb)
		return
	}
	if err != nil {
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

func (this *OAuth2) createEndPoint(evt event.Event) string {
	return this.magicString(this.EndPoint, evt)
}

func (this *OAuth2) getVerb(evt event.Event) string {
	return this.magicString(this.Verb, evt)
}

func (this *OAuth2) Init(conf config.Listener) {

	this.ClientId = conf.ClientId
	this.ClientSecret = conf.ClientSecret
	this.Scope = conf.Scope
	this.TokenUrl = conf.TokenUrl
	this.EndPoint = conf.EndPoint
	this.Verb = conf.Verb
	this.Message = conf.Message
}
