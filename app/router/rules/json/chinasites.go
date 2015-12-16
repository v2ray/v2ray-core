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
	dotIo        = "\\.io$"
	dotLa        = "\\.la$"
	dotNet       = "\\.net$"
	dotOrg       = "\\.org$"
	dotTv        = "\\.tv$"
)

var (
	compiledMatchers []*RegexpDomainMatcher
)

func init() {
	compiledMatchers = make([]*RegexpDomainMatcher, 0, 1024)

	regexpDomains := []string{
		dotCn,
		"\\.xn--fiqs8s$", /* .中国 */

		anySubDomain + "10010" + dotCom,
		anySubDomain + "115" + dotCom,
		anySubDomain + "123u" + dotCom,
		anySubDomain + "126" + dotCom,
		anySubDomain + "126" + dotNet,
		anySubDomain + "127" + dotNet,
		anySubDomain + "163" + dotCom,
		anySubDomain + "17173" + dotCom,
		anySubDomain + "17cdn" + dotCom,
		anySubDomain + "1905" + dotCom,
		anySubDomain + "21cn" + dotCom,
		anySubDomain + "2288" + dotOrg,
		anySubDomain + "3322" + dotOrg,
		anySubDomain + "35" + dotCom,
		anySubDomain + "360doc" + dotCom,
		anySubDomain + "360buy" + dotCom,
		anySubDomain + "360buyimg" + dotCom,
		anySubDomain + "360safe" + dotCom,
		anySubDomain + "36kr" + dotCom,
		anySubDomain + "39" + dotNet,
		anySubDomain + "4399" + dotCom,
		anySubDomain + "51" + dotLa,
		anySubDomain + "51cto" + dotCom,
		anySubDomain + "51job" + dotCom,
		anySubDomain + "51jobcdn" + dotCom,
		anySubDomain + "71" + dotAm,
		anySubDomain + "abchina" + dotCom,
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
		anySubDomain + "allyes" + dotCom,
		anySubDomain + "amap" + dotCom,
		anySubDomain + "anjuke" + dotCom,
		anySubDomain + "appinn" + dotCom,
		anySubDomain + "babytree" + dotCom,
		anySubDomain + "baidu" + dotCom,
		anySubDomain + "baiducontent" + dotCom,
		anySubDomain + "baifendian" + dotCom,
		anySubDomain + "baike" + dotCom,
		anySubDomain + "baixing" + dotCom,
		anySubDomain + "bankcomm" + dotCom,
		anySubDomain + "bankofchina" + dotCom,
		anySubDomain + "bdimg" + dotCom,
		anySubDomain + "bdstatic" + dotCom,
		anySubDomain + "bilibili" + dotCom,
		anySubDomain + "bitauto" + dotCom,
		anySubDomain + "bobo" + dotCom,
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
		anySubDomain + "cloudxns" + dotCom,
		anySubDomain + "cmbchina" + dotCom,
		anySubDomain + "cnbeta" + dotCom,
		anySubDomain + "cnbetacdn" + dotCom,
		anySubDomain + "cnblogs" + dotCom,
		anySubDomain + "cnepub" + dotCom,
		anySubDomain + "cnzz" + dotCom,
		anySubDomain + "coding" + dotNet,
		anySubDomain + "csdn" + dotNet,
		anySubDomain + "ctrip" + dotCom,
		anySubDomain + "dangdang" + dotCom,
		anySubDomain + "daocloud" + dotIo,
		anySubDomain + "diandian" + dotCom,
		anySubDomain + "dianping" + dotCom,
		anySubDomain + "docin" + dotCom,
		anySubDomain + "donews" + dotCom,
		anySubDomain + "douban" + dotCom,
		anySubDomain + "dpfile" + dotCom,
		anySubDomain + "duoshuo" + dotCom,
		anySubDomain + "duowan" + dotCom,
		anySubDomain + "eastday" + dotCom,
		anySubDomain + "emarbox" + dotCom,
		anySubDomain + "etao" + dotCom,
		anySubDomain + "fanli" + dotCom,
		anySubDomain + "fhldns" + dotCom,
		anySubDomain + "getui" + dotCom,
		anySubDomain + "hao123" + dotCom,
		anySubDomain + "hao123img" + dotCom,
		anySubDomain + "haosou" + dotCom,
		anySubDomain + "hexun" + dotCom,
		anySubDomain + "hichina" + dotCom,
		anySubDomain + "huanqiu" + dotCom,
		anySubDomain + "hupu" + dotCom,
		anySubDomain + "iask" + dotCom,
		anySubDomain + "iciba" + dotCom,
		anySubDomain + "idqqimg" + dotCom,
		anySubDomain + "ifanr" + dotCom,
		anySubDomain + "ijinshan" + dotCom,
		anySubDomain + "ipip" + dotNet,
		anySubDomain + "iqiyi" + dotCom,
		anySubDomain + "itjuzi" + dotCom,
		anySubDomain + "jd" + dotCom,
		anySubDomain + "jia" + dotCom,
		anySubDomain + "jianshu" + dotCom,
		anySubDomain + "jiasuhui" + dotCom,
		anySubDomain + "jisuanke" + dotCom,
		anySubDomain + "kaixin001" + dotCom,
		anySubDomain + "kanimg" + dotCom,
		anySubDomain + "kankanews" + dotCom,
		anySubDomain + "kf5" + dotCom,
		anySubDomain + "kouclo" + dotCom,
		anySubDomain + "koudai8" + dotCom,
		anySubDomain + "ku6" + dotCom,
		anySubDomain + "ku6cdn" + dotCom,
		anySubDomain + "ku6img" + dotCom,
		anySubDomain + "lady8844" + dotCom,
		anySubDomain + "leiphone" + dotCom,
		anySubDomain + "letv" + dotCom,
		anySubDomain + "lietou" + dotCom,
		anySubDomain + "lvmama" + dotCom,
		anySubDomain + "lxdns" + dotCom,
		anySubDomain + "mechina" + dotOrg,
		anySubDomain + "mediav" + dotCom,
		anySubDomain + "meika360" + dotCom,
		anySubDomain + "meilishuo" + dotCom,
		anySubDomain + "meishij" + dotNet,
		anySubDomain + "meituan" + dotCom,
		anySubDomain + "meizu" + dotCom,
		anySubDomain + "mi" + dotCom,
		anySubDomain + "miaozhen" + dotCom,
		anySubDomain + "mmstat" + dotCom,
		anySubDomain + "mop" + dotCom,
		anySubDomain + "mydrivers" + dotCom,
		anySubDomain + "netease" + dotCom,
		anySubDomain + "ngacn" + dotCc,
		anySubDomain + "oeeee" + dotCom,
		anySubDomain + "onlinesjtu" + dotCom,
		anySubDomain + "oschina" + dotNet,
		anySubDomain + "paipai" + dotCom,
		anySubDomain + "pchome" + dotNet,
		anySubDomain + "pingplusplus" + dotCom,
		anySubDomain + "pps" + dotTv,
		anySubDomain + "pubyun" + dotCom,
		anySubDomain + "qhimg" + dotCom,
		anySubDomain + "qidian" + dotCom,
		anySubDomain + "qiniu" + dotCom,
		anySubDomain + "qiniudn" + dotCom,
		anySubDomain + "qiniudns" + dotCom,
		anySubDomain + "qiyi" + dotCom,
		anySubDomain + "qiyipic" + dotCom,
		anySubDomain + "qtmojo" + dotCom,
		anySubDomain + "qq" + dotCom,
		anySubDomain + "qqmail" + dotCom,
		anySubDomain + "qunar" + dotCom,
		anySubDomain + "qunarzz" + dotCom,
		anySubDomain + "qzone" + dotCom,
		anySubDomain + "renren" + dotCom,
		anySubDomain + "ruby-china" + dotOrg,
		anySubDomain + "sanwen" + dotNet,
		anySubDomain + "segmentfault" + dotCom,
		anySubDomain + "shutcm" + dotCom,
		anySubDomain + "sina" + dotCom,
		anySubDomain + "sinaapp" + dotCom,
		anySubDomain + "sinaedge" + dotCom,
		anySubDomain + "sinaimg" + dotCom,
		anySubDomain + "sinajs" + dotCom,
		anySubDomain + "smzdm" + dotCom,
		anySubDomain + "sohu" + dotCom,
		anySubDomain + "sogou" + dotCom,
		anySubDomain + "soso" + dotCom,
		anySubDomain + "sspai" + dotCom,
		anySubDomain + "staticfile" + dotOrg,
		anySubDomain + "stockstar" + dotCom,
		anySubDomain + "suning" + dotCom,
		anySubDomain + "tanx" + dotCom,
		anySubDomain + "tao123" + dotCom,
		anySubDomain + "taobao" + dotCom,
		anySubDomain + "taobaocdn" + dotCom,
		anySubDomain + "tencent" + dotCom,
		anySubDomain + "tenpay" + dotCom,
		anySubDomain + "tiexue" + dotNet,
		anySubDomain + "tmall" + dotCom,
		anySubDomain + "tmcdn" + dotNet,
		anySubDomain + "tudou" + dotCom,
		anySubDomain + "tudouui" + dotCom,
		anySubDomain + "unionpay" + dotCom,
		anySubDomain + "unionpaysecure" + dotCom,
		anySubDomain + "upyun" + dotCom,
		anySubDomain + "upaiyun" + dotCom,
		anySubDomain + "v2ex" + dotCom,
		anySubDomain + "vip" + dotCom,
		anySubDomain + "weibo" + dotCom,
		anySubDomain + "weiyun" + dotCom,
		anySubDomain + "wrating" + dotCom,
		anySubDomain + "xiachufang" + dotCom,
		anySubDomain + "xiami" + dotCom,
		anySubDomain + "xiaomi" + dotCom,
		anySubDomain + "xinhuanet" + dotCom,
		anySubDomain + "xinshipu" + dotCom,
		anySubDomain + "xnpic" + dotCom,
		anySubDomain + "xueqiu" + dotCom,
		anySubDomain + "xunlei" + dotCom,
		anySubDomain + "xywy" + dotCom,
		anySubDomain + "yaolan" + dotCom,
		anySubDomain + "yesky" + dotCom,
		anySubDomain + "yigao" + dotCom,
		anySubDomain + "yihaodian" + dotCom,
		anySubDomain + "yihaodianimg" + dotCom,
		anySubDomain + "yingjiesheng" + dotCom,
		anySubDomain + "yhd" + dotCom,
		anySubDomain + "youboy" + dotCom,
		anySubDomain + "youku" + dotCom,
		anySubDomain + "yunba" + dotIo,
		anySubDomain + "yunshipei" + dotCom,
		anySubDomain + "yupoo" + dotCom,
		anySubDomain + "yy" + dotCom,
		anySubDomain + "zbjimg" + dotCom,
		anySubDomain + "zhihu" + dotCom,
		anySubDomain + "zhimg" + dotCom,
		anySubDomain + "zhubajie" + dotCom,
	}

	for _, pattern := range regexpDomains {
		matcher, err := NewRegexpDomainMatcher(pattern)
		if err != nil {
			panic(err)
		}
		compiledMatchers = append(compiledMatchers, matcher)
	}
}
