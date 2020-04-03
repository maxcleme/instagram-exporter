package collector

import "github.com/prometheus/client_golang/prometheus"

var all = []*prometheus.Desc{
	mediaDesc,
	followerDesc,
	followingDesc,
	likeDesc,
	commentDesc,
}

var mediaDesc = prometheus.NewDesc(
	"instagram_media_total",
	"Total media count by username",
	[]string{"username"},
	nil,
)

var followerDesc = prometheus.NewDesc(
	"instagram_follower_total",
	"Total followers count by username",
	[]string{"username"},
	nil,
)

var followingDesc = prometheus.NewDesc(
	"instagram_following_total",
	"Total following count by username",
	[]string{"username"},
	nil,
)

var likeDesc = prometheus.NewDesc(
	"instagram_media_like_total",
	"Total likes count by [username, media_id]",
	[]string{"username", "media_id"},
	nil,
)

var commentDesc = prometheus.NewDesc(
	"instagram_media_comment_total",
	"Total comments count by [username, media_id]",
	[]string{"username", "media_id"},
	nil,
)
