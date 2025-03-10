package dag_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fivetentaylor/pointy/pkg/dag"
	"github.com/fivetentaylor/pointy/pkg/testutils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
	"github.com/fivetentaylor/pointy/rogue/v3/testcases"
)

const essay = `Lifelong learning is essential in today’s rapidly changing world, where new technologies and evolving industries constantly reshape the job market. Engaging in continuous education allows individuals to stay relevant and competitive, ensuring they possess the necessary skills to adapt to new challenges. Moreover, lifelong learning fosters personal growth and intellectual fulfillment, contributing to a more satisfying and meaningful life.

One significant advantage of lifelong learning is the ability to enhance career prospects and job security. As industries advance, employees who continuously update their knowledge and skills are more likely to secure promotions and navigate career transitions successfully. Additionally, employers value workers who demonstrate a commitment to personal development, often leading to increased opportunities and professional recognition. Beyond the workplace, acquiring new skills can open doors to diverse hobbies and interests, enriching one’s personal life.

Furthermore, lifelong learning promotes critical thinking and problem-solving abilities, which are crucial in both professional and personal settings. By engaging with new ideas and perspectives, individuals become more adaptable and resilient in the face of change. This ongoing intellectual stimulation also contributes to mental well-being, reducing the risk of cognitive decline as one ages. In essence, lifelong learning not only supports economic and career advancement but also enhances overall quality of life, making it a vital pursuit for individuals of all ages.`

func TestChunkDocument_simple_three_paragraphs(t *testing.T) {
	doc := v3.NewRogueForQuill("0")

	_, err := doc.Insert(0, essay)
	if err != nil {
		t.Fatal(err)

	}

	chunks, err := dag.ChunkDocument(doc, 100)
	require.NoError(t, err)
	require.Equal(t, 3, len(chunks))

	require.Equal(t,
		"Lifelong learning is essential in today’s rapidly changing world, where new technologies and evolving industries constantly reshape the job market. Engaging in continuous education allows individuals to stay relevant and competitive, ensuring they possess the necessary skills to adapt to new challenges. Moreover, lifelong learning fosters personal growth and intellectual fulfillment, contributing to a more satisfying and meaningful life.\n\n",
		chunks[0].Markdown,
	)

	require.Equal(t,
		"One significant advantage of lifelong learning is the ability to enhance career prospects and job security. As industries advance, employees who continuously update their knowledge and skills are more likely to secure promotions and navigate career transitions successfully. Additionally, employers value workers who demonstrate a commitment to personal development, often leading to increased opportunities and professional recognition. Beyond the workplace, acquiring new skills can open doors to diverse hobbies and interests, enriching one’s personal life.\n\n",
		chunks[1].Markdown,
	)

	require.Equal(t,
		"Furthermore, lifelong learning promotes critical thinking and problem-solving abilities, which are crucial in both professional and personal settings. By engaging with new ideas and perspectives, individuals become more adaptable and resilient in the face of change. This ongoing intellectual stimulation also contributes to mental well-being, reducing the risk of cognitive decline as one ages. In essence, lifelong learning not only supports economic and career advancement but also enhances overall quality of life, making it a vital pursuit for individuals of all ages.\n\n",
		chunks[2].Markdown,
	)
}

const oneReallyBigParagraph = `Lifelong learning is essential in today’s rapidly changing world, where new technologies and evolving industries constantly reshape the job market. Engaging in continuous education allows individuals to stay relevant and competitive, ensuring they possess the necessary skills to adapt to new challenges. Moreover, lifelong learning fosters personal growth and intellectual fulfillment, contributing to a more satisfying and meaningful life. One significant advantage of lifelong learning is the ability to enhance career prospects and job security. As industries advance, employees who continuously update their knowledge and skills are more likely to secure promotions and navigate career transitions successfully. Additionally, employers value workers who demonstrate a commitment to personal development, often leading to increased opportunities and professional recognition. Beyond the workplace, acquiring new skills can open doors to diverse hobbies and interests, enriching one’s personal life.Furthermore, lifelong learning promotes critical thinking and problem-solving abilities, which are crucial in both professional and personal settings. By engaging with new ideas and perspectives, individuals become more adaptable and resilient in the face of change. This ongoing intellectual stimulation also contributes to mental well-being, reducing the risk of cognitive decline as one ages. In essence, lifelong learning not only supports economic and career advancement but also enhances overall quality of life, making it a vital pursuit for individuals of all ages.`

