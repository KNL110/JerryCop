package router

import (
	configparser "JerryCop/internal/configParser"
	"net/url"
	"path"
	"strings"
)

// MatchPrefix checks if the request path matches the given prefix
func MatchPrefix(reqPath, prefix string) bool {
	if prefix == "" || prefix == "/" {
		return true
	}
	prefix = strings.TrimSuffix(prefix, "/")
	return reqPath == prefix || strings.HasPrefix(reqPath, prefix+"/")
}

// FindMatchingServer finds the ServerConfig with the longest matching prefix.
func FindMatchingServer(reqPath string, servers []configparser.ServerConfig) (*configparser.ServerConfig, bool) {
	var bestMatch *configparser.ServerConfig
	longest := -1

	for i := range servers {
		prefix := servers[i].ServerRouting.Prefix
		if MatchPrefix(reqPath, prefix) && len(prefix) > longest {
			longest = len(prefix)
			bestMatch = &servers[i]
		}
	}

	return bestMatch, bestMatch != nil
}

// BuildTargetURL joins the upstream base server URL with the request path
// and preserves the original query parameters.
func BuildTargetURL(upstreamBase, reqPath, rawQuery string) (string, error) {
	targetURL, err := url.Parse(upstreamBase)
	if err != nil {
		return "", err
	}

	targetURL.Path = path.Join(targetURL.Path, reqPath)
	if strings.HasSuffix(reqPath, "/") && !strings.HasSuffix(targetURL.Path, "/") {
		targetURL.Path += "/"
	}
	targetURL.RawQuery = rawQuery

	return targetURL.String(), nil
}
