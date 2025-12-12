package utilx

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

var nonWord = regexp.MustCompile(`[^a-z0-9\-]`)

// Slugify 将字符串转为唯一标识 + 六位 Hash 结果(防止冲突） eg: GoLang学习 -> golang-xue-xi-1fcbf8
func Slugify(name string) string {
	name = strings.Trim(name, " ") // 去空格
	name = strings.ToLower(name)   // 转小写

	args := pinyin.NewArgs()

	// 遍历判断每一个 rune 是汉字还是字母
	var sb strings.Builder
	for k, r := range []rune(name) {
		switch {
		case unicode.Is(unicode.Han, r): // 为汉字
			py := pinyin.Pinyin(string(r), args)
			if len(py) > 0 && len(py[0]) > 0 {
				if k != 0 {
					sb.WriteRune('-')
				}
				sb.WriteString(py[0][0])
			}
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			sb.WriteRune(r)
		}
	}
	slug := nonWord.ReplaceAllString(sb.String(), "")

	h := sha1.Sum([]byte(slug))
	hash := hex.EncodeToString(h[:])[:6]

	return slug + "-" + hash
}
