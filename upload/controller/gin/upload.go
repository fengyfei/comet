package controller

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	mysql "github.com/fengyfei/comet/upload/model/mysql"
	"github.com/gin-gonic/gin"
)

var (
	fileMap = map[string]string{}
	picture = []string{".jpg", ".png", ".jpeg", ".gif", ".bmp"}
	video   = []string{".avi", ".wmv", ".mpg", ".mpeg", ".mpe", ".mov", ".rm", ".ram", ".swf", ".mp4", ".rmvb", ".asf", ".divx", ".vob"}
	fileDir = filePath()
)

var (
	errRequest = errors.New("Request is not post method")
	erruserID  = errors.New("userID invalid")
)

const (
	// InvalidUID - userID invalid
	InvalidUID = 0
	// FileKey - key of the file
	FileKey = "file"
	// FileUploadDir - the root directory of the upload files
	FileUploadDir = "files"
	// PictureDir - save pictures file
	PictureDir = "picture"
	// VideoDir - save videos file
	VideoDir = "video"
	// OtherDir - files other than video and picture
	OtherDir = "other"
)

//UploadController -
type UploadController struct {
	db      *sql.DB
	BaseURL string
	getUID  func(c *gin.Context) (uint32, error)
}

// New -
func New(db *sql.DB, baseURL string, getUID func(c *gin.Context) (uint32, error)) *UploadController {
	return &UploadController{
		db:      db,
		BaseURL: baseURL,
		getUID:  getUID,
	}
}

// RegisterRouter -
func (u *UploadController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := mysql.CreateTable(u.db)
	if err != nil {
		log.Fatal(err)
	}

	err = checkDir(PictureDir, VideoDir, OtherDir)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/upload", u.upload)
}

func checkDir(path ...string) error {
	for _, name := range path {
		_, err := os.Stat(FileUploadDir + "/" + name)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(FileUploadDir+"/"+name, 0777)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Upload single file upload
func (u *UploadController) upload(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Error(errRequest)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	userID, err := u.getUID(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized})
		return
	}

	if userID == InvalidUID {
		c.Error(erruserID)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	file, header, err := c.Request.FormFile(FileKey)
	defer func() {
		file.Close()
		c.Request.MultipartForm.RemoveAll()
	}()

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
		return
	}

	MD5Str, err := MD5(file)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed})
		return
	}

	filePath, err := mysql.QueryByMD5(u.db, MD5Str)
	if err == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	if err != mysql.ErrNoRows {
		c.Error(err)
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict})
		return
	}

	fileSuffix := path.Ext(header.Filename)
	filePath = FileUploadDir + "/" + classifyBySuffix(fileSuffix) + "/" + MD5Str + fileSuffix

	err = copyFile(filePath, file)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	err = mysql.Insert(u.db, userID, filePath, MD5Str)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"status": http.StatusUnsupportedMediaType})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "URL": u.BaseURL + filePath})
}

func filePath() map[string]string {
	for _, suffix := range picture {
		fileMap[suffix] = PictureDir
	}

	for _, suffix := range video {
		fileMap[suffix] = VideoDir
	}

	return fileMap
}

func classifyBySuffix(suffix string) string {

	if dir := fileDir[suffix]; dir != "" {
		return dir
	}
	return OtherDir
}

// MD5 -
func MD5(file io.Reader) (string, error) {
	sum := md5.New()
	_, err := io.Copy(sum, file)
	if err != nil {
		return "", err
	}

	MD5Str := hex.EncodeToString(sum.Sum(nil))
	return MD5Str, nil
}

func copyFile(path string, file io.Reader) error {
	cur, err := os.Create(path)
	defer cur.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(cur, file)
	return err
}
