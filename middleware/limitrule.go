package middleware

import (
	"time"
)

var (
	AuthLimit = FrequentWare(LimitRule{
		Namespace: "auth",
		Limit:     10,
		Window:    1 * time.Minute,
		ByIP:      true,
	})
	QueryLimit = FrequentWare(LimitRule{
		Namespace: "query",
		Limit:     300,
		Window:    1 * time.Minute,
		ByIP:      false,
	})
	ClaimLimit = FrequentWare(LimitRule{
		Namespace: "claim",
		Limit:     30,
		Window:    1 * time.Hour,
		ByIP:      false,
	})
	PublishLimit = FrequentWare(LimitRule{
		Namespace: "publish",
		Limit:     20,
		Window:    1 * time.Hour,
		ByIP:      false,
	})
	UploadLimit = FrequentWare(LimitRule{
		Namespace: "upload",
		Limit:     60,
		Window:    1 * time.Hour,
		ByIP:      false,
	})
)
