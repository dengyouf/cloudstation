package cmd

import (
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dengyouf/cloudstation/store"
	"github.com/dengyouf/cloudstation/store/provider/aliyun"
	"github.com/spf13/cobra"
)

const (
	// BuckName todo
	defaultBuckName = "cloudstations"
	defaultEndpoint = "http://oss-cn-beijing.aliyuncs.com"
	defaultALIAK    = "LTAI5t6BoVne4iGc3gVmvHYZ"
	defaultALISK    = ""
)

var (
	buckName       string
	uploadFilePath string
	bucketEndpoint string
)

// uploadCmd represents the start command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "上传文件到中转站",
	Long:  `上传文件到中转站`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := getProvider()
		if err != nil {
			return err
		}
		if uploadFilePath == "" {
			return fmt.Errorf("upload file path is missing")
		}

		// 为了防止文件都堆在一个文件夹里面 无法查看
		// 我们采用日期进行编码
		day := time.Now().Format("20060102")

		// 为了防止不同用户同一时间上传相同的文件
		// 我们采用用户的主机名作为前置
		hn, err := os.Hostname()
		if err != nil {
			ipAddr := getOutBindIp()
			if ipAddr == "" {
				hn = "unknown"
			} else {
				hn = ipAddr
			}
		}

		fn := path.Base(uploadFilePath)
		ok := fmt.Sprintf("%s/%s/%s", day, hn, fn)
		err = p.UploadFile(buckName, ok, uploadFilePath)
		if err != nil {
			return err
		}
		return nil
	},
}

func getOutBindIp() string {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	addr := strings.Split(conn.LocalAddr().String(), ":")
	if len(addr) == 0 {
		return ""
	}

	return addr[0]
}

// func getSecretKeyFromInput() {
// 	fmt.Printf("请输入access key: ")
// 	fmt.Scanln(&aliSecretKey)
// }
func getSecretKeyFromInputV2() {
	prompt := &survey.Password{
		Message: "请输入access key: ",
	}
	survey.AskOne(prompt, &aliSecretKey)
}
func getProvider() (p store.Uploader, err error) {
	switch ossProvider {
	case "aliyun":
		fmt.Printf("上传云商: 阿里云[%s]\n", defaultEndpoint)
		if aliAccessKey == "" {
			aliAccessKey = defaultALIAK
			fmt.Printf("上传用户: %s\n", aliAccessKey)
			getSecretKeyFromInputV2()
		}
		if aliSecretKey == "" {
			aliSecretKey = defaultALISK
		}
		fmt.Printf("上传用户: %s\n", aliAccessKey)
		p, err = aliyun.NewUploader(bucketEndpoint, aliAccessKey, aliSecretKey)
		return
	case "qcloud":
		return nil, fmt.Errorf("not impl")
	default:
		return nil, fmt.Errorf("unknown oss privier options [aliyun/qcloud]")
	}
}

func init() {
	uploadCmd.PersistentFlags().StringVarP(&uploadFilePath, "file_path", "f", "", "upload file path")
	uploadCmd.PersistentFlags().StringVarP(&buckName, "bucket_name", "b", defaultBuckName, "upload oss bucket name")
	uploadCmd.PersistentFlags().StringVarP(&bucketEndpoint, "bucket_endpoint", "e", defaultEndpoint, "upload oss endpoint")
	RootCmd.AddCommand(uploadCmd)
}
