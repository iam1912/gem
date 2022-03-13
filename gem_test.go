package gem

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupEngine() {
	r := New()
	r.LoadHTMLGlob("./testdata/*")
	r.GET("/", func(c *Context) {
		c.String(200, "HELLO WORLD")
	})
	r.GET("/index", func(c *Context) {
		c.HTML(200, "index.html", H{
			"title": "HELLO WORLD",
		})
	})
	go r.Run(":9090")
}

type EngineTestCase struct {
	Url          string
	ExpectedResp string
}

func TestEngine(t *testing.T) {
	setupEngine()

	testCases := []EngineTestCase{
		{Url: "http://localhost:9090/", ExpectedResp: "HELLO WORLD"},
		{Url: "http://localhost:9090/index", ExpectedResp: `<p>HELLO WORLD</p>`},
	}
	for i, testCase := range testCases {
		resp, err := http.Get(testCase.Url)
		if err != nil {
			t.Errorf("TestEngine #%v get error %v\n", i+1, err.Error())
		}
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		if !assert.Equal(t, testCase.ExpectedResp, string(data)) {
			t.Errorf("TestEngine #%v Expected get %v but get %v\n", i+1, testCase.ExpectedResp, string(data))
		}
	}
}
