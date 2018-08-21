package models

import (
	"errors"
	"regexp"

	"github.com/jmoiron/sqlx"
)

//CheckRule to check match between rule array and receipt
//Return map type and any wrire error encountered
// kiểm tra xem một receipt khi nhận được có trùng khớp với bất kì một rule trong bảng rule không
// Nếu có thì trả về một map kết quả nếu không thì trả về err
func CheckRule(receipt string, rules []string) (map[string]string, error) {
	mapKey := make(map[string]string)
	for i := 0; i < len(rules); i++ {
		re, err := regexp.Compile(rules[i])
		if err != nil {
			return nil, err
		}
		matched := re.FindStringSubmatch(receipt)
		if matched != nil {
			keyName := re.SubexpNames()
			for k, v := range keyName {
				if v != "" {
					mapKey[v] = matched[k]
				}
			}
			return mapKey, nil
		}

	}
	return nil, errors.New("Not pattern matchh")
}

// LoadRegex load all rule from database
// Retuen rule array and any error
// Load tất các các rule trong bảng rule lêns
func LoadRegex(db *sqlx.DB) ([]string, error) {
	var arrayRegex []string
	err := db.Select(&arrayRegex, "SELECT regex FROM rule")
	if err != nil {
		return nil, err
	}
	return arrayRegex, nil
}
