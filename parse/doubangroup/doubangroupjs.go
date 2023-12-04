package doubangroup

import (
	"time"

	"github.com/donghc/crawler/collect"
)

var DoubangroupJSTask = &collect.TaskModel{
	Property: collect.Property{
		Name:     "js_find_douban_sun_room",
		WaitTime: 1 * time.Second,
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
