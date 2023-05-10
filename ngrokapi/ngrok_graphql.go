package ngrokapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	config "com.roger.ngrok.linebot/config"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

var s = `
	schema {
		query: Query
	}

	type Query {
		appusers(): [Appuser]
		#appuser(id: ID!): AppUser
		tunnels(): [Tunnel]!
	}

	type Appuser {
		id: ID!
		username: String
		name: String
		email: String!
	}

	type Tunnel {
		id: ID!
		public_url: String!
		started_at: String!
		proto: String!
		region: String!
		#tunnel_session
		#endpoint
		forwards_to: String!
	}	

`

type appuser struct {
	ID       graphql.ID
	Username string
	Name     string
	Email    string
}

// Schema
type Resolver struct{}

// func (*Resolver) Query() *QueryResolver {
// 	return &QueryResolver{}
// }

// Resolver
// type QueryResolver struct{}

func (*Resolver) Appusers() (*[]*AppuserResolver, error) {
	// do the job
	usersUrl := BaseUrl + "/app/users"
	// GET
	req, err := http.NewRequest("GET", usersUrl, nil)
	if err != nil {
		return nil, err
	}

	apiKey, err := config.GetEnvString("ApiKey")
	if err != nil {
		return nil, err
	}
	req.Header.Set(`Authorization`, `Bearer `+apiKey)
	req.Header.Set(`ngrok-version`, `2`)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	respData := make(map[string]any)
	defer resp.Body.Close()
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(rawBody, &respData); err != nil {
		return nil, err
	}

	// api result
	userData := respData["application_users"]
	ret := []*AppuserResolver{}
	if userArray, ok := userData.([]map[string]any); ok {
		for _, user := range userArray {
			ret = append(ret, &AppuserResolver{
				u: &appuser{
					ID:       graphql.ID(user["id"].(string)),
					Username: user["username"].(string),
					Name:     user["name"].(string),
					Email:    user["email"].(string),
				},
			})
		}
	}
	return &ret, nil
}

type AppuserResolver struct {
	u *appuser
}

func (u *AppuserResolver) ID() graphql.ID {
	return u.u.ID
}
func (u *AppuserResolver) Username() *string {
	return &u.u.Username
}
func (u *AppuserResolver) Name() *string {
	return &u.u.Name
}
func (u *AppuserResolver) Email() string {
	return u.u.Email
}

type tunnel struct {
	ID         graphql.ID
	PublicURL  string
	StartedAt  string
	Proto      string
	Region     string
	ForwardsTo string
}

type tunnelresolver struct {
	t *tunnel
}

func (t *tunnelresolver) ID() graphql.ID {
	return t.t.ID
}
func (t *tunnelresolver) PublicURL() string {
	return t.t.PublicURL
}
func (t *tunnelresolver) StartedAt() string {
	return t.t.StartedAt
}
func (t *tunnelresolver) Proto() string {
	return t.t.Proto
}
func (t *tunnelresolver) Region() string {
	return t.t.Region
}
func (t *tunnelresolver) ForwardsTo() string {
	return t.t.ForwardsTo
}

func (*Resolver) Tunnels() ([]*tunnelresolver, error) {
	data, err := ngrokApiGet("/tunnels")
	if err != nil {
		return nil, err
	}

	tunnels := data["tunnels"].([]any)
	ret := []*tunnelresolver{}
	for _, row := range tunnels {
		tunnel := &tunnel{}
		t := row.(map[string]any)
		fmt.Printf("%+v\n", t)
		tunnel.ID = graphql.ID(t["id"].(string))
		tunnel.PublicURL = t["public_url"].(string)
		tunnel.StartedAt = t["started_at"].(string)
		tunnel.Proto = t["proto"].(string)
		tunnel.Region = t["region"].(string)
		tunnel.ForwardsTo = t["forwards_to"].(string)

		ret = append(ret, &tunnelresolver{tunnel})
	}
	return ret, nil
}

// http get
func ngrokApiGet(path string) (map[string]any, error) {
	// do the job
	usersUrl := BaseUrl + path
	// GET
	req, err := http.NewRequest("GET", usersUrl, nil)
	if err != nil {
		return nil, err
	}

	apiKey, err := config.GetEnvString("ApiKey")
	if err != nil {
		return nil, err
	}
	req.Header.Set(`Authorization`, `Bearer `+apiKey)
	req.Header.Set(`ngrok-version`, `2`)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	respData := make(map[string]any)
	defer resp.Body.Close()
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(rawBody, &respData); err != nil {
		return nil, err
	}
	return respData, nil
}

// http.HttpHandler
var schema = graphql.MustParseSchema(s, &Resolver{})

func GetGraphqlHttpHandler() http.Handler {
	return &relay.Handler{Schema: schema}
}

// test andpoint
func GetNgrokUsers(w http.ResponseWriter, r *http.Request) {
	resp := schema.Exec(context.Background(), `
	{ 
		appusers {
			id
			email
			name
			username
		}
	}`, "", nil)
	fmt.Printf("%+v\n", resp)

	if len(resp.Errors) > 0 {
		// error
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, resp.Errors[0].Error())
	} else {
		// success
		rawData, err := resp.Data.MarshalJSON()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}
		w.Write(rawData)
	}
}

func GetNgrokTunnels(w http.ResponseWriter, r *http.Request) {
	resp := schema.Exec(context.Background(), `
	{ 
		tunnels {
			id
			public_url
			proto
			started_at
			forwards_to
		}
	}`, "", nil)
	fmt.Printf("%+v\n", resp)

	if len(resp.Errors) > 0 {
		// error
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, resp.Errors[0].Error())
	} else {
		// success
		rawData, err := resp.Data.MarshalJSON()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}
		w.Write(rawData)
	}
}
