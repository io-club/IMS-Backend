package main

import (
	// "errors"
	"fmt"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct{
	gorm.Model
	Username string `form:"name" json:"name" binding:"required"`
	Age int	`form:"age" json:"age" binding:"min=8,max=99"`
	Password string `form:"password" json:"password" binding:"required"`
	Hobby string `form:"hobby" json:"hobby"`
}

var db *gorm.DB

func InitDB(){
	var err error
	dsn:="root:123@tcp(127.0.0.1:3306)/homework?charset=utf8mb4&parseTime=True&loc=Local"
	db,err=gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err != nil{
		panic(err)
	}

	db.AutoMigrate(
		&User{},
	)
}


func GetAllHandler(c *gin.Context){
	page:=c.DefaultQuery("page","1")
	limit:=c.DefaultQuery("limit","5")
	pages, _ := strconv.Atoi(page)
	limitsize, _ := strconv.Atoi(limit)
	offset := (pages - 1) * limitsize
	var users []User
	if re:=db.Offset(offset).Limit(limitsize).Find(&users);re.Error!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"msg":re.Error,
		})
		return
	}
	
	

	c.JSON(http.StatusOK,gin.H{
		"data":users,
	})
	
}

func GetOneHandler(c *gin.Context){
	var user User
	id := c.Param("id")
	if re:=db.First(&user,id);re.Error!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"msg":re.Error,
		})
		return
	}
	
	c.JSON(http.StatusOK,gin.H{
		"message":user,
	})
}

func selectHandler(c *gin.Context){
	var users []User
	age:=c.DefaultQuery("age","")
	
	if re:=db.Where("age >?",age).Find(&users);re.Error!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"msg":re.Error,
		})
		return
	}
	
	c.JSON(http.StatusOK,gin.H{
		"data":users,
	})

}
func CreateHandler(c *gin.Context){
	var user User

	if err:=c.ShouldBindJSON(&user);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return
	}
	var re *gorm.DB
	if re =db.Create(&user);re.Error!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":re.Error,
		})
		return
	}

	
	c.JSON(http.StatusOK,gin.H{
		"message":"create successfully",
	})
	

}

func UpdateHandler(c *gin.Context){
	
	id:=c.Param("id")

	var user User

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	fmt.Println(user)

	user.Username=c.Request.FormValue("username")
	user.Password=c.Request.FormValue("password")

	re:=db.Model(&user).Where("id = ?",id).Updates(map[string]interface{}{"username":user.Username,"password":user.Password})
	if re.RowsAffected>0{
		c.JSON(http.StatusOK,gin.H{
			"message":"update successfully",
			"data":user,
		})
	}

	
}

func DeleteHandler(c *gin.Context){

	id := c.Param("id")
	
	re:=db.Delete(&User{},id)
	if re.Error == gorm.ErrRecordNotFound {
          c.JSON(http.StatusNotFound, gin.H{
            "message": "Record not found",
        })
		return
    } 
	
	if re.RowsAffected > 0 {
        
        c.JSON(http.StatusOK, gin.H{
            "message": "Delete successfully",
        })
    } 
}


func main(){
	

	InitDB()

	r:=gin.Default()

	r.GET("/get",GetAllHandler)

	r.GET("/get/:id",GetOneHandler)

	r.GET("/get/select",selectHandler)

	r.POST("/create",CreateHandler)

	r.PUT("/update/:id",UpdateHandler)

	r.DELETE("/delete/:id",DeleteHandler)

	r.Run()
}


