package main

import (
	"bufio"
	"fmt"
	"github.com/struCoder/Go-pinyin"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

type zxT struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	Msg      string         `json:"msg"`
	Children map[string]zxT `json:"children"`
}

const (
	provinceP = `\d{2}0000`
	cityP     = `\d{4}00`
	replaceP  = `(\d{2})(\d{2})(\d{2})`
)

func ParseFile(fileName string) (result map[string]zxT) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = file.Close()
	}()

	replaceRe := regexp.MustCompile(replaceP)

	result = make(map[string]zxT, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		if match, _ := regexp.MatchString(provinceP, words[0]); match {
			result[words[0]] = zxT{
				Id:       words[0],
				Name:     words[1],
				Children: make(map[string]zxT, 0),
			}
		} else if match, _ := regexp.MatchString(cityP, words[0]); match {
			province := replaceRe.ReplaceAllString(words[0], `${1}0000`)
			_, ok := result[province]
			if !ok {
				fmt.Println(province, words)
				continue
			}
			result[province].Children[words[0]] = zxT{
				Id:       words[0],
				Name:     words[1],
				Children: make(map[string]zxT, 0),
			}
		} else {
			province := replaceRe.ReplaceAllString(words[0], `${1}0000`)
			city := replaceRe.ReplaceAllString(words[0], `${1}${2}00`)
			provinceC, ok := result[province]
			if !ok {
				fmt.Println(words)
				continue
			}
			cityC, ok := provinceC.Children[city]
			if !ok {
				result[province].Children[words[0]] = zxT{
					Id:       words[0],
					Name:     words[1],
					Msg:      "save to province",
					Children: make(map[string]zxT, 0),
				}
				continue
			}
			cityC.Children[words[0]] = zxT{
				Id:       words[0],
				Name:     words[1],
				Children: make(map[string]zxT, 0),
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}

var pinyinConv = pinyingo.NewPy(pinyingo.STYLE_NORMAL, pinyingo.NO_SEGMENT)
var initialsConv = pinyingo.NewPy(pinyingo.STYLE_INITIALS, pinyingo.NO_SEGMENT)
var firstLetterConv = pinyingo.NewPy(pinyingo.STYLE_FIRST_LETTER, pinyingo.NO_SEGMENT)

func conv(hz string) []string {
	r1 := pinyinConv.Convert(hz)
	r2 := initialsConv.Convert(hz)
	r3 := firstLetterConv.Convert(hz)

	result := make([]string, 0)
	for index, items := range [][]string{r1, r1, r2, r3} {
		hz1 := make([]string, 0)
		if index == 0 {
			for _, item := range items {
				hz1 = append(hz1, strings.Title(item))
			}
		} else {
			for _, item := range items {
				hz1 = append(hz1, item)
			}
		}
		result = append(result, strings.Join(hz1, ""))
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	fileName := os.Args[1]
	result := ParseFile(fileName)
	items := make([][]string, 0)

	var (
		huaBei   = "1"
		dongBei  = "2"
		huaDong  = "3"
		xiBei    = "4"
		huaZhong = "5"
		huaNan   = "6"
		xiNan    = "7"
		qiTa     = "8"
	)

	areaIds := []string{huaBei, dongBei, huaDong, xiBei, huaZhong, huaNan, xiNan, qiTa}
	areas := map[string]string{"1": "华北", "2": "东北", "3": "华东", "4": "西北", "5": "华中", "6": "华南", "7": "西南", "8": "其他"}

	var regionArea = map[string]string{
		"110000": huaBei,
		"120000": huaBei,
		"130000": huaBei,
		"140000": huaBei,
		"150000": huaBei,
		"210000": dongBei,
		"220000": dongBei,
		"230000": dongBei,
		"310000": huaDong,
		"320000": huaDong,
		"330000": huaDong,
		"340000": huaDong,
		"350000": huaDong,
		"360000": huaZhong,
		"370000": huaDong,
		"410000": huaZhong,
		"420000": huaZhong,
		"430000": huaZhong,
		"440000": huaNan,
		"450000": huaNan,
		"460000": huaNan,
		"500000": xiNan,
		"510000": xiNan,
		"520000": xiNan,
		"530000": xiNan,
		"540000": xiNan,
		"610000": xiBei,
		"620000": xiBei,
		"630000": xiBei,
		"640000": xiBei,
		"650000": xiBei,
		"710000": qiTa,
		"810000": qiTa,
		"820000": qiTa,
	}

	for _, v := range result {
		area, _ := regionArea[v.Id]
		provinceId := v.Id
		items = append(items, format(v, "0", provinceId, "0", area, isLeaf(v)))
		for _, vv := range v.Children {
			cityId := vv.Id
			items = append(items, format(vv, v, provinceId, cityId, area, isLeaf(vv)))
			for _, vvv := range vv.Children {
				items = append(items, format(vvv, vv, provinceId, cityId, area, isLeaf(vvv)))
			}
		}
	}

	values := make([]string, 0)
	for _, areaId := range areaIds {
		values = append(values, fmt.Sprintf("(%s, '%s')", areaId, areas[areaId]))
	}
	fmt.Printf("INSERT INTO fog_addr_area(id, name) VALUES %s;\n", strings.Join(values, ","))

	values = make([]string, 0)
	for _, item := range items {
		values = append(
			values,
			fmt.Sprintf(
				"(%s,'%s',%s,%s,%s,%s,'%s','%s','%s','%s',%s)",
				item[0],
				item[1],
				item[2],
				item[3],
				item[4],
				item[5],
				item[6],
				item[7],
				item[8],
				item[9],
				item[10],
			),
		)
	}
	sort.Strings(values)
	fmt.Printf("INSERT INTO fog_addr_region(id, name, parent_id, province_id, city_id, area_id, pinyin, title_pinyin, initials, first_letter, is_leaf) VALUES \n%s;", strings.Join(values, ",\n"))
}

func isLeaf(v zxT) string {
	isLeaf := "0"
	if len(v.Children) < 1 {
		isLeaf = "1"
	}
	return isLeaf
}

func format(vv zxT, root interface{}, provinceId, cityId, area, isLeaf string) []string {
	pinyins := conv(vv.Name)
	switch v := root.(type) {
	case zxT:
		return []string{vv.Id, vv.Name, v.Id, provinceId, cityId, area, pinyins[0], pinyins[1], pinyins[2], pinyins[3], isLeaf}
	case string:
		return []string{vv.Id, vv.Name, v, provinceId, cityId, area, pinyins[0], pinyins[1], pinyins[2], pinyins[3], isLeaf}
	default:
		return []string{}
	}
}
