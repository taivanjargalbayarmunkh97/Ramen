package utils

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func Filter(filters []FilterObj) func(db *gorm.DB) *gorm.DB {
	where := ""
	for _, filter := range filters {
		operation := "="
		if filter.Operation != "" {
			operation = filter.Operation
		}
		if filter.Value == "" && strings.ToLower(filter.Operation) != "in" && strings.ToLower(filter.Operation) != "notin" {
			continue
		}
		if strings.ContainsAny(filter.FieldName, ";") || strings.ContainsAny(operation, ";") || strings.ContainsAny(filter.Value, ";") {
			continue
		}

		if filter.Values != nil {
			filter.Value = ""
			for _, s := range filter.Values {
				if strings.ToLower(filter.FieldType) == "number" {
					filter.Value = filter.Value + fmt.Sprintf("%v,", s)
				} else {
					filter.Value = filter.Value + fmt.Sprintf("'%s',", s)
				}
			}
			filter.Value = filter.Value[:len(filter.Value)-1]
		}

		if where == "" {
			where = fmt.Sprintf("%s ", FilterColumn(filter.FieldName, filter.FieldType, operation, filter.Value))
		} else {
			where = fmt.Sprintf("%s AND %s ", where, FilterColumn(filter.FieldName, filter.FieldType, operation, filter.Value))
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where(where)
	}
}

func FilterWheres(filters []FilterObj) string {
	where := ""
	for _, filter := range filters {
		operation := "="
		if filter.Operation != "" {
			operation = filter.Operation
		}
		if filter.Value == "" && strings.ToLower(filter.Operation) != "in" && strings.ToLower(filter.Operation) != "notin" {
			continue
		}
		if strings.ContainsAny(filter.FieldName, ";") || strings.ContainsAny(operation, ";") || strings.ContainsAny(filter.Value, ";") {
			continue
		}

		if filter.Values != nil {
			filter.Value = ""
			for _, s := range filter.Values {
				if strings.ToLower(filter.FieldType) == "number" {
					filter.Value = filter.Value + fmt.Sprintf("%v,", s)
				} else {
					filter.Value = filter.Value + fmt.Sprintf("'%s',", s)
				}
			}
			filter.Value = filter.Value[:len(filter.Value)-1]
		}

		if where == "" {
			where = fmt.Sprintf("%s ", FilterColumn(filter.FieldName, filter.FieldType, operation, filter.Value))
		} else {
			where = fmt.Sprintf("%s AND %s ", where, FilterColumn(filter.FieldName, filter.FieldType, operation, filter.Value))
		}
	}

	return where
}

func FilterColumn(columnName string, FieldType string, operation string, value string) string {
	FieldType = strings.ToLower(FieldType)
	operation = strings.ToLower(operation)

	if FieldType == "text" {
		return TextFilter(columnName, operation, value)
	}
	if FieldType == "number" {
		return NumberFilter(columnName, operation, value)
	}
	if FieldType == "date" {
		return DateFilter(columnName, operation, value)
	}
	if operation == "in" {
		return fmt.Sprintf("%s in %s%s%s ", columnName, "(", value, ")")
	}
	if FieldType == "uuid" {
		return UuidFilter(columnName, operation, value)
	}
	return TextFilter(columnName, operation, value)
}

func TextFilter(columnName string, operation string, value string) string {
	operation = strings.ToLower(operation)
	if operation == "like" {
		return fmt.Sprintf("LOWER(%s) %s LOWER('%s%s%s') ", columnName, operation, "%", value, "%")
	}
	if operation == "in" {
		return fmt.Sprintf("%s in %s%s%s ", columnName, "(", value, ")")
	}
	if operation == "notin" {
		return fmt.Sprintf("LOWER(%s) not in %s%s%s ", columnName, "(", value, ")")
	}
	return fmt.Sprintf("LOWER(%s) %s LOWER('%s') ", columnName, operation, value)
}

func UuidFilter(columnName string, operation string, value string) string {
	operation = strings.ToLower(operation)
	if operation == "=" {
		return fmt.Sprintf("%s %s '%s'", columnName, operation, value)
	}
	if operation == "in" {
		return fmt.Sprintf("%s in %s%s%s ", columnName, "(", value, ")")
	}
	if value == "null" {
		return fmt.Sprintf("%s %s null", columnName, operation)
	}
	return fmt.Sprintf("%s %s '%s'", columnName, operation, value)
}

func NumberFilter(columnName string, operation string, value string) string {
	value = strings.ToLower(value)
	if operation == "between" {
		if value == "low" {
			return fmt.Sprintf("%s BETWEEN 0 AND 69 ", columnName)
		} else if value == "medium" {
			return fmt.Sprintf("%s BETWEEN 70 AND 79 ", columnName)
		} else if value == "high" {
			return fmt.Sprintf("%s BETWEEN 80 AND 100 ", columnName)
		}
		return fmt.Sprintf("%s BETWEEN %s AND %s ", columnName, value, value)
	}
	if operation == "like" {
		return fmt.Sprintf("cast(%s as text) %s '%s%s%s' ", columnName, operation, "%", value, "%")
	}
	if operation == "in" {
		return fmt.Sprintf("%s in %s%s%s ", columnName, "(", value, ")")
	}
	if operation == "notin" {
		return fmt.Sprintf("%s not in %s%s%s ", columnName, "(", value, ")")
	}
	return fmt.Sprintf("cast(%s as text) %s '%s' ", columnName, operation, value)
}

func DateFilter(columnName string, operation string, value string) string {
	if operation == "between" {
		dates := strings.Split(value, ",")
		if len(dates) == 2 {
			startDate := strings.TrimSpace(dates[0])
			endDate := strings.TrimSpace(dates[1])
			return fmt.Sprintf("date(%s) BETWEEN '%s' AND '%s' ", columnName, startDate, endDate)
		}
	}
	if operation == "=" {
		return fmt.Sprintf("date(%s) %s '%s' ", columnName, operation, value)
	}
	if operation == "!=" {
		return fmt.Sprintf("date(%s) %s '%s' ", columnName, operation, value)
	}
	if operation == "like" {
		return fmt.Sprintf("to_char(%s , '%s') %s '%s%s%s' ", columnName, "YYYY-mm-dd", operation, "%", value, "%")
	}
	return fmt.Sprintf("%s %s '%s' ", columnName, operation, value)
}

func CheckNilString(text string) bool {
	if text == "" {
		return true
	} else {
		return false
	}
}
