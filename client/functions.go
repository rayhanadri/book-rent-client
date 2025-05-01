package client

import (
	"encoding/json"
	"fmt"
	"io"
	"library-client/api"
	"library-client/model"
	"library-client/variables"
	"net/http"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/olekukonko/tablewriter"
)

func Register() {
	// Register a new user
	fmt.Println("Registering a new account...")
	// craete a new user struct to hold the registration data
	var newUser model.User
	// get user input for name, email, and password
	fmt.Print("Enter your name: ")
	_, err := fmt.Scanf("%[^\n]", &newUser.Name)
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}
	fmt.Print("Enter your email: ")
	_, err = fmt.Scanln(&newUser.Email)
	if err != nil {
		fmt.Println("Error reading email:", err)
		return
	}
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}
	newUser.Password = string(bytePassword)
	fmt.Println() // Print a newline after password input
	newUser.Role = "user"
	newUser.Status = "ACTIVE"
	newUser.Balance = 0
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	// Hit the API endpoint to register the new user
	resp, err := api.SendRequest(http.MethodPost, "users/register", newUser)
	if err != nil {
		fmt.Println("Error registering: Error sending registration request", err.Error())
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var registrationResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&registrationResponse)
	if err != nil {
		fmt.Println("Error registering: Error decoding response")
		return
	}

	// Check if the registration was successful
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		if registrationResponse.Status != 201 {
			fmt.Printf("Registration failed: %s\n", registrationResponse.Message)
		} else {
			fmt.Printf("Registration failed: %s\n", string(body))
		}

		return
	}
}

func Login() {
	// Login to the account
	fmt.Println("Logging into your account...")

	// get user input for email and password
	var email, password string
	fmt.Print("Enter your email: ")
	_, err := fmt.Scanln(&email)
	if err != nil {
		fmt.Println("Error reading email:", err)
		return
	}
	fmt.Print("Enter your password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	password = string(bytePassword)
	fmt.Println() // Print a newline after password input
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}

	// Check if the user is already logged in
	if variables.AccessToken != "" {
		fmt.Println("You are already logged in.")
		return
	}

	// create a new user struct to hold the login data
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest
	loginRequest.Email = email
	loginRequest.Password = password

	// Hit the API endpoint to login the user
	resp, err := api.SendRequest(http.MethodPost, "users/login", loginRequest)
	if err != nil {
		fmt.Println("Error logging in: Error sending login request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var loginResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		fmt.Println("Error logging in: Error decoding response")
		return
	}

	// Check if the login was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Login failed: %s\n", loginResponse.Message)
		return
	}

	// set variables for the logged in user
	variables.AccessToken = loginResponse.Data.(map[string]interface{})["accessToken"].(string)
	variables.RefreshToken = loginResponse.Data.(map[string]interface{})["refreshToken"].(string)
	dataBytes, err := json.Marshal(loginResponse.Data.(map[string]interface{})["user"])
	if err != nil {
		fmt.Println("Error logging in: Error marshalling user data")
		return
	}
	var user model.User
	err = json.Unmarshal(dataBytes, &user)
	if err != nil {
		fmt.Println("Error logging in: Error unmarshalling user data")
		return
	}

	variables.CurrentUser = user
	variables.IsLoggedIn = true
	fmt.Println("Login successful!")
}

func Logout() {
	// Logout from the account
	fmt.Println("Logging out of your account...")

	// Check if the user is already logged in
	if variables.AccessToken == "" {
		fmt.Println("You are not logged in.")
		return
	}

	// Remove the access token, refresh token, and user data
	variables.AccessToken = ""
	variables.RefreshToken = ""
	variables.CurrentUser = model.User{}
	variables.IsLoggedIn = false
	fmt.Println("Logout successful!")
}

