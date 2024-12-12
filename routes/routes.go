package routes

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Vansh3140/GOVault/drivers"
	"github.com/gofiber/fiber/v2"
)

// Address represents the structure for address details
type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

// User represents the structure for user details
type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

var (
	db, err = drivers.Connect() // Initialize the database connection
)

// copyNonEmptyFields copies non-zero fields from src to dest
func copyNonEmptyFields(src, dest interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	destVal := reflect.ValueOf(dest).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		destField := destVal.Field(i)

		if !srcField.IsZero() {
			destField.Set(srcField)
		}
	}
}

// GetAll retrieves all records from a specific collection
func GetAll(c *fiber.Ctx) error {
	collection := c.Params("collection")

	if collection == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Collection parameter is required",
		})
	}

	records, err := db.ReadAll(collection)
	if err != nil {
		return c.Status(500).SendString(string(err.Error()))
	}

	allRecords := []User{}

	for _, value := range records {
		record := User{}
		if err := json.Unmarshal([]byte(value), &record); err != nil {
			fmt.Printf("Error while unmarshaling data %v", err)
			continue
		}
		allRecords = append(allRecords, record)
	}

	return c.Status(200).JSON(allRecords)
}

// GetOne retrieves a specific record by resource name
func GetOne(c *fiber.Ctx) error {
	collection := c.Params("collection")
	resource := c.Params("resource")

	if collection == "" || resource == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing collection or resource",
		})
	}

	newUser := User{}

	err := db.Read(collection, resource, &newUser)
	if err != nil {
		return c.Status(500).SendString(string(err.Error()))
	}

	return c.Status(200).JSON(newUser)
}

// CreateOne creates a new record in the specified collection
func CreateOne(c *fiber.Ctx) error {
	collection := c.Params("collection")

	if collection == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing collection name",
		})
	}

	employeeRecord := User{}

	if err := json.Unmarshal(c.Body(), &employeeRecord); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Error while unmarshaling data",
		})
	}

	if err := db.Write(collection, employeeRecord.Name, User{
		Name:    employeeRecord.Name,
		Contact: employeeRecord.Contact,
		Company: employeeRecord.Company,
		Age:     employeeRecord.Age,
		Address: employeeRecord.Address,
	}); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while writing data in DBVault",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":        "created",
		"resource_name": employeeRecord.Name,
	})
}

// UpdateOne updates a specific record in a collection
func UpdateOne(c *fiber.Ctx) error {
	collection := c.Params("collection")
	resource := c.Params("resource")

	if collection == "" || resource == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Missing collection or resource values",
			"error":   err.Error(),
		})
	}

	newEmployee := User{}

	if err := json.Unmarshal(c.Body(), &newEmployee); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while unmarshaling data",
			"error":   err.Error(),
		})
	}

	oldEmployee := User{}

	err := db.Read(collection, resource, &oldEmployee)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while reading data from DBVault",
			"error":   err.Error(),
		})
	}

	if newEmployee.Name != "" && newEmployee.Name != oldEmployee.Name {
		if err := db.Delete(collection, oldEmployee.Name); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Unable to update the records",
				"error":   err.Error(),
			})
		}
	}

	copyNonEmptyFields(&newEmployee, &oldEmployee)

	if err := db.Write(collection, oldEmployee.Name, User{
		Name:    oldEmployee.Name,
		Contact: oldEmployee.Contact,
		Company: oldEmployee.Company,
		Age:     oldEmployee.Age,
		Address: oldEmployee.Address,
	}); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while updating data in DBVault",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":        "updated",
		"resource_name": oldEmployee.Name,
	})
}

// DeleteAll removes all records from a specific collection
func DeleteAll(c *fiber.Ctx) error {
	collection := c.Params("collection")

	if collection == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Missing collection name",
			"error":   err.Error(),
		})
	}

	if err := db.Delete(collection, ""); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Unable to delete the records",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":          "deleted",
		"collection_name": collection,
	})
}

// DeleteOne removes a specific record from a collection
func DeleteOne(c *fiber.Ctx) error {
	collection := c.Params("collection")
	resource := c.Params("resource")

	if collection == "" || resource == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Missing collection or resource name",
			"error":   err.Error(),
		})
	}

	if err := db.Delete(collection, resource); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Unable to delete the record",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":        "deleted",
		"resource_name": resource,
	})
}
