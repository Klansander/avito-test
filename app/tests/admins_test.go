package tests

import (
	"fmt"
	json "github.com/json-iterator/go"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (s *APITestSuite) TestAdminCreateCourse() {

	r := s.Require()

	name := fmt.Sprintf(`{
    "tag_id":[1,2],
    "feature_id":6,
    "is_active": false,
    "content": {"title": "some_title", "text": "some_text", "url": "some_url"}
}`)

	req, _ := http.NewRequest("POST", "/banner", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp := httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	//---------2-----------
	req, _ = http.NewRequest("POST", "/banner", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)

	//---------2--------------

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела ответа:", err)
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Ошибка чтения тела ответа:", err)
		return
	}

	req, _ = http.NewRequest("GET", "/user_banner?tag_id=1&&feature_id=6", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела ответа:", err)
		return
	}
	var res1 map[string]interface{}
	var res2 map[string]interface{}

	err = json.Unmarshal([]byte(`{
  "content": {
    "title": "some_title",
    "text": "some_text",
    "url": "some_url"
  }
}`), &res1)
	if err != nil {
		fmt.Println("Ошибка Unmarshal res1 :", err)
		return
	}
	err = json.Unmarshal(body, &res2)
	if err != nil {
		fmt.Println("Ошибка Unmarshal res2:", err)
		return
	}
	r.Equal(len(res1), len(res2))
	for k, v := range res1 {
		r.Equal(v, res2[k])
	}
	//--------3---------

	req, _ = http.NewRequest("GET", "/user_banner?tag_id=1&&feature_id=6&&use_last_revision=true", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела ответа:", err)
		return
	}
	clear(res2)

	err = json.Unmarshal(body, &res2)
	if err != nil {
		fmt.Println("Ошибка Unmarshal res2:", err)
		return
	}

	r.Equal(len(res1), len(res2))
	for k, v := range res1 {
		r.Equal(v, res2[k])
	}
	//--------4------------
	req, _ = http.NewRequest("GET", "/user_banner?tag_id=1&&feature_id=1", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)

	//--------5------------
	req, _ = http.NewRequest("GET", "/user_banner?tag_id=1", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)

	//--------6------------
	req, _ = http.NewRequest("GET", "/user_banner?feature_id=1", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)

	//--------7------------

	req, _ = http.NewRequest("GET", "/banner", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)

	//--------8------------
	req, _ = http.NewRequest("GET", "/banner", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)

	//--------9------------
	req, _ = http.NewRequest("GET", "/banner", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	//--------9------------
	req, _ = http.NewRequest("PATCH", "/banner/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)

	//--------9------------
	req, _ = http.NewRequest("PATCH", "/banner/2", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)

	//--------9------------
	name = fmt.Sprintf(`{
    "tag_id":[1,2],
    "feature_id":2,
    "is_active": false,
    "content": {"title": "some_1itle", "text": "some_text", "url": "some_url"}
}`)

	req, _ = http.NewRequest("PATCH", "/banner/1", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	//--------10------------

	req, _ = http.NewRequest("GET", "/user_banner?tag_id=1&&feature_id=2&&use_last_revision=true", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	clear(res1)

	err = json.Unmarshal([]byte(`{"content": {"title": "some_1itle", "text": "some_text", "url": "some_url"}}`), &res1)
	if err != nil {
		fmt.Println("Ошибка Unmarshal res2:", err)
		return
	}
	r.Equal(http.StatusOK, resp.Result().StatusCode)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела ответа:", err)
		return
	}
	clear(res2)

	err = json.Unmarshal(body, &res2)
	if err != nil {
		fmt.Println("Ошибка Unmarshal res2:", err)
		return
	}

	r.Equal(len(res1), len(res2))
	for k, v := range res1 {
		r.Equal(v, res2[k])
	}

	//--------9------------
	req, _ = http.NewRequest("PATCH", "/banner/2", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)

	//--------10------------
	req, _ = http.NewRequest("PATCH", "/banner/2", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)

	//--------9------------
	req, _ = http.NewRequest("DELETE", "/banner/2", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)

	//--------10------------
	req, _ = http.NewRequest("DELETE", "/banner/2", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)

	//--------9------------
	req, _ = http.NewRequest("DELETE", "/banner/2", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)

	//--------10------------
	req, _ = http.NewRequest("DELETE", "/banner/s", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)

	//--------10------------
	req, _ = http.NewRequest("DELETE", "/banner/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNoContent, resp.Result().StatusCode)

}
