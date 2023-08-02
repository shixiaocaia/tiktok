package service

import (
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

// grpc的具体逻辑

type UserService struct {
	pb.UnimplementedUserServiceServer
}

//// UpdateUserFavoritedCount 更新用户 获赞数，ActionType 1：表示+1 2：-1
//func (u UserService) UpdateUserFavoritedCount(ctx context.Context, req *pb.UpdateUserFavoritedCountReq) (*pb.UpdateUserFavoritedCountRsp, error) {
//	err := repository.UpdateUserFavoritedNum(req.UserId, req.ActionType)
//	if err != nil {
//		log.Errorf("UpdateUserFavoritedCount err", req.UserId)
//		return nil, err
//	}
//	return &pb.UpdateUserFavoritedCountRsp{}, nil
//}
//
//// UpdateUserFollowCount 更新用户 喜爱的视频数，ActionType 1：表示+1 2：-1
//func (u UserService) UpdateUserFavoriteCount(ctx context.Context, req *pb.UpdateUserFavoriteCountReq) (*pb.UpdateUserFavoriteCountRsp, error) {
//	err := repository.UpdateUserFavoriteNum(req.UserId, req.ActionType)
//	if err != nil {
//		log.Errorf("UpdateUserFavoriteCount err", req.UserId)
//		return nil, err
//	}
//	return &pb.UpdateUserFavoriteCountRsp{}, nil
//}
//
//// UpdateUserFollowCount 更新用户 关注数，ActionType 1：表示+1 2：-1
//func (u UserService) UpdateUserFollowCount(ctx context.Context, req *pb.UpdateUserFollowCountReq) (*pb.UpdateUserFollowCountRsp, error) {
//	err := repository.UpdateUserFollowNum(req.UserId, req.ActionType)
//	if err != nil {
//		log.Errorf("UpdateUserFollowCount err", req.UserId)
//		return nil, err
//	}
//	return &pb.UpdateUserFollowCountRsp{}, nil
//}
//
//// UpdateUserFollowerCount 更新用户 粉丝数，ActionType 1：表示+1 2：-1
//func (u UserService) UpdateUserFollowerCount(ctx context.Context, req *pb.UpdateUserFollowerCountReq) (*pb.UpdateUserFollowerCountRsp, error) {
//	err := repository.UpdateUserFollowerNum(req.UserId, req.ActionType)
//	if err != nil {
//		log.Errorf("UpdateUserFollowerCount err", req.UserId)
//		return nil, err
//	}
//	return &pb.UpdateUserFollowerCountRsp{}, nil
//}
//
//func (u UserService) GetUserInfoDict(ctx context.Context, req *pb.GetUserInfoDictRequest) (*pb.GetUserInfoDictResponse, error) {
//	userList, err := repository.GetUserList(req.UserIdList)
//	if err != nil {
//		log.Errorf("GetUserInfoDict err", req.UserIdList)
//		return nil, err
//	}
//	resp := &pb.GetUserInfoDictResponse{UserInfoDict: make(map[int64]*pb.UserInfo)}
//
//	for _, user := range userList {
//		resp.UserInfoDict[user.Id] = &pb.UserInfo{
//			Id:              user.Id,
//			Name:            user.Name,
//			Avatar:          user.Avatar,
//			FollowCount:     user.Follow,
//			FollowerCount:   user.Follower,
//			BackgroundImage: user.BackgroundImage,
//			Signature:       user.Signature,
//			TotalFavorited:  user.TotalFav,
//			FavoriteCount:   user.FavCount,
//		}
//	}
//
//	return resp, nil
//}
//
//func (u UserService) CacheChangeUserCount(ctx context.Context, req *pb.CacheChangeUserCountReq) (*pb.CacheChangeUserCountRsp, error) {
//	uid := strconv.FormatInt(req.UserId, 10)
//	mutex := lock.GetLock("user_" + uid)
//	defer lock.UnLock(mutex)
//	user, err := repository.CacheGetUser(req.UserId)
//	if err != nil {
//		log.Infof("CacheChangeUserCount err", req.UserId)
//		return nil, err
//	}
//
//	switch req.CountType {
//	case "follow":
//		user.Follow += req.Op
//	case "follower":
//		user.Follower += req.Op
//	case "like":
//		user.FavCount += req.Op
//	case "liked":
//		user.TotalFav += req.Op
//	}
//	repository.CacheSetUser(user)
//
//	return &pb.CacheChangeUserCountRsp{}, nil
//}
//
//func (u UserService) CacheGetAuthor(ctx context.Context, req *pb.CacheGetAuthorReq) (*pb.CacheGetAuthorRsp, error) {
//	key := strconv.FormatInt(req.VideoId, 10)
//	data, err := repository.CacheHGet("video", key)
//	if err != nil {
//		log.Errorf("CacheGetAuthor err", req.VideoId)
//		return nil, err
//	}
//
//	uid := int64(0)
//	err = json.Unmarshal(data, &uid)
//	if err != nil {
//		return nil, err
//	}
//
//	return &pb.CacheGetAuthorRsp{UserId: uid}, nil
//}
//
//func (u UserService) GetUserInfoList(ctx context.Context, req *pb.GetUserInfoListRequest) (*pb.GetUserInfoListResponse, error) {
//	response := new(pb.GetUserInfoListResponse)
//	log.Infof("GetUserInfoList req", req.IdList)
//	for _, userId := range req.IdList {
//		info, err := repository.GetUserInfo(userId)
//		if err != nil {
//			log.Errorf("GetUserInfoList err", userId, err)
//			return nil, err
//		}
//		response.UserInfoList = append(response.UserInfoList, UserToUserInfo(info))
//	}
//
//	return response, nil
//}
//
//func (u UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
//	user, err := repository.GetUserInfo(req.Id)
//	if err != nil {
//		return nil, err
//	}
//	response := &pb.GetUserInfoResponse{
//		UserInfo: UserToUserInfo(user),
//	}
//
//	return response, nil
//}
//
//func (u UserService) CheckPassWord(ctx context.Context, req *pb.CheckPassWordRequest) (*pb.CheckPassWordResponse, error) {
//	info, err := repository.GetUserInfo(req.Username)
//	if err != nil {
//		return nil, err
//	}
//	// 验证密码是否正确
//	err = bcrypt.CompareHashAndPassword([]byte(info.Password), []byte(req.Password))
//	if err != nil {
//		return nil, errors.New("password error")
//	}
//	token, err := GenToken(info.Id, req.Username)
//	if err != nil {
//		return nil, err
//	}
//	response := &pb.CheckPassWordResponse{
//		UserId: info.Id,
//		Token:  token,
//	}
//	return response, nil
//}
//
//func (u UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
//	sign, err := repository.UserNameIsExist(req.Username)
//	if err != nil {
//		log.Error("UserNameIsExist err ", err)
//		return nil, err
//	}
//	if sign {
//		return nil, fmt.Errorf("user %s exists", req.Username)
//	}
//	info, err := repository.InsertUser(req.Username, req.Password)
//	if err != nil {
//		return nil, err
//	}
//	token, err := GenToken(info.Id, req.Username)
//	if err != nil {
//		return nil, err
//	}
//	registerResponse := &pb.RegisterResponse{
//		UserId: info.Id,
//		Token:  token,
//	}
//
//	return registerResponse, nil
//}
