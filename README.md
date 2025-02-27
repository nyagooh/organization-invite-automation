# GitHub Organization Inviter

This Go project automates the process of inviting GitHub users to join a GitHub organization. It reads a list of usernames from a file, checks if they are already members or have pending invitations, and sends an invitation if they are not.

## Features

- **Environment Variable Support:** Loads configuration from a `.env` file using [godotenv](https://github.com/joho/godotenv).

- **GitHub API Integration:** Utilizes the [go-github](https://github.com/google/go-github) library for interacting with the GitHub API.

- **OAuth2 Authentication:** Authenticates requests to GitHub using an OAuth2 token.

- **Prevent Duplicate Invitations:** Checks if a user is already a member or has a pending invitation before sending a new invite.

- **Command-Line Interface:** Accepts the GitHub organization name and a file path (containing GitHub usernames) as command-line arguments.

- **Error Handling:** Gracefully handles errors at each step, including file access, API calls, and scanning input.

## Prerequisites

- **Go:** Version 1.17 or later.

- **GitHub Personal Access Token:** Ensure you have a token with the appropriate scopes (e.g., `admin:org` `read.user`) to invite users.

- **Dependencies:** The project uses Go modules to manage dependencies.

## Setup

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/nyagooh/organization-invite-automation.git
   cd organization-invite-automation
   ```
2. **Create a .env File:**
    In the root of your project, create a file named .env and add your GitHub token:
    ```bash
    GITHUB_TOKEN=your_github_personal_access_token
    ```
3. **Download Dependencies:**
   Run the following command to download all dependencies:
  
    ```bash
    go mod init <input go mod name>
    ```
## Usage
Compile and run the program by providing the GitHub organization name and the path to the file containing usernames:
  ```bash
  go run . <org-name> <txtfile>
  ```
  example
  ```bash
  go run . kisumujs username.txt
  ```
  The program will:

- Load environment variables.
-  Read the file line by line to get GitHub usernames.

-  Check if each user is already a member or has a pending invitation.

-  Send an invitation only if needed.

- Log the invitation status and any errors encountered.

## Code Overview
### main.go

1. Environment Setup:

     -  Loads environment variables from a .env file.

     - Reads command-line arguments for the organization name and file path.
2. File Reading:

    - Opens the specified file and reads it line by line using a buffered scanner.
3.  GitHub API Client:

    - Sets up OAuth2 authentication with the provided token.

     - Creates a GitHub client to perform API operations.
4.  Invitation Process:

     -  For each username, checks if the user is already a member or has a pending invitation using the helper function shouldInviteUser.

    -  If not, it fetches the user details and sends an invitation to the organization.
    -  Logs the status of each invitation and any errors that occur.

### Helper Function: shouldInviteUser

The function shouldInviteUser is used to verify whether a user should be invited. It performs two checks:

   - Membership Check: Uses client.Organizations.IsMember to determine if the user is already a member of the organization.

   - Pending Invitation Check: Uses client.Organizations.ListPendingOrgInvitations to check if there's already an invitation pending for the user.

## Error Handling

- Environment & File Errors: Logs a message if the .env file is not found or if the specified file cannot be opened.

- API Errors: Logs detailed errors if there are issues with fetching user details or sending invitations.

- Scanner Errors: Checks and logs errors that occur while reading the input file.

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.