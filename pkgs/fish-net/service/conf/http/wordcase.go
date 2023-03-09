package http

import (
	"fishnet/domain"
	"fishnet/glb"
	"fishnet/service/conf/pack"
	"fishnet/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateWordcase(c *gin.Context) {
	var req pack.CreateWordcaseRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	wordcases, err := _wordcaseUsecase.QueryWordcase(nil, &req.Group, &req.Key, 1, 0)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	if len(wordcases) > 0 {
		c.JSON(400, gin.H{
			"msg": "wordcase already exists",
		})
		return
	}
	wordcase := &domain.Wordcase{
		GroupName: req.Group,
		Key:       req.Key,
		Value:     req.Value,
		Order:     req.Order,
		Disable:   req.Disable,
		Remark:    req.Remark,
	}
	if err := _wordcaseUsecase.CreateWordcase([]*domain.Wordcase{wordcase}); err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
		"data": gin.H{
			"id": wordcase.ID,
		},
	})
}

func DeleteWordcase(c *gin.Context) {
	wordcaseIDStr := c.Param("wordcaseId")
	if wordcaseIDStr == "" {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	wordcaseId, err := strconv.ParseInt(wordcaseIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	err = _wordcaseUsecase.DeleteWordcase(wordcaseId)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func UpdateWordcase(c *gin.Context) {
	wordcaseIDStr := c.Param("wordcaseId")
	if wordcaseIDStr == "" {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	wordcaseId, err := strconv.ParseInt(wordcaseIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	var req pack.UpdateWordcaseRequest
	if err := c.BindJSON(&req); err != nil {
		glb.LOG.Info("bind json failed")
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	glb.LOG.Info(util.SPrettyLog(req))
	err = _wordcaseUsecase.UpdateWordcase(wordcaseId, &req.Value, &req.Order, &req.Disable, &req.Remark)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func QueryWordcase(c *gin.Context) {
	if wordcaseIDStr := c.Param("wordcaseId"); wordcaseIDStr != "" {
		wordcaseID, err := strconv.ParseInt(wordcaseIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"msg": "bad request",
			})
			return
		}
		wordcases, err := _wordcaseUsecase.QueryWordcase(&wordcaseID, nil, nil, 1, 0)
		if err != nil {
			c.JSON(500, gin.H{
				"msg": "internal server error",
			})
			return
		}
		if len(wordcases) == 0 {
			c.JSON(404, gin.H{
				"msg": "wordcase not found",
			})
			return
		}
		c.JSON(200, gin.H{
			"msg": "ok",
			"data": gin.H{
				"wordcases": pack.Value(wordcases[0]),
			},
		})
		return
	}

	glb.LOG.Info("aaa multi")
	var req pack.QueryWordcaseRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	wordcases, err := _wordcaseUsecase.QueryWordcase(nil, &req.Group, &req.Key, req.Limit, req.Offset)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
		"data": gin.H{
			"list":  pack.Groups(wordcases),
			"total": len(wordcases),
		},
	})
}
