package cim

import (
	"strings"
	"unicode"
)

// PropertyMapper 属性名映射器（CIM 属性名 <-> Iceberg 字段名）
type PropertyMapper struct {
	// CIM 属性名到 Iceberg 字段名的映射
	cimToIceberg map[string]string
	// Iceberg 字段名到 CIM 属性名的反向映射
	icebergToCim map[string]string
}

// NewPropertyMapper 创建属性名映射器
func NewPropertyMapper() *PropertyMapper {
	pm := &PropertyMapper{
		cimToIceberg: make(map[string]string),
		icebergToCim: make(map[string]string),
	}
	
	// 初始化特殊映射（与 mapping/rules.py 中的 ATTRIBUTE_NAME_MAPPING 保持一致）
	specialMappings := map[string]string{
		"mRID":            "m_rid",
		"assetID":         "asset_id",
		"workID":          "work_id",
		"taskID":          "task_id",
		"nominalVoltage":  "nominal_voltage_kv",
		"ratedCurrent":    "rated_current",
		"normalOpen":      "normal_open",
		"installDate":     "install_date",
		"serialNumber":    "serial_number",
		"serviceCategory": "service_category",
		"readingType":     "reading_type",
		"readingTime":     "reading_time",
		"constructionKind": "construction_kind",
		"sequenceNumber":  "sequence_number",
		"switchOnCount":   "switch_on_count",
		"switchOnDate":    "switch_on_date",
	}
	
	for cim, iceberg := range specialMappings {
		pm.cimToIceberg[cim] = iceberg
		pm.icebergToCim[iceberg] = cim
	}
	
	return pm
}

// CIMToIceberg 将 CIM 属性名转换为 Iceberg 字段名
func (pm *PropertyMapper) CIMToIceberg(cimName string) string {
	// 先检查特殊映射
	if iceberg, exists := pm.cimToIceberg[cimName]; exists {
		return iceberg
	}
	
	// 使用通用转换规则：驼峰转下划线
	return camelToSnake(cimName)
}

// IcebergToCIM 将 Iceberg 字段名转换为 CIM 属性名
func (pm *PropertyMapper) IcebergToCIM(icebergName string) string {
	// 先检查反向映射
	if cim, exists := pm.icebergToCim[icebergName]; exists {
		return cim
	}
	
	// 使用通用转换规则：下划线转驼峰
	return snakeToCamel(icebergName)
}

// camelToSnake 将驼峰命名转换为下划线命名
// 例如：mRID -> m_rid, normalOpen -> normal_open
func camelToSnake(s string) string {
	if s == "" {
		return s
	}
	
	var result strings.Builder
	runes := []rune(s)
	
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			// 如果前一个字符不是下划线，添加下划线
			if i > 0 && runes[i-1] != '_' {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	
	return result.String()
}

// snakeToCamel 将下划线命名转换为驼峰命名
// 例如：m_rid -> mRID, normal_open -> normalOpen
func snakeToCamel(s string) string {
	if s == "" {
		return s
	}
	
	parts := strings.Split(s, "_")
	var result strings.Builder
	
	for i, part := range parts {
		if part == "" {
			continue
		}
		
		if i == 0 {
			// 第一个部分：首字母小写（除非是特殊映射）
			if len(part) > 0 {
				runes := []rune(part)
				if len(runes) > 1 && unicode.IsUpper(runes[1]) {
					// 如果第二个字符是大写，保持原样（如 mRID）
					result.WriteString(part)
				} else {
					result.WriteRune(unicode.ToLower(runes[0]))
					if len(runes) > 1 {
						result.WriteString(string(runes[1:]))
					}
				}
			}
		} else {
			// 后续部分：首字母大写
			if len(part) > 0 {
				runes := []rune(part)
				result.WriteRune(unicode.ToUpper(runes[0]))
				if len(runes) > 1 {
					result.WriteString(string(runes[1:]))
				}
			}
		}
	}
	
	return result.String()
}

// MapPropertiesToIceberg 将 Actor 的属性映射转换为 Iceberg 字段格式
func (pm *PropertyMapper) MapPropertiesToIceberg(properties map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for cimName, value := range properties {
		icebergName := pm.CIMToIceberg(cimName)
		result[icebergName] = value
	}
	return result
}

// MapIcebergToProperties 将 Iceberg 字段格式转换为 Actor 的属性映射
func (pm *PropertyMapper) MapIcebergToProperties(icebergFields map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for icebergName, value := range icebergFields {
		cimName := pm.IcebergToCIM(icebergName)
		result[cimName] = value
	}
	return result
}

// DefaultPropertyMapper 默认的属性映射器实例
var DefaultPropertyMapper = NewPropertyMapper()
