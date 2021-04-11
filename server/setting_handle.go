package server

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xchan/mod"
	"sync/atomic"
)

func (h *Handler) Setting(c *gin.Context) {
	// 获取当前的上传类型

	c.HTML(200, "setting.html", gin.H{
		"storeType":    h.s.GetCurrentStoreType(), // 获取当前存储类型
		"local_domain": h.s.GetLocalDomain(),      // 如果有本地存储则获取本地域名+协议
		"qiniu_key":    h.s.OSS.Key,
		"qiniu_secret": h.s.OSS.Secret,
		"qiniu_bucket": h.s.OSS.Bucket,
		"qiniu_domain": h.s.OSS.Domain,
	})
}

// UpdateSetting 修改上传设置
func (h *Handler) UpdateSetting(c *gin.Context) {
	position := c.PostForm("position")
	localDomain := c.PostForm("local_domain")
	key := c.PostForm("qiniu_key")
	secret := c.PostForm("qiniu_secret")
	bucket := c.PostForm("qiniu_bucket")
	qiNiuDomain := c.PostForm("qiniu_domain")

	// 后端防抖,防止手快直接点了两下安装,修改了两次,为 1 时代表正在安装
	if atomic.LoadUint32(atomicLock) == 1 {
		// 已经开始了安装程序
		c.JSON(200, Response(400, "正在修改配置中,请稍后呢~", nil))
		return
	}
	// 修改中
	atomic.StoreUint32(atomicLock, 1)
	// 解锁
	defer atomic.StoreUint32(atomicLock, 0)

	// 参数校验
	switch position {

	case "local": // 本地
		if err := h.s.UpdateOrAddLocalSetting(localDomain); err != nil {
			c.JSON(200, Response(400, "修改本地存储失败"+err.Error(), nil))
			return
		}
		c.JSON(200, Response(200, "修改成功", nil))
		return
	case "qiniu": // 七牛
		if err := h.s.UpdateOrAddQiNiuSetting(&mod.QiNiuOSSConfig{
			Key:    key,
			Secret: secret,
			Bucket: bucket,
			Domain: qiNiuDomain,
		}); err != nil {
			c.JSON(200, Response(400, "修改本地存储失败"+err.Error(), nil))
			return
		}
		c.JSON(200, Response(200, "修改成功", nil))
		return
	}

}
