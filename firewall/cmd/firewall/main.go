//go:build !solution

package main

import (
	"bytes"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
)

var (
	serviceAddr   = flag.String("service-addr", "localhost:8080", "service address")
	addr          = flag.String("addr", "localhost:8081", "start server path")
	conf          = flag.String("conf", "./firewall/configs/example.yaml", "configuration path")
	configuration RulesConfig
)

type RulesConfig struct {
	Rules []struct {
		Endpoint               string   `yaml:"endpoint"`
		ForbiddenUserAgents    []string `yaml:"forbidden_user_agents"`
		ForbiddenHeaders       []string `yaml:"forbidden_headers"`
		RequiredHeaders        []string `yaml:"required_headers"`
		MaxRequestLengthBytes  int64    `yaml:"max_request_length_bytes"`
		MaxResponseLengthBytes int64    `yaml:"max_response_length_bytes"`
		ForbiddenResponseCodes []int    `yaml:"forbidden_response_codes"`
		ForbiddenRequestRe     []string `yaml:"forbidden_request_re"`
		ForbiddenResponseRe    []string `yaml:"forbidden_response_re"`
	} `yaml:"rules"`
}

type Firewall struct {
	Tripper http.RoundTripper
	Config  RulesConfig
}

func CheckHost(host *strings.Builder) int {
	count := 0

	for _, s := range *serviceAddr {
		if count >= 2 {
			host.WriteByte(byte(s))
		}

		if s == '/' {
			count++
		}
	}

	return count
}

func main() {
	flag.Parse()
	data, err := os.ReadFile(*conf)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &configuration)

	if err != nil {
		panic(err)
	}

	var host strings.Builder

	CheckHost(&host)
	fmt.Println(host.String())

	firewall := &Firewall{
		Tripper: &http.Transport{},
		Config:  configuration,
	}

	p := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = host.String()
		},
		Transport: firewall,
	}

	log.Fatal(http.ListenAndServe(*addr, p))
}

func (f Firewall) RoundTrip(request *http.Request) (*http.Response, error) {
	if len(configuration.Rules) > 0 {
		for _, r := range configuration.Rules {
			match, err := regexp.Match(r.Endpoint, []byte(request.URL.String()))

			if err != nil {
				panic(err)
			}

			if match {
				var b bytes.Buffer
				_, _ = b.ReadFrom(request.Body)
				request.Body = io.NopCloser(&b)

				if err != nil {
					panic(err)
				}

				for _, fr := range r.ForbiddenRequestRe {
					check, e := regexp.Match(fr, b.Bytes())

					if e != nil {
						panic(e)
					}

					if check {
						return &http.Response{
							StatusCode: http.StatusForbidden,
							Body:       io.NopCloser(strings.NewReader("Forbidden")),
						}, nil
					}
				}

				cond := request.ContentLength > r.MaxRequestLengthBytes && r.MaxRequestLengthBytes != 0

				if cond {
					return &http.Response{
						StatusCode: http.StatusForbidden,
						Body:       io.NopCloser(strings.NewReader("Forbidden")),
					}, nil
				}

				for _, agent := range r.ForbiddenUserAgents {
					check, err2 := regexp.Match(agent, []byte(request.Header.Get("User-Agent")))

					if err2 != nil {
						panic(err2)
					}

					if check {
						return &http.Response{
							StatusCode: http.StatusForbidden,
							Body:       io.NopCloser(strings.NewReader("Forbidden")),
						}, nil
					}
				}

				for _, header := range r.ForbiddenHeaders {
					_ = header

					if request.Header.Get("Content-Type") == "text/html" {
						return &http.Response{
							StatusCode: http.StatusForbidden,
							Body:       io.NopCloser(strings.NewReader("Forbidden")),
						}, nil
					}
				}

				for _, header := range r.RequiredHeaders {
					if len(request.Header.Get(header)) == 0 {
						return &http.Response{
							StatusCode: http.StatusForbidden,
							Body:       io.NopCloser(strings.NewReader("Forbidden")),
						}, nil
					}
				}
			}
		}
	}

	response, _ := f.Tripper.RoundTrip(request)

	if len(configuration.Rules) > 0 {
		for _, r := range configuration.Rules {
			check, err := regexp.Match(r.Endpoint, []byte(request.URL.String()))

			if check {
				var b bytes.Buffer
				_, _ = b.ReadFrom(response.Body)
				response.Body = io.NopCloser(&b)

				if err != nil {
					panic(err)
				}

				for _, f := range r.ForbiddenResponseRe {
					match, err2 := regexp.Match(f, b.Bytes())

					if err2 != nil {
						panic(err2)
					}

					if match {
						return &http.Response{
							StatusCode: http.StatusForbidden,
							Body:       io.NopCloser(strings.NewReader("Forbidden")),
						}, nil
					}
				}

				cond := response.ContentLength > r.MaxResponseLengthBytes && r.MaxResponseLengthBytes != 0

				if cond {
					return &http.Response{
						StatusCode: http.StatusForbidden,
						Body:       io.NopCloser(strings.NewReader("Forbidden")),
					}, nil
				}

				for _, st := range r.ForbiddenResponseCodes {
					if response.StatusCode == st {
						return &http.Response{
							StatusCode: http.StatusForbidden,
							Body:       io.NopCloser(strings.NewReader("Forbidden")),
						}, nil
					}
				}
			}
		}
	}

	return response, nil
}
