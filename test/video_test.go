package test

import "net/http"
import "testing"

// 测试 “视频流” 功能
func TestFeed(t *testing.T) {
	e := newExpect(t)
	feedResp := e.GET("/douyin/feed/").Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_msg").String().IsEqual("success")
	feedResp.Value("video_list").Array().Length().Gt(0)
	for _, element := range feedResp.Value("video_list").Array().Iter() {
		video := element.Object()
		video.ContainsKey("id")
		video.ContainsKey("author")
		video.Value("play_url").String().NotEmpty()
		video.Value("cover_url").String().NotEmpty()
	}
}

// 测试 “视频投稿、发布列表” 功能
func TestPublishAndList(t *testing.T) {
	e := newExpect(t)

	// 随机得到一个用户名以及密码
	var publish_username = generateRandomUsername()
	var publish_password = generateRandomPassword()

	// 获取测试用户的 userId 和 token
	userId, token := getTestUserIdAndToken(publish_username, publish_password, e)

	// 发布视频请求并进行断言
	publishResp := e.POST("/douyin/publish/action/").
		WithMultipart().
		WithFile("data", "test.mp4"). // 添加要上传的视频文件路径
		WithFormField("token", token).
		WithFormField("title", "用户视频上传测试"). // 设置视频标题
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// 验证发布响应的状态码
	publishResp.Value("status_msg").String().IsEqual("success")

	// 获取发布列表请求并进行断言
	publishListResp := e.GET("/douyin/publish/list/").
		WithQuery("user_id", userId).WithQuery("token", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// 验证发布列表响应的状态码
	publishResp.Value("status_msg").String().IsEqual("success")

	// 验证视频列表不为空并遍历每个视频元素
	publishListResp.Value("video_list").Array().Length().Gt(0)

	for _, element := range publishListResp.Value("video_list").Array().Iter() {
		video := element.Object()

		// 验证视频对象中是否包含指定字段
		video.ContainsKey("id")
		video.ContainsKey("author")

		// 验证视频播放链接和封面链接不为空
		video.Value("play_url").String().NotEmpty()
		video.Value("cover_url").String().NotEmpty()
	}
}
