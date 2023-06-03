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
		"position": true,
		"customer": true,
		"order":    true,
		"income":   true,
		"trade":    true,
	},
	"User": {
		"account":  true,
		"captcha":  true,
		"home":     true,
		"password": true,
		"profile":  true,
		"position": true,
		"customer": true,
		"order":    true,
		"income":   true,
		"trade":    true,
	},
}
