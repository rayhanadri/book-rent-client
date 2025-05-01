package cli

import (
	"fmt"
	"library-client/client"
	"library-client/model"
	"library-client/variables"
	"strconv"
	"strings"

	"bufio"

	"os"

	"github.com/olekukonko/tablewriter"
)

func Execute() {
	for {
		menu := ShowMenu()
		InputMenu(menu)
	}
}

func ShowMenu() []model.Menu {
	fmt.Println("===========================================")
	fmt.Println("Welcome to the Library Book Rent System")
	fmt.Println("===========================================")
	fmt.Println(variables.CurrentUser.Name)

	var tableMenuLoggedOut = []model.Menu{
		{No: 1, Name: "Register", Desc: "Register a new account"},
		{No: 2, Name: "Login", Desc: "Log into your account"},
		{No: 3, Name: "Get All Books", Desc: "Get all books"},
		{No: 4, Name: "Get Book By ID", Desc: "Get book by ID"},
		{No: 5, Name: "Exit", Desc: "Exit the application"},
	}

	var tableMenuLoggedIn = []model.Menu{
		{No: 1, Name: "User Details", Desc: "View your details"},
		{No: 2, Name: "Topup Balance", Desc: "Topup your balance"},

		{No: 3, Name: "Get All Books", Desc: "Get all books"},
		{No: 4, Name: "Get Book By ID", Desc: "Get book by ID"},

		{No: 5, Name: "Rent a Book", Desc: "Rent a book"},
		{No: 6, Name: "Return a Book", Desc: "Return a book by Rent ID"},
		{No: 7, Name: "Get All Rent History", Desc: "Get all rent history"},
		{No: 8, Name: "Get Rent History Detail", Desc: "Get a rent history by Rent ID"},

		{No: 9, Name: "Get All Transactions", Desc: "Get all users transactions"},
		{No: 10, Name: "Get Transaction Detail By ID", Desc: "Get a transaction by ID"},
		{No: 11, Name: "Confirm Pending Transaction by ID", Desc: "Confirm a pending transaction by ID"},

		// {No: 11, Name: "Refresh Token", Desc: "Refresh your token"},
		{No: 12, Name: "Logout", Desc: "Log out of your account"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"No", "Menu", "Description"})

	if variables.IsLoggedIn {
		for _, row := range tableMenuLoggedIn {
			table.Append([]string{fmt.Sprintf("%d", row.No), row.Name, row.Desc})
		}
		table.Render()

		return tableMenuLoggedIn
	} else {
		for _, row := range tableMenuLoggedOut {
			table.Append([]string{fmt.Sprintf("%d", row.No), row.Name, row.Desc})
		}
		table.Render()
		return tableMenuLoggedOut
	}

}

func InputMenu(Menu []model.Menu) {
	reader := bufio.NewReader(os.Stdin)
	var choice int = 0
	fmt.Println("===========================================")
	fmt.Print("Enter your choice: ")

	choiceStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	choice, err = strconv.Atoi(strings.TrimSpace(choiceStr))
	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return
	}

	validChoice := false
	for _, item := range Menu {
		if choice == item.No {
			validChoice = true
			break
		}
	}

	if !validChoice {
		fmt.Println("Invalid choice. Please select a valid menu option.")
		return
	}

	fmt.Printf("You selected: %s\n", Menu[choice-1].Name)

	if variables.IsLoggedIn {
		//handle input menu logged in
		switch choice {
		case 1:
			client.GetUserDetails()
		case 2:
			var AmoutTopup int
			fmt.Print("Enter Amount to Topup: ")

			AmoutTopupStr, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
			AmoutTopupStr = strings.TrimSpace(AmoutTopupStr)
			AmoutTopup, err = strconv.Atoi(AmoutTopupStr)
			if err == nil {
				client.TopupBalance(AmoutTopup)
			} else {
				fmt.Println("Invalid input. Please enter a valid amount.")
			}
		case 3:
			client.GetAllBooks()
		case 4:
			var bookID int
			fmt.Print("Enter Book ID: ")

			bookIDStr, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
			bookIDStr = strings.TrimSpace(bookIDStr)
			bookID, err = strconv.Atoi(bookIDStr)
			if err != nil {
				fmt.Println("Invalid input. Please enter a valid Book ID.")
				return
			}
			client.GetBookByID(bookID)
		case 5:
			client.RentABook()
		case 6:
			client.ReturnABook()
		case 7:
			client.GetAllRentHistory()
		case 8:
			client.GetRentHistoryByID()
		case 9:
			client.GetAllTransactions()
		case 10:
			client.GetTransactionByID()
		case 11:
			client.ConfirmPendingTransaction()
		case 12:
			client.Logout()
			fmt.Println("You have logged out successfully.")
		default:
			fmt.Println("Invalid choice. Please select a valid menu option.")
		}

	} else {
		// handle input menu logged out
		switch choice {
		case 1:
			client.Register()
			Execute()
		case 2:
			client.Login()
			Execute()
		case 3:
			client.GetAllBooks()
		case 4:
			var bookID int
			fmt.Print("Enter Book ID: ")
			_, err := fmt.Scanln(&bookID)
			if err == nil {
				client.GetBookByID(bookID)
			} else {
				fmt.Println("Invalid input. Please enter a valid Book ID.")
			}
		case 5:
			fmt.Println("Exiting the application...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please select a valid menu option.")
		}
	}
}
