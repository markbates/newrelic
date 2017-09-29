package newrelic

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	nr "github.com/newrelic/go-agent"
)

func New(name string) buffalo.MiddlewareFunc {
	env := envy.Get("GO_ENV", "development")

	config := nr.NewConfig(name, os.Getenv("NEW_RELIC_LICENSE_KEY"))
	config.Enabled = env == "production"

	na, _ := nr.NewApplication(config)
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			req := c.Request()
			txn := na.StartTransaction(req.URL.String(), c.Response(), req)
			ri := c.Value("current_route").(buffalo.RouteInfo)
			txn.AddAttribute("PathName", ri.PathName)
			txn.AddAttribute("RequestID", c.Value("request_id"))
			defer txn.End()
			err := next(c)
			if err != nil {
				txn.NoticeError(err)
				return err
			}
			return nil
		}
	}
}
