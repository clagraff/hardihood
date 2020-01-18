package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"
)

func getCSS() string {
	return `
		body {
			background-color: #fdfdfd;
			text-align: center;
			font-family: Sans;
		}

		section {
			display: block;
			background: #ffffff;
			border: 1px solid #F0F0F0;
			border-radius: 9px;
			text-align: left;
			margin: 2rem auto 2rem auto;
			width: 1200px;
			padding: 1rem;

			-webkit-box-shadow: 0px 0px 5px 0px rgba(0,0,0,0.30);
			-moz-box-shadow: 0px 0px 5px 0px rgba(0,0,0,0.30);
			box-shadow: 0px 0px 5px 0px rgba(0,0,0,0.30);
		}

		section header {
			font-size: 2rem;
		}

		section ul {
			list-style: none;
			padding: 0px;
			margin: 0 1rem;
		}

		section li {
			border: 1px solid #F0F0F0;
			border-radius: 9px;	
			background-color: #FBFBFB;
			padding: 1rem;
			margin: 1rem 0 1rem 0;
		}

		section li span {
			float: right;
		}

		span.healthy {
			color: green;
		}

		span.sick {
			color: red;
		}

		span.healthy .icon {
			background-color: green;
		}

		span.sick .icon {
			background-color: red;
		}

		span.icon {
			border-radius: 4rem;
			width: 1.2rem;
			text-align: center;
			margin-left: 1rem;
			color: white;
			font-weight: bold;
		}

`
}

func listingHTML() string {
	return `
<html>
	<head>
		<title>Status Page</title>
		<style>
			{{.CSS}}
		</style>
	</head>
	<body>
	{{range .Services}}
		<section>
			<header>{{.Name}}</header>
			<ul>
				{{range .Checks}}
				<li>
					{{.Description}}
					{{$status := .Status}}
					<span class="{{$status.CSSIdent}}">{{$status.Name}}<span class="icon">{{$status.HTMLChar}}</span></span>
				</li>
				{{end}}
			</ul>
		</section>
	{{end}}
	</body>
</html>
`
}

type Status interface {
	Name() string
	CSSIdent() string
	HTMLChar() string
}

type status struct {
	name     string
	cssIdent string
	htmlChar string
}

func (s status) Name() string     { return s.name }
func (s status) CSSIdent() string { return s.cssIdent }
func (s status) HTMLChar() string { return s.htmlChar }

func MakeStatus(name, cssIdent, htmlChar string) Status {
	return status{
		name:     name,
		cssIdent: cssIdent,
		htmlChar: htmlChar,
	}
}

var Healthy Status = MakeStatus("Healthy", "healthy", "✓")
var Sick Status = MakeStatus("Sick", "sick", "✕")

type Check interface {
	Description() string
	Status() Status
}

type check struct {
	description string
	statusFn    func() Status
}

func (c check) Description() string { return c.description }
func (c check) Status() Status      { return c.statusFn() }

func MakeCheck(desc string, fn func() Status) Check {
	return check{
		description: desc,
		statusFn:    fn,
	}
}

type Service interface {
	Name() string
	Checks() []Check
}

type service struct {
	name   string
	checks []Check
}

func (s service) Name() string    { return s.name }
func (s service) Checks() []Check { return s.checks }

func MakeService(name string, checks []Check) Service {
	return service{
		name:   name,
		checks: checks,
	}
}

func main() {
	exampleService := MakeService(
		"example",
		[]Check{
			MakeCheck(
				"Random number is even",
				func() Status {
					if rand.Intn(100)%2 == 0 {
						return Healthy
					}
					return Sick
				},
			),
			MakeCheck(
				"Random number is odd",
				func() Status {
					if rand.Intn(100)%2 != 0 {
						return Healthy
					}
					return Sick
				},
			),
		},
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := struct {
			CSS      string
			Services []Service
		}{
			CSS:      getCSS(),
			Services: []Service{exampleService},
		}

		html := listingHTML()

		tmpl, err := template.New("page").Parse(html)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		err = tmpl.Execute(w, page)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
