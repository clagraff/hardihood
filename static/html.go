package static

func HTML() string {
	return `
<html>
	<head>
		<title>{{.Config.Title}}</title>
		<meta http-equiv="refresh" content="{{.Config.Refresh}}">
		<link rel="icon" 
			type="image/png" 
			href="{{.Config.Favicon}}">
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
