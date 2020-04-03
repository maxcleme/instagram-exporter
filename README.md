# instagram-exporter

**WARNING: This exporter doesn't use Instagram official API because Facebook restrict them for Business/Creator account, moreover some metrics are not available in [Instagram Basic Display API](https://developers.facebook.com/docs/instagram-basic-display-api). Instead, this exporter use [goinsta](https://github.com/ahmdrz/goinsta) which bypass Instagram API, use at your own risk.**

## Configuration

Put the following file into your home folder : 

**.instagram-exporter.yaml**
```yaml
http_port: 2112
http_path: /metrics
log_level: debug

ig_login: XXXXXXXXXXXXX
ig_password: XXXXXXXXXXXXX
ig_token_path: /tmp/instagram-exporter/.ig_token
ig_target_users: AAAAAA,BBBBBB,CCCCCC
```

**Note: `ig_login` and `ig_password` will only be used once in order to create the token file, once done you can remove them from the config file.**

#### Metrics

| Name | Tags | Description |
| -- | -- | -- |
| `instagram_media_total` | `[username]` |  Total media count |
| `instagram_follower_total` | `[username]` | Total followers count |
| `instagram_following_total` | `[username]` | Total following count |
| `instagram_media_like_total` | `[username, media_id]` | Total likes |
| `instagram_media_comment_total` | `[username, media_id]` | Total comments |
