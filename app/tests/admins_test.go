package tests

import (
	"fmt"
	json "github.com/json-iterator/go"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

func (s *APITestSuite) TestAdminCreateCourse() {

	r := s.Require()
	countTest := 1
	name := `{
    "tag_id":[1,2],
    "feature_id":6,
    "is_active": false,
    "content": {"title": "some_title", "text": "some_text", "url": "some_url"}
}`

	req, err := http.NewRequest("POST", "/banner", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp := httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)
	r.Equal(resp.Body.String(), "{\"banner_id\":1}")

	//---------2-----------
	req, err = http.NewRequest("POST", "/banner", strings.NewReader(name))
	printErr(err, &countTest)
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

	req, err = http.NewRequest("GET", "/user_banner?tag_id=1&&feature_id=6", strings.NewReader(name))
	printErr(err, &countTest)
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

	req, err = http.NewRequest("GET", "/user_banner?tag_id=1&&feature_id=6&&use_last_revision=true", strings.NewReader(name))
	printErr(err, &countTest)
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
	req, err = http.NewRequest("GET", "/user_banner?tag_id=1&&feature_id=1", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)

	//--------5------------
	req, err = http.NewRequest("GET", "/user_banner?tag_id=1", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)

	//--------6------------
	req, err = http.NewRequest("GET", "/user_banner?feature_id=1", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)

	//--------7------------

	req, err = http.NewRequest("GET", "/banner", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)

	//--------8------------
	req, err = http.NewRequest("GET", "/banner", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)

	//--------9------------
	req, err = http.NewRequest("GET", "/banner", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	//--------9------------
	req, err = http.NewRequest("PATCH", "/banner/1", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)

	//--------9------------
	req, err = http.NewRequest("PATCH", "/banner/2", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)

	//--------9------------
	name = `{
    "tag_id":[1,2],
    "feature_id":2,
    "is_active": false,
    "content": {"title": "some_1itle", "text": "some_text", "url": "some_url"}
}`

	req, err = http.NewRequest("PATCH", "/banner/1", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	//--------10------------

	req, err = http.NewRequest("GET", "/user_banner?tag_id=1&&feature_id=2&&use_last_revision=true", strings.NewReader(name))
	printErr(err, &countTest)
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
	req, err = http.NewRequest("PATCH", "/banner/2", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)

	//--------10------------
	req, err = http.NewRequest("PATCH", "/banner/2", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)

	//--------9------------
	req, err = http.NewRequest("DELETE", "/banner/2", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)

	//--------10------------
	req, err = http.NewRequest("DELETE", "/banner/2", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)

	//--------9------------
	req, err = http.NewRequest("DELETE", "/banner/2", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)

	//--------10------------
	req, err = http.NewRequest("DELETE", "/banner/s", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)

	//--------10------------
	req, err = http.NewRequest("DELETE", "/banner/1", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNoContent, resp.Result().StatusCode)

	//--------10------------

	name = `{
    "tag_id":[1,2],
    "feature_id":3,
    "is_active": false,
    "content": {"title": "some_title", "text": "some_text", "url": "some_url"}
}`

	req, err = http.NewRequest("POST", "/banner", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	r.Equal(resp.Body.String(), "{\"banner_id\":3}")

	name = `{
    "tag_id":[1,2],
    "feature_id":4,
    "is_active": false,
    "content": {"title": "some_title", "text": "some_text", "url": "some_url"}
}`

	req, err = http.NewRequest("POST", "/banner", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	r.Equal(resp.Body.String(), "{\"banner_id\":4}")

	name = `{
    "tag_id":[3],
    "feature_id":5,
    "is_active": false,
    "content": {"title": "some_title", "text": "some_text", "url": "some_url"}
}`

	req, err = http.NewRequest("POST", "/banner", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	r.Equal(resp.Body.String(), "{\"banner_id\":5}")

	name = `{
    "tag_id":[2],
    "feature_id":6,
    "is_active": false,
    "content": {"title": "some_title", "text": "some_text", "url": "some_url"}
}`

	req, err = http.NewRequest("POST", "/banner", strings.NewReader(name))
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	r.Equal(resp.Body.String(), "{\"banner_id\":6}")

	req, err = http.NewRequest("DELETE", "/banner/?tag_id=1", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)

	req, err = http.NewRequest("DELETE", "/banner/?tag_id=1", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)

	req, err = http.NewRequest("DELETE", "/banner/?tag_id=1", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNoContent, resp.Result().StatusCode)

	req, err = http.NewRequest("DELETE", "/banner/?feature_id=5", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNoContent, resp.Result().StatusCode)

	req, err = http.NewRequest("DELETE", "/banner/?feature_id=3", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)

	req, err = http.NewRequest("GET", "/banner/version?banner_id=1&&version=1", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	req, err = http.NewRequest("GET", "/banner/version?banner_id=1&&version=1", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)

	req, err = http.NewRequest("GET", "/banner/version?banner_id=1&&version=1", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "user_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)

	req, err = http.NewRequest("GET", "/banner/version?banner_id=1", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	req, err = http.NewRequest("GET", "/banner/version?banner_id=2", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)

	req, err = http.NewRequest("GET", "/banner/version?banner_id=s", nil)
	printErr(err, &countTest)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("token", "admin_token")

	resp = httptest.NewRecorder()
	s.router.Router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)

}
func printErr(err error, countTest *int) {
	if err != nil {
		fmt.Println("Ошибка запроса теста "+strconv.Itoa(*countTest)+":", err)
		return
	}
	*countTest++
}
