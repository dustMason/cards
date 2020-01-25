package actions

import (
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
)

const remarkTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <title>Title</title>
    <meta charset="utf-8">
    <style>
			@import url('https://fonts.googleapis.com/css?family=Heebo:300,400,500,700|Inconsolata');

			/* Global Styles */
			body {
				font-family: Heebo, "Helvetica Neue", Helvetica, Arial, sans-serif;
			}
			.remark-slide-content h1, .remark-slide-content h2 {
				font-weight: 500;
			}
			.remark-slide-content p, .remark-slide-content li, .remark-slide-content td, .remark-slide-content th {
				font-size: 24pt;
				line-height: 1.6;
			}
			.remark-code, .remark-inline-code {
				font-family: Inconsolata, monospace;
			}
			th {
				border-bottom: 1px solid black;
			}
			td, th {
				padding: 8px;
			}
			img {
				display: block;
				max-width:100%;
				max-height:100%;
				width: auto;
				height: auto;
			}
			.center img {
				margin-right: auto;
				margin-left: auto;
			}
			.smaller p, .smaller div, .smaller li, .smaller th, .smaller td {
				font-size: 18pt;
			}
			.footnote p {
				position: absolute;
				bottom: 3em;
				font-size: 8pt !important;
			}
			.remark-slide-number {
				background-color: white;
				padding: 0 5px;
				border-radius: 5px;
				font-size: 20px !important;
			}
			.no-number .remark-slide-number {
				display: none;
			}
			#qrcode {
				width: 384px;
				height: 384px;
			}

			/* Title Slide Layout. Use .smokescreen[...] to contain the h1/h2. */
			.remark-slide-content.title h1, .remark-slide-content.title h2, .remark-slide-content.title h3 {
				color: white;
				font-size: 50pt;
				margin: 30pt;
				font-weight: 300;
			}
			.remark-slide-content.title h2 {
				font-size: 30pt;
			}
			.remark-slide-content.title h3 {
				font-size: 22pt;
			}
			.smokescreen {
				width: 100%;
				position: absolute;
				left: 0px;
				top: 33%;
				background-color: rgba(0,0,0,.7);
				vertical-align: middle;
				text-align: center;
			}

			/* Columnar Layouts. Two- and three-column layouts use .col classnames and float
			 * next to each other. The img-right uses .col and .rc (for right-column) and
			 * they're not equal-width or height. For convenience, two-column layouts also
			 * allow you to name the right column with .rc classname, so you can switch
			 * between layouts without changing the markup, just the slide's class.
			 */

			/* Two-Column Layout */
			.two-column .rc, .two-column .col {
				width: 48%;
				float: left;
				margin-right: 1%;
			}
			.two-column .rc img, .two-column .col img, .three-column .col img, .three-column .rc img {
				display: block;
				min-width: 100%;
				min-height: 100%;
				max-width: 100%;
				max-height: 100%;
				width: auto;
				height: auto;
				margin: 0;
				padding: 0;
			}

			/* Two-Column Layout, Text Left, Image Right */
			.img-right .col {
				width: 62.5%;
				 padding-right: 1em;
			}
			.img-right .rc {
				position: absolute;
				top: 0;
				left: 62.5%;
				width: 37.5%;
				height: 100%;
				margin: 0;
				padding: 0;
			}
			.img-right .rc p { /* Remove empty line above image wrapped in <p> */
				padding: 0;
				margin: 0;
			}

			/* Three-Column Layout */
			.three-column .col, .three-column .rc {
				width: 32%;
				float: left;
			}

			/* Shrink Images To Fit In A Vertical Space */
			.img-450h img {
				display: block;
				max-height: 450px !important;
				width: auto;
				margin: 0;
				padding: 0;
			}
			.img-300h img {
				display: block;
				max-height: 300px !important;
				width: auto;
				margin: 0;
				padding: 0;
			}
			.center img, .img-center img {
				display: block;
				margin-left: auto;
				margin-right: auto;
			}
    </style>
  </head>
  <body>
    <textarea id="source">{{.Source}}</textarea>
    <script src="https://remarkjs.com/downloads/remark-latest.min.js">
    </script>
    <script>
      var slideshow = remark.create();
    </script>
  </body>
</html>
`

type remarkContents struct {
	Source string
}

func Present(source string) {
	tmpl, err := template.New("remark").Parse(remarkTemplate)
	if err != nil {
		panic(err)
	}

	file, err := ioutil.TempFile("", "cards-remark-*.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_ = tmpl.Execute(file, &remarkContents{Source: source})

	open, err := exec.LookPath("open")
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(open, file.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
