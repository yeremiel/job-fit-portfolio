package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"job-fit/internal/evaluator"
	"job-fit/internal/failure"
)

const MaxBodyBytes = 128 * 1024

type JobMetadata struct {
	Company   string `json:"company,omitempty"`
	Title     string `json:"title,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
}
type input struct {
	job         JobMetadata
	description string
}

func parseRequest(w http.ResponseWriter, r *http.Request) (input, *apiError) {
	bad := input{}
	media, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" || (params["charset"] != "" && !strings.EqualFold(params["charset"], "utf-8")) || len(params) > 1 || (len(params) == 1 && params["charset"] == "") || r.Header.Get("Content-Encoding") != "" {
		return bad, problem(415)
	}
	if r.URL.RawQuery != "" || r.URL.ForceQuery {
		return bad, problem(400)
	}
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	b, err := io.ReadAll(r.Body)
	if err != nil {
		if _, ok := err.(*http.MaxBytesError); ok {
			return bad, problem(413)
		}
		return bad, malformed()
	}
	if !utf8.Valid(b) || !json.Valid(b) || !validSurrogates(b) {
		return bad, malformed()
	}
	d := json.NewDecoder(bytes.NewReader(b))
	if !uniqueKeys(d, 0) {
		return bad, problem(400)
	}
	var top map[string]json.RawMessage
	if json.Unmarshal(b, &top) != nil || top == nil || len(top) != 1 || top["job"] == nil {
		return bad, problem(400)
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(top["job"], &fields) != nil || fields == nil {
		return bad, problem(400)
	}
	values := map[string]string{}
	for key, raw := range fields {
		switch key {
		case "company", "title", "sourceUrl", "description":
		default:
			return bad, problem(400)
		}
		var value string
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) || json.Unmarshal(raw, &value) != nil {
			return bad, problem(400)
		}
		values[key] = value
	}
	if _, ok := values["description"]; !ok {
		return bad, validation("job.description", "required", "Description is required.")
	}
	for _, key := range []string{"description", "company", "title", "sourceUrl"} {
		value, ok := values[key]
		if !ok {
			continue
		}
		limit := 16 * 1024
		switch key {
		case "company":
			limit = 256
		case "title":
			limit = 512
		case "sourceUrl":
			limit = 2048
		}
		if len(value) > limit {
			return bad, problem(413)
		}
		if strings.ContainsRune(value, '\x00') {
			return bad, validation("job."+key, "invalid_character", "NUL is not allowed.")
		}
		if strings.TrimSpace(value) == "" {
			return bad, validation("job."+key, "blank", "Field must not be blank.")
		}
	}
	source := strings.TrimSpace(values["sourceUrl"])
	if source != "" {
		u, err := url.Parse(source)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil {
			return bad, validation("job.sourceUrl", "invalid_format", "Source URL must be an absolute HTTP or HTTPS URL without credentials.")
		}
	}
	if _, err := evaluator.PrepareJob(values["description"]); err != nil {
		if failure.Classify(err) == failure.TooLarge {
			return bad, problem(413)
		}
		return bad, validation("job.description", "invalid_format", "Description is invalid.")
	}
	return input{JobMetadata{strings.TrimSpace(values["company"]), strings.TrimSpace(values["title"]), source}, values["description"]}, nil
}

// encoding/json accepts duplicate keys and replaces isolated UTF-16 surrogates.
// Check both explicitly before interpreting the request object.
func uniqueKeys(d *json.Decoder, depth int) bool {
	if depth > 100 {
		return false
	}
	t, err := d.Token()
	if err != nil {
		return false
	}
	delim, ok := t.(json.Delim)
	if !ok {
		return true
	}
	if delim == '{' {
		seen := map[string]bool{}
		for d.More() {
			key, err := d.Token()
			if err != nil {
				return false
			}
			k, ok := key.(string)
			if !ok || seen[k] {
				return false
			}
			seen[k] = true
			if !uniqueKeys(d, depth+1) {
				return false
			}
		}
	} else if delim == '[' {
		for d.More() {
			if !uniqueKeys(d, depth+1) {
				return false
			}
		}
	} else {
		return false
	}
	_, err = d.Token()
	return err == nil
}
func validSurrogates(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] != '"' {
			continue
		}
		i++
		for ; i < len(b) && b[i] != '"'; i++ {
			if b[i] != '\\' {
				continue
			}
			i++
			if b[i] != 'u' {
				continue
			}
			n, _ := strconv.ParseUint(string(b[i+1:i+5]), 16, 16)
			i += 4
			if n >= 0xDC00 && n <= 0xDFFF {
				return false
			}
			if n >= 0xD800 && n <= 0xDBFF {
				if i+6 >= len(b) || b[i+1] != '\\' || b[i+2] != 'u' {
					return false
				}
				low, err := strconv.ParseUint(string(b[i+3:i+7]), 16, 16)
				if err != nil || low < 0xDC00 || low > 0xDFFF {
					return false
				}
				i += 6
			}
		}
	}
	return true
}
