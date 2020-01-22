package static

func DefaultCSS() string {
	return `
html {
    height: 100%;
}

body {
	height: 100%;
	margin: 0;
	padding: 0;
	text-align: center;
	font-family: Sans;

	background: rgb(238,238,224);
	background: linear-gradient(0deg, rgba(238,238,224,1) 0%, rgba(238,238,224,1) 68%, rgba(0,188,208,1) 68%);

	background-repeat: no-repeat;
	background-attachment: fixed;
}

section {
	display: block;
	background: #ffffff;
	border: 1px solid #F0F0F0;
	text-align: left;
	margin: 2rem auto 2rem auto;
	width: 960px;
	padding: 1rem;

	box-shadow: 0 2px 5px 0 rgba(0,0,0,0.16),0 2px 10px 0 rgba(0,0,0,0.12);
	/*-webkit-box-shadow: 0px 0px 5px 0px rgba(0,0,0,0.30);
	-moz-box-shadow: 0px 0px 5px 0px rgba(0,0,0,0.30);
	box-shadow: 0px 0px 5px 0px rgba(0,0,0,0.30);*/

	color: #ffffff;

}

section header {
	font-size: 2rem;
	color: #212121;
}

section ul {
	list-style: none;
	padding: 0px;
	margin: 0 1rem;
}

section li {
	border: 1px solid #F0F0F0;
	background-color: #FBFBFB;
	padding: 1rem;
	margin: 1rem 0 1rem 0;
}

section li span {
	float: right;
}

span.healthy {
	color: #0D8F4F;
}

span.sick {
	color: #D70461;
}

li.healthy {
	background-color: #0D8F4F;
}

li.sick {
	background-color: #d50000;
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
