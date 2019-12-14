package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func web() {

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {

		ss, _ := json.Marshal(cfg.Pages)

		fmt.Fprintf(w, string(ss))

	})

	http.HandleFunc("/do", func(w http.ResponseWriter, r *http.Request) {

		pageType := r.FormValue("pt")
		fromStr := r.FormValue("from")
		toStr := r.FormValue("to")

		from, err := strconv.Atoi(fromStr)
		if err != nil || from < 1 {
			from = 1
		}
		to, err := strconv.Atoi(toStr)
		if err != nil || to < 1 {
			to = 1
		}

		var nums []int
		for i := from; i <= to; i++ {
			nums = append(nums, i)
		}

		fmt.Fprintf(w, spider(pageType, nums...))

	})

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)

}
