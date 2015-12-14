package json

import (
	"strings"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

type ChinaSitesRule struct {
	Rule
}

func (this *ChinaSitesRule) Apply(dest v2net.Destination) bool {
	address := dest.Address()
	if !address.IsDomain() {
		return false
	}
	domain := strings.ToLower(address.Domain())
	for _, matcher := range compiledMatchers {
		if matcher.Match(domain) {
			return true
		}
	}
	return false
}

const (
	anySubDomain = "^(.*\\.)?"
	dotAm        = "\\.am$"
	dotCc        = "\\.cc$"
	dotCn        = "\\.cn$"
	dotCom       = "\\.com$"
	dotLa        = "\\.la$"
	dotNet       = "\\.net$"
	dotOrg       = "\\.org$"
	dotTv        = "\\.tv$"
)

var (
	compiledMatchers = make([]*RegexpDomainMatcher, 0, 1024)

	regexpDomains = []string{
		dotCn,

		anySubDomain + "10010" + dotCom,
		anySubDomain + "115" + dotCom,
		anySubDomain + "123u" + dotCom,
		anySubDomain + "126" + dotCom,
		anySubDomain + "126" + dotNet,
		anySubDomain + "163" + dotCom,
		anySubDomain + "17173" + dotCom,
		anySubDomain + "17cdn" + dotCom,
		anySubDomain + "1905" + dotCom,
		anySubDomain + "21cn" + dotCom,
		anySubDomain + "2288" + dotOrg,
		anySubDomain + "3322" + dotOrg,
		anySubDomain + "360doc" + dotCom,
		anySubDomain + "360buy" + dotCom,
		anySubDomain + "360buyimg" + dotCom,
		anySubDomain + "360safe" + dotCom,
		anySubDomain + "36kr" + dotCom,
		anySubDomain + "39" + dotNet,
		anySubDomain + "4399" + dotCom,
		anySubDomain + "51" + dotLa,
		anySubDomain + "51job" + dotCom,
		anySubDomain + "51jobcdn" + dotCom,
		anySubDomain + "71" + dotAm,
		anySubDomain + "acfun" + dotTv,
		anySubDomain + "aicdn" + dotCom,
		anySubDomain + "alibaba" + dotCom,
		anySubDomain + "alicdn" + dotCom,
		anySubDomain + "aliimg.com" + dotCom,
		anySubDomain + "alipay" + dotCom,
		anySubDomain + "alipayobjects" + dotCom,
		anySubDomain + "aliyun" + dotCom,
		anySubDomain + "aliyuncdn" + dotCom,
		anySubDomain + "aliyuncs" + dotCom,
		anySubDomain + "anjuke" + dotCom,
		anySubDomain + "babytree" + dotCom,
		anySubDomain + "baidu" + dotCom,
		anySubDomain + "baiducontent" + dotCom,
		anySubDomain + "baifendian" + dotCom,
		anySubDomain + "baike" + dotCom,
		anySubDomain + "baixing" + dotCom,
		anySubDomain + "bdimg" + dotCom,
		anySubDomain + "bdstatic" + dotCom,
		anySubDomain + "bilibili" + dotCom,
		anySubDomain + "bitauto" + dotCom,
		anySubDomain + "ccb" + dotCom,
		anySubDomain + "cctv" + dotCom,
		anySubDomain + "cctvpic" + dotCom,
		anySubDomain + "cdn20" + dotCom,
		anySubDomain + "ch" + dotCom,
		anySubDomain + "china" + dotCom,
		anySubDomain + "chinacache" + dotCom,
		anySubDomain + "chinacache" + dotNet,
		anySubDomain + "chinamobile" + dotCom,
		anySubDomain + "chinaz" + dotCom,
		anySubDomain + "chuangxin" + dotCom,
		anySubDomain + "clouddn" + dotCom,
		anySubDomain + "cmbchina" + dotCom,
		anySubDomain + "cnbeta" + dotCom,
		anySubDomain + "cnbetacdn" + dotCom,
		anySubDomain + "cnblogs" + dotCom,
		anySubDomain + "cnepub" + dotCom,
		anySubDomain + "cnzz" + dotCom,
		anySubDomain + "csdn" + dotNet,
		anySubDomain + "ctrip" + dotCom,
		anySubDomain + "dangdang" + dotCom,
		anySubDomain + "diandian" + dotCom,
		anySubDomain + "dianping" + dotCom,
		anySubDomain + "docin" + dotCom,
		anySubDomain + "donews" + dotCom,
		anySubDomain + "douban" + dotCom,
		anySubDomain + "duowan" + dotCom,
		anySubDomain + "eastday" + dotCom,
		anySubDomain + "fanli" + dotCom,
		anySubDomain + "fhldns" + dotCom,
		anySubDomain + "hao123" + dotCom,
		anySubDomain + "hao123img" + dotCom,
		anySubDomain + "haosou" + dotCom,
		anySubDomain + "hexun" + dotCom,
		anySubDomain + "huanqiu" + dotCom,
		anySubDomain + "hupu" + dotCom,
		anySubDomain + "iask" + dotCom,
		anySubDomain + "iqiyi" + dotCom,
		anySubDomain + "jd" + dotCom,
		anySubDomain + "jia" + dotCom,
		anySubDomain + "kouclo" + dotCom,
		anySubDomain + "ku6" + dotCom,
		anySubDomain + "ku6cdn" + dotCom,
		anySubDomain + "ku6img" + dotCom,
		anySubDomain + "lady8844" + dotCom,
		anySubDomain + "leiphone" + dotCom,
		anySubDomain + "letv" + dotCom,
		anySubDomain + "lvmama" + dotCom,
		anySubDomain + "meika360" + dotCom,
		anySubDomain + "meituan" + dotCom,
		anySubDomain + "mi" + dotCom,
		anySubDomain + "miaozhen" + dotCom,
		anySubDomain + "ngacn" + dotCc,
		anySubDomain + "oeeee" + dotCom,
		anySubDomain + "oschina" + dotNet,
		anySubDomain + "paipai" + dotCom,
		anySubDomain + "pps" + dotTv,
		anySubDomain + "qidian" + dotCom,
		anySubDomain + "qiniu" + dotCom,
		anySubDomain + "qiniudn" + dotCom,
		anySubDomain + "qiniudns" + dotCom,
		anySubDomain + "qiyi" + dotCom,
		anySubDomain + "qiyipic" + dotCom,
		anySubDomain + "qq" + dotCom,
		anySubDomain + "qqmail" + dotCom,
		anySubDomain + "qunar" + dotCom,
		anySubDomain + "qzone" + dotCom,
		anySubDomain + "renren" + dotCom,
		anySubDomain + "smzdm" + dotCom,
		anySubDomain + "sohu" + dotCom,
		anySubDomain + "sogou" + dotCom,
		anySubDomain + "soso" + dotCom,
		anySubDomain + "stockstar" + dotCom,
		anySubDomain + "suning" + dotCom,
		anySubDomain + "tanx" + dotCom,
		anySubDomain + "tao123" + dotCom,
		anySubDomain + "taobao" + dotCom,
		anySubDomain + "taobaocdn" + dotCom,
		anySubDomain + "tencent" + dotCom,
		anySubDomain + "tenpay" + dotCom,
		anySubDomain + "tmall" + dotCom,
		anySubDomain + "tudou" + dotCom,
		anySubDomain + "unionpay" + dotCom,
		anySubDomain + "unionpaysecure" + dotCom,
		anySubDomain + "upyun" + dotCom,
		anySubDomain + "upaiyun" + dotCom,
		anySubDomain + "vip" + dotCom,
		anySubDomain + "weibo" + dotCom,
		anySubDomain + "weiyun" + dotCom,
		anySubDomain + "xiami" + dotCom,
		anySubDomain + "xinhuanet" + dotCom,
		anySubDomain + "xueqiu" + dotCom,
		anySubDomain + "xunlei" + dotCom,
		anySubDomain + "xywy" + dotCom,
		anySubDomain + "yaolan" + dotCom,
		anySubDomain + "yingjiesheng" + dotCom,
		anySubDomain + "yhd" + dotCom,
		anySubDomain + "youboy" + dotCom,
		anySubDomain + "youku" + dotCom,
		anySubDomain + "zhihu" + dotCom,
	}
)

func init() {
	for _, pattern := range regexpDomains {
		matcher, err := NewRegexpDomainMatcher(pattern)
		if err != nil {
			panic(err)
		}
		compiledMatchers = append(compiledMatchers, matcher)
	}
}
