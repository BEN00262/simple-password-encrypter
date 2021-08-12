package main

import "encoding/json"

type Credential struct {
	Target   string
	Password string
}

func parseJSONFile(fileData []byte) []Credential {
	var credentials []Credential

	json.Unmarshal(fileData, &credentials)

	return credentials
}

func searchCredential(Target string, credentials []Credential) string {

	for i := 0; i < len(credentials); i++ {
		if credentials[i].Target == Target {
			return credentials[i].Password
		}
	}

	return ""
}

func insertCredential(credential Credential, credentials []Credential) []Credential {
	// check if the credential exists first of all if it does just update it :)
	for i := 0; i < len(credentials); i++ {
		if credentials[i].Target == credential.Target {
			credentials[i].Password = credential.Password
			return credentials
		}
	}

	credentials = append(credentials, credential)
	return credentials
}

func removeCredential(target string, credentials []Credential) []Credential {
	for i := 0; i < len(credentials); i++ {
		if credentials[i].Target == target {
			return append(credentials[:i], credentials[i+1:]...)
		}
	}
	return credentials
}

func dumpToJSONFile(credentials *[]Credential) ([]byte, error) {
	convertedBytes, err := json.Marshal(*credentials)

	if err != nil {
		return nil, err
	}

	return convertedBytes, nil
}
