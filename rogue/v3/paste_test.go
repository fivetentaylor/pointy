package v3_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

// function to load html from file big_paste.html
func loadHtml(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		panic(err)
	}

	// read the file
	b := make([]byte, fInfo.Size())
	n, err := f.Read(b)
	if err != nil {
		panic(err)
	}

	return string(b[:n])
}

func TestPasteHtml(t *testing.T) {
	testCases := []struct {
		name             string
		fn               func(r *v3.Rogue)
		index            int
		htmlContent      string
		expectedMarkdown string
		expectedHtml     string
	}{
		{
			name:             "empty string",
			htmlContent:      "",
			expectedMarkdown: "",
			expectedHtml:     "",
		},
		{
			name:             "p tag",
			htmlContent:      "<p>This is some text</p>",
			expectedMarkdown: "This is some text\n\n",
			expectedHtml:     "<p data-rid=\"q_1\"><span data-rid=\"auth0_3\">This is some text</span></p>",
		},
		{
			name:             "simple bold text",
			htmlContent:      `<p class="p1">This is some <b>bold</b> text</p>`,
			expectedMarkdown: "This is some **bold** text\n\n\n\n",
			expectedHtml:     "<p data-rid=\"auth0_25\"><span data-rid=\"auth0_3\">This is some </span><strong data-rid=\"auth0_16\">bold</strong><span data-rid=\"auth0_20\"> text</span></p><p data-rid=\"q_1\"></p>",
		},
		{
			name: "simple list",
			htmlContent: `
<ul class="ul1">
<li class="li4"><span class="s1"></span>Bullet 1</li>
<li class="li4"><span class="s1"></span>Bullet 2</li>
<li class="li4"><span class="s1"></span>Bullet 3</li>
</ul>
`,
			expectedMarkdown: "- Bullet 1\n- Bullet 2\n- Bullet 3\n\n\n",
		},
		{
			name: "in the middle of a doc",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "Hello World!\n\nGoodbye moon.")
				require.NoError(t, err)
			},
			index:            13,
			htmlContent:      `<p class="p1">This is some <b>bold</b> text</p>`,
			expectedMarkdown: "Hello World!\n\nThis is some **bold** text\n\n\n\nGoodbye moon.\n\n",
		},

		{
			name:             "Paste from apple notes",
			htmlContent:      _appleNote,
			expectedMarkdown: "## **Title**\n\n\n\n### **Heading**\n\n\n\n**Subheading**\n\n\n\n- Bullet 1\n- Bullet 2\n- Bullet 3\n\n\n1. List 1\n1. List 2\n1. List 3\n\n\n",
		},
		{
			name:         "Single sentence with special characters",
			htmlContent:  `<p class="p1">Hö'elün (fl. 1162–1210) <b>was</b> a Mongolian noblewoman 🚺 and the <span style="text-decoration: underline">mother of Temüjin, better known as Genghis</span>.</p>`,
			expectedHtml: "<p data-rid=\"auth0_108\"><span data-rid=\"auth0_3\">Hö&#39;elün (fl.\u20091162–1210) </span><strong data-rid=\"auth0_27\">was</strong><span data-rid=\"auth0_30\"> a Mongolian noblewoman 🚺 and the </span><u data-rid=\"auth0_65\">mother of Temüjin, better known as Genghis</u><span data-rid=\"auth0_107\">.</span></p><p data-rid=\"q_1\"></p>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("auth0")
			if tc.fn != nil {
				tc.fn(r)
			}

			_, err := r.PasteHtml(tc.index, tc.htmlContent)
			require.NoError(t, err)

			if tc.expectedHtml != "" {
				html, err := r.GetHtml(v3.RootID, v3.LastID, true, false)
				require.NoError(t, err)
				require.Equal(t, tc.expectedHtml, html)
			}

			if tc.expectedMarkdown != "" {
				mkd, err := r.GetMarkdownBeforeAfter(v3.RootID, v3.LastID)
				require.NoError(t, err)
				require.Equal(t, tc.expectedMarkdown, mkd)
			}

		})
	}
}

