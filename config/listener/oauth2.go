package listener

type OAuth2 struct {
	ListenerBase
}

func (this *OAuth2) GetClientId() string {
	return this.getString("clientId")
}

func (this *OAuth2) GetClientSecret() string {
	return this.getString("clientSecret")
}

func (this *OAuth2) GetTokenUrl() string {
	return this.getString("tokenUrl")
}

func (this *OAuth2) GetScope() []string {
	return this.getStringSlice("scope")
}

func (this *OAuth2) GetEndPoint() string {
	return this.getString("endPoint")
}

func (this *OAuth2) GetVerb() string {
	return this.getString("verb")
}

func (this *OAuth2) GetBody() string {
	return this.getString("body")
}
