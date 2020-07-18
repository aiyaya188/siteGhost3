package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	//"finace-bank-receipt/utils"

	"github.com/pquerna/otp/totp"
)

//Create2FAQR 生成2fa二维码
//params 账号,域名
//return secret秘钥,base64url,error
func Create2FAQR(account string) (string, string, error) {
	issuer := "后台管理系统"
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
	})
	if err != nil {
		fmt.Println("创建用户2FA Secret失败", err.Error())
		return "", "", err
	}
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		fmt.Println("创建用户2FA Secret失败", err.Error())
		return "", "", err
	}
	err = png.Encode(&buf, img)
	if err != nil {
		fmt.Println("创建用户2FA Secret失败", err.Error())
		return "", "", err
	}
	pngbase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return key.Secret(), pngbase64, nil
}

//ValidTotp 校验totp验证码
func ValidTotp(passcord, secret string) bool {
	return totp.Validate(passcord, secret)
}

//GetTotp 获取用户TOTP二维码
func GetTotp(userName string) (map[string]string, error) {
	//新增totp secret code
	secret, png, err := Create2FAQR(userName)
	if err != nil {
		return nil, err
	}
	//创建totp二维码
	totp := make(map[string]string)
	totp["qrcode"] = "data:image/png;base64," + png
	totp["secret"] = secret
	return totp, nil
}
