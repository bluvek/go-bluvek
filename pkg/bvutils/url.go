package bvutils

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

// GetRequestPath 获取请求路径
func GetRequestPath(path, prefix string) (uri string, id int64) {
	uri = strings.TrimPrefix(path, prefix)
	re := regexp.MustCompile(`^(.*)/(\d+)$`)
	matches := re.FindStringSubmatch(uri)
	if len(matches) == 3 {
		uri = matches[1]
		id = cast.ToInt64(matches[2])
	}

	return
}

// ConvertToRestfulURL 将URI转换为REST ful URL
func ConvertToRestfulURL(url string) string {
	re := regexp.MustCompile(`(^.+?/[^/]+)/\d+$`)
	return re.ReplaceAllString(url, `$1/:id`)
}

// ConvertRestfulURLToUri 将REST ful URL转换为URI
func ConvertRestfulURLToUri(url string) (string, string) {
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return url, ""
	}
	last := parts[len(parts)-1]
	if strings.HasPrefix(last, ":") {
		path := strings.Join(parts[:len(parts)-1], "/")
		param := strings.TrimPrefix(last, ":")
		return path, param
	}

	return url, ""
}

// RemoveDomain 移除URL中的域名
func RemoveDomain(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	result := parsedURL.EscapedPath()
	if parsedURL.RawQuery != "" {
		result += "?" + parsedURL.RawQuery
	}

	if parsedURL.Fragment != "" {
		result += "#" + parsedURL.Fragment
	}

	return result
}

// JoinDomain 将域名和 URL 路径拼接为完整 URL
func JoinDomain(domain string, path string) string {
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}

	parsedDomain, err := url.Parse(domain)
	if err != nil {
		return ""
	}

	parsedPath, err := url.Parse(path)
	if err != nil {
		return ""
	}

	return parsedDomain.ResolveReference(parsedPath).String()
}
