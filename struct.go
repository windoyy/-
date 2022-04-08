package main
type classInfo struct {
	Result int `json:"result"`
	Msg interface{} `json:"msg"`
	Data struct {
		Ext struct {
			From string `json:"_from_"`
		} `json:"ext"`
		ReadingDuration int `json:"readingDuration"`
		ActiveList []struct {
			UserStatus int `json:"userStatus"`
			GroupID int `json:"groupId"`
			IsLook int `json:"isLook"`
			Type int `json:"type"`
			ReleaseNum int `json:"releaseNum"`
			Content string `json:"content,omitempty"`
			AttendNum int `json:"attendNum"`
			ActiveType int `json:"activeType"`
			Logo string `json:"logo"`
			NameOne string `json:"nameOne"`
			StartTime int64 `json:"startTime"`
			ID int `json:"id"`
			EndTime int64 `json:"endTime"`
			Status int `json:"status"`
			NameFour string `json:"nameFour"`
			ExtraInfo struct {
				NoticeID int `json:"noticeId"`
			} `json:"extraInfo,omitempty"`
			NameTwo string `json:"nameTwo,omitempty"`
			OtherID string `json:"otherId,omitempty"`
			Source int `json:"source,omitempty"`
		} `json:"activeList"`
	} `json:"data"`
	ErrorMsg interface{} `json:"errorMsg"`
}

type tokenRes struct {
	Result bool `json:"result"`
	Token string `json:"_token"`
}

type uploadPicRes struct {
	Result bool `json:"result"`
	Msg string `json:"msg"`
	Puid int `json:"puid"`
	Data struct {
		Preview string `json:"preview"`
		Filetype string `json:"filetype"`
		PreviewURL string `json:"previewUrl"`
		Suffix string `json:"suffix"`
		Resid int64 `json:"resid"`
		Duration int `json:"duration"`
		Pantype string `json:"pantype"`
		Puid int `json:"puid"`
		Filepath string `json:"filepath"`
		Crc string `json:"crc"`
		Isfile bool `json:"isfile"`
		Residstr string `json:"residstr"`
		ObjectID string `json:"objectId"`
		Extinfo string `json:"extinfo"`
		Thumbnail string `json:"thumbnail"`
		Creator int `json:"creator"`
		ModifyDate int64 `json:"modifyDate"`
		ResTypeValue int `json:"resTypeValue"`
		DisableOpt bool `json:"disableOpt"`
		Sort int `json:"sort"`
		Topsort int `json:"topsort"`
		Restype string `json:"restype"`
		Size int `json:"size"`
		UploadDate int64 `json:"uploadDate"`
		Name string `json:"name"`
	} `json:"data"`
	Crc string `json:"crc"`
	Resid int64 `json:"resid"`
	ObjectID string `json:"objectId"`
}