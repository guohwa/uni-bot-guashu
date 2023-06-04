package web

var auth map[string]map[string]bool = map[string]map[string]bool{
	"Public": {
		"account": true,
		"captcha": true,
	},
	"Admin": {
		"account":  true,
		"captcha":  true,
		"password": true,
		"profile":  true,
		"home":     true,
		"command":  true,
		"position": true,
		"order":    true,
		"income":   true,
		"trade":    true,
		"customer": true,
		"user":     true,
	},
	"User": {
		"account":  true,
		"captcha":  true,
		"password": true,
		"profile":  true,
		"home":     true,
		"command":  true,
		"position": true,
		"order":    true,
		"income":   true,
		"trade":    true,
		"customer": true,
	},
}
