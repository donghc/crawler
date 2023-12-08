package doubangroupjs

import (
	"github.com/donghc/crawler/collect"
)

var (
	cookie = "bid=tC5jgShcUIU; ll=\"108288\"; __utmc=30149280; douban-fav-remind=1; ct=y; ap_v=0,6.0; dbcl2=\"167037167:j9y6vSfLpN8\"; ck=iZmW; push_doumail_num=0; __utma=30149280.1262193242.1695303769.1700921222.1700924224.9; __utmz=30149280.1700924224.9.2.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmt=1; __utmv=30149280.16703; push_noty_num=0; __utmb=30149280.11.5.1700924241581; frodotk_db=\"f912931b88658561eb37b95299309023\"; ps=y"
)
var DoubangroupJSTask = &collect.TaskModel{
	Property: collect.Property{
		Name:     "js_find_douban_sun_room",
		WaitTime: 1,
		MaxDepth: 5,
		Cookie:   cookie,
	},
	Root: `
		var arr = new Array();
 		for (var i = 0; i <= 50; i+=25) {
			var obj = {
               URL: "https://www.douban.com/group/beijingzufang/discussion?start=" + i, 
			   Priority: 1,
			   RuleName: "解析网站URL",
			   Method: "GET",
		   };
		   arr.push(obj);
		};
		AddJsReq(arr);
			`,
	Rules: []collect.RuleModel{
		{
			Name: "解析网站URL",
			ParseFunc: `
			ctx.ParseJSReg("解析阳台房","(https://www.douban.com/group/topic/[0-9a-z]+/)\"[^>]*>([^<]+)</a>");
			`,
		},
		{
			Name: "解析阳台房",
			ParseFunc: `
			ctx.OutputJS("<div class=\"topic-content\">[\\s\\S]*?阳台[\\s\\S]*?<div class=\"aside\">");
			`,
		},
	},
}
