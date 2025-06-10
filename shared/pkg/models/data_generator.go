package models

import (
	"log"
	"strings"

	"gorm.io/gorm"
)

func GenerateCategories(db *gorm.DB) error {
	categories := []string{
		"Web Development",
		"Mobile Development",
		"Data Science",
		"Machine Learning",
		"Blockchain",
		"Cloud Computing",
		"Cybersecurity",
		"DevOps",
		"Artificial Intelligence",
		"Internet of Things (IoT)",
		"Game Development",
		"Augmented Reality (AR) and Virtual Reality (VR)",
		"API Development",
		"Software Testing",
		"Database Management",
		"UI/UX Design",
		"Content Management Systems (CMS)",
		"Search Engine Optimization (SEO)",
		"Digital Marketing",
		"Business Intelligence",
		"Project Management",
		"Agile Methodologies",
		"Version Control Systems",
		"Microservices Architecture",
		"Serverless Computing",
		"Edge Computing",
	}

	var existingCategories []Category
	if err := db.Find(&existingCategories).Error; err != nil {
		return err
	}

	existingCategoryNames := make(map[string]bool)
	for _, category := range existingCategories {
		existingCategoryNames[strings.ToLower(category.Name)] = true
	}
	for _, name := range categories {
		slug := GenerateSlug(name)
		if !existingCategoryNames[strings.ToLower(name)] {
			newCategory := Category{
				Name: name,
				Slug: slug,
			}
			if err := db.Create(&newCategory).Error; err != nil {
				return err
			}
		}
	}

	log.Println("✅ categories generated successfully")
	return nil
}

func GenerateSlug(name string) string {
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	return slug
}
