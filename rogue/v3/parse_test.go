package v3_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func TestParseHtml(t *testing.T) {
	testCases := []struct {
		name          string
		htmlContent   string
		expectedText  string
		expectedSpans []v3.TextSpan
	}{
		{
			name:         "simple bold text",
			htmlContent:  `<p class="p1">This is some <b>bold </b>text</p>`,
			expectedText: "This is some bold text\n",
			expectedSpans: []v3.TextSpan{
				{StartIndex: 13, EndIndex: 18, Format: v3.FormatV3Span{"b": "true"}},
			},
		},
		{
			name: "simple list with text after",
			htmlContent: `
<ul class="ul1">
<li class="li4"><span class="s1"></span>Bullet 1</li>
<li class="li4"><span class="s1"></span>Bullet 2</li>
<li class="li4"><span class="s1"></span>Bullet 3</li>
</ul>
<p class="p1">after the list</p>
`,
			expectedText: "Bullet 1\nBullet 2\nBullet 3\nafter the list\n",
			expectedSpans: []v3.TextSpan{
				{StartIndex: 0, EndIndex: 8, Format: v3.FormatV3BulletList(0)},
				{StartIndex: 9, EndIndex: 17, Format: v3.FormatV3BulletList(0)},
				{StartIndex: 18, EndIndex: 26, Format: v3.FormatV3BulletList(0)},
			},
		},
		{
			name:          "bold with font weight normal",
			htmlContent:   `<p class="p1">This is some not <b style="font-weight:normal;">bold</b> text</p>`,
			expectedText:  "This is some not bold text",
			expectedSpans: []v3.TextSpan{},
		},
		{
			name:         "p with large font size",
			htmlContent:  `<p style="font-size:24px;">A header with some <b>bold</b> in it</p>`,
			expectedText: "A header with some bold in it\n",
			expectedSpans: []v3.TextSpan{
				{StartIndex: 0, EndIndex: 29, Format: v3.FormatV3Header(2)},
				{StartIndex: 19, EndIndex: 23, Format: v3.FormatV3Span{"b": "true"}},
			},
		},
		{
			name: "all span formats with Apple Notes",
			htmlContent: `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
<meta http-equiv="Content-Style-Type" content="text/css">
<title></title>
<meta name="Generator" content="Cocoa HTML Writer">
<meta name="CocoaVersion" content="2487.5">
<style type="text/css">
p.p1 {margin: 0.0px 0.0px 0.0px 0.0px; font: 13.0px 'Helvetica Neue'}
span.s1 {text-decoration: underline}
span.s2 {text-decoration: line-through}
</style>
</head>
<body>
<p class="p1"><b>Just</b> a <i>paragraph</i> <span class="s1">with</span> <span class="s2">some <b>formatting</b></span> in it</p>
</body>
</html>`,
			expectedText: "Just a paragraph with some formatting in it\n",
			expectedSpans: []v3.TextSpan{
				{StartIndex: 0, EndIndex: 4, Format: v3.FormatV3Span{"b": "true"}},
				{StartIndex: 7, EndIndex: 16, Format: v3.FormatV3Span{"i": "true"}},
				{StartIndex: 17, EndIndex: 21, Format: v3.FormatV3Span{"u": "true"}},
				{StartIndex: 22, EndIndex: 37, Format: v3.FormatV3Span{"s": "true"}},
				{StartIndex: 27, EndIndex: 37, Format: v3.FormatV3Span{"b": "true"}},
			},
		},
		{

			name: "Many styles - apple notes",
			htmlContent: `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
<meta http-equiv="Content-Style-Type" content="text/css">
<title></title>
<meta name="Generator" content="Cocoa HTML Writer">
<meta name="CocoaVersion" content="2487.5">
<style type="text/css">
p.p1 {margin: 0.0px 0.0px 0.0px 0.0px; font: 20.0px 'Helvetica Neue'}
p.p2 {margin: 0.0px 0.0px 0.0px 0.0px; font: 13.0px 'Helvetica Neue'; min-height: 15.0px}
p.p3 {margin: 0.0px 0.0px 2.0px 0.0px; font: 16.0px 'Helvetica Neue'}
p.p4 {margin: 0.0px 0.0px 0.0px 0.0px; font: 13.0px 'Helvetica Neue'}
li.li4 {margin: 0.0px 0.0px 0.0px 0.0px; font: 13.0px 'Helvetica Neue'}
span.s1 {font: 9.0px Menlo}
ol.ol1 {list-style-type: decimal}
ul.ul1 {list-style-type: disc}
</style>
</head>
<body>
<p class="p1"><b>Title</b></p>
<p class="p2"><br></p>
<p class="p3"><b>Heading</b></p>
<p class="p2"><br></p>
<p class="p4"><b>Subheading</b></p>
<p class="p2"><br></p>
<ul class="ul1">
<li class="li4"><span class="s1"></span>Bullet 1</li>
<li class="li4"><span class="s1"></span>Bullet 2</li>
<li class="li4"><span class="s1"></span>Bullet 3</li>
</ul>
<p class="p2"><br></p>
<ol class="ol1">
<li class="li4">List 1</li>
<li class="li4">List 2</li>
<li class="li4">List 3</li>
</ol>
</body>
</html>`,
			expectedText: "Title\n\nHeading\n\nSubheading\n\nBullet 1\nBullet 2\nBullet 3\n\nList 1\nList 2\nList 3\n",
			expectedSpans: []v3.TextSpan{
				{StartIndex: 0, EndIndex: 5, Format: v3.FormatV3Header(2)},
				{StartIndex: 0, EndIndex: 5, Format: v3.FormatV3Span{"b": "true"}},
				{StartIndex: 7, EndIndex: 14, Format: v3.FormatV3Header(3)},
				{StartIndex: 7, EndIndex: 14, Format: v3.FormatV3Span{"b": "true"}},
				{StartIndex: 16, EndIndex: 26, Format: v3.FormatV3Span{"b": "true"}},
				{StartIndex: 28, EndIndex: 36, Format: v3.FormatV3BulletList(0)},
				{StartIndex: 37, EndIndex: 45, Format: v3.FormatV3BulletList(0)},
				{StartIndex: 46, EndIndex: 54, Format: v3.FormatV3BulletList(0)},
				{StartIndex: 56, EndIndex: 62, Format: v3.FormatV3OrderedList(0)},
				{StartIndex: 63, EndIndex: 69, Format: v3.FormatV3OrderedList(0)},
				{StartIndex: 70, EndIndex: 76, Format: v3.FormatV3OrderedList(0)},
			},
		},
		{
			name:         "Single sentence with special characters",
			htmlContent:  `<p class="p1">Hö'elün (fl. 1162–1210) <b>was</b> a Mongolian noblewoman 🚺 and the <span style="text-decoration: underline">mother of Temüjin, better known as Genghis</span>.</p>`,
			expectedText: "Hö'elün (fl. 1162–1210) was a Mongolian noblewoman 🚺 and the mother of Temüjin, better known as Genghis.\n",
			expectedSpans: []v3.TextSpan{
				{StartIndex: 24, EndIndex: 27, Format: v3.FormatV3Span{"b": "true"}},
				{StartIndex: 62, EndIndex: 104, Format: v3.FormatV3Span{"u": "true"}},
			},
		},
		{
			name:         "Single sentence - from wikipedia",
			htmlContent:  `<meta charset='utf-8'><b style="color: rgb(32, 33, 34); font-family: sans-serif; font-size: 14px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; white-space: normal; background-color: rgb(245, 255, 250); text-decoration-thickness: initial; text-decoration-style: initial; text-decoration-color: initial;"><a href="https://en.wikipedia.org/wiki/H%C3%B6%27el%C3%BCn" title="Hö'elün" style="text-decoration: none; color: var(--color-progressive,#36c); background: none; overflow-wrap: break-word;">Hö'elün</a></b><span style="color: rgb(32, 33, 34); font-family: sans-serif; font-size: 14px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; white-space: normal; background-color: rgb(245, 255, 250); text-decoration-thickness: initial; text-decoration-style: initial; text-decoration-color: initial; display: inline !important; float: none;"><span> </span>(</span><abbr title="floruit ('flourished'&nbsp;– known to have been active at a particular time or during a particular period)" style="border-bottom: 0px; cursor: help; text-decoration: underline dotted; color: rgb(32, 33, 34); font-family: sans-serif; font-size: 14px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; white-space: normal; background-color: rgb(245, 255, 250);">fl.</abbr><span style="color: rgb(32, 33, 34); font-family: sans-serif; font-size: 14px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; white-space: nowrap; background-color: rgb(245, 255, 250); text-decoration-thickness: initial; text-decoration-style: initial; text-decoration-color: initial;"> 1162–1210</span><span style="color: rgb(32, 33, 34); font-family: sans-serif; font-size: 14px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; white-space: normal; background-color: rgb(245, 255, 250); text-decoration-thickness: initial; text-decoration-style: initial; text-decoration-color: initial; display: inline !important; float: none;">) was a<span> </span></span><a href="https://en.wikipedia.org/wiki/Mongol_Empire" title="Mongol Empire" style="text-decoration: none; color: var(--color-progressive,#36c); background: none rgb(245, 255, 250); overflow-wrap: break-word; font-family: sans-serif; font-size: 14px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; white-space: normal;">Mongolian</a><span style="color: rgb(32, 33, 34); font-family: sans-serif; font-size: 14px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; white-space: normal; background-color: rgb(245, 255, 250); text-decoration-thickness: initial; text-decoration-style: initial; text-decoration-color: initial; display: inline !important; float: none;"><span> </span>noblewoman and the mother of<span> </span></span><a href="https://en.wikipedia.org/wiki/Genghis_Khan" title="Genghis Khan" style="text-decoration: none; color: var(--color-progressive,#36c); background: none rgb(245, 255, 250); overflow-wrap: break-word; font-family: sans-serif; font-size: 14px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; white-space: normal;">Temüjin</a><span style="color: rgb(32, 33, 34); font-family: sans-serif; font-size: 14px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; white-space: normal; background-color: rgb(245, 255, 250); text-decoration-thickness: initial; text-decoration-style: initial; text-decoration-color: initial; display: inline !important; float: none;">, better known as Genghis Khan</span>`,
			expectedText: "Hö'elün\u00a0(fl.\u20091162–1210) was a\u00a0Mongolian\u00a0noblewoman and the mother of\u00a0Temüjin, better known as Genghis Khan",
			expectedSpans: []v3.TextSpan{
				{StartIndex: 0, EndIndex: 7, Format: v3.FormatV3Span{"b": "true"}},
				{StartIndex: 0, EndIndex: 7, Format: v3.FormatV3Span{"a": "https://en.wikipedia.org/wiki/H%C3%B6%27el%C3%BCn"}},
				{StartIndex: 30, EndIndex: 39, Format: v3.FormatV3Span{"a": "https://en.wikipedia.org/wiki/Mongol_Empire"}},
				{StartIndex: 69, EndIndex: 76, Format: v3.FormatV3Span{"a": "https://en.wikipedia.org/wiki/Genghis_Khan"}},
			},
		},

		{
			name:         "Headers - Google docs",
			htmlContent:  `<meta charset='utf-8'><meta charset="utf-8"><b style="font-weight:normal;" id="docs-internal-guid-a6df7dee-7fff-f2a6-0e1b-8955ba59dc4a"><p dir="ltr" style="line-height:1.38;margin-top:0pt;margin-bottom:3pt;"><span style="font-size:26pt;font-family:Arial,sans-serif;color:#000000;background-color:transparent;font-weight:400;font-style:normal;font-variant:normal;text-decoration:none;vertical-align:baseline;white-space:pre;white-space:pre-wrap;">Title</span></p><br /><p dir="ltr" style="line-height:1.38;margin-top:0pt;margin-bottom:16pt;"><span style="font-size:15pt;font-family:Arial,sans-serif;color:#666666;background-color:transparent;font-weight:400;font-style:normal;font-variant:normal;text-decoration:none;vertical-align:baseline;white-space:pre;white-space:pre-wrap;">Subtitle</span></p><h1 dir="ltr" style="line-height:1.38;margin-top:20pt;margin-bottom:6pt;"><span style="font-size:20pt;font-family:Arial,sans-serif;color:#000000;background-color:transparent;font-weight:400;font-style:normal;font-variant:normal;text-decoration:none;vertical-align:baseline;white-space:pre;white-space:pre-wrap;">Heading</span></h1><h2 dir="ltr" style="line-height:1.38;margin-top:18pt;margin-bottom:6pt;"><span style="font-size:16pt;font-family:Arial,sans-serif;color:#000000;background-color:transparent;font-weight:400;font-style:normal;font-variant:normal;text-decoration:none;vertical-align:baseline;white-space:pre;white-space:pre-wrap;">Heading 2</span></h2><h3 dir="ltr" style="line-height:1.38;margin-top:16pt;margin-bottom:4pt;"><span style="font-size:13.999999999999998pt;font-family:Arial,sans-serif;color:#434343;background-color:transparent;font-weight:400;font-style:normal;font-variant:normal;text-decoration:none;vertical-align:baseline;white-space:pre;white-space:pre-wrap;">Heading 3</span></h3></b>`,
			expectedText: "Title\n\nSubtitle\nHeading\nHeading 2\nHeading 3\n",
			expectedSpans: []v3.TextSpan{
				{StartIndex: 0, EndIndex: 5, Format: v3.FormatV3Header(1)},
				{StartIndex: 7, EndIndex: 15, Format: v3.FormatV3Header(3)},
				{StartIndex: 16, EndIndex: 23, Format: v3.FormatV3Header(2)},
				{StartIndex: 24, EndIndex: 33, Format: v3.FormatV3Header(3)},
				{StartIndex: 34, EndIndex: 43, Format: v3.FormatV3Header(3)},
			},
		},
		{
			name: "Headers - Apple Notes",
			htmlContent: `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
<meta http-equiv="Content-Style-Type" content="text/css">
<title></title>
<meta name="Generator" content="Cocoa HTML Writer">
<meta name="CocoaVersion" content="2487.5">
<style type="text/css">
p.p1 {margin: 0.0px 0.0px 0.0px 0.0px; font: 20.0px 'Helvetica Neue'}
p.p2 {margin: 0.0px 0.0px 0.0px 0.0px; font: 13.0px 'Helvetica Neue'; min-height: 15.0px}
p.p3 {margin: 0.0px 0.0px 2.0px 0.0px; font: 16.0px 'Helvetica Neue'}
p.p4 {margin: 0.0px 0.0px 0.0px 0.0px; font: 13.0px 'Helvetica Neue'}
</style>
</head>
<body>
<p class="p1"><b>Title</b></p>
<p class="p2"><br></p>
<p class="p3"><b>Heading</b></p>
<p class="p2"><br></p>
<p class="p4"><b>Subheading</b></p>
</body>
</html>`,
			expectedText: "Title\n\nHeading\n\nSubheading\n",
			expectedSpans: []v3.TextSpan{
				{StartIndex: 0, EndIndex: 5, Format: v3.FormatV3Header(2)},
				{StartIndex: 0, EndIndex: 5, Format: v3.FormatV3Span{"b": "true"}},
				{StartIndex: 7, EndIndex: 14, Format: v3.FormatV3Header(3)},
				{StartIndex: 7, EndIndex: 14, Format: v3.FormatV3Span{"b": "true"}},
				{StartIndex: 16, EndIndex: 26, Format: v3.FormatV3Span{"b": "true"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plainText, spans, err := v3.ParseHtml(tc.htmlContent)
			require.NoError(t, err)
			require.Equal(t, tc.expectedText, plainText, "Text does not match")
			require.Equal(t, tc.expectedSpans, spans)
		})
	}
}
