package forms

import (
	"net/http"
	"net/url"
)

type Form struct {
    url.Values
    Errors errors
}

func New(data url.Values) *Form {
    return &Form{
        data, 
        errors(map[string][]string{}),
    }
}

// checks if the form is empty
func (f *Form)Has(filed string, r *http.Request) bool {
    x := r.FormValue(filed)
    if x == ""{
        return false
    }
    return true
}
