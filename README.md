# GoSQL

GoSQL is a terminal-based Minimalist PostgresSQL database management tool written in Go. It provides an intuitive interface for managing database connections and executing PostgresSQL queries.

## Features

- Database connection profile management
- SQL query execution
- Table data viewing and editing
- Schema browser
- Real-time data updates
- Secure password storage

## Installation

1. Make sure you have Go installed on your system
2. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/GoSQL.git
   ```
3. Navigate to the project directory:
   ```bash
   cd GoSQL
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Build the project:
   ```bash
   go build
   ```

## Usage

1. Run the application:
   ```bash
   ./GoSQL
   ```

2. Create a new database connection profile:
   - Fill in the connection details (Profile Name, Host, Port, Username, Password, Database Name)
   - Click "Connect" or press Enter to save the profile

3. Select a profile from the list to connect to the database

4. Use the query interface to execute SQL commands

## Keyboard Shortcuts

### Profile Management
- `Tab` / `Down Arrow`: Move to next field
- `Shift + Tab` / `Up Arrow`: Move to previous field
- `Ctrl + X`: Jump to buttons
- `n`: Create a new Profile
- `u`: Update a Profile
- `p`: Return Profile list
- `Enter`: Submit form

### Query Interface
- `Ctrl + T`: Focus on table tree
- `Ctrl + K`: Focus on query input
- `Ctrl + U`: Focus on data table
- `Ctrl + A`: Execute current query
- `Ctrl + D`: Toggle database list visibility
- `Ctrl + R`: Reset changes
- `Ctrl + S`: Submit data changes
- `Escape`: Close status modal
- `e`: Edit selected cell
- `s`: Save changes

### Table Navigation
- `Up/Down Arrow`: Move between rows
- `Left/Right Arrow`: Move between columns
- `1-9`: Quick navigation (type number + direction)
- `d`: Move down specified number of rows
- `u`: Move up specified number of rows
- `Ctrl + X`: Toggle node expansion in tree view

