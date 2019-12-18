package apnian

import (
	"encoding/json"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

const alert = "Bawk bawk, bitches! 🐔"
const sound = "default"
const linkUrl = "https://electronicpanopticon.com"

func TestAPS(t *testing.T) {
	t.Run("GenerateAPS()", func(t *testing.T) {
		expected := testAPS()

		aps := GenerateAPS(alert, sound, linkUrl)

		assert.Equal(t, expected, aps)
	})

	t.Run("ToJsonBytes()", func(t *testing.T) {
		expected := testAPS()

		b := expected.ToJsonBytes()
		var aps APS
		err := json.Unmarshal(b, &aps)

		assert.Nil(t, err)
		assert.Equal(t, expected, &aps)
	})
}

func testAPS() *APS {
	payload := Payload{
		Alert:   alert,
		Sound:   sound,
		LinkUrl: linkUrl,
	}
	return &APS{
		APS: payload,
	}
}

//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\
// Benchmark!

var result *APS

func benchmarkGenerateAPSFact(b *testing.B, num int) {
	var aps *APS
	for x := 0; x < b.N; x++ {
		alert := faker.Word()
		sound := faker.URL()
		linkUrl := faker.URL()
		aps = GenerateAPS(alert, sound, linkUrl)

	}
	result = aps
}

func BenchmarkFact20(b *testing.B) {
	benchmarkGenerateAPSFact(b, 20)
}