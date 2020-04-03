package collector

import (
	"github.com/ahmdrz/goinsta/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type collector struct {
	insta *goinsta.Instagram

	login     string
	password  string
	tokenPath string
	usernames []string
}

type Option func(c *collector)

func WithLogin(login, password string) Option {
	return func(c *collector) {
		c.login = login
		c.password = password
	}
}

func WithTokenPath(path string) Option {
	return func(c *collector) {
		c.tokenPath = path
	}
}

func WithTargets(usernames ...string) Option {
	return func(c *collector) {
		c.usernames = usernames
	}
}

func Instagram(opts ...Option) (*collector, error) {
	c := &collector{}
	for _, opt := range opts {
		opt(c)
	}

	insta, err := goinsta.Import(c.tokenPath)
	if err != nil {
		logrus.WithError(err).Debug("token file not found")
		logrus.Debug("trying to login with provided credentials")
		insta = goinsta.New(c.login, c.password)
		if err := insta.Login(); err != nil {
			return nil, err
		}
		if err := insta.Export(c.tokenPath); err != nil {
			return nil, err
		}
		logrus.Debug("login success")
	} else {
		logrus.Debug("token read successfully")
	}

	c.insta = insta
	return c, nil
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debug("collecting metrics")

	for _, username := range c.usernames {
		user, err := c.insta.Profiles.ByName(username)
		if err != nil {
			logrus.WithError(err).Fatal("cannot fetch user profile")
		}

		ch <- prometheus.MustNewConstMetric(
			mediaDesc,
			prometheus.GaugeValue,
			float64(user.MediaCount),
			username,
		)
		ch <- prometheus.MustNewConstMetric(
			followerDesc,
			prometheus.GaugeValue,
			float64(user.FollowerCount),
			username,
		)
		ch <- prometheus.MustNewConstMetric(
			followingDesc,
			prometheus.GaugeValue,
			float64(user.FollowingCount),
			username,
		)

		feed := user.Feed()
		if feed == nil {
			logrus.Fatal("collector: cannot fetch user feed")
		}

		for feed.Next() {
			for _, item := range feed.Items {
				ch <- prometheus.MustNewConstMetric(
					likeDesc,
					prometheus.GaugeValue,
					float64(item.Likes),
					username, item.ID,
				)
				ch <- prometheus.MustNewConstMetric(
					commentDesc,
					prometheus.GaugeValue,
					float64(item.CommentCount),
					username, item.ID,
				)
			}
		}
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	logrus.Debug("collecting descriptions")
	for _, desc := range all {
		ch <- desc
	}
}
