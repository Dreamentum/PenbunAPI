package controller

import (
	"log"
	"net/http"

	"penbun.com/api/src/config"
	"penbun.com/api/src/model"

	"github.com/gin-gonic/gin"
)

//  Publisher
func GetPublishers(ctx *gin.Context) {
	rows, err := config.DB.Query("SELECT PublisherId, PublisherName, PublisherContactName, PublisherContactPhone, PublisherAddress, PublisherDistrict, PublisherProvince, PublisherZipCode, Description, Status, UpdateDate, UpdateBy FROM tb_Publishers ORDER BY PublisherId")
	if err != nil {
		log.Fatalln("[!][Publisher] SELECT failed:", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var publishers []model.PUBLISHER

	for rows.Next() {
		var p model.PUBLISHER

		err := rows.Scan(&p.ID, &p.Name, &p.ContactName, &p.ContactPhone, &p.Address, &p.District, &p.Province, &p.ZipCode,
			&p.Description, &p.Status, &p.UpdatedAt, &p.UpdateBy)
		if err != nil {
			log.Fatalln("[!][Publisher] Scan failed:", err.Error())
			return
		}
		publishers = append(publishers, p)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": publishers})
	log.Println("[=][Publisher] Selected all Publishers")
}
