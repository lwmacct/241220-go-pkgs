// utils.go
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cobra"
)

type FlagTag struct {
	Bind     string
	Default  string
	Group    string
	Required string
	Note     string
	Flag     string
}

// 绑定命令行标志到结构体字段
func bind(cmd *cobra.Command, val reflect.Value, prefix string, groupNames []string) {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)

		// 仅处理导出的字段
		if !fieldVal.CanSet() {
			continue
		}

		tag := &FlagTag{
			Bind:     fieldTyp.Tag.Get("bind"),
			Default:  fieldTyp.Tag.Get("default"),
			Group:    fieldTyp.Tag.Get("group"),
			Required: fieldTyp.Tag.Get("required"),
			Note:     fieldTyp.Tag.Get("note"),
			Flag:     fieldTyp.Tag.Get("flag"),
		}

		if tag.Bind == "false" {
			// 如果字段 bind 标签为 false，则跳过
			continue
		}

		// 获取字段的组信息
		if tag.Group != "" && !contains(groupNames, tag.Group) {
			// 如果字段的组信息不在指定的组中，则跳过
			continue
		}

		// 获取标志名称
		flagName := fieldTyp.Name
		if prefix != "" {
			flagName = prefix + "." + flagName
		}

		// 获取自定义的标志名称
		if tag.Flag != "" {
			flagName = tag.Flag
		} else {
			// 将点号替换为连字符，并转换为 kebab-case
			flagName = toKebabCase(flagName)
		}

		usage := tag.Note
		defaultValue, err := parseDefaultValue(fieldVal, tag.Default)
		if err != nil {
			fmt.Printf("解析字段 %s 的默认值时出错: %v\n", flagName, err)
			continue
		}

		// 根据字段类型绑定标志
		fieldKind := fieldTyp.Type.Kind()
		switch {
		case fieldTyp.Type == reflect.TypeOf(time.Duration(0)):
			defaultDuration := defaultValue.(time.Duration)
			fieldVal.Set(reflect.ValueOf(defaultDuration))
			cmd.Flags().DurationVar(fieldVal.Addr().Interface().(*time.Duration), flagName, defaultDuration, usage)

		case fieldKind == reflect.Bool:
			defaultBool := defaultValue.(bool)
			fieldVal.SetBool(defaultBool)
			cmd.Flags().BoolVar(fieldVal.Addr().Interface().(*bool), flagName, defaultBool, usage)

		// int 系列
		case fieldTyp.Type == reflect.TypeOf(int(0)):
			defaultInt := defaultValue.(int)
			fieldVal.SetInt(int64(defaultInt))
			cmd.Flags().IntVar(fieldVal.Addr().Interface().(*int), flagName, defaultInt, usage)
		case fieldTyp.Type == reflect.TypeOf(int8(0)):
			defaultInt := defaultValue.(int8)
			fieldVal.SetInt(int64(defaultInt))
			cmd.Flags().Int8Var(fieldVal.Addr().Interface().(*int8), flagName, defaultInt, usage)
		case fieldTyp.Type == reflect.TypeOf(int16(0)):
			defaultInt := defaultValue.(int16)
			fieldVal.SetInt(int64(defaultInt))
			cmd.Flags().Int16Var(fieldVal.Addr().Interface().(*int16), flagName, defaultInt, usage)
		case fieldTyp.Type == reflect.TypeOf(int32(0)):
			defaultInt := defaultValue.(int32)
			fieldVal.SetInt(int64(defaultInt))
			cmd.Flags().Int32Var(fieldVal.Addr().Interface().(*int32), flagName, defaultInt, usage)
		case fieldTyp.Type == reflect.TypeOf(int64(0)):
			defaultInt := defaultValue.(int64)
			fieldVal.SetInt(defaultInt)
			cmd.Flags().Int64Var(fieldVal.Addr().Interface().(*int64), flagName, defaultInt, usage)

		// float 系列
		case fieldTyp.Type == reflect.TypeOf(float32(0)):
			defaultFloat := defaultValue.(float32)
			fieldVal.SetFloat(float64(defaultFloat))
			cmd.Flags().Float32Var(fieldVal.Addr().Interface().(*float32), flagName, defaultFloat, usage)
		case fieldTyp.Type == reflect.TypeOf(float64(0)):
			defaultFloat := defaultValue.(float64)
			fieldVal.SetFloat(defaultFloat)
			cmd.Flags().Float64Var(fieldVal.Addr().Interface().(*float64), flagName, defaultFloat, usage)

		// uint 系列
		case fieldTyp.Type == reflect.TypeOf(uint(0)):
			defaultUint := defaultValue.(uint)
			fieldVal.SetUint(uint64(defaultUint))
			cmd.Flags().UintVar(fieldVal.Addr().Interface().(*uint), flagName, defaultUint, usage)
		case fieldTyp.Type == reflect.TypeOf(uint8(0)):
			defaultUint := defaultValue.(uint8)
			fieldVal.SetUint(uint64(defaultUint))
			cmd.Flags().Uint8Var(fieldVal.Addr().Interface().(*uint8), flagName, defaultUint, usage)
		case fieldTyp.Type == reflect.TypeOf(uint16(0)):
			defaultUint := defaultValue.(uint16)
			fieldVal.SetUint(uint64(defaultUint))
			cmd.Flags().Uint16Var(fieldVal.Addr().Interface().(*uint16), flagName, defaultUint, usage)
		case fieldTyp.Type == reflect.TypeOf(uint32(0)):
			defaultUint := defaultValue.(uint32)
			fieldVal.SetUint(uint64(defaultUint))
			cmd.Flags().Uint32Var(fieldVal.Addr().Interface().(*uint32), flagName, defaultUint, usage)
		case fieldTyp.Type == reflect.TypeOf(uint64(0)):
			defaultUint := defaultValue.(uint64)
			fieldVal.SetUint(defaultUint)
			cmd.Flags().Uint64Var(fieldVal.Addr().Interface().(*uint64), flagName, defaultUint, usage)

		// string
		case fieldKind == reflect.String:
			defaultStr := defaultValue.(string)
			fieldVal.SetString(defaultStr)
			cmd.Flags().StringVar(fieldVal.Addr().Interface().(*string), flagName, defaultStr, usage)

		// 切片处理
		case fieldKind == reflect.Slice:
			elemKind := fieldTyp.Type.Elem().Kind()
			switch {
			// 字符串切片
			case elemKind == reflect.String:
				defaultSlice := defaultValue.([]string)
				fieldVal.Set(reflect.ValueOf(defaultSlice))
				cmd.Flags().StringSliceVar(fieldVal.Addr().Interface().(*[]string), flagName, defaultSlice, usage)
			// int 系列切片
			case elemKind >= reflect.Int && elemKind <= reflect.Int64:
				defaultSlice := reflect.ValueOf(defaultValue)
				intSlice := make([]int, defaultSlice.Len())
				for j := 0; j < defaultSlice.Len(); j++ {
					intSlice[j] = int(defaultSlice.Index(j).Int())
				}
				fieldVal.Set(reflect.ValueOf(intSlice))
				cmd.Flags().IntSliceVar(fieldVal.Addr().Interface().(*[]int), flagName, intSlice, usage)
			case elemKind == reflect.Int32:
				defaultSlice := defaultValue.([]int32)
				fieldVal.Set(reflect.ValueOf(defaultSlice))
				cmd.Flags().Int32SliceVar(fieldVal.Addr().Interface().(*[]int32), flagName, defaultSlice, usage)
			case elemKind == reflect.Int64:
				defaultSlice := defaultValue.([]int64)
				fieldVal.Set(reflect.ValueOf(defaultSlice))
				cmd.Flags().Int64SliceVar(fieldVal.Addr().Interface().(*[]int64), flagName, defaultSlice, usage)
			// float 系列切片
			case elemKind == reflect.Float64:
				defaultSlice := defaultValue.([]float64)
				fieldVal.Set(reflect.ValueOf(defaultSlice))
				cmd.Flags().Float64SliceVar(fieldVal.Addr().Interface().(*[]float64), flagName, defaultSlice, usage)
			case elemKind == reflect.Bool:
				defaultSlice := defaultValue.([]bool)
				fieldVal.Set(reflect.ValueOf(defaultSlice))
				cmd.Flags().BoolSliceVar(fieldVal.Addr().Interface().(*[]bool), flagName, defaultSlice, usage)
			// uint 系列切片
			case elemKind >= reflect.Uint && elemKind <= reflect.Uint64:
				defaultSlice := reflect.ValueOf(defaultValue)
				uintSlice := make([]uint, defaultSlice.Len())
				for j := 0; j < defaultSlice.Len(); j++ {
					uintSlice[j] = uint(defaultSlice.Index(j).Uint())
				}
				fieldVal.Set(reflect.ValueOf(uintSlice))
				cmd.Flags().UintSliceVar(fieldVal.Addr().Interface().(*[]uint), flagName, uintSlice, usage)
			default:
				log.Println("不支持的切片元素类型", elemKind)
			}
		case fieldKind == reflect.Map:
			elemKind := fieldTyp.Type.Elem().Kind()
			switch {
			// map[string]string
			case elemKind == reflect.String:
				defaultMap := defaultValue.(map[string]string)
				fieldVal.Set(reflect.ValueOf(defaultMap))
				cmd.Flags().StringToStringVar(fieldVal.Addr().Interface().(*map[string]string), flagName, defaultMap, usage)
			// map[string]int
			case elemKind == reflect.Int:
				defaultMap := defaultValue.(map[string]int)
				fieldVal.Set(reflect.ValueOf(defaultMap))
				cmd.Flags().StringToIntVar(fieldVal.Addr().Interface().(*map[string]int), flagName, defaultMap, usage)
			default:
				log.Println("utils.go", "不支持的 map 元素类型", elemKind)
			}
		case fieldKind == reflect.Struct:
			// 递归处理嵌套的结构体
			bind(cmd, fieldVal, flagName, groupNames)
		default:
			log.Println("utils.go", "不支持的类型", fieldKind)
		}

		// 检查环境变量，优先使用命令行参数
		envPrefix := os.Getenv("PREFIX_ACF")
		if envPrefix == "" {
			envPrefix = "ACF_"
		}
		envVar := envPrefix + strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
		envValue := os.Getenv(envVar)

		if envValue != "" && !cmd.Flags().Changed(flagName) {
			if err := cmd.Flags().Set(flagName, envValue); err != nil {
				fmt.Printf("从环境变量 %s 设置标志 %s 时出错: %v\n", envVar, flagName, err)
			}
		}

		// 如果标志是必选的，标记为必选
		if tag.Required == "true" {
			if err := cmd.MarkFlagRequired(flagName); err != nil {
				fmt.Printf("标记标志 %s 为必选时出错: %v\n", flagName, err)
			}
		}
	}
}

