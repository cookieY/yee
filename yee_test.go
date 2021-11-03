package yee

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func indexHandle(c Context) (err error) {
	return c.JSON(http.StatusOK, "ok")
}

func addRouter(y *Core) {
	y.GET("/", indexHandle)
}

func TestYee(t *testing.T) {
	y := New()
	addRouter(y)
	y.Run(":9999")
}

func TestRestApi(t *testing.T) {
	y := New()
	y.Restful("/", RestfulAPI{
		Get: func(c Context) (err error) {
			return c.String(http.StatusOK, "updated")
		},
		Post: func(c Context) (err error) {
			return c.String(http.StatusOK, "get it")
		},
	})
}

func TestDownload(t *testing.T) {
	y := New()
	y.GET("/", func(c Context) (err error) {
		return c.File("args.go")
	})
	y.Run(":9999")
}

func TestStatic(t *testing.T) {
	y := New()
	y.Static("/front", "dist")
	y.GET("/", func(c Context) error {
		return c.HTMLTpl(http.StatusOK, "./dist/index.html")
	})
	y.Run(":9999")
}

const ver = `alt-svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000,h3-Q050=":443"; ma=2592000,h3-Q046=":443"; ma=2592000,h3-Q043=":443"; ma=2592000,quic=":443"; ma=2592000; v="46,43"`

func TestH3(t *testing.T) {
	y := New()
	y.GET("/", func(c Context) (err error) {
		return c.JSON(http.StatusOK, "hello")
	})
	y.RunH3(":443", "henry.com+4.pem", "henry.com+4-key.pem")
}

// Setup a bare-bones TLS config for the server
func TestGenerateTLSConfig(t *testing.T) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)   //把 1 左移 128 位，返回给 big.Int
	serialNumber, _ := rand.Int(rand.Reader, max) //返回在 [0, max) 区间均匀随机分布的一个随机值
	subject := pkix.Name{                         //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
		Organization:       []string{"Manning Publications Co."},
		OrganizationalUnit: []string{"Books"},
		CommonName:         "Go Web Programming",
	}
	template := x509.Certificate{
		SerialNumber: serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature, //KeyUsage 与 ExtKeyUsage 用来表明该证书是用来做服务器认证的
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},               // 密钥扩展用途的序列
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048) //生成一对具有指定字位数的RSA密钥

	//CreateCertificate基于模板创建一个新的证书
	//第二个第三个参数相同，则证书是自签名的
	//返回的切片是DER编码的证书
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk) //DER 格式
	certOut, _ := os.Create("cert.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICAET", Bytes: derBytes})
	certOut.Close()
	keyOut, _ := os.Create("key.pem")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}
