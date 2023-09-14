package forms

type errors map[string][]string

// adds an error message in a given form field
func (e errors)Add(field, message string){
    e[field] = append(e[field], message)
}

// get the first error message
func (e errors)Get(field string) string {
    errorString := e[field]
    if len(errorString) == 0 {
        return ""
    }
    return errorString[0]
}

