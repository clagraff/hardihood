package static

func DefaultCSS() string {
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
	width: 960px;
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
}`
}
