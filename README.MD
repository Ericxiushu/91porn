## Depend

go packages:
    
    "github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/robertkrimen/otto"
	"github.com/zyxar/argo/rpc"
	"golang.org/x/net/proxy"

download tool：aria2

## Config

config file: /91porn/conf/conf.json

```json
{
    "port": 8888,
    "domain": "91porn.com",
    "md5js": "https://91porn.com/js/md5.js",
    "aria2":{
        "uri":"http://127.0.0.1:6800/jsonrpc",
        "token":"abcd"
    },
    "pageType": "rf",
    "pages": {
        "hot": "当前最热",
        "rp": "最近得分",
        "long": "10分钟以上",
        "md": "本月讨论",
        "tf": "本月收藏",
        "mf": "收藏最多",
        "rf": "最近加精",
        "top": "本月最热"
    }
}
```

## API

http 1.1 GET

1、/map


2、/do?pt=hot&from=1&to=10


## Reference

[tx7do/91porn_spider](https://github.com/tx7do/91porn_spider)

[zzjzz9266a/91porn_php](https://github.com/zzjzz9266a/91porn_php)