package auth

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"
	"wq_submitter/configs"
)

var (
	BrainClient *BrainAuth
	once        sync.Once
	conf        *configs.GlobalConfig
)

func init() {
	conf = configs.GetGlobalConfig()
}

type user struct {
	ID string `json:"id"`
}

type token struct {
	Expiry float64 `json:"expiry"`
}

type loginResponse struct {
	User        user     `json:"user"`
	Token       token    `json:"token"`
	Permissions []string `json:"permissions"`
}
type BrainAuth struct {
	HttpClient  *http.Client
	expireTimer time.Timer
	mutex       sync.Mutex
}

func GetBrainAuth() *BrainAuth {
	once.Do(func() {
		BrainClient = newBrainAuth()
	})
	return BrainClient
}

func newBrainAuth() *BrainAuth {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Errorf("Failed to create cookie jar: %s", err.Error())
		return nil
	}

	brain := &BrainAuth{
		HttpClient: &http.Client{Jar: jar},
		mutex:      sync.Mutex{},
	}

	expireTimeNum, err := brain.login()
	if err != nil {
		log.Errorf("newBrainAuth Failed {%s}", err.Error())
		return nil
	}
	if expireTimeNum == -1 {
		log.Error("newBrainAuth Failed in Login")
		return nil
	}

	brain.expireTimer = *time.NewTimer(time.Duration(0.9*expireTimeNum) * time.Second)
	return brain
}

func (auth *BrainAuth) login() (float64, error) {

	username := conf.CredentialConfig.UserName
	password := conf.CredentialConfig.Password

	// 创建一个新的 POST 请求
	req, err := http.NewRequest("POST", "https://api.worldquantbrain.com/authentication", nil)
	if err != nil {
		return -1, fmt.Errorf("创建请求失败: %w", err)
	}
	// 设置基本认证
	req.SetBasicAuth(username, password)

	// 发送请求 (使用带有 cookie jar 的 httpClient)
	resp, err := auth.HttpClient.Do(req)
	if err != nil {
		log.Errorf("发送认证请求失败: %s", err.Error())
		return -1, fmt.Errorf("发送认证请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 打印响应状态码
	log.Infof("响应状态码: %d\n", resp.StatusCode)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("读取响应体失败: %s", err.Error())
	}

	if resp.StatusCode >= 400 {
		log.Errorf("Code: %d, Message: %s", resp.StatusCode, string(body))
		return -1, fmt.Errorf(string(body))
	}
	var responseData loginResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Errorf("解析Resp Body失败: %s", err.Error())
		return -1, fmt.Errorf("解析Resp Body失败: %w", err)
	}

	// 打印响应内容
	log.Infof("响应内容: %s\n", string(body))

	return responseData.Token.Expiry, nil

}
func (auth *BrainAuth) freshToken() error {
	auth.mutex.Lock()
	defer auth.mutex.Unlock()

	oldTimer := auth.expireTimer
	defer oldTimer.Stop()

	expireTimeNum, err := auth.login()
	if err != nil {
		log.Errorf("FreshBrainAuth Failed {%s}", err.Error())
		return err
	}
	auth.expireTimer = *time.NewTimer(time.Duration(0.9*expireTimeNum) * time.Second)
	return nil
}
func (auth *BrainAuth) CheckFreshToken() bool {

	select {
	case <-auth.expireTimer.C:
		err := auth.freshToken()
		if err != nil {
			log.Errorf("NeedFreshToken Failed {%s}", err.Error())
			return true
		}
	default:

	}
	return false
}
