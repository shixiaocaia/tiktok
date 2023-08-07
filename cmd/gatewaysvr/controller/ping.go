package controller

//func Feed(ctx *gin.Context) {
//	//todo 判断用户是否登录，先默认不登陆
//
//	// 根据时间戳获取视频列表
//	currentTime, err := strconv.ParseInt(ctx.Query("latest_time"), 10, 64)
//	if err != nil || currentTime == int64(0) {
//		currentTime = utils.GetCurrentTime()
//	}
//
//	userId, _ := ctx.Get("UserId")
//	tokenId := userId.(int64)
//
//	// 根据视频id获取该视频的作者信息，点赞信息，评论数
//
//	// 是否点赞，是否关注这个视频
//
//}

//func Ping(ctx *gin.Context) {
//	ctx.JSON(200, gin.H{
//		"message": config.GetGlobalConfig().Ping,
//	})
//}
//
//func Greet(ctx *gin.Context) {
//	resp, err := utils.GetGreeterClient().SayHello(ctx, &pb.HelloRequest{
//		Name: "tiktok",
//	})
//
//	if err != nil {
//		log.Error("Greet error", err)
//		ctx.JSON(http.StatusOK, gin.H{
//			"message": err.Error(),
//		})
//	}
//
//	ctx.JSON(200, gin.H{
//		"message": resp.Message,
//	})
//}
