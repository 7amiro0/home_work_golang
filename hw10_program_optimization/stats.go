package hw10programoptimization

import (
	"bufio"
	"github.com/mailru/easyjson"
	"io"
	"strings"
)

type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domain = strings.ToLower(domain)
	content := bufio.NewScanner(r)
	result := make(DomainStat)
	user := User{}

	for content.Scan() {
		if err := easyjson.Unmarshal(content.Bytes(), &user); err != nil {
			return nil, err
		} else if email := strings.SplitN(user.Email, "@", 2); len(email) == 2 {
			if strings.Contains(strings.ToLower(email[1]), domain) {
				result[strings.ToLower(email[1])]++
			}
		}
	}

	return result, nil
}
