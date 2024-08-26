package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Team struct {
    ID           uint   `gorm:"primaryKey" json:"id"`
    TeamName     string `json:"teamName"`
    TeamTelegram string `json:"teamTelegram"`
    MembersCount int    `json:"membersCount"`
    EventID      uint   `json:"-"`
}

type Event struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    Name        string `json:"name"`
    Date        time.Time `json:"date"`
    Description string    `json:"description"`
    TeamCount   int       `json:"teamCount"`
    Teams       []Team    `gorm:"foreignKey:EventID" json:"teams"`
}

func getEvents(c *gin.Context) {
    var events []Event
    result := db.Preload("Teams").Find(&events)
    if result.Error != nil {
        c.AbortWithError(http.StatusNotFound, result.Error)
        return
    }
    c.JSON(http.StatusOK, events)
}

func addEvent(c *gin.Context) {
    var event Event
    if err := c.ShouldBindJSON(&event); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.Create(&event)
    c.JSON(http.StatusOK, event)
}

func addTeam(c *gin.Context) {
    var team Team
    if err := c.ShouldBindJSON(&team); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.Create(&team)
    c.JSON(http.StatusOK, team)
}

func deleteEvent(c *gin.Context) {
    id := c.Params.ByName("id")
    var event Event
    db.Where("id = ?", id).Delete(event)
    c.JSON(http.StatusOK, gin.H{"id #" + id: "deleted"})
}

func deleteTeam(c *gin.Context) {
    id := c.Params.ByName("id")
    var team Team
    db.Where("id = ?", id).Delete(team)
    c.JSON(http.StatusOK, gin.H{"id #" + id: "deleted"})
}

func main() {
    db, err := gorm.Open("sqlite3", "test.db")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()

    // Миграция базы данных
    db.AutoMigrate(&Event{}, &Team{})
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.GET("/events", getEvents)
    r.POST("/events", addEvent)
    r.POST("/teams", addTeam)
    r.DELETE("/events/:id", deleteEvent)
    r.DELETE("/teams/:id", deleteTeam)
    r.Run()
}