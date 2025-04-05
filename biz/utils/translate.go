package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"meeting_agent/conf"
	"strings"

	alimt20181012 "github.com/alibabacloud-go/alimt-20181012/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var LanguageMap = map[uint32]string{
	0: "zh",
	1: "en",
	2: "ja",
	3: "ko",
	4: "fr",
	5: "de",
	6: "es",
	7: "pt",
	8: "ru",
}

func CreateClient() (_result *alimt20181012.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(conf.GetConf().Aliyun.AccessKey),
		AccessKeySecret: tea.String(conf.GetConf().Aliyun.AccessSecret),
	}
	config.Endpoint = tea.String("mt.ap-southeast-1.aliyuncs.com")
	_result = &alimt20181012.Client{}
	_result, _err = alimt20181012.NewClient(config)
	return _result, _err
}

func Translate(text string, lang uint32) (translated string, err error) {
	client, err := CreateClient()
	if err != nil {
		return "", fmt.Errorf("创建翻译客户端失败: %w", err)
	}

	translateGeneralRequest := &alimt20181012.TranslateGeneralRequest{
		FormatType:     tea.String("text"),
		SourceLanguage: tea.String("auto"),
		TargetLanguage: tea.String(LanguageMap[lang]),
		Scene:          tea.String("general"),
		SourceText:     tea.String(text),
	}
	runtime := &util.RuntimeOptions{}

	s, tryErr := func() (t string, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, err := client.TranslateGeneralWithOptions(translateGeneralRequest, runtime)
		if err != nil {
			return "", fmt.Errorf("调用翻译API失败: %w", err)
		}
		return *result.Body.Data.Translated, nil
	}()

	if tryErr != nil {
		var e = &tea.SDKError{}
		var _t *tea.SDKError
		if errors.As(tryErr, &_t) {
			e = _t
		}
		hlog.Error("翻译服务错误: %v", tea.StringValue(e.Message))

		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(e.Data)))
		if err := d.Decode(&data); err != nil {
			return "", fmt.Errorf("解析错误信息失败: %w", err)
		}

		if m, ok := data.(map[string]interface{}); ok {
			if recommend, ok := m["Recommend"].(string); ok {
				hlog.Error("建议: %s", recommend)
			}
		}
		return "", fmt.Errorf("翻译失败: %s", tea.StringValue(e.Message))
	}
	return s, nil
}
