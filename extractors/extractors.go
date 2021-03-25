package extractors

import (
	"net/url"
	"strings"

	"github.com/sallyciel/annie/extractors/acfun"
	"github.com/sallyciel/annie/extractors/bcy"
	"github.com/sallyciel/annie/extractors/bilibili"
	"github.com/sallyciel/annie/extractors/douyin"
	"github.com/sallyciel/annie/extractors/douyu"
	"github.com/sallyciel/annie/extractors/facebook"
	"github.com/sallyciel/annie/extractors/geekbang"
	"github.com/sallyciel/annie/extractors/haokan"
	"github.com/sallyciel/annie/extractors/instagram"
	"github.com/sallyciel/annie/extractors/iqiyi"
	"github.com/sallyciel/annie/extractors/mgtv"
	"github.com/sallyciel/annie/extractors/miaopai"
	"github.com/sallyciel/annie/extractors/netease"
	"github.com/sallyciel/annie/extractors/pixivision"
	"github.com/sallyciel/annie/extractors/pornhub"
	"github.com/sallyciel/annie/extractors/qq"
	"github.com/sallyciel/annie/extractors/tangdou"
	"github.com/sallyciel/annie/extractors/tiktok"
	"github.com/sallyciel/annie/extractors/tumblr"
	"github.com/sallyciel/annie/extractors/twitter"
	"github.com/sallyciel/annie/extractors/types"
	"github.com/sallyciel/annie/extractors/udn"
	"github.com/sallyciel/annie/extractors/universal"
	"github.com/sallyciel/annie/extractors/vimeo"
	"github.com/sallyciel/annie/extractors/weibo"
	"github.com/sallyciel/annie/extractors/xvideos"
	"github.com/sallyciel/annie/extractors/yinyuetai"
	"github.com/sallyciel/annie/extractors/youku"
	"github.com/sallyciel/annie/extractors/youtube"
	"github.com/sallyciel/annie/utils"
)

var extractorMap map[string]types.Extractor

func init() {
	douyinExtractor := douyin.New()
	youtubeExtractor := youtube.New()

	extractorMap = map[string]types.Extractor{
		"": universal.New(), // universal extractor

		"douyin":     douyinExtractor,
		"iesdouyin":  douyinExtractor,
		"bilibili":   bilibili.New(),
		"bcy":        bcy.New(),
		"pixivision": pixivision.New(),
		"youku":      youku.New(),
		"youtube":    youtubeExtractor,
		"youtu":      youtubeExtractor, // youtu.be
		"iqiyi":      iqiyi.New(),
		"mgtv":       mgtv.New(),
		"tangdou":    tangdou.New(),
		"tumblr":     tumblr.New(),
		"vimeo":      vimeo.New(),
		"facebook":   facebook.New(),
		"douyu":      douyu.New(),
		"miaopai":    miaopai.New(),
		"163":        netease.New(),
		"weibo":      weibo.New(),
		"instagram":  instagram.New(),
		"twitter":    twitter.New(),
		"qq":         qq.New(),
		"yinyuetai":  yinyuetai.New(),
		"geekbang":   geekbang.New(),
		"pornhub":    pornhub.New(),
		"xvideos":    xvideos.New(),
		"udn":        udn.New(),
		"tiktok":     tiktok.New(),
		"haokan":     haokan.New(),
		"acfun":      acfun.New(),
	}
}

// Extract is the main function to extract the data.
func Extract(u string, option types.Options) ([]*types.Data, error) {
	u = strings.TrimSpace(u)
	var domain string

	bilibiliShortLink := utils.MatchOneOf(u, `^(av|BV|ep)\w+`)
	if len(bilibiliShortLink) > 1 {
		bilibiliURL := map[string]string{
			"av": "https://www.bilibili.com/video/",
			"BV": "https://www.bilibili.com/video/",
			"ep": "https://www.bilibili.com/bangumi/play/",
		}
		domain = "bilibili"
		u = bilibiliURL[bilibiliShortLink[1]] + u
	} else {
		u, err := url.ParseRequestURI(u)
		if err != nil {
			return nil, err
		}
		if u.Host == "haokan.baidu.com" {
			domain = "haokan"
		} else {
			domain = utils.Domain(u.Host)
		}
	}
	extractor := extractorMap[domain]
	videos, err := extractor.Extract(u, option)
	if err != nil {
		return nil, err
	}
	for _, v := range videos {
		v.FillUpStreamsData()
	}
	return videos, nil
}
