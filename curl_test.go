package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_isValidURL(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "voyagegroup url with http://", args: args{u: "http://voyagegroup.com"}, want: true},
		{name: "voyagegroup url", args: args{u: "voyagegroup.com"}, want: false},
		{name: "invalid url", args: args{u: "     \\\\"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidURL(tt.args.u); got != tt.want {
				t.Errorf("isValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidMethod(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "GET", args: args{method: "GET"}, want: true},
		{name: "POST", args: args{method: "POST"}, want: true},
		{name: "invalid method", args: args{method: "PIYOPIYO"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidMethod(tt.args.method); got != tt.want {
				t.Errorf("isValidMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHeader(t *testing.T) {
	type args struct {
		rawHeader string
	}
	tests := []struct {
		name      string
		args      args
		wantKey   string
		wantValue string
	}{
		{name: "X-Treasure", args: args{rawHeader: "X-Treasure: 🍺"}, wantKey: "X-Treasure", wantValue: "🍺"},
		{name: "Content-Type", args: args{rawHeader: "Content-Type: application/json"}, wantKey: "Content-Type", wantValue: "application/json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotValue := parseHeader(tt.args.rawHeader)
			if gotKey != tt.wantKey {
				t.Errorf("parseHeader() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotValue != tt.wantValue {
				t.Errorf("parseHeader() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		u       string
		headers []string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "success", args: args{u: "", headers: []string{"X-Treasure: 🍺"}}, want: "🍺"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		headerJSON, _ := json.MarshalIndent(r.Header, "", "  ")
		w.Write(headerJSON)
	}))
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Get(server.URL, tt.args.headers)
			defer resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			var header http.Header
			json.Unmarshal(body, &header)

			if expected := header.Get("X-Treasure"); expected != tt.want {
				t.Errorf("Get() = %v, want %v", expected, tt.want)
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		u       string
		headers []string
		body    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "with json", args: args{u: "", headers: []string{"Content-Type: application/json"}, body: `{"ajito":"🍺"}`}, want: `{"ajito":"🍺"}`},
		{name: "with json and multi header", args: args{u: "", headers: []string{"Content-Type: application/json", "X-Treasure: 🍺"}, body: `{"ajito":"🍺"}`}, want: `{"ajito":"🍺"}`},
		{name: "with urlencoded data", args: args{u: "", headers: []string{"Content-Type:application/x-www-form-urlencoded"}, body: `text=hello`}, want: `text=hello`},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		b, _ := ioutil.ReadAll(r.Body)
		w.Write(b)
	}))
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Post(server.URL, tt.args.headers, tt.args.body)
			defer resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			if expected := string(body); expected != tt.want {
				t.Errorf("Get() = %v, want %v", expected, tt.want)
			}
		})
	}
}
