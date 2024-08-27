package main

import (
    "net/http"
    "github.com/gin-contrib/cors"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Login    string `json:"userLogin"`
	Password string `json:"userPassword"`
	IsAdmin  bool   `json:"isUserAdmin"`
}

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

var db *gorm.DB

func main() {
    db, err := gorm.Open("sqlite3", "test.db")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()

    // Миграция базы данных
    db.AutoMigrate(&Event{}, &Team{})
    r := gin.Default()
    r.Use(cors.Default())
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.POST("/add_event", func(c *gin.Context) {
        var event Event
        if err := c.ShouldBindJSON(&event); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    
        if err := db.Create(&event).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    
        c.JSON(http.StatusCreated, event)
    })
    r.GET("/get_events", func(c *gin.Context) {
        var events []Event
        if err := db.Preload("Teams").Find(&events).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    
        c.JSON(http.StatusOK, events)
    })
    
    r.Run()
}