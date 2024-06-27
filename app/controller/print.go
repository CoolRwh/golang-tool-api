package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jadefox10200/goprint"
	"golang-tool-api/app/models/response"
	"net/http"
)

var Print = new(printer)

type printer struct {
}

func (p printer) List(c *gin.Context) {
	printers, err := goprint.GetPrinterNames()
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err.Error()))
	}
	c.JSON(http.StatusOK, response.Ok(printers))
	return
}

func (p printer) PrintFile(c *gin.Context) {
	type Data struct {
		PrinterName string `json:"printerName"`
		Num         int    `json:"num"`
		FilePath    string `json:"filePath"`
	}
	var Param Data
	err := c.ShouldBindJSON(&Param)
	if err != nil {
		c.JSON(http.StatusOK, response.ResultCustom(&response.BusinessError{
			Code: response.RequestParamError,
			Msg:  err.Error(),
		}))
		return
	}
	printerHarder, err := goprint.GoOpenPrinter(Param.PrinterName)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail("打印机状态异常！"+err.Error()))
		return
	}
	for i := 0; i <= Param.Num; i++ {
		//go func() {
		//	err = goprint.GoPrint(printerHarder, Param.FilePath)
		//	if err != nil {
		//		c.JSON(http.StatusOK, response.Fail(struct {
		//			Msg string `json:"msg"`
		//		}{Msg: err.Error()}))
		//		return
		//	}
		//}()
		err = goprint.GoPrint(printerHarder, Param.FilePath)
		if err != nil {
			c.JSON(http.StatusOK, response.Fail(err.Error()))
			return
		}
	}
	c.JSON(http.StatusOK, response.Ok("success"))
	return
}