func TestChunkDocument_oneReallyBigParagraph(t *testing.T) {
	doc := v3.NewRogueForQuill("0")

	_, err := doc.Insert(0, oneReallyBigParagraph)
	if err != nil {
		t.Fatal(err)

	}

	chunks, err := dag.ChunkDocument(doc, 100)
	require.NoError(t, err)
	require.Equal(t, 1, len(chunks))

	require.Equal(t,
		"Lifelong learning is essential in today’s rapidly changing world, where new technologies and evolving industries constantly reshape the job market. Engaging in continuous education allows individuals to stay relevant and competitive, ensuring they possess the necessary skills to adapt to new challenges. Moreover, lifelong learning fosters personal growth and intellectual fulfillment, contributing to a more satisfying and meaningful life. One significant advantage of lifelong learning is the ability to enhance career prospects and job security. As industries advance, employees who continuously update their knowledge and skills are more likely to secure promotions and navigate career transitions successfully. Additionally, employers value workers who demonstrate a commitment to personal development, often leading to increased opportunities and professional recognition. Beyond the workplace, acquiring new skills can open doors to diverse hobbies and interests, enriching one’s personal life.Furthermore, lifelong learning promotes critical thinking and problem-solving abilities, which are crucial in both professional and personal settings. By engaging with new ideas and perspectives, individuals become more adaptable and resilient in the face of change. This ongoing intellectual stimulation also contributes to mental well-being, reducing the risk of cognitive decline as one ages. In essence, lifelong learning not only supports economic and career advancement but also enhances overall quality of life, making it a vital pursuit for individuals of all ages.\n\n",
		chunks[0].Markdown,
	)
}

const extendedEssay = `Lifelong learning is essential in today’s rapidly changing world. With new technologies and evolving industries, the job market is constantly being reshaped. To stay relevant, individuals must adapt by continuously educating themselves.

Engaging in lifelong learning helps people remain competitive in their careers. By keeping their skills up to date, they can adapt to new challenges more effectively. This ensures they are always prepared for changes in their field.

A key benefit of lifelong learning is enhanced career prospects. Employees who invest in their education are more likely to receive promotions. They are also better positioned to navigate career transitions successfully.

Employers value workers who show a commitment to personal development. This dedication often leads to more opportunities and greater professional recognition. It demonstrates a proactive approach to growth and improvement.

Beyond career advancement, lifelong learning enriches personal life as well. Acquiring new skills can lead to new hobbies and interests. These pursuits can add joy and satisfaction to everyday life.

Lifelong learning also sharpens critical thinking and problem-solving abilities. These skills are valuable in both professional and personal settings. By engaging with new ideas, individuals can better handle unexpected situations.

Exposure to different perspectives fosters adaptability. Lifelong learners are more resilient when faced with change. This flexibility is essential in an unpredictable world.

Intellectual stimulation through learning benefits mental well-being. It keeps the mind active and reduces the risk of cognitive decline. This is particularly important as individuals grow older.

Lifelong learning contributes to a more fulfilling life overall. It provides a sense of achievement and purpose. By constantly growing, individuals can lead a more meaningful and rewarding life.

In essence, lifelong learning supports both economic and personal advancement. It enhances career opportunities and enriches one's quality of life. Therefore, it remains a vital pursuit for people of all ages.`

