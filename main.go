package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/phanletrunghieu/dkhp-uit/util"
)

func main() {
	// get flag args
	idPtr := flag.String("id", "", "Mã số sinh viên")
	passPtr := flag.String("pass", "", "Mật khẩu tài khoản chứng thực")
	loopPtr := flag.Int("loop", 1, "Số vòng lặp (-1 để chạy vô tận)")
	flag.Parse()
	classes := strings.Join(flag.Args(), "\r\n")

	// create client
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	for i := 0; i < *loopPtr || *loopPtr <= -1; i++ {
		if *loopPtr > 1 || *loopPtr <= -1 {
			fmt.Printf("Loop %d\n", i+1)
		}

		// login
		err := Login(client, idPtr, passPtr)
		if err != nil {
			fmt.Println("Login fail")
			continue
		}

		// get form token
		data, err := util.HttpGet(client, "https://dkhp.uit.edu.vn/sinhvien/hocphan/dangky")
		if err != nil {
			fmt.Println("Login fail")
			continue
		}
		form_token, err := util.GetFormToken(data)
		if err != nil {
			fmt.Println("It's not time")
			continue
		}

		//dkhp
		dataDKHP, err := Dkhp(client, form_token, classes)
		if err != nil {
			fmt.Println("DKHP fail")
			continue
		}

		//get result
		result, err := util.GetResultOfDKHP(dataDKHP)
		if err != nil {
			fmt.Println("Can't show result")
			continue
		}

		fmt.Printf("%v\n", result)
		break
	}
}

func Login(client *http.Client, idPtr *string, passPtr *string) error {
	v := url.Values{}
	v.Set("name", *idPtr)
	v.Set("pass", *passPtr)
	v.Set("form_id", "user_login")
	_, err := util.HttpPost(client, "https://dkhp.uit.edu.vn", v.Encode())
	if err != nil {
		return err
	}
	return nil
}

func Dkhp(client *http.Client, form_token string, classes string) (string, error) {
	v := url.Values{}
	v.Set("dsmalop", classes)
	v.Set("form_id", "uit_dkhp_dangky_form")
	v.Set("form_token", form_token)
	v.Set("op", "Đăng ký")
	dataDKHP, err := util.HttpPost(client, "https://dkhp.uit.edu.vn/sinhvien/hocphan/dangky", v.Encode())
	if err != nil {
		return "", err
	}
	return dataDKHP, nil
}
