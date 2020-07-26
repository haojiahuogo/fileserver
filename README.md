# fileserver
## 漫画Go语言 文件服务


## 运行
go mod init fileserver

## 请求参数
formData = multipart/form-data
file = 图片

zipstr="[{"mode":"Fill", "fillType":"TopRight","w":160,"h":160,"name":"simg"},{"mode":"Fit","w":800,"h":10000,"name":"bimg"}]"

mode字段为Fill时, 可用新参数fillType=Center|TopLeft|Top|TopRight|Left|Right|BottomLeft|Bottom|BottomRight 如果不传入fillType,默认为Center

name=simg|bimg  返回的(缩略图/大图)名称

![](http://www.baidu.com/img/bdlogo.gif)  

# nodejs使用示例
文件服务器
//1.html页面部分使用
``` jQuery
jQuery('#imgUpload').click().fileupload({
    dataType: 'json',
    formData:{"zipstr":'[{"mode":"Fit","w":100,"h":100}]'}, //图片上传参数设置
    url: "/tools/imgupload/upload?type=user&t=" + new Date().getTime(),//文件上传地址
    done: function (e, result) {
        if (result.result.errno == 0) {
            //..业务逻辑
        } else {
            jQuery.messager.alert('提示', result.result.errmsg);
        }
    }
});

//2.js处理部分
async uploadtonetAction() {
    let fPath = file.path;
    //以下请求地址建议配置到config
    var postImgUrl = "http://127.0.0.1:8080/v1/image"; //请求上传地址（内网）
    var params = this.post("zipImage");
    let req = think.promisify(request.post);
    let options = {
        url: postImgUrl,
        method: "post",
        headers: {"zipstr": params},
        formData: {
            file: fs.createReadStream(fPath)
        }
    };
    let res = await req(options);
    let result = JSON.parse(res.body);
    //..业务逻辑
}
```

# Net使用示例

``` NET
//以下请求地址建议配置到web.config
string reqUrl = "http://127.0.0.1:8080/v1/image";   //请求上传地址（内网）
string localPath = @"D:\Image\6.jpg";                   //本地要上传的图片地址

WebClient web = new WebClient();
//请求参数设置
JArray reqParams = new JArray();
JObject param1 = new JObject();
param1["mode"] = "Fit";
param1["w"] = 700;
param1["h"] = 10000;
reqParams.Add(param1);
JObject param2 = new JObject();
param2["mode"] = "Fill";
param2["w"] = 100;
param2["h"] = 100;
reqParams.Add(param2);
web.Headers.Add("zipstr", reqParams.ToString(Newtonsoft.Json.Formatting.None));
byte[] res = web.UploadFile(reqUrl, localPath); //执行请求
//返回结果展示
JObject resJson = JObject.Parse(System.Text.Encoding.UTF8.GetString(res));
StringBuilder sb = new StringBuilder();
for (int i = 0; i < resJson.Count; i++)
{
    sb.Append(string.Format("{0}{1}", imgUrl, resJson[i.ToString()]["Uri"])+"<br>");
}
Response.Write(sb.ToString());
```