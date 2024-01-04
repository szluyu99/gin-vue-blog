package ginblog

import (
	"bytes"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCookieSessions(t *testing.T) {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/hello", func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("hello") != "world" {
			session.Set("hello", "world")
			session.Save()
		}

		c.JSON(200, gin.H{"hello": session.Get("hello")})
	})

	c := NewTestClient(r)
	{
		w := c.Get("/hello")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"hello":"world"}`, w.Body.String())
	}

}

// func TestRedisSessions(t *testing.T) {
// 	r := gin.Default()
// 	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
// 	r.Use(sessions.Sessions("mysession", store))

// 	r.GET("/incr", func(c *gin.Context) {
// 		session := sessions.Default(c)
// 		var count int
// 		v := session.Get("count")
// 		if v == nil {
// 			count = 0
// 		} else {
// 			count = v.(int)
// 			count++
// 		}
// 		session.Set("count", count)
// 		session.Save()
// 		c.JSON(200, gin.H{"count": count})
// 	})

// 	c1 := NewTestClient(r)
// 	{
// 		w := c1.Get("/incr")
// 		assert.Equal(t, http.StatusOK, w.Code)
// 		assert.Equal(t, `{"count":0}`, w.Body.String())

// 		w = c1.Get("/incr")
// 		assert.Equal(t, http.StatusOK, w.Code)
// 		assert.Equal(t, `{"count":1}`, w.Body.String())
// 	}

// 	c2 := NewTestClient(r)
// 	{
// 		w := c2.Get("/incr")
// 		assert.Equal(t, http.StatusOK, w.Code)
// 		assert.Equal(t, `{"count":0}`, w.Body.String())

// 		w = c2.Get("/incr")
// 		assert.Equal(t, http.StatusOK, w.Code)
// 		assert.Equal(t, `{"count":1}`, w.Body.String())
// 	}
// }

type TestClient struct {
	r         http.Handler
	CookieJar http.CookieJar
	Scheme    string
	Host      string
}

func NewTestClient(r http.Handler) (c *TestClient) {
	jar, _ := cookiejar.New(nil)
	return &TestClient{
		r:         r,
		CookieJar: jar,
		Scheme:    "http",
		Host:      "1.2.3.4",
	}
}

func (c *TestClient) SendReq(path string, req *http.Request) *httptest.ResponseRecorder {
	req.URL.Scheme = "http"
	req.URL.Host = "MOCK_SERVER"
	req.RemoteAddr = "127.0.0.1:1234"

	currentUrl := &url.URL{
		Scheme: c.Scheme,
		Host:   c.Host,
		Path:   path,
	}

	cookies := c.CookieJar.Cookies(currentUrl)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	w := httptest.NewRecorder()
	c.r.ServeHTTP(w, req)
	c.CookieJar.SetCookies(currentUrl, w.Result().Cookies())
	return w
}

// Get return *httptest.ResponseRecorder
func (c *TestClient) Get(path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	return c.SendReq(path, req)
}

// Post return *httptest.ResponseRecorder
func (c *TestClient) Post(path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return c.SendReq(path, req)
}
