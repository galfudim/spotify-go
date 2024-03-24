package main

import (
	"encoding/json"
	. "spotify-go/common/spotify"
)

func getScopeString() string {
	allScopes := GetAllScopes()
	scopeString := ""
	for i := 0; i < len(allScopes); i++ {
		if i != len(allScopes)-1 {
			scopeString += string(allScopes[i]) + "+"
		} else {
			scopeString += string(allScopes[i])
		}
	}

	return scopeString
}

func jsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil
	}
	return result
}