// contains 检查字符串是否在列表中
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func toKebabCase(s string) string {
	// 替换所有非字母数字字符为 '-'
	var result []rune
	prevDash := false
	for _, r := range s {
		if r == '.' || r == ' ' || r == '_' {
			if !prevDash {
				result = append(result, '-')
				prevDash = true
			}
			continue
		}
		if unicode.IsUpper(r) {
			// 如果不是第一个字符且前一个字符不是 '-'
			if len(result) > 0 && result[len(result)-1] != '-' {
				// 仅在当前大写字母后面跟随小写字母时插入 '-'
				// 以避免在连续大写字母之间插入 '-'
				nextIndex := len(result)
				if nextIndex < len(s)-1 {
					nextRune := rune(s[nextIndex+1])
					if unicode.IsLower(nextRune) {
						result = append(result, '-')
					}
				}
			}
			result = append(result, unicode.ToLower(r))
			prevDash = false
		} else {
			result = append(result, r)
			prevDash = false
		}
	}

	// 转换为字符串并移除可能的首尾 '-'
	kebab := strings.Trim(string(result), "-")
	return kebab
}

// parseDefaultValue 解析字段的默认值
func parseDefaultValue(field reflect.Value, defaultStr string) (any, error) {
	fieldType := field.Type()

	// 如果没有默认值，返回类型的零值
	if defaultStr == "" {
		return reflect.Zero(fieldType).Interface(), nil
	}

	// 针对 time.Duration 类型处理
	if fieldType == reflect.TypeOf(time.Duration(0)) {
		duration, err := time.ParseDuration(defaultStr)
		if err != nil {
			return nil, fmt.Errorf("解析 time.Duration 类型的默认值失败: %v", err)
		}
		return duration, nil
	}

	// 根据类型解析默认值
	switch fieldType.Kind() {
	case reflect.Bool:
		return strconv.ParseBool(defaultStr)
	case reflect.Int:
		i, err := strconv.ParseInt(defaultStr, 10, 0)
		return int(i), err
	case reflect.Int8:
		i, err := strconv.ParseInt(defaultStr, 10, 8)
		return int8(i), err
	case reflect.Int16:
		i, err := strconv.ParseInt(defaultStr, 10, 16)
		return int16(i), err
	case reflect.Int32:
		i, err := strconv.ParseInt(defaultStr, 10, 32)
		return int32(i), err
	case reflect.Int64:
		return strconv.ParseInt(defaultStr, 10, 64)
	case reflect.Uint:
		u, err := strconv.ParseUint(defaultStr, 10, 0)
		return uint(u), err
	case reflect.Uint8:
		u, err := strconv.ParseUint(defaultStr, 10, 8)
		return uint8(u), err
	case reflect.Uint16:
		u, err := strconv.ParseUint(defaultStr, 10, 16)
		return uint16(u), err
	case reflect.Uint32:
		u, err := strconv.ParseUint(defaultStr, 10, 32)
		return uint32(u), err
	case reflect.Uint64:
		return strconv.ParseUint(defaultStr, 10, 64)
	case reflect.Float32:
		f, err := strconv.ParseFloat(defaultStr, 32)
		return float32(f), err
	case reflect.Float64:
		return strconv.ParseFloat(defaultStr, 64)
	case reflect.String:
		return defaultStr, nil
	case reflect.Slice:
		return parseSliceDefaultValue(fieldType, defaultStr)
	default:
		return nil, fmt.Errorf("不支持的类型: %v", fieldType.Kind())
	}
}

// parseSliceDefaultValue 处理切片类型的默认值解析
func parseSliceDefaultValue(fieldType reflect.Type, defaultStr string) (interface{}, error) {
	elemType := fieldType.Elem()
	elements := strings.Split(defaultStr, ",")
	resultSlice := reflect.MakeSlice(fieldType, len(elements), len(elements))

	for i, elem := range elements {
		elem = strings.TrimSpace(elem)
		parsedElem, err := parseDefaultValue(reflect.New(elemType).Elem(), elem)
		if err != nil {
			return nil, fmt.Errorf("解析切片元素 '%s' 时出错: %v", elem, err)
		}
		resultSlice.Index(i).Set(reflect.ValueOf(parsedElem).Convert(elemType))
	}
	return resultSlice.Interface(), nil
}

func isInputFromPipe() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("获取标准输入状态失败:", err)
		return false
	}

	// 检查标准输入是否是管道输入或者重定向
	return (fileInfo.Mode() & os.ModeCharDevice) == 0
}

func getPackageName(skip int) string {
	_, filename, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return path.Base(path.Dir(filename))
}
