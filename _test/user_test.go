package _test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	UserID     int    `json:"user_id"`
	Token      string `json:"token"`
}

type Response2 struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	User       User   `json:"user"`
}

type User struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	FollowCount     int    `json:"follow_count"`
	FollowerCount   int    `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  int    `json:"total_favorited"`
	WorkCount       int    `json:"work_count"`
	FavoriteCount   int    `json:"favorite_count"`
}

var url = "124.220.156.104:8005"

func TestRegister(t *testing.T) {
	res, err := http.Post("http://124.220.156.104:8005/douyin/user/login/?username=test_xiaocai&password=test", "POST", nil)
	assert.Nil(t, err)
	var response Response
	returnedJson, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(returnedJson, &response)
	assert.Nil(t, err)
	assert.Equal(t, 0, response.StatusCode)
	assert.NotNil(t, response.UserID)
	assert.NotNil(t, response.Token)
}

func TestLogin(t *testing.T) {
	res, err := http.Post("http://124.220.156.104:8005/douyin/user/login/?username=test&password=123456", "POST", nil)
	assert.Nil(t, err)
	var response Response
	returnedJson, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(returnedJson, &response)
	assert.Nil(t, err)
	assert.Equal(t, 0, response.StatusCode)
	assert.NotNil(t, response.UserID)
	assert.NotNil(t, response.Token)
}

func TestGetUserInfo(t *testing.T) {
	res, err := http.Get("http://124.220.156.104:8005/douyin/user/?user_id=20044&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJ0ZXN0IiwidXNlcl9pZCI6MjAwNDQsImlzcyI6InNlcnZlciJ9.JCMk9hJ8LtaAk0EDzyImeX_wm_pUQAjTYywaP3AaR9E")
	assert.Nil(t, err)
	var response Response2
	returnedJson, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(returnedJson, &response)
	assert.Nil(t, err)
	assert.Equal(t, 0, response.StatusCode)
	assert.NotNil(t, response.User)
	assert.Equal(t, 20044, response.User.ID)
	assert.Equal(t, "test", response.User.Name)
	assert.Equal(t, 0, response.User.FollowCount)
	assert.Equal(t, 0, response.User.FollowerCount)
	assert.Equal(t, false, response.User.IsFollow)
	assert.NotNil(t, response.User.Avatar)
	assert.NotNil(t, response.User.BackgroundImage)
	assert.NotNil(t, response.User.Signature)
	assert.Equal(t, 0, response.User.TotalFavorited)
	assert.Equal(t, 0, response.User.WorkCount)
	assert.Equal(t, 0, response.User.FavoriteCount)
}
