package web

var auth map[string]map[string]bool = map[string]map[string]bool{
	"Public": {
		"account": true,
		"captcha": true,
	},
	"Admin": {
		"user":     true,
		"account":  true,
		"captcha":  true,
		"home":     true,
		"password": true,
		"profile":  true,
		"command":  true,
		"customer": true,
	},
	"User": {
		"account":  true,
		"captcha":  true,
		"home":     true,
		"password": true,
		"profile":  true,
		"command":  true,
		"customer": true,
	},
}
