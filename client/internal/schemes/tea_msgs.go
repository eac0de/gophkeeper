package schemes

type SaveTextDataMsg struct {
	Err        error
	StatusCode int
	TextData   UserTextData
}

type SaveFileDataMsg struct {
	Err        error
	StatusCode int
	FileData   UserFileData
}

type SaveBankCardMsg struct {
	Err        error
	StatusCode int
	BankCard   UserBankCard
}

type SaveAuthInfoMsg struct {
	Err        error
	StatusCode int
	AuthInfo   UserAuthInfo
}

type UpdateListItemMsg struct {
	ItemIndex int
	ActiveTab int
	Item      interface{}
}

type DeleteListItemMsg struct {
	ItemIndex  int
	ActiveTab  int
	Err        error
	StatusCode int
}

type GetUserAuthInfoListMsg struct {
	List       []UserAuthInfo
	Err        error
	StatusCode int
}

type GetUserFileDataListMsg struct {
	List       []UserFileData
	Err        error
	StatusCode int
}

type GetUserTextDataListMsg struct {
	List       []UserTextData
	Err        error
	StatusCode int
}

type GetUserBankCardListMsg struct {
	List       []UserBankCard
	Err        error
	StatusCode int
}

type DownloadFileMsg struct {
	Err        error
	StatusCode int
}
