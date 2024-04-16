package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func transformationCriteria(inputData map[string]interface{}) map[string]interface{} {
	resultMap := make(map[string]interface{})
	for key, value := range inputData {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if _, ok := value.(map[string]interface{}); !ok {
			continue
		}
		tranforedValue := criteriaChecks(value.(map[string]interface{}))
		if tranforedValue != nil {
			resultMap[key] = tranforedValue
		}
	}
	return resultMap
}

func criteriaChecks(value map[string]interface{}) interface{} {
	for key := range value {
		temp := value[key]
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		value[key] = temp
	}
	for key, val := range value {
		key = strings.TrimSpace(key)
		switch key {
		case "S":
			return transformString(val.(string))
		case "N":
			return transformNumeric(val.(string))
		case "BOOL":
			return transformBool(val.(string))
		case "NULL":
			nullValue := strings.TrimSpace(val.(string))
			if nullValue == "" {
				return nil
			}
			if _, ok := val.(string); !ok {
				continue
			}
			return transformNull(val.(string))
		case "L":
			if _, ok := val.([]interface{}); !ok {

				return nil
			}
			res := transformList(val.([]interface{}))
			if res != nil {
				return res
			}
		case "M":
			return transformMap(value)
		default:
			continue
		}
	}
	return nil
}

func transformMap(inputMap map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range inputMap {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if valMap, ok := value.(map[string]interface{}); ok {
			for k, val := range valMap {
				if transformedValue := criteriaChecks(val.(map[string]interface{})); transformedValue != nil {
					result[k] = transformedValue
				}
			}
		}
	}
	return result
}

func transformList(list []interface{}) []interface{} {
	var result []interface{}
	for _, item := range list {
		if _, ok := item.(map[string]interface{}); !ok {
			continue
		}
		if transformedItem := criteriaChecks(item.(map[string]interface{})); transformedItem != nil {
			result = append(result, transformedItem)
		}

	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func transformNull(nullValue string) interface{} {
	nullValue = strings.TrimSpace(nullValue)
	if nullValue == "1" ||
		strings.ToLower(nullValue) == "true" || strings.ToLower(nullValue) == "t" {
		nullValue = "null"

		return nullValue
	} else {
		return nil
	}

}

func transformBool(boolValue string) interface{} {
	boolValue = strings.TrimSpace(boolValue)
	if boolValue == "" {
		return nil
	}
	result, err := strconv.ParseBool(boolValue)
	if err != nil {
		return nil
	}
	return result
}
func transformNumeric(numericValue string) interface{} {
	numericValue = strings.TrimSpace(numericValue)
	numericValue = strings.TrimLeft(numericValue, "0")

	if numericValue == "" {
		return nil
	}
	if result, err := strconv.ParseFloat(numericValue, 64); err == nil {
		return result
	}
	if result, err := strconv.ParseInt(numericValue, 10, 64); err == nil {
		return result
	}

	return nil
}

func transformString(stringValue string) interface{} {
	stringValue = strings.TrimSpace(stringValue)
	if stringValue == "" {
		return nil
	}
	if timestamp, err := time.Parse(time.RFC3339, stringValue); err == nil {
		return timestamp.Unix()
	}
	return stringValue
}

func main() {
	start := time.Now()
	inputData, err := os.ReadFile("input.json")
	if err != nil {
		fmt.Println("unable to read file due to %s\n", err.Error())
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(inputData, &jsonData)
	if err != nil {
		fmt.Println("unable to Unmarshall : %v", err.Error())
	}
	resOutput := transformationCriteria(jsonData)
	result, err := json.Marshal(resOutput)
	if err != nil {
		fmt.Println("unable to marshal : %v", err.Error())
	}

	fmt.Println(string(result))
	elapsed := time.Since(start)
	fmt.Println("Execution time : ", elapsed)
}