func GetAllBooks() {
	// Fetch all books from the API
	fmt.Println("Fetching all books...")

	// Hit the API endpoint to fetch all books
	resp, err := api.SendRequest(http.MethodGet, "books", nil)
	if err != nil {
		fmt.Println("Error fetching books: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var booksResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&booksResponse)
	if err != nil {
		fmt.Println("Error fetching books: Error decoding response")
		return
	}

	// Check if the fetch was successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if booksResponse.Status != 200 {
			fmt.Printf("Failed to fetch books: %s\n", booksResponse.Message)
		} else {
			fmt.Printf("Failed to fetch books: %s\n", string(body))
		}
		return
	}

	var books []model.Book
	dataBytes, err := json.Marshal(booksResponse.Data)
	if err != nil {
		fmt.Println("Error fetching books: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &books)
	if err != nil {
		fmt.Println("Error fetching books: Error unmarshalling response data")
		return
	}

	// Display the books in a table format
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Author", "Year", "Category", "Stock", "Price"})
	for _, book := range books {
		table.Append([]string{
			fmt.Sprintf("%d", book.ID),
			book.Title,
			book.Author,
			book.PublishedAt.Format("2006"),
			book.Category,
			fmt.Sprintf("%d", book.Stock),
			fmt.Sprintf("%d", book.Price),
		})
	}
	table.Render()
	fmt.Println("All books displayed successfully!")
}

func GetBookByID(book_id int) {
	// Fetch book by ID from the API
	fmt.Println("Fetching book data...")

	// Hit the API endpoint to fetch book by ID
	resp, err := api.SendRequest(http.MethodGet, fmt.Sprintf("books/%d", book_id), nil)
	if err != nil {
		fmt.Println("Error fetching book: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var bookResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&bookResponse)
	if err != nil {
		fmt.Println("Error fetching book: Error decoding response")
		return
	}

	var book model.Book
	dataBytes, err := json.Marshal(bookResponse.Data)
	if err != nil {
		fmt.Println("Error fetching book: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &book)
	if err != nil {
		fmt.Println("Error fetching book: Error unmarshalling response data")
		return
	}

	// Check if the fetch was successful
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Author", "Year", "Category", "Stock", "Price"})
	table.Append([]string{
		fmt.Sprintf("%d", book.ID),
		book.Title,
		book.Author,
		book.PublishedAt.Format("2006"),
		book.Category,
		fmt.Sprintf("%d", book.Stock),
		fmt.Sprintf("%d", book.Price),
	})
	table.Render()
}

// Users
func GetUserDetails() {
	// Fetch user details from the API
	fmt.Println("Fetching user details...")

	// Hit the API endpoint to fetch user details
	resp, err := api.SendRequest(http.MethodGet, "users/me", nil)
	if err != nil {
		fmt.Println("Error fetching user: Error sending request")
		return
	}

	// Get the response from the API
	var userResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		fmt.Println("Error fetching user details: Error decoding response")
		return
	}

	// Check if the fetch was successful
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Token expired, refreshing token...")
			RefreshToken()
			GetUserDetails()
			return
		}

		if userResponse.Status != 200 {
			fmt.Printf("Failed to fetch user details: %s\n", userResponse.Message)
			return
		} else {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Failed to fetch user details: %s\n", string(body))
			return
		}
	}

	var user model.User
	dataBytes, err := json.Marshal(userResponse.Data)
	if err != nil {
		fmt.Println("Error fetching user: Error marshalling response data")
		return
	}

	err = json.Unmarshal(dataBytes, &user)
	if err != nil {
		fmt.Println("Error fetching user: Error unmarshalling response data")
		return
	}
	defer resp.Body.Close()

	// Display the user details in a table format
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Email", "Status", "Balance"})
	table.Append([]string{
		fmt.Sprintf("%d", user.ID),
		user.Name,
		user.Email,
		user.Status,
		fmt.Sprintf("%d", user.Balance),
	})
	table.Render()
}

func TopupBalance(AmountTopup int) {
	// Topup balance using payment gateway
	fmt.Println("Creating topup transaction...")

	// get user input for amount to topup
	var transaction model.Transaction
	transaction.TransactionType = "Topup"
	transaction.PaymentMethod = "Payment Gateway"
	transaction.Amount = AmountTopup
	transaction.Description = "Topup balance"

	// Hit the API endpoint to create a transaction
	resp, err := api.SendRequest(http.MethodPost, "transactions/create-transaction", transaction)
	if err != nil {
		fmt.Println("Error topping up balance: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var transResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&transResponse)
	if err != nil {
		fmt.Println("Error fetching transaction: Error decoding response")
		return
	}

	// Check if the transaction was successful
	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized && variables.AccessToken != "" {
			fmt.Println("Token expired, refreshing token...")
			RefreshToken()
			TopupBalance(AmountTopup)
			return
		}
		body, _ := io.ReadAll(resp.Body)
		if transResponse.Status != 201 {
			fmt.Printf("Failed to create transaction: %s\n", transResponse.Message)
		} else {
			fmt.Printf("Failed to create transaction: %s\n", string(body))
		}
		return
	}

	var trans model.Transaction
	dataBytes, err := json.Marshal(transResponse.Data)
	if err != nil {
		fmt.Println("Error fetching transaction: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &trans)
	if err != nil {
		fmt.Println("Error fetching transaction: Error unmarshalling response data")
		return
	}

	// Display the transaction details in a table format
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Transaction Type", "Payment Method", "Amount", "Status", "Description", "Invoice ID", "Invoice URL"})
	table.Append([]string{
		fmt.Sprintf("%d", trans.ID),
		trans.TransactionType,
		trans.PaymentMethod,
		fmt.Sprintf("%d", trans.Amount),
		trans.Status,
		trans.Description,
		trans.InvoiceID,
		trans.InvoiceURL,
	})
	table.Render()

	// Continue transaction to payment using the invoice URL
	fmt.Println("Topup transaction created, please continue to payment using the invoice URL")
	fmt.Println("Please confirm if you already paid the invoice")
	// Prompt user for confirmation
	fmt.Print("Invoice Paid? (y/n): ")
	var confirm string
	_, err = fmt.Scanln(&confirm)
	if err != nil {
		fmt.Println("Error reading confirmation:", err)
		return
	}
	// Check if the user confirmed the payment
	if confirm == "y" || confirm == "Y" {
		UpdateTransaction(trans.ID)
	} else {
		fmt.Println("Transaction not confirmed. Please try again later using Confirm Pending Transaction Menu.")
	}
}

func UpdateTransaction(transactionID int) {
	// Update transaction status
	fmt.Println("Checking transaction status...")

	// Hit the API endpoint to update transaction status
	resp, err := api.SendRequest(http.MethodPut, fmt.Sprintf("transactions/update-transaction/%d", transactionID), nil)
	if err != nil {
		fmt.Println("Error updating transaction: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var transResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&transResponse)
	if err != nil {
		fmt.Println("Error fetching transaction: Error decoding response")
	}

	// Check if the update was successful
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized && variables.AccessToken != "" {
			fmt.Println("Token expired, refreshing token...")
			RefreshToken()
			UpdateTransaction(transactionID)
			return
		}
		body, _ := io.ReadAll(resp.Body)
		if transResponse.Status != 200 {
			fmt.Printf("Failed to update transaction: %s\n", transResponse.Message)
		} else {
			fmt.Printf("Failed to update transaction: %s\n", string(body))
		}
		return
	}

	var trans model.Transaction
	dataBytes, err := json.Marshal(transResponse.Data)
	if err != nil {
		fmt.Println("Error fetching transaction: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &trans)
	if err != nil {
		fmt.Println("Error fetching transaction: Error unmarshalling response data")
		return
	}

	// Display the transaction details in a table format
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Transaction Type", "Payment Method", "Amount", "Status", "Description", "Invoice ID", "Invoice URL"})
	table.Append([]string{
		fmt.Sprintf("%d", trans.ID),
		trans.TransactionType,
		trans.PaymentMethod,
		fmt.Sprintf("%d", trans.Amount),
		trans.Status,
		trans.Description,
		trans.InvoiceID,
		trans.InvoiceURL,
	})
	table.Render()
}

func ConfirmPendingTransaction() {
	// Confirm pending transaction
	fmt.Println("Get Pending Transaction Data...")

	// Prompt user for transaction ID
	var transactionID int
	for {
		fmt.Print("Enter Transaction ID: ")
		_, err := fmt.Scanln(&transactionID)
		if err != nil {
			fmt.Println("Error reading Transaction ID:", err)
			continue
		}
		if transactionID <= 0 {
			fmt.Println("Invalid Transaction ID. Please enter a positive number.")
			continue
		}
		break
	}

	// Hit the API endpoint to fetch transaction by ID
	resp, err := api.SendRequest(http.MethodPut, fmt.Sprintf("transactions/update-transaction/%d", transactionID), nil)
	if err != nil {
		fmt.Println("Error fetching transaction: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var transResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&transResponse)
	if err != nil {
		fmt.Println("Error fetching transaction: Error decoding response")
		return
	}

	// Check if the fetch was successful
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized && variables.AccessToken != "" {
			fmt.Println("Token expired, refreshing token...")
			RefreshToken()
			ConfirmPendingTransaction()
			return
		}
		body, _ := io.ReadAll(resp.Body)

		if transResponse.Status != 200 {
			fmt.Printf("Failed to fetch transaction: %s\n", transResponse.Message)
		} else {
			fmt.Printf("Failed to fetch transaction: %s\n", string(body))
		}

		return
	}

	var trans model.Transaction
	dataBytes, err := json.Marshal(transResponse.Data)
	if err != nil {
		fmt.Println("Error fetching transaction: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &trans)
	if err != nil {
		fmt.Println("Error fetching transaction: Error unmarshalling response data")
		return
	}

	// Display the transaction details in a table format
	tablewriter := tablewriter.NewWriter(os.Stdout)
	tablewriter.SetHeader([]string{"ID", "Transaction Type", "Payment Method", "Amount", "Status", "Description", "Invoice ID", "Invoice URL"})
	tablewriter.Append([]string{
		fmt.Sprintf("%d", trans.ID),
		trans.TransactionType,
		trans.PaymentMethod,
		fmt.Sprintf("%d", trans.Amount),
		trans.Status,
		trans.Description,
		trans.InvoiceID,
		trans.InvoiceURL,
	})
	tablewriter.Render()
}

func RefreshToken() {
	// Refresh the access token
	fmt.Println("Refreshing token...")

	// Hit the API endpoint to refresh the token
	resp, err := api.RefreshToken()
	if err != nil {
		fmt.Println("Error refreshing token: Error sending request", err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to refresh token: %s\n", string(body))
		return
	}

	defer resp.Body.Close()

	var tokenResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		fmt.Println("Error refreshing token: Error decoding response")
		return
	}

	// Assign the new access token and refresh token to the variables
	variables.AccessToken = tokenResponse.Data.(map[string]interface{})["accessToken"].(string)
	variables.RefreshToken = tokenResponse.Data.(map[string]interface{})["refreshToken"].(string)
	dataBytes, err := json.Marshal(tokenResponse.Data.(map[string]interface{})["user"])
	if err != nil {
		fmt.Println("Error logging in: Error marshalling user data")
		return
	}
	var user model.User
	err = json.Unmarshal(dataBytes, &user)
	if err != nil {
		fmt.Println("Error logging in: Error unmarshalling user data")
		return
	}

	variables.CurrentUser = user
	variables.IsLoggedIn = true
}

func RentABook() {
	// Rent a book
	fmt.Println("Renting a book...")

	// Prompt user for book ID
	var book_id int
	for {
		fmt.Print("Enter Book ID: ")
		_, err := fmt.Scanln(&book_id)
		if err != nil {
			fmt.Println("Error reading Book ID:", err)
			continue
		}
		if book_id <= 0 {
			fmt.Println("Invalid Book ID. Please enter a positive number.")
			continue
		}
		break
	}

	// Hit the API endpoint to fetch book by ID
	resp, err := api.SendRequest(http.MethodGet, fmt.Sprintf("books/%d", book_id), nil)
	if err != nil {
		fmt.Println("Error fetching book: Error sending request")
		return
	}
	defer resp.Body.Close()

	var bookResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&bookResponse)
	if err != nil {
		fmt.Println("Error fetching book: Error decoding response")
		return
	}

	// Check if the fetch was successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if bookResponse.Status != 200 {
			fmt.Printf("Failed to fetch book: %s\n", bookResponse.Message)
		} else {
			fmt.Printf("Failed to fetch book: %s\n", string(body))
		}
		return
	}

	var book model.Book
	dataBytes, err := json.Marshal(bookResponse.Data)
	if err != nil {
		fmt.Println("Error fetching book: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &book)
	if err != nil {
		fmt.Println("Error fetching book: Error unmarshalling response data")
		return
	}

	// Prompt user for quantity and number of days
	var quantity int
	for {
		fmt.Print("Enter Quantity: ")
		_, err = fmt.Scanln(&quantity)
		if err != nil {
			fmt.Println("Error reading Quantity:", err)
			continue
		}
		if quantity <= 0 {
			fmt.Println("Invalid Quantity. Please enter a positive number.")
			continue
		}
		break
	}
	var numberOfDays int
	for {
		fmt.Print("Enter Number of Days: ")
		_, err = fmt.Scanln(&numberOfDays)
		if err != nil {
			fmt.Println("Error reading Number of Days:", err)
			continue
		}
		if numberOfDays <= 0 {
			fmt.Println("Invalid Number of Days. Please enter a positive number.")
			continue
		}
		break
	}
	var unitPrice int = book.Price
	if quantity > book.Stock {
		fmt.Println("Quantity exceeds available stock.")
		return
	}
	var totalPrice int = unitPrice * quantity * numberOfDays
	rentStartDate := time.Now()
	rentEndDate := time.Now().AddDate(0, 0, numberOfDays)
	createdAt := time.Now()
	updatedAt := time.Now()

	// create a new rent struct to hold the rent data
	var rent model.Rent
	rent.BookID = book.ID
	rent.UserID = variables.CurrentUser.ID
	rent.Quantity = quantity
	rent.TotalPrice = totalPrice
	rent.RentStartDate = rentStartDate
	rent.RentEndDate = rentEndDate
	rent.CreatedAt = createdAt
	rent.UpdatedAt = updatedAt

	// Hit the API endpoint to create a rent
	resp, err = api.SendRequest(http.MethodPost, "rents", rent)
	if err != nil {
		fmt.Println("Error renting book: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized && variables.AccessToken != "" {
			fmt.Println("Token expired, refreshing token...")
			RefreshToken()
			RentABook()
			return
		}
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to rent book: %s\n", string(body))
		return
	}

	var rentResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&rentResponse)
	if err != nil {
		fmt.Println("Error renting book: Error decoding response")
		return
	}
	var rentResp model.Rent
	dataBytes, err = json.Marshal(rentResponse.Data)
	if err != nil {
		fmt.Println("Error renting book: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &rentResp)
	if err != nil {
		fmt.Println("Error renting book: Error unmarshalling response data")
		return
	}

	// Check if the rent was successful
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Book ID", "Title", "Author", "Quantity", "Total Price", "Status", "Rent Start Date", "Rent End Date"})
	table.Append([]string{
		fmt.Sprintf("%d", rentResp.ID),
		rentResp.Book.Title,
		rentResp.Book.Author,
		fmt.Sprintf("%d", rentResp.BookID),
		fmt.Sprintf("%d", rentResp.Quantity),
		fmt.Sprintf("%d", rentResp.TotalPrice),
		rentResp.RentStatus,
		rentResp.RentStartDate.Format("2006-01-02"),
		rentResp.RentEndDate.Format("2006-01-02"),
	})
	table.Render()

	// Prompt user for payment method
	fmt.Println("Please select your payment method:")
	fmt.Println("1. App Balance")
	fmt.Println("2. Payment Gateway")
	var paymentChoice int
	var paymentMethod string
	for {
		fmt.Print("Your payment choice: ")
		_, err = fmt.Scanln(&paymentChoice)
		if err != nil {
			fmt.Println("Error reading payment choice:", err)
			continue
		}
		if paymentChoice == 1 {
			paymentMethod = "App Balance"
			break
		} else if paymentChoice == 2 {
			paymentMethod = "Payment Gateway"
			break
		} else {
			fmt.Println("Invalid payment choice. Please select 1 for App Balance or 2 for Payment Gateway.")
		}
	}

	// Create a new transaction struct to hold the transaction data
	transaction := model.Transaction{
		TransactionType: "Rent",
		PaymentMethod:   paymentMethod,
		Amount:          totalPrice,
		Description:     "Rent book",
		RentID:          rentResp.ID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Hit the API endpoint to create a transaction
	resp, err = api.SendRequest(http.MethodPost, "transactions/create-transaction", transaction)
	if err != nil {
		fmt.Println("Error saving transaction: Error sending request")
		return
	}
	defer resp.Body.Close()

	var transResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&transResponse)
	if err != nil {
		fmt.Println("Error saving transaction: Error decoding response")
		return
	}

	// Check if the transaction was successful
	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized && variables.AccessToken != "" {
			fmt.Println("Token expired, refreshing token...")
			RefreshToken()
			RentABook()
			return
		}
		body, _ := io.ReadAll(resp.Body)
		if transResponse.Status != 201 {
			fmt.Printf("Failed to save transaction: %s\n", transResponse.Message)
		} else {
			fmt.Printf("Failed to save transaction: %s\n", string(body))
		}
		return
	}

	var transResp model.Transaction
	dataBytes, err = json.Marshal(transResponse.Data)
	if err != nil {
		fmt.Println("Error saving transaction: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &transResp)
	if err != nil {
		fmt.Println("Error saving transaction: Error unmarshalling response data")
		return
	}

	// Display the transaction details in a table format
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Transaction Type", "Payment Method", "Amount", "Status", "Description", "Invoice ID", "Invoice URL"})
	table.Append([]string{
		fmt.Sprintf("%d", transResp.ID),
		transResp.TransactionType,
		transResp.PaymentMethod,
		fmt.Sprintf("%d", transResp.Amount),
		transResp.Status,
		transResp.Description,
		transResp.InvoiceID,
		transResp.InvoiceURL,
	})
	table.Render()

	// Prompt user for payment confirmation
	if transResp.PaymentMethod == "Payment Gateway" {
		fmt.Println("Rent transaction created, please continue to payment using the invoice URL")
		fmt.Println("Please confirm if you already paid the invoice")
		fmt.Print("Invoice Paid? (y/n): ")
		var confirm string
		_, err = fmt.Scanln(&confirm)
		if err != nil {
			fmt.Println("Error reading confirmation:", err)
			fmt.Println("Transaction not confirmed. Please try again later using Confirm Pending Transaction Menu.")
			return
		}
		if confirm == "y" || confirm == "Y" {
			UpdateTransaction(transResp.ID)
		} else {
			fmt.Println("Transaction not confirmed. Please try again later using Confirm Pending Transaction Menu.")
		}

	}

	// Hit the API endpoint to fetch rent by ID
	resp, err = api.SendRequest(http.MethodGet, fmt.Sprintf("rents/%d", rentResp.ID), nil)
	if err != nil {
		fmt.Println("Error fetching rent: Error sending request")
		return
	}
	defer resp.Body.Close()

	var rentResponse2 model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&rentResponse2)
	if err != nil {
		fmt.Println("Error fetching rent: Error decoding response")
		return
	}

	// Check if the fetch was successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if rentResponse2.Status != 200 {
			fmt.Printf("Failed to fetch rent: %s\n", rentResponse2.Message)
		} else {
			fmt.Printf("Failed to fetch rent: %s\n", string(body))
		}
		return
	}

	var rentResp2 model.Rent
	dataBytes, err = json.Marshal(rentResponse2.Data)
	if err != nil {
		fmt.Println("Error fetching rent: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &rentResp2)
	if err != nil {
		fmt.Println("Error fetching rent: Error unmarshalling response data")
		return
	}

	// Display the rent details in a table format
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Book ID", "Title", "Author", "Quantity", "Total Price", "Status", "Rent Start Date", "Rent End Date"})
	table.Append([]string{
		fmt.Sprintf("%d", rentResp2.ID),
		rentResp2.Book.Title,
		rentResp2.Book.Author,
		fmt.Sprintf("%d", rentResp2.BookID),
		fmt.Sprintf("%d", rentResp2.Quantity),
		fmt.Sprintf("%d", rentResp2.TotalPrice),
		rentResp2.RentStatus,
		rentResp2.RentStartDate.Format("2006-01-02"),
		rentResp2.RentEndDate.Format("2006-01-02"),
	})
	table.Render()
}

func ReturnABook() {
	// Return a book
	fmt.Println("Starting return book...")

	// Prompt user for rent ID
	var rent_id int
	for {
		fmt.Print("Enter Rent ID: ")
		_, err := fmt.Scanln(&rent_id)
		if err != nil {
			fmt.Println("Error reading Rent ID:", err)
			continue
		}
		break
	}

	// Hit the API endpoint to fetch rent by ID
	resp, err := api.SendRequest(http.MethodPut, fmt.Sprintf("rents/return/%d", rent_id), nil)
	if err != nil {
		fmt.Println("Error fetching rent: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var rentResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&rentResponse)
	if err != nil {
		fmt.Println("Error fetching rent: Error decoding response")
		return
	}

	// Check if the fetch was successful
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized && variables.AccessToken != "" {
			fmt.Println("Token expired, refreshing token...")
			RefreshToken()
			ReturnABook()
			return
		}
		body, _ := io.ReadAll(resp.Body)
		if rentResponse.Status != 200 {
			fmt.Printf("Failed to fetch rent: %s\n", rentResponse.Message)
		} else {
			fmt.Printf("Failed to fetch rent: %s\n", string(body))
		}
		return
	}

	var rent model.Rent
	dataBytes, err := json.Marshal(rentResponse.Data)
	if err != nil {
		fmt.Println("Error fetching rent: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &rent)
	if err != nil {
		fmt.Println("Error fetching rent: Error unmarshalling response data")
		return
	}

	// Display the rent details in a table format
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Book ID", "Title", "Author", "Quantity", "Total Price", "Status", "Rent Start Date", "Rent End Date"})
	table.Append([]string{
		fmt.Sprintf("%d", rent.ID),
		rent.Book.Title,
		rent.Book.Author,
		fmt.Sprintf("%d", rent.BookID),
		fmt.Sprintf("%d", rent.Quantity),
		fmt.Sprintf("%d", rent.TotalPrice),
		rent.RentStatus,
		rent.RentStartDate.Format("2006-01-02"),
		rent.RentEndDate.Format("2006-01-02"),
	})
	table.Render()
}

func GetAllRentHistory() {
	// Fetch all rent history from the API
	fmt.Println("Fetching all rent history...")

	// Hit the API endpoint to get all rent history
	resp, err := api.SendRequest(http.MethodGet, "rents", nil)
	if err != nil {
		fmt.Println("Error fetching rent history: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var rentResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&rentResponse)
	if err != nil {
		fmt.Println("Error fetching rent history: Error decoding response")
		return
	}

	// Check if the fetch was successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if rentResponse.Status != 200 {
			fmt.Printf("Failed to fetch rent history: %s\n", rentResponse.Message)
		} else {
			fmt.Printf("Failed to fetch rent history: %s\n", string(body))
		}
		return
	}

	var rents []model.Rent
	dataBytes, err := json.Marshal(rentResponse.Data)
	if err != nil {
		fmt.Println("Error fetching rent history: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &rents)
	if err != nil {
		fmt.Println("Error fetching rent history: Error unmarshalling response data")
		return
	}

	// Display the rent history in a table format
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Book ID", "Title", "Author", "Quantity", "Total Price", "Status", "Rent Start Date", "Rent End Date"})
	for _, rent := range rents {
		table.Append([]string{
			fmt.Sprintf("%d", rent.ID),
			rent.Book.Title,
			rent.Book.Author,
			fmt.Sprintf("%d", rent.BookID),
			fmt.Sprintf("%d", rent.Quantity),
			fmt.Sprintf("%d", rent.TotalPrice),
			rent.RentStatus,
			rent.RentStartDate.Format("2006-01-02"),
			rent.RentEndDate.Format("2006-01-02"),
		})
	}
	table.Render()
}

func GetRentHistoryByID() {
	// Fetch rent history by ID from the API
	fmt.Println("Fetching rent history by ID...")

	// Prompt user for rent ID
	var rent_id int
	for {
		fmt.Print("Enter Rent ID: ")
		_, err := fmt.Scanln(&rent_id)
		if err != nil {
			fmt.Println("Error reading Rent ID:", err)
			continue
		}
		if rent_id <= 0 {
			fmt.Println("Invalid Rent ID. Please enter a positive number.")
			continue
		}
		break
	}

	// Hit the API endpoint to fetch rent by ID
	resp, err := api.SendRequest(http.MethodGet, fmt.Sprintf("rents/%d", rent_id), nil)
	if err != nil {
		fmt.Println("Error fetching rent: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var rentResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&rentResponse)
	if err != nil {
		fmt.Println("Error fetching rent: Error decoding response")
		return
	}

	// Check if the fetch was successful
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized && variables.AccessToken != "" {
			fmt.Println("Token expired, refreshing token...")
			RefreshToken()
			GetRentHistoryByID()
			return
		}
		body, _ := io.ReadAll(resp.Body)
		if rentResponse.Status != 200 {
			fmt.Printf("Failed to fetch rent: %s\n", rentResponse.Message)
		} else {
			fmt.Printf("Failed to fetch rent: %s\n", string(body))
		}
		return
	}

	var rent model.Rent
	dataBytes, err := json.Marshal(rentResponse.Data)
	if err != nil {
		fmt.Println("Error fetching rent: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &rent)
	if err != nil {
		fmt.Println("Error fetching rent: Error unmarshalling response data")
		return
	}

	// Display the rent details in a table format
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Book ID", "Title", "Author", "Quantity", "Total Price", "Status", "Rent Start Date", "Rent End Date"})
	table.Append([]string{
		fmt.Sprintf("%d", rent.ID),
		fmt.Sprintf("%d", rent.BookID),
		rent.Book.Title,
		rent.Book.Author,
		fmt.Sprintf("%d", rent.Quantity),
		fmt.Sprintf("%d", rent.TotalPrice),
		rent.RentStatus,
		rent.RentStartDate.Format("2006-01-02"),
		rent.RentEndDate.Format("2006-01-02"),
	})
	table.Render()
	fmt.Println("Rent history displayed successfully!")
}

func GetAllTransactions() {
	// Fetch all transactions from the API
	fmt.Println("Fetching all transactions...")

	// Hit the API endpoint to get all transactions
	resp, err := api.SendRequest(http.MethodGet, "transactions", nil)
	if err != nil {
		fmt.Println("Error fetching transactions: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var transResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&transResponse)
	if err != nil {
		fmt.Println("Error fetching transactions: Error decoding response")
		return
	}

	// Check if the fetch was successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusUnauthorized && variables.AccessToken != "" {
			fmt.Println("Token expired, refreshing token...")
			RefreshToken()
			GetAllTransactions()
			return
		}
		if transResponse.Status != 200 {
			fmt.Printf("Failed to fetch transactions: %s\n", transResponse.Message)
		} else {
			fmt.Printf("Failed to fetch transactions: %s\n", string(body))
		}
		return
	}

	var transactions []model.Transaction
	dataBytes, err := json.Marshal(transResponse.Data)
	if err != nil {
		fmt.Println("Error fetching transactions: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &transactions)
	if err != nil {
		fmt.Println("Error fetching transactions: Error unmarshalling response data")
		return
	}

	// Display the transactions in a table format
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Transaction Type", "Payment Method", "Amount", "Status", "Description", "Invoice ID", "Invoice URL"})
	for _, transaction := range transactions {
		table.Append([]string{
			fmt.Sprintf("%d", transaction.ID),
			transaction.TransactionType,
			transaction.PaymentMethod,
			fmt.Sprintf("%d", transaction.Amount),
			transaction.Status,
			transaction.Description,
			transaction.InvoiceID,
			transaction.InvoiceURL,
		})
	}
	table.Render()
}

func GetTransactionByID() {
	// Fetch transaction by ID from the API
	fmt.Println("Fetching transaction by ID...")

	// Prompt user for transaction ID
	var transaction_id int
	for {
		fmt.Print("Enter Transaction ID: ")
		_, err := fmt.Scanln(&transaction_id)
		if err != nil {
			fmt.Println("Error reading Transaction ID:", err)
			continue
		}
		break
	}

	// Hit the API endpoint to fetch transaction by ID
	resp, err := api.SendRequest(http.MethodGet, fmt.Sprintf("transactions/%d", transaction_id), nil)
	if err != nil {
		fmt.Println("Error fetching transaction: Error sending request")
		return
	}
	defer resp.Body.Close()

	// Get the response from the API
	var transResponse model.ResponseApi
	err = json.NewDecoder(resp.Body).Decode(&transResponse)
	if err != nil {
		fmt.Println("Error fetching transaction: Error decoding response")
		return
	}

	// Check if the fetch was successful
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized && variables.AccessToken != "" {
			fmt.Println("Token expired, refreshing token...")
			RefreshToken()
			GetTransactionByID()
			return
		}
		body, _ := io.ReadAll(resp.Body)
		if transResponse.Status != 200 {
			fmt.Printf("Failed to fetch transaction: %s\n", transResponse.Message)
		} else {
			fmt.Printf("Failed to fetch transaction: %s\n", string(body))
		}
		return
	}

	var transaction model.Transaction
	dataBytes, err := json.Marshal(transResponse.Data)
	if err != nil {
		fmt.Println("Error fetching transaction: Error marshalling response data")
		return
	}
	err = json.Unmarshal(dataBytes, &transaction)
	if err != nil {
		fmt.Println("Error fetching transaction: Error unmarshalling response data")
		return
	}

	// Display the transaction details in a table format
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Transaction Type", "Payment Method", "Amount", "Status", "Description", "Invoice ID", "Invoice URL"})
	table.Append([]string{
		fmt.Sprintf("%d", transaction.ID),
		transaction.TransactionType,
		transaction.PaymentMethod,
		fmt.Sprintf("%d", transaction.Amount),
		transaction.Status,
		transaction.Description,
		transaction.InvoiceID,
		transaction.InvoiceURL,
	})
	table.Render()
}
