package urlutil

import (
	"net/url"
	"path"
	"sort"
	"strings"
)

func Normalize(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	parsed.Scheme = strings.ToLower(parsed.Scheme)
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", nil
	}

	parsed.Host = strings.ToLower(parsed.Host)

	parsed.Host = removeDefaultPort(parsed.Host, parsed.Scheme)

	parsed.Fragment = ""

	if parsed.Path == "" {
		parsed.Path = "/"
	} else {
		parsed.Path = path.Clean(parsed.Path)

		if parsed.Path == "." {
			parsed.Path = "/"
		}
	}

	if parsed.RawQuery != "" {
		parsed.RawQuery = sortQuery(parsed.RawQuery)
	}

	return parsed.String(), nil
}

func ResolveReference(base, ref string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	refURL, err := url.Parse(ref)
	if err != nil {
		return "", err
	}
	resolved := baseURL.ResolveReference(refURL)
	return Normalize(resolved.String())
}

func ExtractDomain(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := parsed.Hostname()
	return strings.ToLower(host)
}

func IsValidCrawlURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return false
	}

	if parsed.Host == "" {
		return false
	}

	ext := strings.ToLower(path.Ext(parsed.Path))
	skipExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".svg": true, ".webp": true, ".ico": true, ".bmp": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true,
		".xlsx": true, ".ppt": true, ".pptx": true,
		".zip": true, ".tar": true, ".gz": true, ".rar": true,
		".mp3": true, ".mp4": true, ".avi": true, ".mov": true,
		".wmv": true, ".flv": true, ".wav": true, ".ogg": true,
		".exe": true, ".dmg": true, ".apk": true,
		".css": true, ".js": true, ".json": true, ".xml": true,
		".woff": true, ".woff2": true, ".ttf": true, ".eot": true,
	}

	return !skipExts[ext]
}

func removeDefaultPort(host, scheme string) string {
	if strings.HasSuffix(host, ":80") && scheme == "http" {
		return strings.TrimSuffix(host, ":80")
	}
	if strings.HasSuffix(host, ":443") && scheme == "https" {
		return strings.TrimSuffix(host, ":443")
	}
	return host
}

func sortQuery(rawQuery string) string {
	params, err := url.ParseQuery(rawQuery)
	if err != nil {
		return rawQuery
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		vals := params[k]
		sort.Strings(vals)
		for _, v := range vals {
			parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
	}

	return strings.Join(parts, "&")
}
