package handlers

import (
	"bbsgo/utils"
	"net/http"
)

func GetQiniuUploadToken(w http.ResponseWriter, r *http.Request) {
	token := utils.GetQiniuUploadToken()
	utils.Success(w, map[string]string{
		"token": token,
	})
}
