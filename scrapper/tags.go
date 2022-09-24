package scrapper

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/utils"
)

const (
	tagName = "manga-reader"

	optionErrorList              = "error-list"
	optionError                  = "error"
	optionMangaList              = "manga-list"
	optionMangaListCurrentPage   = "manga-list.current-page"
	optionMangaListLastPage      = "manga-list.last-page"
	optionMangaListNextPageURL   = "manga-list.next-page-url"
	optionMangaTitle             = "manga.title"
	optionMangaURL               = "manga.url"
	optionMangaImageURL          = "manga.image-url"
	optionMangaSlug              = "manga.slug"
	optionMangaOtherID           = "manga.id"
	optionMangaSynopsis          = "manga.synopsis"
	optionChapterList            = "manga.chapter-list"
	optionChapterListCurrentPage = "manga.chapter-list.current-page"
	optionChapterListLastPage    = "manga.chapter-list.last-page"
	optionChapterListNextPageURL = "manga.chapter-list.next-page-url"
	optionChapterTitle           = "manga.chapter.title"
	optionChapterOtherID         = "manga.chapter.id"
	optionChapterNumber          = "manga.chapter.number"
	optionChapterUploadDate      = "manga.chapter.upload-date"
	optionChapterURL             = "manga.chapter.url"
	optionChapterPageList        = "manga.chapter.pages-list"
	optionChapterPage            = "manga.chapter.page"
)

func isTagOptionEmpty(val string) bool {
	return val == "" || val == "-"
}

func hasOrderPreference(option string) (string, int, bool) {
	pattern, err := regexp.Compile("\\[[^]]+\\]")
	if err != nil {
		return option, 0, false
	}

	matches := pattern.FindAllString(option, -1)
	if len(matches) < 1 {
		return option, 0, false
	}

	strOrder := strings.Trim(matches[len(matches)-1], "[]")
	order, err := strconv.Atoi(strOrder)
	if err != nil {
		return option, 0, false
	}

	option = pattern.ReplaceAllString(option, "")

	return option, order, true
}

func hasDateFormat(option string) (string, string, bool) {
	opts := strings.Split(option, "|")
	format := ""
	if len(opts) > 1 {
		format = opts[1]
	}

	return opts[0], format, len(opts) > 1
}

func areOptionsEqual(expectedOption, actualOption string) bool {
	actualOptions := strings.Split(actualOption, ",")
	for _, opt := range actualOptions {
		trimmedOption, _, _ := hasOrderPreference(opt)
		opt, _, _ := hasDateFormat(trimmedOption)
		if opt == expectedOption {
			return true
		}
	}
	return false
}

func getInterfaceValue(data interface{}, option string) (interface{}, bool) {
	if data == nil {
		return nil, false
	}

	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Struct {
		return nil, false
	}

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tagOption := typeField.Tag.Get(tagName)

		if isTagOptionEmpty(tagOption) {
			if value, ok := getInterfaceValue(v.Field(i).Interface(), option); ok {
				return value, ok
			}
			continue
		}

		if areOptionsEqual(tagOption, option) {
			if _, f, ok := hasDateFormat(tagOption); ok {
				return parseDateByFormat(v.Field(i).Interface(), f), true
			}
			return v.Field(i).Interface(), true
		}
	}

	return nil, false
}

type orderedValue struct {
	Order int
	Value interface{}
}

func parseDateByFormat(val interface{}, format string) interface{} {
	format = strings.Split(format, "[")[0]
	if strVal, ok := val.(string); ok {
		if date, err := utils.ParseDate(strVal, format); err == nil {
			return date
		}
	}

	return val
}

func getInterfaceValueWithOrder(data interface{}, option string, oVal *orderedValue) {
	if data == nil {
		return
	}

	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tagOption := typeField.Tag.Get(tagName)

		if isTagOptionEmpty(tagOption) {
			getInterfaceValueWithOrder(v.Field(i).Interface(), option, oVal)
		}

		if areOptionsEqual(option, tagOption) {
			if _, order, ok := hasOrderPreference(tagOption); ok {
				if oVal.Value == nil || (oVal.Order >= order && !isEmpty(v.Field(i).Interface())) {
					oVal.Order = order
					oVal.Value = v.Field(i).Interface()
				}
			} else {
				oVal.Value = v.Field(i).Interface()
			}
			if _, f, ok := hasDateFormat(tagOption); ok {
				oVal.Value = parseDateByFormat(oVal.Value, f)
			}
		}
	}
}

func isEmpty(data interface{}) bool {
	if data == nil {
		return true
	}

	switch v := data.(type) {
	case string:
		return len(v) == 0
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	}

	return false
}

func getString(data interface{}) (string, bool) {
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.String:
		return fmt.Sprintf("%v", v.Interface()), true
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		val := fmt.Sprintf("%d", v.Interface())
		return val, true
	case reflect.Float32:
		if fVal, ok := v.Interface().(float32); ok {
			strVal := utils.GetString(float64(fVal))
			return strVal, true
		}
	case reflect.Float64:
		if fVal, ok := v.Interface().(float64); ok {
			strVal := utils.GetString(fVal)
			return strVal, true
		}
	}

	return "", false
}

func getInt(data interface{}) (int, bool) {
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.String:
		if strVal, ok := v.Interface().(string); ok {
			if val, err := strconv.Atoi(strVal); err == nil {
				return val, true
			}
		}
		return 0, false
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		val, _ := v.Interface().(int)
		return val, true
	}

	return 0, false
}

