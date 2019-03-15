/*
 * Revision History:
 *     Initial: 2019/03/14        Yang ChengKai
 */

package controller

import (
	"errors"
	"net/http"

	permission "github.com/fengyfei/comet/permission/model/mysql"
	"github.com/gin-gonic/gin"
)

var (
	//FirstURL -
	FirstURL      = "/api/v1/permission/addurl"
	errPermission = errors.New("Admin permission is wrong")
)

//CheckPermission -
func CheckPermission(c *Controller, getUID func(ctx *gin.Context) (uint32, error)) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		var check = false
		URLL := ctx.Request.URL.Path

		a, err := getUID(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusBadGateway, err)
			return
		}

		adRole, err := permission.AdminGetRoleMap(c.db, a)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		urlRole, err := permission.URLPermissions(c.db, &URLL)
		if err != nil {
			ctx.AbortWithError(http.StatusFailedDependency, err)
			return
		}

		lenthrole, err := permission.GetRoleIDMap(c.db)
		if err != nil {
			ctx.AbortWithError(http.StatusNotExtended, err)
			return
		}

		urlroleid, err := permission.URLPermissions(c.db, &FirstURL)
		if err != nil {
			ctx.AbortWithError(http.StatusVariantAlsoNegotiates, err)
			return
		}

		if len(lenthrole) != 0 && len(urlroleid) != 0 {
			for urlkey := range urlRole {
				for adkey := range adRole {
					if urlkey == adkey {
						check = true
					}
				}
			}

			if !check {
				ctx.AbortWithError(http.StatusForbidden, errPermission)
			}
		}
	}
}
