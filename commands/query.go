package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crufter/borg/conf"
	"github.com/crufter/borg/types"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func Query(q string) error {
	client := &http.Client{Timeout: time.Duration(10 * time.Second)}
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/v1/query?l=%v&p=%v&q=%v", host(), *conf.L, *conf.P, url.QueryEscape(q)), nil)
	if err != nil {
		fmt.Println("Failed to create request: " + err.Error())
	}
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while making request: " + err.Error())
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}
	problems := []types.Problem{}
	err = json.Unmarshal(body, &problems)
	if err != nil {
		return errors.New("Malformed response from server")
	}
	renderQuery(problems)
	return nil
}

func renderQuery(problems []types.Problem) {
	const padding = 4
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight)
	for i, prob := range problems {
		if i > 0 {
			fmt.Fprintln(w, "")
		}
		fmt.Fprintln(w, fmt.Sprintf("(%v)", i+1), prob.Title)
		line := 0
	Loop:
		for x, sol := range prob.Solutions {
			fmt.Fprintf(w, "\t[%v%v]", i+1, x+1)
			for i, bodyPart := range sol.Body {
				if i > 0 {
					fmt.Fprintln(w, "\t\t", "-")
				}
				bodyPartLines := strings.Split(bodyPart, "\n")
				for j, bodyPartLine := range bodyPartLines {
					t := "\t\t"
					if i == 0 && j == 0 {
						t = "\t"
					}
					if len(strings.TrimSpace(bodyPartLine)) == 0 {
						continue
					}
					fmt.Fprintln(w, t, strings.Trim(bodyPartLine, "\n"))
					line++
					if line == 10 && *conf.F == false {
						fmt.Fprintln(w, "\t", "...", "\t")
						break Loop
					}
				}
			}
		}
	}
	w.Flush()
}

func host() string {
	return fmt.Sprintf("http://%v:9992", *conf.H)
}