func getStringValueWithOrder(data interface{}, option string) (string, bool) {
	oVal := &orderedValue{}
	getInterfaceValueWithOrder(data, option, oVal)

	if oVal.Value == nil {
		return "", false
	}

	return getString(oVal.Value)
}

func getStringValue(data interface{}, option string) (string, bool) {
	if val, ok := getInterfaceValue(data, option); ok {
		if strVal, ok := getString(val); ok {
			return strVal, true
		}
	}

	return "", false
}

func getIntValue(data interface{}, option string) (int, bool) {
	if data, ok := getInterfaceValue(data, option); ok {
		return getInt(data)
	}

	return 0, false
}

func getManga(data interface{}) models.Manga {
	manga := models.Manga{}

	manga.Title, _ = getStringValueWithOrder(data, optionMangaTitle)
	manga.OtherID, _ = getStringValueWithOrder(data, optionMangaOtherID)
	manga.URL, _ = getStringValueWithOrder(data, optionMangaURL)
	manga.Slug, _ = getStringValueWithOrder(data, optionMangaSlug)
	manga.ImageURL, _ = getStringValueWithOrder(data, optionMangaImageURL)
	manga.Synopsis, _ = getStringValueWithOrder(data, optionMangaSynopsis)

	return manga
}

func getChapter(data interface{}) models.MangaChapter {
	chapter := models.MangaChapter{}

	chapter.OtherID, _ = getStringValueWithOrder(data, optionChapterOtherID)
	chapter.Title, _ = getStringValueWithOrder(data, optionChapterTitle)
	chapter.URL, _ = getStringValueWithOrder(data, optionChapterURL)

	chapter.Number, _ = getStringValueWithOrder(data, optionChapterNumber)
	chapter.Number = GetChapterNumber(chapter.Number)

	chapter.UploadDate, _ = getStringValueWithOrder(data, optionChapterUploadDate)

	return chapter
}

func getMangas(data interface{}) ([]models.Manga, error) {
	mangas := []models.Manga{}
	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return mangas, errors.Errorf("manga list interface is not of kind Array but of kind %s", v.Kind())
	}

	for i := 0; i < v.Len(); i++ {
		mangas = append(mangas, getManga(v.Index(i).Interface()))
	}

	return mangas, nil
}

func getChapters(data interface{}) ([]models.MangaChapter, bool) {
	chapters := []models.MangaChapter{}
	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return chapters, false
	}

	for i := 0; i < v.Len(); i++ {
		chapters = append(chapters, getChapter(v.Index(i).Interface()))
	}

	return chapters, true
}

func getPageURLs(data interface{}) ([]string, bool) {
	imageURLs := []string{}
	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return imageURLs, false
	}

	for i := 0; i < v.Len(); i++ {
		url, _ := getString(v.Index(i).Interface())
		imageURLs = append(imageURLs, url)
	}

	return imageURLs, true
}

func GetMangaListFromTags(data interface{}) ([]models.Manga, error) {
	mangas := []models.Manga{}

	iMangaList, ok := getInterfaceValue(data, optionMangaList)
	if !ok {
		return mangas, nil
	}

	return getMangas(iMangaList)
}

func GetMangaListPageData(data interface{}) (int, int, string) {
	currentPage, _ := getIntValue(data, optionMangaListCurrentPage)
	lastPage, _ := getIntValue(data, optionMangaListLastPage)
	nextPageURL, _ := getStringValue(data, optionMangaListNextPageURL)

	return currentPage, lastPage, nextPageURL
}

func GetChapterListPageData(data interface{}) (int, int, string) {
	currentPage, _ := getIntValue(data, optionChapterListCurrentPage)
	lastPage, _ := getIntValue(data, optionChapterListLastPage)
	nextPageURL, _ := getStringValue(data, optionChapterListNextPageURL)

	return currentPage, lastPage, nextPageURL
}

func GetMangaInfoFromTags(data interface{}) (models.Manga, error) {
	manga := getManga(data)

	iChapterList, ok := getInterfaceValue(data, optionChapterList)
	if !ok {
		return manga, nil
	}

	manga.Chapters, ok = getChapters(iChapterList)
	if !ok {
		return manga, errors.New("could not get chapters")
	}

	return manga, nil
}

func GetChapterPagesFromTags(data interface{}) (models.Pages, error) {
	pages := models.Pages{}
	oVal := orderedValue{}
	getInterfaceValueWithOrder(data, optionChapterPageList, &oVal)

	var ok bool
	pages.URLs, ok = getPageURLs(oVal.Value)
	if !ok {
		return pages, errors.New("could not get chapter pages")
	}

	return pages, nil
}

func GetErrorsFromTags(data interface{}) []string {
	val, ok := getInterfaceValue(data, optionErrorList)
	if !ok {
		val, _ := getStringValue(data, optionError)
		return []string{val}
	}

	errorMessages := []string{}
	v := reflect.ValueOf(val)

	if v.Kind() != reflect.Slice {
		return errorMessages
	}

	for i := 0; i < v.Len(); i++ {
		message, ok := getString(v.Index(i).Interface())
		if ok && len(message) > 0 {
			errorMessages = append(errorMessages, message)
		}
	}

	return errorMessages
}
