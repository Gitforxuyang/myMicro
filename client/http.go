package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type HttpClient interface {
	Do(ctx context.Context,method string,url string,data interface{},options ...HttpOptions) (string,error)
	DoRpc(ctx context.Context,method string,url string,data interface{},options ...HttpOptions)
}

type httpClient struct {
	client *http.Client
}

type HttpOptions struct {
	Headers map[string]string
}
func NewHttpClient() HttpClient{
	return &httpClient{client:&http.Client{Transport:&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   1 * time.Second, //连接建立时1s超时
			KeepAlive: 5 * time.Second, //每5s发送一个心跳包
		}).DialContext,
		MaxIdleConns:          100, //最大空闲连接数
		MaxIdleConnsPerHost: 80, //跟每个主机之间设置的最大连接数。 设置为80. 防止与高频server之间的连接占满了名额，其它连接无法使用到连接池
		MaxConnsPerHost:100, //跟每个主机之间建立的最大连接数
		IdleConnTimeout:       90 * time.Second, //空闲连接超时时间
		ExpectContinueTimeout: 1 * time.Second,
	},Timeout:time.Second*5}}
}
func (m *httpClient) Do(ctx context.Context,method string,url string,data interface{},options ...HttpOptions) (string,error){
	bts,err:=json.Marshal(data)
	if err!=nil{
		return "",err
	}
	body:=bytes.NewBuffer(bts)
	req,err:=http.NewRequest(method,url,body)
	req.Header.Set("Content-Type","application/json")
	for _,option:=range options{
		if option.Headers!=nil{
			for key,value:=range option.Headers{
				req.Header.Set(key,value)
			}
		}
	}
	if err!=nil{
		return "",err
	}
	resp,err:=m.client.Do(req.WithContext(ctx))
	if err!=nil{
		return "",err
	}
	bytes,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return "",err
	}
	str:=string(bytes)
	return str,nil
}
func (m *httpClient) DoRpc(ctx context.Context, method string, url string, data interface{}, options ...HttpOptions) {
	panic("implement me")
}
