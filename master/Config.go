package master

import(
    "io/ioutil"
    "encoding/json"
)

var(
    G_config *Config
)

type Config struct{
    ApiPort string `json:"apiPort"`
    ApiReadTimeout int `json:"apiReadTimeout"`
    ApiWriteTimeout int `json:"apiWriteTimeout"`
    EtcdEndpoints []string `json:"etcdEndpoints"`
    EtcdDialTimeout int `json:"etcdDialTimeout"`
}

func InitConfig(filename string)(err error){
    var(
        content []byte
        config Config
    )

    if content,err = ioutil.ReadFile(filename);err != nil{
        return
    }

    if err = json.Unmarshal(content, &config);err != nil{
        return
    }

    G_config = &config

    return
}
