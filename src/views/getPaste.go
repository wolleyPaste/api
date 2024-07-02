package views

import (
	"wolley-api/src/db"

	"github.com/gin-gonic/gin"
)

func GetPaste(c *gin.Context) {
	id := c.Request.Header.Get("id")
	if len(id) != 8 {
		c.JSON(406, gin.H{
			"message": "unvalid id scheme",
		})
		return
	}

	paste, exists := db.GetPaste(id)
	if !exists {
		c.JSON(404, gin.H{
			"message": "paste doesn't exist",
		})
		return
	}

	if paste.Password != "" {

		if c.Request.Header.Get("password") == paste.Password {
			paste.Password = ""
			c.JSON(200, gin.H{
				"message": "sucessfully fetched secured paste",
				"paste":   paste,
			})
			return
		} else {
			c.JSON(401, gin.H{
				"message": "incorrect password",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "sucessfully fetched paste",
		"paste":   paste,
	})
}
