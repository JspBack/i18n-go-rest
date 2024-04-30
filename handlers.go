package main

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func addFAQ(c *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Question   string `json:"question" validate:"required"`
		TrQuestion string `json:"trQuestion" validate:"required"`
		Answers    []struct {
			Title    string `json:"title"`
			TrTitle  string `json:"trTitle"`
			Answer   string `json:"answer"`
			TrAnswer string `json:"trAnswer"`
		} `json:"answers" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	// Validate the input using the validator package
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMsg string
		for _, e := range validationErrors {
			errorMsg += e.Field() + " is required. "
		}
		return c.Status(fiber.StatusBadRequest).SendString(errorMsg)
	}

	// Generate UUID for FAQ entry
	faqID := uuid.New()

	// Create FAQ entry
	faq := FAQ{
		ID:         faqID,
		Question:   req.Question,
		TrQuestion: req.TrQuestion,
	}

	// Parse and assign answers
	for _, ans := range req.Answers {
		faq.Answers = append(faq.Answers, Answer{
			FAQID:    faqID,
			Title:    ans.Title,
			TrTitle:  ans.TrTitle,
			Answer:   ans.Answer,
			TrAnswer: ans.TrAnswer,
		})
	}

	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create FAQ entry and related answers in a transaction
	if err := tx.Create(&faq).Error; err != nil {
		log.Println("Error creating FAQ:", err)
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Translate the response message based on the language header
	lang := c.Get("x-custom-lang")
	if lang == "" {
		lang = "en-US" // Default to English if no language header is provided
	}
	
	T := i18n.NewLocalizer(bundle, lang)
	message, err := TranslateMessage(T, "faq_added")
	if err != nil {
		log.Println("Error translating message:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	if err := createBackup(); err != nil {
		fmt.Println("Error creating backup:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.SendString(message)
}

func deleteFAQ(c *fiber.Ctx) error {
	// Get the ID from the URL parameters
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ID parameter is required")
	}

	// Parse the ID string into a UUID object
	faqID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID parameter")
	}

	// Delete the FAQ using the FAQ ID
	result := db.Where("id = ?", faqID).Delete(&FAQ{})
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).SendString("FAQ not found")
	}

	// Translate the response message based on the language header
	lang := c.Get("x-custom-lang")
	if lang == "" {
		lang = "en-US" // Default to English if no language header is provided
	}
	T := i18n.NewLocalizer(bundle, lang)
	message, err := TranslateMessage(T, "faq_deleted")
	if err != nil {
		log.Println("Error translating message:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// if err := createBackup(); err != nil {
	// 	fmt.Println("Error creating backup:", err)
	// 	return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	// }

	return c.SendString(message)
}

func getFAQ(c *fiber.Ctx) error {
	var faqs []FAQ
	if err := db.Preload("Answers").Find(&faqs).Error; err != nil {
		log.Println("Error retrieving FAQs:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	if err := createBackup(); err != nil {
		fmt.Println("Error creating backup:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(faqs)
}

func patchFAQ(c *fiber.Ctx) error {
	// Get the ID from the URL parameters
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ID parameter is required")
	}

	// Parse the ID string into a UUID object
	faqID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID parameter")
	}

	// Parse request body
	var req struct {
		Question   string `json:"question" validate:"required"`
		TrQuestion string `json:"trQuestion" validate:"required"`
		Answers    []struct {
			Title    string `json:"title"`
			TrTitle  string `json:"trTitle"`
			Answer   string `json:"answer"`
			TrAnswer string `json:"trAnswer"`
		} `json:"answers" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	// Validate the input using the validator package
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMsg string
		for _, e := range validationErrors {
			errorMsg += e.Field() + " is required. "
		}
		return c.Status(fiber.StatusBadRequest).SendString(errorMsg)
	}

	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update FAQ entry
	if err := tx.Model(&FAQ{}).Where("id = ?", faqID).Updates(FAQ{
		Question:   req.Question,
		TrQuestion: req.TrQuestion,
	}).Error; err != nil {
		log.Println("Error updating FAQ:", err)
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Delete existing answers
	if err := tx.Where("faq_id = ?", faqID).Delete(&Answer{}).Error; err != nil {
		log.Println("Error deleting existing answers:", err)
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Parse and assign answers
	for _, ans := range req.Answers {
		if err := tx.Create(&Answer{
			FAQID:    faqID,
			Title:    ans.Title,
			TrTitle:  ans.TrTitle,
			Answer:   ans.Answer,
			TrAnswer: ans.TrAnswer,
		}).Error; err != nil {
			log.Println("Error creating answer:", err)
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Translate the response message based on the language header
	lang := c.Get("x-custom-lang")
	if lang == "" {
		lang = "en-US" // Default to English if no language header is provided
	}
	T := i18n.NewLocalizer(bundle, lang)
	message, err := TranslateMessage(T, "faq_updated")
	if err != nil {
		log.Println("Error translating message:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	if err := createBackup(); err != nil {
		fmt.Println("Error creating backup:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.SendString(message)

}
