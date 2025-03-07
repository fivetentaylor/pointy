package utils_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/pkg/utils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func TestFindSimilarSubstringsUint16(t *testing.T) {
	t.Run("find similar substrings", func(t *testing.T) {
		s := "Hello I am a doc, Jello I am a doc"
		pattern := "hello"

		substrings := utils.FindSimilarSubstringsUint16(
			v3.StrToUint16(s),
			v3.StrToUint16(pattern),
			1,
		)

		require.Equal(t, substrings, []int{0, 18})
	})

	t.Run("no similar substrings", func(t *testing.T) {
		s := "Hello I am a doc, Jello I am a doc"
		pattern := "goodbye"

		substrings := utils.FindSimilarSubstringsUint16(
			v3.StrToUint16(s),
			v3.StrToUint16(pattern),
			1,
		)

		require.Equal(t, substrings, []int{})
	})

	t.Run("empty compare", func(t *testing.T) {
		s := "Hello I am a doc, Jello I am a doc"
		pattern := ""

		substrings := utils.FindSimilarSubstringsUint16(
			v3.StrToUint16(s),
			v3.StrToUint16(pattern),
			1,
		)

		require.Equal(t, substrings, []int{})
	})

	t.Run("rw example", func(t *testing.T) {
		s := colleensDoc
		pattern := colleensSelection

		length := len(pattern)
		md := int(math.Floor(0.2*float64(length))) + 1

		substrings := utils.FindSimilarSubstringsUint16(
			v3.StrToUint16(s),
			v3.StrToUint16(pattern),
			md,
		)

		for i := 0; i < len(substrings); i++ {
			substring := s[substrings[i] : substrings[i]+length]
			fmt.Println(substring)
		}

		require.Len(t, substrings, 1)
	})
}

const colleensDoc = `Designing Tools for AI: Navigating the Organic Landscape of Generative AI
The Challenge of Designing for AI
In the rapidly evolving world of technology, designing tools for artificial intelligence (AI) presents a unique set of challenges. Unlike traditional software development, where structures and inputs are more concrete, AI—particularly generative AI—introduces a level of fluidity and unpredictability that demands a new approach to design.
The Shift from Structured to Organic
Traditional Software Design:
Relied on structured databases
Clear input/output relationships
Predictable form fields and data entry
AI and Generative AI Design:
Organic and fluid nature
Less control over inputs and outputs
Conversational user experiences add complexity
Key Challenges in AI Tool Design
Unpredictable Outputs: Generative AI can produce a wide range of results, making it difficult to design interfaces that accommodate all possibilities.
Conversational Interfaces: Natural language interactions require more flexible and adaptive user interfaces.
Balancing User Control and AI Autonomy: Determining the right level of user input versus AI-driven decisions is crucial.
Explaining AI Decision-Making: Designing interfaces that make AI processes transparent and understandable to users.
Adapting to Rapid AI Advancements: Tools must be flexible enough to incorporate new AI capabilities as they emerge.
Strategies for Effective AI Tool Design
User-Centric Approach: Focus on user needs and expectations when interacting with AI systems.
Iterative Design Process: Embrace frequent testing and refinement to address the organic nature of AI interactions.
Flexible UI Components: Design modular interface elements that can adapt to various AI outputs.
Clear Communication: Implement effective ways to convey AI capabilities and limitations to users.
Ethical Considerations: Incorporate ethical guidelines and safeguards into the design process.
The Future of AI Tool Design
As AI continues to evolve, designers must stay adaptable and innovative. The challenge lies in creating intuitive, powerful tools that harness the potential of AI while maintaining a cohesive and satisfying user experience.
By embracing the organic nature of AI and focusing on user-centered design principles, we can create tools that not only meet the current needs but also pave the way for future advancements in human-AI interaction.
`

const colleensSelection = `Key Challenges in AI Tool Design
**Unpredictable Outputs: **Generative AI can produce a wide range of results, making it difficult to design interfaces that accommodate all possibilities.
**Conversational Interfaces: **Natural language interactions require more flexible and adaptive user interfaces.
**Balancing User Control and AI Autonomy: **Determining the right level of user input versus AI-driven decisions is crucial.
**Explaining AI Decision-Making: **Designing interfaces that make AI processes transparent and understandable to users.
**Adapting to Rapid AI Advancements: **Tools must be flexible enough to incorporate new AI capabilities as they emerge.`
