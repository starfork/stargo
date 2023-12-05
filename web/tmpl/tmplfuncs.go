package tmpl

import (
	"fmt"
	"html"
	"html/template"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// HTML2str returns escaping text convert from html.
func HTML2str(html string) string {
	re := regexp.MustCompile(`\<[\S\s]+?\>`)
	html = re.ReplaceAllStringFunc(html, strings.ToLower)

	// remove STYLE
	re = regexp.MustCompile(`\<style[\S\s]+?\</style\>`)
	html = re.ReplaceAllString(html, "")

	// remove SCRIPT
	re = regexp.MustCompile(`\<script[\S\s]+?\</script\>`)
	html = re.ReplaceAllString(html, "")

	re = regexp.MustCompile(`\<[\S\s]+?\>`)
	html = re.ReplaceAllString(html, "\n")

	re = regexp.MustCompile(`\s{2,}`)
	html = re.ReplaceAllString(html, "\n")

	return strings.TrimSpace(html)
}

// Substr returns the substr from start to length.
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

func Htmlquote(text string) string {
	// HTML编码为实体符号
	/*
	   Encodes `text` for raw use in HTML.
	       >>> htmlquote("<'&\\">")
	       '&lt;&#39;&amp;&quot;&gt;'
	*/

	text = html.EscapeString(text)
	text = strings.NewReplacer(
		`“`, "&ldquo;",
		`”`, "&rdquo;",
		` `, "&nbsp;",
	).Replace(text)

	return strings.TrimSpace(text)
}

// Htmlunquote returns unquoted html string.
func Htmlunquote(text string) string {
	// 实体符号解释为HTML
	/*
	   Decodes `text` that's HTML quoted.
	       >>> htmlunquote('&lt;&#39;&amp;&quot;&gt;')
	       '<\\'&">'
	*/

	text = html.UnescapeString(text)

	return strings.TrimSpace(text)
}

func AssetsJs(text string) template.HTML {
	text = "<script src=\"" + text + "\"></script>"

	return template.HTML(text)
}

// AssetsCSS returns stylesheet link tag with src string.
func AssetsCSS(text string) template.HTML {
	text = "<link href=\"" + text + "\" rel=\"stylesheet\" />"

	return template.HTML(text)
}

func MapGet(arg1 interface{}, arg2 ...interface{}) (interface{}, error) {
	arg1Type := reflect.TypeOf(arg1)
	arg1Val := reflect.ValueOf(arg1)

	if arg1Type.Kind() == reflect.Map && len(arg2) > 0 {
		// check whether arg2[0] type equals to arg1 key type
		// if they are different, make conversion
		arg2Val := reflect.ValueOf(arg2[0])
		arg2Type := reflect.TypeOf(arg2[0])
		if arg2Type.Kind() != arg1Type.Key().Kind() {
			// convert arg2Value to string
			var arg2ConvertedVal interface{}
			arg2String := fmt.Sprintf("%v", arg2[0])

			// convert string representation to any other type
			switch arg1Type.Key().Kind() {
			case reflect.Bool:
				arg2ConvertedVal, _ = strconv.ParseBool(arg2String)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				arg2ConvertedVal, _ = strconv.ParseInt(arg2String, 0, 64)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				arg2ConvertedVal, _ = strconv.ParseUint(arg2String, 0, 64)
			case reflect.Float32, reflect.Float64:
				arg2ConvertedVal, _ = strconv.ParseFloat(arg2String, 64)
			case reflect.String:
				arg2ConvertedVal = arg2String
			default:
				arg2ConvertedVal = arg2Val.Interface()
			}
			arg2Val = reflect.ValueOf(arg2ConvertedVal)
		}

		storedVal := arg1Val.MapIndex(arg2Val)

		if storedVal.IsValid() {
			var result interface{}

			switch arg1Type.Elem().Kind() {
			case reflect.Bool:
				result = storedVal.Bool()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				result = storedVal.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				result = storedVal.Uint()
			case reflect.Float32, reflect.Float64:
				result = storedVal.Float()
			case reflect.String:
				result = storedVal.String()
			default:
				result = storedVal.Interface()
			}

			// if there is more keys, handle this recursively
			if len(arg2) > 1 {
				return MapGet(result, arg2[1:]...)
			}
			return result, nil
		}
		return nil, nil

	}
	return nil, nil
}
