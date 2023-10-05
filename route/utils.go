package route

import (
	"Lightnovel/model"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

func PasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func Unhex(s string) ([]byte, error) {
	return hex.DecodeString(s[:model.IDHexLength])
}

func PasswordVerify(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

func IsUsernameValid(username string) bool {
	matched, _ := regexp.MatchString(
		fmt.Sprintf("^[a-zA-Z]{%v,%v}$", model.UsernameMinLength, model.UsernameMaxLength),
		username,
	)
	return matched
}

func IsPasswordValid(password string) bool {
	return !(len(password) < model.PasswordMinLength || len(password) > model.PasswordMaxLength)
}

func checkUserMetadata(data model.UserMetadata) (bool, ErrorCode) {
	if IsUsernameValid(data.Username) == false {
		return false, BadUsername
	}

	if utf8.RuneCountInString(data.Displayname) < model.DisplayNameMinLength ||
		utf8.RuneCountInString(data.Displayname) > model.DisplayNameMaxLength {
		return false, BadDisplayname
	}

	for _, ch := range []rune(data.Displayname) {
		if !unicode.IsPrint(ch) {
			return false, BadDisplayname
		}
	}

	if _, err := mail.ParseAddress(data.Email); err != nil ||
		utf8.RuneCountInString(data.Email) > model.EmailMaxLength {
		return false, BadEmail
	}

	return true, BadInput
}

func parseIntArray(nums string) []int {
	var res []int
	numStr := strings.Split(nums, ",")
	for _, str := range numStr {
		num, err := strconv.Atoi(str)
		if err == nil {
			res = append(res, num)
		}
	}
	return res
}

func getFiltersAndSort(c *fiber.Ctx) model.FiltersAndSortNovel {
	fromDate, err := time.Parse(time.DateOnly, c.Query("from", ""))
	if err != nil {
		fromDate = model.DefaultFiltersAndSort.FromDate
	}
	toDate, err := time.Parse(time.DateOnly, c.Query("to", ""))
	if err != nil {
		toDate = model.DefaultFiltersAndSort.ToDate
	}
	pageQuery := c.QueryInt("page", 1)
	page := model.DefaultFiltersAndSort.Page
	if pageQuery > 1 {
		page = uint(pageQuery)
	}
	orderBy := model.OrderBy(c.Query("orderBy", ""))
	if !orderBy.Validate() {
		orderBy = model.DefaultFiltersAndSort.OrderBy
	}
	sortOrder := model.SortOrder(c.Query("sortOrder", ""))
	if !sortOrder.Validate() {
		sortOrder = model.DefaultFiltersAndSort.SortOrder
	}

	return model.FiltersAndSortNovel{
		SortOrder:  sortOrder,
		OrderBy:    orderBy,
		Adult:      c.QueryBool("adult", model.DefaultFiltersAndSort.Adult),
		Language:   c.Query("language", model.DefaultFiltersAndSort.Language),
		Tag:        parseIntArray(c.Query("tag", "")),
		TagExclude: parseIntArray(c.Query("tagExclude", "")),
		Search:     c.Query("search", model.DefaultFiltersAndSort.Search),
		Page:       page,
		FromDate:   fromDate,
		ToDate:     toDate,
		Status: model.NovelStatusID(
			c.QueryInt("status", int(model.DefaultFiltersAndSort.Status)),
		),
	}
}
