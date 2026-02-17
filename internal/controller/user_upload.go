package controller

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	filePath = "/books/epub"
	sshIp    = "182.92.1.221:22"
)

func UploadFile(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)
	// 获取表单文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}

	// 安全处理文件名
	filename := filepath.Base(file.Filename)

	// 校验文件类型
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".png" && ext != ".webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型"})
		return
	}

	// 保存文件到服务器指定目录
	saveDir := "" // 你的服务器目录
	savePath := filepath.Join(saveDir, filename)

	fmt.Printf("savePath: %s\n", savePath)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 返回文件路径（前端可以拼公网访问 URL）
	c.JSON(http.StatusOK, gin.H{
		"path": savePath,
	})
}

type SftpResonse struct {
	Data      []byte `json:"data"`
	RemoteUrl string `json:"remote_url"`
}

// 使用sftp 将获取到的文件保存到远程服务器
func UploadFileSFTP(ctx *gin.Context) (SftpResonse, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return SftpResonse{}, fmt.Errorf("获取文件信息失败")
	}
	filename := filepath.Base(file.Filename)
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".epub" {
		return SftpResonse{}, fmt.Errorf("不支持的文件格式")
	}

	srcFile, _ := file.Open()

	data, err := io.ReadAll(srcFile)
	if err != nil {
		return SftpResonse{}, fmt.Errorf("读取文件失败")
	}

	defer srcFile.Close()

	// SSH 客户端配置
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("207710cwL"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshConn, _ := ssh.Dial("tcp", sshIp, config)
	defer sshConn.Close()

	sftpClient, _ := sftp.NewClient(sshConn)
	defer sftpClient.Close()

	// 保存文件到远程路径
	remotePath := filepath.Join(filePath, filepath.Base(file.Filename))
	dstFile, _ := sftpClient.Create(remotePath)
	defer dstFile.Close()

	io.Copy(dstFile, srcFile)
	return SftpResonse{data, remotePath}, nil
}