func TestChunkDocument_extendedEssay(t *testing.T) {
	doc := v3.NewRogueForQuill("0")

	_, err := doc.Insert(0, extendedEssay)
	require.NoError(t, err)

	chunks, err := dag.ChunkDocument(doc, 100)
	require.NoError(t, err)
	require.Equal(t, 5, len(chunks))

	require.Equal(t,
		"Lifelong learning is essential in today’s rapidly changing world. With new technologies and evolving industries, the job market is constantly being reshaped. To stay relevant, individuals must adapt by continuously educating themselves.\n\n\n\nEngaging in lifelong learning helps people remain competitive in their careers. By keeping their skills up to date, they can adapt to new challenges more effectively. This ensures they are always prepared for changes in their field.\n\n",
		chunks[0].Markdown,
	)

	require.Equal(t,
		"Lifelong learning contributes to a more fulfilling life overall. It provides a sense of achievement and purpose. By constantly growing, individuals can lead a more meaningful and rewarding life.\n\n\n\nIn essence, lifelong learning supports both economic and personal advancement. It enhances career opportunities and enriches one&#39;s quality of life. Therefore, it remains a vital pursuit for people of all ages.\n\n",
		chunks[4].Markdown,
	)
}

func TestChunkDocument_with_checks(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()

	err := os.Chdir("../..")
	require.NoError(t, err)

	checks, err := dag.ListFuncationalChecks(ctx, "threadv2")
	require.NoError(t, err)

	for _, check := range checks {
		t.Run(check.CheckName, func(t *testing.T) {
			var doc v3.Rogue
			err := json.Unmarshal([]byte(check.SerializedRogue), &doc)
			require.NoError(t, err)

			chunks, err := dag.ChunkDocument(&doc, 250)
			require.NoError(t, err)

			for i, chunk := range chunks {
				fmt.Printf("%d: %q\n", i, chunk.Markdown)
			}

			if len(chunks) == 0 {
				t.Fatal("no chunks")
			}

			mkdown, err := doc.GetFullMarkdown()
			require.NoError(t, err)

			var chunkMarkdowns []string
			for _, chunk := range chunks {
				// fmt.Printf("%d: %q\n", i, chunk)
				chunkMarkdowns = append(chunkMarkdowns, chunk.Markdown)
				mkdown = strings.Replace(mkdown, strings.TrimSpace(chunk.Markdown), "", 1)
			}

			require.Equal(t, "", strings.TrimSpace(mkdown))
			// require.Equal(t, "", strings.TrimSpace(mkdown), fmt.Sprintf("%+v", chunkMarkdowns))

		})
	}

}