func TestPaste(t *testing.T) {
	testCases := []struct {
		name          string
		fn            func(r *v3.Rogue)
		visIndex      int
		selLen        int
		curSpanFormat v3.FormatV3Span
		items         []v3.PasteItem
		expectedHtml  string
	}{
		{
			name: "empty paste",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "Hello World!\nGoodbye moon.")
				require.NoError(t, err)
			},
			visIndex: 13,
			selLen:   0,
			items: []v3.PasteItem{
				{Kind: "string", Mime: "text/html", Data: ``},
			},
			expectedHtml: "<p><span>Hello World!</span></p><p><span>Goodbye moon.</span></p>",
		},
		{
			name: "apple note",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "Hello World!\nGoodbye moon.")
				require.NoError(t, err)
			},
			visIndex: 13,
			selLen:   0,
			items: []v3.PasteItem{
				{Kind: "string", Mime: "text/html", Data: _appleNote},
			},
			expectedHtml: "<p><span>Hello World!</span></p><h2><strong>Title</strong></h2><p></p><h3><strong>Heading</strong></h3><p></p><p><strong>Subheading</strong></p><p></p><ul><li><span>Bullet 1</span></li><li><span>Bullet 2</span></li><li><span>Bullet 3</span></li></ul><p></p><ol><li><span>List 1</span></li><li><span>List 2</span></li><li><span>List 3</span></li></ol><p><span>Goodbye moon.</span></p>",
		},
		{
			name: "paste simple html keep formatting",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(0, 1, v3.FormatV3Header(1))
				require.NoError(t, err)
			},
			visIndex: 6,
			selLen:   0,
			items: []v3.PasteItem{
				{Kind: "string", Mime: "text/html", Data: "<p>cruel </p>"},
			},
			expectedHtml: "<h1><span>Hello cruel World!</span></h1>",
		},
		{
			name:     "big paste from wikipedia html",
			visIndex: 0,
			selLen:   0,
			items: []v3.PasteItem{
				{Kind: "string", Mime: "text/html", Data: loadHtml("big_paste.html")},
			},
			expectedHtml: `<h3><span> From Wikipedia, the free encyclopedia </span></h3><h3><span>       The Founding Ceremony of the Nation </span></h3><p><span>     1967 revision     Artist   </span><a href="https://en.wikipedia.org/wiki/Dong_Xiwen">Dong Xiwen</a><span>     Year   1953; revised 1954, 1967     Type   </span><a href="https://en.wikipedia.org/wiki/Oil_painting">Oil</a><span> on canvas     Dimensions   229 cm × 400 cm (90 in × 160 in)     Location   </span><a href="https://en.wikipedia.org/wiki/National_Museum_of_China">National Museum of China</a><span>, </span><a href="https://en.wikipedia.org/wiki/Beijing">Beijing</a><span>        Chinese name     </span><a href="https://en.wikipedia.org/wiki/Simplified_Chinese_characters">Simplified Chinese</a><span>   开国大典     </span><a href="https://en.wikipedia.org/wiki/Traditional_Chinese_characters">Traditional Chinese</a><span>   開國大典        </span><a href="https://en.wikipedia.org/wiki/Hanyu_Pinyin">Hanyu Pinyin</a><span>   Kāiguó dàdiǎn        </span><strong> Transcriptions </strong><span>         </span><strong><em>The Founding Ceremony of the Nation</em></strong><span> (or </span><strong><em>The Founding of the Nation</em></strong><span>) is a 1953 </span><a href="https://en.wikipedia.org/wiki/Oil_painting">oil painting</a><span> by Chinese artist </span><a href="https://en.wikipedia.org/wiki/Dong_Xiwen">Dong Xiwen</a><span>. It depicts </span><a href="https://en.wikipedia.org/wiki/Mao_Zedong">Mao Zedong</a><span> and other </span><a href="https://en.wikipedia.org/wiki/Chinese_Communist_Party">Communist Party</a><span> officials </span><a href="https://en.wikipedia.org/wiki/Proclamation_of_the_People%27s_Republic_of_China">proclaiming the People&#39;s Republic of China</a><span> at </span><a href="https://en.wikipedia.org/wiki/Tiananmen_Square">Tiananmen Square</a><span> on October 1, 1949. A prominent example of </span><a href="https://en.wikipedia.org/wiki/Socialist_realism">socialist realism</a><span>, it is one of the most celebrated works of official Chinese art. The painting was repeatedly revised, and a replica painting made to accommodate further changes, as some of the leaders it depicted fell from power and later were </span><a href="https://en.wikipedia.org/wiki/Political_rehabilitation#China">rehabilitated</a><span>. </span></p><p><span>  After the </span><a href="https://en.wikipedia.org/wiki/Chinese_Communist_Revolution">Chinese Communist Revolution</a><span>, the Party sought to memorialize their achievements through artworks. Dong was commissioned to create a visual representation of the October 1 ceremony, which he had attended. He viewed it as essential that the painting show both the people and their leaders. After working for three months, he completed an oil painting in a </span><a href="https://en.wikipedia.org/wiki/Folk_art">folk art</a><span> style, drawing upon Chinese art history for the contemporary subject. The success of the painting was assured when Mao viewed it and liked it, and it was reproduced in large numbers for display in the home. </span></p><p><span>  The 1954 purge of </span><a href="https://en.wikipedia.org/wiki/Gao_Gang">Gao Gang</a><span> from the government resulted in Dong being ordered to remove him from the painting. Gao&#39;s departure was not the last; Dong was forced to remove then-Chinese president </span><a href="https://en.wikipedia.org/wiki/Liu_Shaoqi">Liu Shaoqi</a><span> in 1967. The winds of political fortune continued to shift during the </span><a href="https://en.wikipedia.org/wiki/Cultural_Revolution">Cultural Revolution</a><span>, and a reproduction was painted by other artists in 1972, to accommodate another deletion. The replica was modified in 1979 to include the purged individuals, who had been rehabilitated. Both canvases are in the </span><a href="https://en.wikipedia.org/wiki/National_Museum_of_China">National Museum of China</a><span> in Beijing. </span></p><h2><span>   Background </span></h2><p><span>   Following the establishment of the People&#39;s Republic in 1949, Communists quickly took control of </span><a href="https://en.wikipedia.org/wiki/Chinese_art#Communist_and_socialist_art_(1950%E2%80%931980s)">art in China</a><span>. The </span><a href="https://en.wikipedia.org/wiki/Socialist_realism">socialist realism</a><span>that was characteristic of </span><a href="https://en.wikipedia.org/wiki/Soviet_art">Soviet art</a><span> came to be highly influential in the People&#39;s Republic. The new government proposed a series of paintings, preferably in oil, to memorialize the history of the Party, and its triumph in 1949. To this end, in December 1950, arts official Wang Yeqiu proposed to Deputy Minister of Culture </span><a href="https://en.wikipedia.org/wiki/Zhou_Yang_(literary_theorist)">Zhou Yang</a><span> that there be an art exhibition the following year to commemorate the 30th anniversary of the founding of the Party in China. Wang had toured the Soviet Union and observed its art, with which he was greatly impressed, and he proposed that sculptures and paintings be exhibited depicting the Party&#39;s history, for eventual inclusion in the planned Museum of the Chinese Revolution. Even before gaining full control of the country, the Party had used art as propaganda, a technique especially effective as much of the Chinese population was then illiterate. Wang&#39;s proposal was preliminarily approved in March 1951, and a committee, including the art critic and official </span><a href="https://en.wikipedia.org/wiki/Jiang_Feng_(artist)">Jiang Feng</a><span>, was appointed to seek suitable artists.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007785%E2%80%93786-1">[1]</a><span> Although, nearly 100 paintings were produced for the 1951 exhibition, not enough were found to be suitable, and it was cancelled.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn19-2">[2]</a><span> </span></p><p><span>    </span><em>The Founding Ceremony of the Nation</em><span> on exhibit along with artifacts of the October 1, 1949, ceremony. National Museum of China. Seen in 2018.    The use of oil paintings to memorialize events and make a political statement was not new; 19th-century examples include </span><a href="https://en.wikipedia.org/wiki/John_Trumbull">John Trumbull</a><span>&#39;s paintings for the </span><a href="https://en.wikipedia.org/wiki/United_States_Capitol#Art">United States Capitol</a><span> (1817–1821) and </span><a href="https://en.wikipedia.org/wiki/Jacques-Louis_David">Jacques-Louis David</a><span>&#39;s </span><a href="https://en.wikipedia.org/wiki/The_Coronation_of_Napoleon"><em>The Coronation of Napoleon</em></a><span> (1807).</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007785-3">[3]</a><span> Oil painting allowed for a blending of tones to produce a wide range of realistic, attractive colors, in a way not possible with traditional Chinese ink and brush painting.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007789-4">[4]</a><span> Wang admired how, in Moscow museums, </span><a href="https://en.wikipedia.org/wiki/Lenin">Lenin</a><span>&#39;s career was chronicled and made accessible to the masses through artifacts accompanied by oil paintings showing crucial moments in the Communist leader&#39;s career. He and higher-level officials decided to use a similar technique as they planned the Museum of the Chinese Revolution. Thus, they sought to chronicle the Party&#39;s history and showcase its accomplishments. Paintings were commissioned, even though the museum did not yet exist and would not open until 1961.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007790%E2%80%93792-5">[5]</a><span> Chinese leaders were eager to be presented in paintings, wanting to be immortalized as central characters in the nation&#39;s historical drama.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Hung_2005171%E2%80%93172-6">[6]</a><span> </span></p><p><span>  None of the works initially obtained for the museum depicted the crowning moment of the revolution, the ceremony at </span><a href="https://en.wikipedia.org/wiki/Tiananmen_Square">Tiananmen Square</a><span> on October 1, 1949, when </span><a href="https://en.wikipedia.org/wiki/Mao_Zedong">Mao Zedong</a><span> proclaimed the People&#39;s Republic. Officials deemed such a work essential.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn1-7">[7]</a><span> </span><a href="https://en.wikipedia.org/wiki/Dong_Xiwen">Dong Xiwen</a><span>, a professor at the </span><a href="https://en.wikipedia.org/wiki/Central_Academy_of_Fine_Arts">Central Academy of Fine Arts</a><span> (CAFA) in Beijing, was accomplished and politically reliable, and had been present at the October 1 ceremony: he was an obvious candidate.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007790%E2%80%93792-5">[5]</a><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn3-8">[8]</a><span> Although Dong later complained that never in his career had he had full freedom of choice as to his paintings</span><span>’</span><span> subjects,</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007783-9">[9]</a><span> </span><em>The Founding of the Nation</em><span> would make him famous.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn19-2">[2]</a><span> </span></p><h2><span>   Subject and techniques </span></h2><p><span>     Dong&#39;s original painting    The painting depicts the inaugural ceremony of the People&#39;s Republic of China on October 1, 1949. The focus is on Mao, who stands on </span><a href="https://en.wikipedia.org/wiki/Tiananmen_Gate">Tiananmen Gate</a><span>&#39;s balcony, reading his proclamation into (originally) two microphones.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews81-10">[10]</a><span> Dong took some liberties with the appearance of Tiananmen Gate, opening up the space in front of Mao to grant the chairman a more direct connection with his people,</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007809-11">[11]</a><span> something that architect </span><a href="https://en.wikipedia.org/wiki/Liang_Sicheng">Liang Sicheng</a><span> deemed a mistake for a builder, but artistically brilliant.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn19-2">[2]</a><span> Five </span><a href="https://en.wikipedia.org/wiki/Doves_as_symbols">doves</a><span> fly into the sky to Mao&#39;s right. Before him on Tiananmen Square, honor guards and members of patriotic organizations are assembled in orderly ranks, with some holding banners. </span><a href="https://en.wikipedia.org/wiki/Qianmen">Qianmen</a><span>, the gate at the south end of the square, is visible, as is </span><a href="https://en.wikipedia.org/wiki/Yongdingmen">Yongdingmen</a><span>gate (seen to the left of Mao).</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews81-10">[10]</a><span> Beyond the old city walls that at the time enclosed the square (they were torn down in the 1950s), the city of Beijing is visible, and in green is represented the nation of China, with those further scenes under bright sunlight and sharply defined clouds.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007810-12">[12]</a><span> October 1 had been an overcast day in Beijing; Dong took artistic license with the weather.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn3-8">[8]</a><span> </span></p><p><span>  To Mao&#39;s left are seen his lieutenants in the Communist takeover. In the original painting, the front row, which is ordered by rank, consisted of (from left) General </span><a href="https://en.wikipedia.org/wiki/Zhu_De">Zhu De</a><span>, </span><a href="https://en.wikipedia.org/wiki/Liu_Shaoqi">Liu Shaoqi</a><span>, Madame </span><a href="https://en.wikipedia.org/wiki/Song_Qingling">Song Qingling</a><span> (the widow of </span><a href="https://en.wikipedia.org/wiki/Sun_Zhongshan">Sun Zhongshan</a><span>), </span><a href="https://en.wikipedia.org/wiki/Li_Jishen">Li Jishen</a><span>, </span><a href="https://en.wikipedia.org/wiki/Zhang_Lan">Zhang Lan</a><span>(with beard), and at far right, General </span><a href="https://en.wikipedia.org/wiki/Gao_Gang">Gao Gang</a><span>. </span><a href="https://en.wikipedia.org/wiki/Zhou_Enlai">Zhou Enlai</a><span> was furthest left in the second row, and beside him were </span><a href="https://en.wikipedia.org/wiki/Dong_Biwu">Dong Biwu</a><span>, two men whose identities are uncertain, and furthest right, </span><a href="https://en.wikipedia.org/wiki/Guo_Moruo">Guo Moruo</a><span>. </span><a href="https://en.wikipedia.org/wiki/Lin_Boqu">Lin Boqu</a><span> was furthest left in the third row.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews81-10">[10]</a><span> The leaders congregate close to each other, while remaining a respectful distance from Mao. This emphasizes his primacy, as does the fact that he is depicted as taller than his lieutenants.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Hung_2005172%E2%80%93173-13">[13]</a><span> The point of view of the observer is well back on the balcony, from which most of the square would be obscured by the floor. Rather than show only Mao and sky, Dong manipulated the </span><a href="https://en.wikipedia.org/wiki/Perspective_(graphical)">perspective</a><span>, raising the horizon and intensifying the foreshortening of the balcony. Necessarily, only the officials and not the crowd below are represented as individuals; art historian </span><a href="https://en.wikipedia.org/wiki/Wu_Hung">Wu Hung</a><span> wrote that </span><span>“</span><span>the parading masses in the Square derive strength from a collective anonymity. The combination of the two—above and below, the leaders and the people—constitutes a comprehensive representation of New China.&#34;</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Hung_2005173-14">[14]</a><span> </span></p><p><span>  Mao and his officials are surrounded by </span><a href="https://en.wikipedia.org/wiki/Paper_lantern">lanterns</a><span>, symbols of prosperity;</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007783-9">[9]</a><span> the </span><a href="https://en.wikipedia.org/wiki/Chrysanthemum">chrysanthemums</a><span> on either side symbolize longevity. The doves represent peace restored to a nation long wracked by war.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007809-11">[11]</a><span> The new five-star </span><a href="https://en.wikipedia.org/wiki/Flag_of_China">flag of China</a><span>, rising over the people, represents the end of the </span><a href="https://en.wikipedia.org/wiki/Fengjian">feudal system</a><span> and the rebirth of the nation as the People&#39;s Republic.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn19-2">[2]</a><span> Mao, who is presented as a statesman, not as the revolutionary leader he was during the conflict, faces Qianmen, aligning himself along Beijing&#39;s old imperial North-South Axis, symbolizing his authority. The chairman is at the center of multiple, concentric circles in the painting, with the innermost formed by the front row of his comrades, another by the people in the square, and the outermost the old city walls. Surrounding them are the sunlit scenes, envisioning a glorious future for China with Mao the heart of the nation.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007809%E2%80%93810-15">[15]</a><span> </span></p><p><span>    Dong borrowed techniques from the </span><a href="https://en.wikipedia.org/wiki/Mogao_Caves">Dunhuang murals</a><span> of the Tang dynasty.    Although Dong had been trained in </span><a href="https://en.wikipedia.org/wiki/Western_painting">Western painting</a><span>, he chose a </span><a href="https://en.wikipedia.org/wiki/Folk_art">folk art</a><span> style for </span><em>The Founding of the Nation</em><span>, using bright, contrasting colors in a manner similar to that of </span><a href="https://en.wikipedia.org/wiki/Chinese_New_Year">New Year</a><span>&#39;s prints popular in China. He stated in 1953, </span><span>“</span><span>the Chinese people like bright, intense colors. This convention is in line with the theme of </span><em>The Founding Ceremony of the Nation</em><span>. In my choice of colors I did not hesitate to put aside the complex colors commonly adopted in Western painting as well as the conventional rules for oil painting.&#34;</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007809-11">[11]</a><span> Artists in the early years of the People&#39;s Republic, including Dong, sought to satisfy Chinese aesthetic tastes in their works and so depicted their subjects in bright original color, avoiding complex use of light and shadow on faces as in many Western paintings.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Bing65-16">[16]</a><span> By European standards, the painting&#39;s colors are overly intense and saturated. The color </span><a href="https://en.wikipedia.org/wiki/Vermilion">vermilion</a><span> was used for large areas of the columns, the carpets, and the lanterns, setting a tone for the work. The blooming flowers, the flags and banners, and the blue and white sky all give the painting a happy atmosphere</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Bing66-17">[17]</a><span>—a joyful, festive air, as well as giving </span><span>“</span><span>cultural sublimity&#34;, appropriate for a work depicting the founding of a nation.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-china_museum_2-18">[18]</a><span> </span></p><p><span>  Dong drew upon Chinese art history, using techniques from </span><a href="https://en.wikipedia.org/wiki/Mogao_Caves">Dunhuang murals</a><span> of the </span><a href="https://en.wikipedia.org/wiki/Tang_dynasty">Tang dynasty</a><span>, </span><a href="https://en.wikipedia.org/wiki/Ming_dynasty">Ming dynasty</a><span> portraits, and ancient figure paintings. Patterns on the carpet, columns, lanterns, and railing evoke cultural symbols.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-china_museum_2-18">[18]</a><span> The colors of the painting are reminiscent of crudely printed rural </span><a href="https://en.wikipedia.org/wiki/Woodcut">woodcuts</a><span>; this is emphasized by the black outlines of a number of objects, including the pillars and stone railing, as those outlines are characteristic of such woodcuts.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews82-19">[19]</a><span>Dong noted, </span><span>“</span><span>If this painting is rich in national styles, it is largely because I adopted these [native] approaches.&#34;</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007809-11">[11]</a><span> </span></p><h2><span>   Composition </span></h2><p><span>     </span><a href="https://en.wikipedia.org/wiki/Tiananmen_Gate">Tiananmen Gate</a><span> in 2009    </span><em>The Founding of the Nation</em><span> was one of several paintings commissioned for the new Museum of the Chinese Revolution from faculty members at CAFA. Two of these, </span><a href="https://en.wikipedia.org/wiki/Luo_Gongliu">Luo Gongliu</a><span>&#39;s </span><em>Tunnel Warfare</em><span> and Wang Shikuo&#39;s </span><em>Sending Him Off to the Army</em><span>, were completed in 1951; </span><em>The Founding of the Nation</em><span> was finished the following year.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews76%E2%80%9379-20">[20]</a><span> These commissions were regarded as from the government and were highly prestigious. State assistance, such as access to archives, was available.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007791%E2%80%93792-21">[21]</a><span> </span></p><p><span>  At the time CAFA chose Dong, he was painting workers at the </span><a href="https://en.wikipedia.org/wiki/Shijingshan">Shijingshan</a><span> power plant outside Beijing. Dong reviewed the photographs of the event, but found them unsatisfactory as none showed both the leaders and the people gathered in the square below, which he felt was necessary. He created a postcard-size sketch, but was dissatisfied with it, feeling it did not capture the grandeur of the occasion. Taking advice from other artists, Dong made adjustments to his plan.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn3-8">[8]</a><span> </span></p><p><span>  Dong rented a small room in Beijing above a store selling soy sauce.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn3-8">[8]</a><span> Jiang intervened to give Dong time and space to create the painting;</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews76%E2%80%9379-20">[20]</a><span> the artist needed three months to complete his work. The room was smaller than the painting, which is four meters wide, and Dong would affix part of the canvas to the ceiling, working on his back. To save commuting time, he slept in a chair. He smoked cigarettes constantly as he worked. His daughter brought meals, but he was often unable to eat. Once the painting was under way, several of Dong&#39;s colleagues, including oil painter Ai Zhongxin, came to visit. They decided that the figure of Mao, the central one of the painting, was not tall enough. Dong removed the figure of Mao from the canvas, and painted him again, increasing his height by just under an inch (2.54 cm).</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn3-8">[8]</a><span> </span></p><p><span>  In painting the sky and the pillars, Dong used a pen and brush, as if doing a traditional Chinese painting. He depicted the clothing in detail; Madame Song wears gloves showing flowers, while Zhang Lan&#39;s silk robe appears carefully ironed for the momentous day.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn7-22">[22]</a><span>Dong used sawdust to enhance the texture of the carpet on which Mao stands;</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007809-11">[11]</a><span> he painted the marble railing as yellowish rather than white, thus emphasizing the age of the Chinese nation.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn19-2">[2]</a><span> The leaders in the painting were asked to examine their portraits for accuracy.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Hung_2005172-23">[23]</a><span> </span></p><h2><span>   Reception and prominence </span></h2><p><span>   When the painting was unveiled in 1953, most Chinese critics were enthusiastic. </span><a href="https://en.wikipedia.org/wiki/Xu_Beihong">Xu Beihong</a><span>, the president of CAFA and a pioneer in using realism in oil painting, admired the manner in which the work fulfilled its political mission, but complained that because of the colors, it barely resembled an oil painting.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007810-12">[12]</a><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Bing65-16">[16]</a><span> He and others, though, saw that the painting opened a new chapter in Chinese art development.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Bing66-17">[17]</a><span> Zhu Dan, head of the People&#39;s Fine Arts Publishing House, which would reproduce the painting for the masses, argued that it was more a poster than an oil painting. Other artists stated that Dong&#39;s earlier works, such as </span><em>Kazakh Shepherdess</em><span>(1947) and </span><em>Liberation</em><span> (1949), were better examples of the new national style of art.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn7-22">[22]</a><span> Senior Party leaders, though, approved of the painting, as art historian Chang-Tai Hung put it, </span><span>“</span><span>seeing it as a testament to the young nation&#39;s evolving identity and growing confidence&#34;.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007810-12">[12]</a><span> </span></p><p><span>  Soon after the unveiling, Jiang wanted to arrange an exhibition at which government officials, including Mao, could view and publicly endorse the new Chinese art. He had connections in Mao&#39;s inner circle, and Dong and others organized it to be in conjunction with meetings at </span><a href="https://en.wikipedia.org/wiki/Zhongnanhai">Zhongnanhai</a><span> that Mao led. This was, most likely, the only time Mao attended an art exhibition after 1949. Mao visited the exhibition three times in between meetings and especially liked </span><em>The Founding of the Nation</em><span>—the official photograph of the event shows Mao and Zhou Enlai viewing the canvas with Dong.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews80-24">[24]</a><span> The chairman stared at the painting for a long time and finally said, </span><span>“</span><span>It is a great nation. It really is a great nation.&#34;</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Hung_2005172-23">[23]</a><span> Mao also stated that the portrayal of Dong Biwu was particularly well rendered. As Dong Biwu was in the second row, mostly hidden by the large Zhu De, Mao was most likely joking, but the favorable reaction by the country&#39;s leader assured the success of the painting.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews80-24">[24]</a><span> </span></p><p><span>  </span><em>The Founding of the Nation</em><span> was hailed as one of the greatest oil paintings ever by a Chinese artist by reviewers in that country, and more than 500,000 reproductions were sold in three months.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007783-9">[9]</a><span> Mao&#39;s praise helped boost the painting and its painter. Dong&#39;s techniques were seen as bridging the gap between the elitist medium of oil painting and popular art, and as a boost to Jiang&#39;s position that realistic art could be politically desirable.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews82-19">[19]</a><span> It was reproduced in primary and secondary school textbooks.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn19-2">[2]</a><span> The painting appeared on the front page of </span><a href="https://en.wikipedia.org/wiki/People%27s_Daily"><em>People&#39;s Daily</em></a><span> in September 1953, and became an officially approved interior decoration. One English-language magazine published by the Chinese government for distribution abroad showed a model family in a modern apartment, with a large poster of </span><em>The Founding of the Nation</em><span> on the wall.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews80%E2%80%9381-25">[25]</a><span> According to Chang-Tai Hung, the painting </span><span>“</span><span>became a celebrated propaganda piece&#34;.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2005920-26">[26]</a><span> </span></p><h2><span>   Later history and political changes </span></h2><p><span>     1954 revision with </span><a href="https://en.wikipedia.org/wiki/Gao_Gang">Gao Gang</a><span>deleted      The replica painting (1979 revision)    In February 1954 </span><a href="https://en.wikipedia.org/wiki/Gao_Gang">Gao Gang</a><span>, the head of the State Planning Council, was purged from government; he killed himself only months later. His presence in the painting immediately on Mao&#39;s left placed arts officials in a quandary. Given its popularity among officials and the people, </span><em>The Founding of the Nation</em><span> had to be shown at the Second National Arts Exhibition (1955), but it was unthinkable that Gao, deemed a traitor, should be depicted. Accordingly, Dong was ordered to remove Gao from the painting, which he did.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews82%E2%80%9383-27">[27]</a><span> </span></p><p><span>  In erasing Gao, Dong expanded the basket of pink chrysanthemums which stands at the officials</span><span>’</span><span> feet, and completed the depiction of Yongdingmen gate, in the original seen only in part behind Gao. He was forced to expand the section of sky seen above the people assembled in Tiananmen Square, which affected the placement of Mao as the center of attention. He compensated for this, to some extent, by adding two more microphones to Mao&#39;s right. Julia Andrews, in her book on the art of the People&#39;s Republic, suggested that Dong&#39;s solution was not entirely satisfactory as the microphones dominate the center of the painting, and Mao is diminished by the expanded space around him. The modified painting was shown in the 1955 exhibition, and in 1958 in Moscow. Although the painting was later altered again and does not exist in this form, this version is the one most commonly reproduced.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews82%E2%80%9383-27">[27]</a><span> </span></p><p><span>  When the Museum of the Chinese Revolution opened on Tiananmen Square in 1961, the painting was displayed on a huge wall in the gallery devoted to the Communist triumph, but in 1966, during the </span><a href="https://en.wikipedia.org/wiki/Cultural_Revolution">Cultural Revolution</a><span>, radicals shut down the museum, and it remained closed until 1969.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2005931-28">[28]</a><span> During that time, Chinese president Liu Shaoqi, accused of taking a </span><span>“</span><span>capitalist road&#34;, was purged from government. His removal from the painting was ordered in 1967, and Dong was tasked to carry it out.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007784-29">[29]</a><span> Dong had suffered during the Cultural Revolution: accused of being a rightist, he was expelled from the party for two years, sent to a rural </span><a href="https://en.wikipedia.org/wiki/Labor_camp">work camp</a><span>, and then was </span><a href="https://en.wikipedia.org/wiki/Political_rehabilitation#China">rehabilitated</a><span> by being made to labor as a steelworker.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn1-7">[7]</a><span> Dong&#39;s task was difficult, as Liu was one of the most prominent figures in the first row, standing to the left of Madame Song. Officials wanted Liu replaced with </span><a href="https://en.wikipedia.org/wiki/Lin_Biao">Lin Biao</a><span>, much in Mao&#39;s favor at the time. Dong was unwilling to give Lin prominence he had not then had, and though he could not refuse outright at the dangerous time of the Cultural Revolution, he eventually got permission to merely remove Liu. The figure was too large to simply delete, so Liu was repainted as Dong Biwu, and made to appear as if in the second row. According to Andrews, the attempt was a failure: </span><span>“</span><span>the new Dong Biwu does not recede into the second row as intended. Instead, he appears as a leering, glowing figure, a strangely malevolent character in the midst of an otherwise stately group&#34;.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews84-30">[30]</a><span> Officials deemed the revised work unexhibitable. Andrews speculated that Dong may have been trying to sabotage the change, or may have been affected by the stress of the years of the Cultural Revolution.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews84-30">[30]</a><span> </span></p><p><span>  In 1972, as part of a renovation of the Museum of the Chinese Revolution, officials wanted to exhibit Dong&#39;s painting again, but they decreed that Lin Boqu, whose white-haired figure was furthest left, must be removed.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews84-30">[30]</a><span> This was because the </span><a href="https://en.wikipedia.org/wiki/Gang_of_Four">Gang of Four</a><span>, then in control of China, blamed Lin Boqu (who had died in 1960) for opposing the marriage, in the revolutionary days, of Mao with </span><a href="https://en.wikipedia.org/wiki/Jiang_Qing">Jiang Qing</a><span> (one of the Four). Sources differ on what took place regarding the painting: Chang-Tai Hung related that Dong, terminally ill with cancer, could not make the changes, so his student </span><a href="https://en.wikipedia.org/wiki/Jin_Shangyi">Jin Shangyi</a><span> and another artist, Zhao Yu, were assigned to do the work. The two feared damaging the original canvas, so made an exact replica but for the required changes, with Dong brought forth from his hospital for consultations.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007784-29">[29]</a><span> According to Andrews, Jin and Zhao created the new version because Dong would not let anyone else alter his painting.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews84-30">[30]</a><span> Jin later stated that the painting, while effective politically, also shows Dong&#39;s inner world.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-china_museum_2-18">[18]</a><span> </span></p><p><span>  With the end of the Cultural Revolution in 1976 and the subsequent accession of </span><a href="https://en.wikipedia.org/wiki/Deng_Xiaoping">Deng Xiaoping</a><span>, many of the purged figures of earlier years were rehabilitated, and the authorities in 1979 decided to bring more historical accuracy to the painting. Dong had died in 1973; his family strongly opposed anyone altering the original painting, and the government respected their wishes. Jin was on a tour outside China, so the government assigned Yan Zhenduo to make changes to the replica. He placed Liu, Lin Boqu and Gao in the painting</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEC._Hung_2007784-29">[29]</a><span> and made other changes: a previously unidentifiable man in the back row now resembles Deng Xiaoping. The replica painting was restored to the Museum of the Chinese Revolution.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews85-31">[31]</a><span> </span></p><p><span>    Visitors take photographs of </span><em>The Founding Ceremony of the Nation</em><span>, National Museum of China, in 2018.    The painting was reproduced on </span><a href="https://en.wikipedia.org/wiki/Postage_stamps_and_postal_history_of_China#People's_Republic_of_China">Chinese postage stamps</a><span> in 1959 and 1999, for the tenth and fiftieth anniversaries of the founding of the People&#39;s Republic.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-scott_2007-32">[32]</a><span> Also in 1999, the museum authorized a private company to make small-scale gold foil reproductions of the painting. Dong&#39;s family sued, and in 2002 the courts found that Dong&#39;s heirs held the copyright to the painting, and that the museum only had the right to exhibit it.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-33">[33]</a><span> Joe McDonald of the </span><a href="https://en.wikipedia.org/wiki/Associated_Press">Associated Press</a><span> deemed the upholding of the copyright </span><span>“</span><span>a triumph for China&#39;s capitalist ambitions over its leftist history&#34;.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-ap-34">[34]</a><span> In 2014, the art museum at CAFA held a retrospective of Dong&#39;s works, exhibiting the small-scale draft of the painting, possessed by Dong&#39;s family, for the first time. </span><a href="https://en.wikipedia.org/wiki/Fan_Di%27an">Fan Di&#39;an</a><span>, curator of the exhibition, stated, </span><span>“</span><span>The changes to the painting tell a bitter story, reflecting the political influences on art. But it didn&#39;t affect Dong Xiwen&#39;s love of art.&#34;</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cd2014-35">[35]</a><span> </span></p><p><span>  Wu Hung described </span><em>The Founding of the Nation</em><span> as </span><span>“</span><span>arguably the most celebrated work of official Chinese art&#34;.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Hung_200869-36">[36]</a><span> He noted that the painting is the only </span><span>“</span><span>canonized</span><span>”</span><span> one depicting the October 1 ceremony, and that other artists have tended to give the people&#39;s perspective, subjecting themselves to Mao&#39;s gaze.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Hung_2005274-37">[37]</a><span> The painting is a modern-day example of </span><a href="https://en.wikipedia.org/wiki/Damnatio_memoriae"><em>damnatio memoriae</em></a><span>, the alteration of artworks or other objects to remove the image or name of a disfavored person.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEUnverzagt220-38">[38]</a><span> Deng Zhangyu, in a 2014 article, called the painting </span><span>“</span><span>the most significant historical image of China&#39;s founding&#34;.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cd2014-35">[35]</a><span> Wu Hung suggested that the alterations to it over the years, while always showing Mao proclaiming the new government, parallel the changes that have come to China&#39;s leadership during the years of Communist governance.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Hung_2005274-37">[37]</a><span> Andrews wrote that </span><span>“</span><span>its greatest importance to the art world was its elevation as a model of party-approved oil painting&#34;.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews77,_80-39">[39]</a><span> Writer Wu Bing in 2009 called it </span><span>“</span><span>a milestone in Chinese oil painting, boldly incorporating national styles&#34;.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEWu_Bing66-17">[17]</a><span> The painting has never been as highly regarded in the West as in China; according to Andrews, </span><span>“</span><span>art history students have been known to roar with laughter when slides of it appear on the screen&#34;.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews80-24">[24]</a><span> Art historian </span><a href="https://en.wikipedia.org/wiki/Michael_Sullivan_(art_historian)">Michael Sullivan</a><span>dismissed it as mere propaganda.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-FOOTNOTEAndrews80-24">[24]</a><span> Today, following a merger of museums, both paintings are in the </span><a href="https://en.wikipedia.org/wiki/National_Museum_of_China">National Museum of China</a><span>, on Tiananmen Square.</span><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-cn3-8">[8]</a><a href="https://en.wikipedia.org/wiki/The_Founding_Ceremony_of_the_Nation#cite_note-china_museum_1-40">[40]</a><span> </span></p><p><span> </span></p>`,
		},
		{
			name:     "complex rogue doc",
			visIndex: 0,
			selLen:   0,
			items: []v3.PasteItem{
				{Kind: "string", Mime: "text/html", Data: _complexRogueDoc},
			},
			expectedHtml: _complexRogueDoc + "<p></p>", // tack on the newline from the rogue
		},
		{
			name:     "list city",
			visIndex: 0,
			selLen:   0,
			items: []v3.PasteItem{
				{Kind: "string", Mime: "text/html", Data: "<h1><span>List Play</span></h1><ol><li><span>one</span></li><p><span>testing</span></p><ol><li><span>niceooh</span></li><p><span>hmmmmm</span></p><ol><li><span>cool</span></li><li><span>wow works great!</span></li><p><span>hello world</span></p><p><span>nice</span></p><li><span>cool beans</span></li></ol></ol><li><span>nice bullet</span></li><li><span>and number</span></li></ol>"},
			},
			expectedHtml: "<h1><span>List Play</span></h1><ol><li><span>one</span></li><p><span>testing</span></p><ol><li><span>niceooh</span></li><p><span>hmmmmm</span></p><ol><li><span>cool</span></li><li><span>wow works great!</span></li><p><span>hello world</span></p><p><span>nice</span></p><li><span>cool beans</span></li></ol></ol><li><span>nice bullet</span></li><li><span>and number</span></li></ol><p></p>",
		},
		{
			name:     "plaintext",
			visIndex: 0,
			selLen:   0,
			items: []v3.PasteItem{
				{Kind: "string", Mime: "text/plain", Data: "hello world"},
			},
			expectedHtml: "<p><span>hello world</span></p>",
		},
		{
			name: "paste link",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "hello world")
				require.NoError(t, err)
			},
			visIndex: 0,
			selLen:   11,
			items: []v3.PasteItem{
				{Kind: "string", Mime: "text/plain", Data: "https://www.google.com"},
			},
			expectedHtml: "<p><a href=\"https://www.google.com\">hello world</a></p>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("auth0")
			if tc.fn != nil {
				tc.fn(r)
			}

			_, _, err := r.Paste(tc.visIndex, tc.selLen, v3.FormatV3Span{}, tc.items)
			require.NoError(t, err)

			firstID, err := r.GetFirstID()
			require.NoError(t, err)

			lastID, err := r.GetLastID()
			require.NoError(t, err)

			html, err := r.GetHtml(firstID, lastID, false, true)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHtml, html)
		})
	}
}

var _appleNote = `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
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
</html>`

var _complexRogueDoc = `<h1><span>h1</span></h1><h2><span>h2</span></h2><h3><span>h3</span></h3><h4><span>h4</span></h4><h5><span>h5</span></h5><h6><span>h6</span></h6><ul><li><span>Soft paws on carpet</span></li><li><span>Whiskers twitch, eyes gleam with light</span></li><p><span>soft return 🙌🌎</span></p><li><span>Feline grace unfolds</span></li></ul><hr/><ol><li><span>one</span></li><li><span>two</span></li><ol><li><span>three</span></li><ul><li><span>four</span></li></ul></ol></ol><pre><code class="language-python" data-language="python"># this is some python
# with a multiple lines

for i in range(10):
	print(f&#34;hello {i}&#34;)
</code></pre><p></p><p><code>a code snippet</code></p><p><span>Some </span><strong>forma</strong><strong><s>tting</s></strong><s> </s><em><s>with</s></em><em> diffe</em><em><u>rent </u></em><u>stuff</u></p><blockquote><p><span>and a block quote about amazing things</span></p><p><span>neato!</span></p></blockquote><p></p>`
