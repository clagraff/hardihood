package main

import (
	"fmt"
	"html"
	"log"
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
		<section>
			<header>hardihood</header>
			<ul>
				<li>
					CDN
					<span class="healthy">Healthy <span class="icon">✓</span></span>
				</li>
				<li>
					Conversions
					<span class="sick">Outage <span class="icon">✕</span></span>
				</li>
				<li>
					Site delivery
					<span class="sick">Outage <span class="icon">✕</span></span>
				</li>
				<li>
					API
					<span class="healthy">Healthy <span class="icon">✓</span></span>
				</li>
			</ul>
		</section>
		<section>
			<header>hardihood</header>
			<ul>
				<li>
					CDN
					<span class="healthy">Healthy</span>
				</li>
				<li>
					Conversions
					<span class="sick">Outage</span>
				</li>
				<li>
					Site delivery
					<span class="sick">Outage</span>
				</li>
				<li>
					API
					<span class="healthy">Healthy</span>
				</li>
			</ul>
		</section>

	</body>
</html>
`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		page := struct {
			CSS string
		}{
			CSS: getCSS(),
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
