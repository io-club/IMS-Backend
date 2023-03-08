package http

import (
	"fishnet/domain"
	"fishnet/service/conf/pack"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateWordcase(c *gin.Context) {
	var req pack.CreateWordcaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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
	wordcaseId, err := strconv.ParseInt(c.Param("wordcaseId"), 10, 64)
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
	wordcaseId, err := strconv.ParseInt(c.Param("wordcaseId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	var req *pack.UpdateWordcaseRequest
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
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
	wordcaseIDStr := c.Query("wordcaseId")
	if wordcaseIDStr != "" {
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

	var req *pack.QueryWordcaseRequest
	if err := c.ShouldBindJSON(req); err != nil {
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
			"wordcases": pack.Groups(wordcases),
		},
	})
}
