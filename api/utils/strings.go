package utils

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	rand2 "math/rand"

	"golang.org/x/crypto/sha3"
)

// RandString generate rand string with specified length
func RandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	data := []byte(str)
	var result []byte
	r := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	for range length {
		result = append(result, data[r.Intn(len(data))])
	}
	return string(result)
}

func RandomNumber(bit int) int {
	minNum := intPow(10, bit-1)
	maxNum := intPow(10, bit) - 1

	rand2.NewSource(time.Now().UnixNano())
	return rand2.Intn(maxNum-minNum+1) + minNum
}

func intPow(x, y int) int {
	result := 1
	for i := 0; i < y; i++ {
		result *= x
	}
	return result
}

func Contains(slice []string, item string) bool {
	for _, e := range slice {
		if e == item {
			return true
		}
	}
	return false
}

// Stamp2str 时间戳转字符串
func Stamp2str(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// Str2stamp 字符串转时间戳
func Str2stamp(str string) int64 {
	if len(str) == 0 {
		return 0
	}

	var layout string
	if strings.Contains(str, "T") {
		layout = "2006-01-02T15:04:05-07:00"
	} else {
		if len(str) < 12 {
			str = str + " 00:00:00"
		}
		layout = "2006-01-02 15:04:05"
	}

	t, err := time.ParseInLocation(layout, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

func GenPassword(pass string, salt string) string {
	data := []byte(pass + salt)
	hash := sha3.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

func JsonEncode(value interface{}) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func JsonDecode(src string, dest interface{}) error {
	return json.Unmarshal([]byte(src), dest)
}

func InterfaceToString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return JsonEncode(value)
}

// CutWords 截取前 N 个单词
func CutWords(str string, num int) string {
	// 按空格分割字符串为单词切片
	words := strings.Fields(str)

	// 如果单词数量超过指定数量，则裁剪单词；否则保持原样
	if len(words) > num {
		return strings.Join(words[:num], " ") + " ..."
	} else {
		return str
	}
}

// HasChinese 判断文本是否含有中文
func HasChinese(text string) bool {
	for _, char := range text {
		if unicode.Is(unicode.Scripts["Han"], char) {
			return true
		}
	}
	return false
}

func GenRedeemCode(codeLength int) (string, error) {
	bytes := make([]byte, codeLength/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsValidEmail 检查给定的字符串是否是有效的电子邮件地址
func IsValidEmail(email string) bool {
	// 这个正则表达式匹配大多数常见的邮箱格式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidMobile 检查给定的字符串是否是有效的中国大陆手机号
func IsValidMobile(phone string) bool {
	// 支持 13x, 14x, 15x, 16x, 17x, 18x, 19x 开头的号码
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}

// TruncateString 截断字符串
func TruncateString(str string, length int) string {
	if utf8.RuneCountInString(str) > length {
		return string([]rune(str)[:length]) + "..."
	}
	return str
}

// generateNickname 生成有趣的中文昵称，包含武侠、玄幻角色名
func GenerateNickname() string {
	adjectives := []string{
		"神秘的", "快乐的", "暴躁的", "优雅的", "闪亮的", "呆萌的", "机智的", "勇敢的", "慵懒的", "疯狂的",
		"潇洒的", "逍遥的", "无敌的", "绝世的", "孤独的", "冷酷的", "热血的", "天真的", "霸气的", "深情的",
	}
	nouns := []string{
		// 金庸
		"令狐冲", "杨过", "小龙女", "张无忌", "黄蓉", "郭靖", "段誉", "萧峰", "虚竹", "王重阳",
		"东方不败", "任我行", "周芷若", "赵敏", "韦小宝", "阿青", "独孤求败", "黄药师", "欧阳锋", "洪七公",
		// 古龙
		"小李飞刀", "李寻欢", "楚留香", "陆小凤", "花满楼", "西门吹雪", "叶孤城", "阿飞", "孙小红", "铁中棠",
		// 玄幻
		"萧炎", "药老", "美杜莎", "林动", "小貂", "唐三", "小舞", "比比东", "霍雨浩", "王林", "李穆婉",
		"孟浩", "韩立", "王宝乐", "石昊", "叶凡", "荒天帝", "无始大帝", "狠人大帝", "女帝", "帝尊",
		"炎帝", "辰东", "李七夜", "李长生", "徐凤年", "李淳罡", "陈长生", "苏离", "苏醒", "林雷", "李穆婉",
		//
		"郭靖", "黄蓉", "杨过", "小龙女", "张无忌", "赵敏", "乔峰", "段誉", "虚竹", "王语嫣", "令狐冲", "任盈盈", "东方不败", "韦小宝", "陈家洛", "胡斐", "苗人凤", "李寻欢", "楚留香", "陆小凤", "西门吹雪", "叶孤城", "沈浪", "熊猫儿", "萧十一郎", "连城璧", "张丹枫", "云蕾", "厉胜男", "金世遗", "狄云", "丁典", "袁承志", "温青青", "霍青桐", "文泰来", "谢逊", "周芷若", "灭绝师太", "洪七公", "欧阳锋", "黄药师", "一灯大师", "周伯通", "梅超风", "阿朱", "阿紫", "萧炎", "林动", "唐三", "霍雨浩", "韩立", "徐凤年", "宁缺", "陈平安", "石昊", "叶凡", "秦羽", "林雷", "萧玄", "药尘", "绫清竹", "小舞", "唐舞麟", "南宫婉", "姜泥", "桑桑", "宁姚", "云曦", "姬紫月", "姜立", "迪莉娅", "方寒", "牧尘", "洛璃", "白小纯", "陈长生", "唐三十六", "苏铭", "王林", "萧晨", "林辰", "龙辰", "叶锋", "沈翔", "凌尘", "东邪", "西毒", "南帝", "北丐", "中神通", "小李飞刀", "盗帅", "陆小凤", "玉面飞狐", "雪山飞狐", "金蛇郎君", "铁掌水上漂", "西狂", "小东邪", "光明左使", "光明右使", "紫衫龙王", "白眉鹰王", "金毛狮王", "青翼蝠王", "毒手药王", "潇湘夜雨", "君子剑", "淑女剑", "玉面孟尝", "九指神丐", "北侠", "中顽童", "赤练仙子", "混元霹雳手", "炎帝", "武祖", "海神", "大主宰", "荒天帝", "叶天帝", "剑尊", "血刀老祖", "万佛之主", "妖帝", "剑神", "刀神", "医仙", "李白", "杜甫", "苏轼", "辛弃疾", "屈原", "司马迁", "诸葛亮", "曹操", "刘备", "关羽", "张飞", "岳飞", "文天祥", "王阳明", "朱元璋", "王羲之", "颜真卿", "柳公权", "张旭", "吴道子", "李时珍", "孙思邈", "华佗", "扁鹊", "祖冲之", "张衡", "蔡伦", "毕昇", "沈括",
	}

	r := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	adj := adjectives[r.Intn(len(adjectives))]
	noun := nouns[r.Intn(len(nouns))]
	number := RandomNumber(6)
	return fmt.Sprintf("%s%s@%d", adj, noun, number)
}

// GenerateAvatar 生成随机头像
func GenerateAvatar() string {
	r := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("/images/avatar/user%02d.png", r.Intn(40)+1)
}