func TestChunkDocument_realworld_large_doc(t *testing.T) {
	doc := testcases.Load(t, "large_doc.json")

	chunks, err := dag.ChunkDocument(doc, 2048)
	require.NoError(t, err)

	for i, chunk := range chunks {
		fmt.Printf("%d: %q\n", i, chunk.Markdown)
	}

	require.Equal(t, 3, len(chunks))

	require.Equal(t,
		"When the Museum of the Chinese Revolution opened on Tiananmen Square in 1961, the painting was displayed on a huge wall in the gallery devoted to the Communist triumph, but in 1966, during the Cultural Revolution, radicals shut down the museum, and it remained closed until 1969.[28] During that time, Chinese president Liu Shaoqi, accused of taking a &#34;capitalist road&#34;, was purged from government. His removal from the painting was ordered in 1967, and Dong was tasked to carry it out.[29] Dong had suffered during the Cultural Revolution: accused of being a rightist, he was expelled from the party for two years, sent to a rural work camp, and then was rehabilitated by being made to labor as a steelworker.[7] Dong&#39;s task was difficult, as Liu was one of the most prominent figures in the first row, standing to the left of Madame Song. Officials wanted Liu replaced with Lin Biao, much in Mao&#39;s favor at the time. Dong was unwilling to give Lin prominence he had not then had, and though he could not refuse outright at the dangerous time of the Cultural Revolution, he eventually got permission to merely remove Liu. The figure was too large to simply delete, so Liu was repainted as Dong Biwu, and made to appear as if in the second row. According to Andrews, the attempt was a failure: &#34;the new Dong Biwu does not recede into the second row as intended. Instead, he appears as a leering, glowing figure, a strangely malevolent character in the midst of an otherwise stately group&#34;.[30] Officials deemed the revised work unexhibitable. Andrews speculated that Dong may have been trying to sabotage the change, or may have been affected by the stress of the years of the Cultural Revolution.[30]\n\n\n\nIn 1972, as part of a renovation of the Museum of the Chinese Revolution, officials wanted to exhibit Dong&#39;s painting again, but they decreed that Lin Boqu, whose white-haired figure was furthest left, must be removed.[30] This was because the Gang of Four, then in control of China, blamed Lin Boqu (who had died in 1960) for opposing the marriage, in the revolutionary days, of Mao with Jiang Qing (one of the Four). Sources differ on what took place regarding the painting: Chang-Tai Hung related that Dong, terminally ill with cancer, could not make the changes, so his student Jin Shangyi and another artist, Zhao Yu, were assigned to do the work. The two feared damaging the original canvas, so made an exact replica but for the required changes, with Dong brought forth from his hospital for consultations.[29] According to Andrews, Jin and Zhao created the new version because Dong would not let anyone else alter his painting.[30] Jin later stated that the painting, while effective politically, also shows Dong&#39;s inner world.[18]\n\n\n\nWith the end of the Cultural Revolution in 1976 and the subsequent accession of Deng Xiaoping, many of the purged figures of earlier years were rehabilitated, and the authorities in 1979 decided to bring more historical accuracy to the painting. Dong had died in 1973; his family strongly opposed anyone altering the original painting, and the government respected their wishes. Jin was on a tour outside China, so the government assigned Yan Zhenduo to make changes to the replica. He placed Liu, Lin Boqu and Gao in the painting[29] and made other changes: a previously unidentifiable man in the back row now resembles Deng Xiaoping. The replica painting was restored to the Museum of the Chinese Revolution.[31]\n\n\n\n\n\nVisitors take photographs of The Founding Ceremony of the Nation, National Museum of China, in 2018.\n\nThe painting was reproduced on Chinese postage stamps in 1959 and 1999, for the tenth and fiftieth anniversaries of the founding of the People&#39;s Republic.[32] Also in 1999, the museum authorized a private company to make small-scale gold foil reproductions of the painting. Dong&#39;s family sued, and in 2002 the courts found that Dong&#39;s heirs held the copyright to the painting, and that the museum only had the right to exhibit it.[33] Joe McDonald of the Associated Press deemed the upholding of the copyright &#34;a triumph for China&#39;s capitalist ambitions over its leftist history&#34;.[34] In 2014, the art museum at CAFA held a retrospective of Dong&#39;s works, exhibiting the small-scale draft of the painting, possessed by Dong&#39;s family, for the first time. Fan Di&#39;an, curator of the exhibition, stated, &#34;The changes to the painting tell a bitter story, reflecting the political influences on art. But it didn&#39;t affect Dong Xiwen&#39;s love of art.&#34;[35]\n\n\n\nWu Hung described The Founding of the Nation as &#34;arguably the most celebrated work of official Chinese art&#34;.[36] He noted that the painting is the only &#34;canonized&#34; one depicting the October 1 ceremony, and that other artists have tended to give the people&#39;s perspective, subjecting themselves to Mao&#39;s gaze.[37] The painting is a modern-day example of damnatio memoriae, the alteration of artworks or other objects to remove the image or name of a disfavored person.[38] Deng Zhangyu, in a 2014 article, called the painting &#34;the most significant historical image of China&#39;s founding&#34;.[35] Wu Hung suggested that the alterations to it over the years, while always showing Mao proclaiming the new government, parallel the changes that have come to China&#39;s leadership during the years of Communist governance.[37] Andrews wrote that &#34;its greatest importance to the art world was its elevation as a model of party-approved oil painting&#34;.[39] Writer Wu Bing in 2009 called it &#34;a milestone in Chinese oil painting, boldly incorporating national styles&#34;.[17] The painting has never been as highly regarded in the West as in China; according to Andrews, &#34;art history students have been known to roar with laughter when slides of it appear on the screen&#34;.[24] Art historian Michael Sullivan dismissed it as mere propaganda.[24] Today, following a merger of museums, both paintings are in the National Museum of China, on Tiananmen Square.[8][40]\n\n",
		chunks[2].Markdown,
	)
}
