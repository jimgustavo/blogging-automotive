package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type BlogPost struct {
	ID         int       `json:"id"`
	Category   string    `json:"category"`
	Title      string    `json:"title"`
	Picture    string    `json:"picture"`
	Summary    string    `json:"summary"`
	Author     string    `json:"author"`
	EditorData string    `json:"editor_data"`
	CreatedAt  time.Time `json:"created_at"`
}

var (
	db *sql.DB
)

func init() {
	// Initialize the database connection in an init function
	var err error
	db, err = sql.Open("postgres", "postgres://tavito:mamacita@localhost:5432/blog_automotive?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Check the database connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping the database:", err)
	}
}

func main() {
	defer db.Close()

	router := mux.NewRouter()

	// Serve static files from the "static" directory
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.HandleFunc("/blogposts", CreateBlogPost).Methods("POST")
	router.HandleFunc("/blogposts", GetBlogPosts).Methods("GET")
	router.HandleFunc("/blogposts/{id}", GetBlogPost).Methods("GET")       // New endpoint to get a single blog post
	router.HandleFunc("/blogposts/{id}", UpdateBlogPost).Methods("PUT")    // Endpoint to update a blog post
	router.HandleFunc("/blogposts/{id}", DeleteBlogPost).Methods("DELETE") // Endpoint to delete a blog post
	router.HandleFunc("/blogposts/filter", FilterBlogPosts).Methods("GET")
	router.HandleFunc("/reading-page.html", ReadingPageHandler).Methods("GET")
	//router.HandleFunc("/blogposts/filter", FilterBlogPosts).Methods("GET") // Route for filtering blog posts

	log.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", router)
}

func ReadingPageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract blog post ID from URL parameters
	postID := r.URL.Query().Get("id")

	// Check if postID is empty
	if postID == "" {
		http.Error(w, "Post ID is missing in the URL", http.StatusBadRequest)
		return
	}

	// Fetch blog post data based on the ID from the database
	_, err := getBlogPostByID(postID)
	if err != nil {
		// Handle error (e.g., blog post not found)
		http.Error(w, "Failed to fetch blog post", http.StatusNotFound)
		return
	}

	// Render the reading page with the blog post content
	// You may pass the blog post data to the template
	// for rendering the content in the HTML page
	// For example:
	// renderTemplate(w, "reading-page.html", post)
}

func getBlogPostByID(postID string) (BlogPost, error) {
	// Fetch blog post data from the database based on the ID
	// You'll need to implement this function to interact with your database
	// Return the fetched blog post and any error that occurred

	// For now, let's return an empty blog post and nil error
	return BlogPost{}, nil
}

func CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the JSON request body into a BlogPost struct
	var newPost BlogPost
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		log.Println("Failed to decode JSON:", err)
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	// Insert the new post data into the database
	_, err := db.Exec(
		"INSERT INTO posts (category, title, picture, summary, author, editor_data, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		newPost.Category, newPost.Title, newPost.Picture, newPost.Summary, newPost.Author, newPost.EditorData, time.Now(),
	)
	if err != nil {
		log.Println("Failed to insert data into the database:", err)
		http.Error(w, "Failed to insert data into the database", http.StatusInternalServerError)
		return
	}

	// Optionally, you can return the newly created post as JSON response
	json.NewEncoder(w).Encode(newPost)
}

func GetBlogPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Execute a SELECT query to retrieve all blog posts from the database
	rows, err := db.Query("SELECT id, category, title, picture, summary, author, editor_data, created_at FROM posts")
	if err != nil {
		log.Println("Failed to execute query:", err)
		http.Error(w, "Failed to retrieve data from the database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the retrieved blog posts
	var retrievedPosts []BlogPost

	// Iterate through the result rows and scan them into BlogPost structs
	for rows.Next() {
		var post BlogPost
		if err := rows.Scan(
			&post.ID,
			&post.Category,
			&post.Title,
			&post.Picture,
			&post.Summary,
			&post.Author,
			&post.EditorData,
			&post.CreatedAt,
		); err != nil {
			log.Println("Failed to scan row:", err)
			// Skip this row and continue to the next one
			continue
		}
		retrievedPosts = append(retrievedPosts, post)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Println("Error while iterating over database rows:", err)
		http.Error(w, "Error while iterating over database rows", http.StatusInternalServerError)
		return
	}

	// Encode and send the retrieved blog posts as JSON response
	if len(retrievedPosts) == 0 {
		// If no blog posts were retrieved, return an empty array
		w.WriteHeader(http.StatusOK) // Set HTTP status code to 200
		w.Write([]byte("[]"))        // Write an empty array to the response
		return
	}

	json.NewEncoder(w).Encode(retrievedPosts)
}

func GetBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the ID parameter from the URL
	params := mux.Vars(r)
	postID := params["id"]

	// Execute a SELECT query to retrieve the blog post from the database
	var post BlogPost
	err := db.QueryRow("SELECT id, category, title, picture, summary, author, editor_data, created_at FROM posts WHERE id = $1", postID).
		Scan(&post.ID, &post.Category, &post.Title, &post.Picture, &post.Summary, &post.Author, &post.EditorData, &post.CreatedAt)
	if err != nil {
		log.Println("Failed to retrieve data from the database:", err)
		http.Error(w, "Failed to retrieve data from the database", http.StatusInternalServerError)
		return
	}

	// Encode and send the retrieved blog post as JSON response
	json.NewEncoder(w).Encode(post)
}

func UpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the ID parameter from the URL
	params := mux.Vars(r)
	postID := params["id"]

	// Decode the JSON request body into a BlogPost struct
	var updatedPost BlogPost
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		log.Println("Failed to decode JSON:", err)
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	// Update the blog post in the database
	_, err := db.Exec(
		"UPDATE posts SET category = $1, title = $2, picture = $3, summary = $4, author = $5, editor_data = $6 WHERE id = $7",
		updatedPost.Category, updatedPost.Title, updatedPost.Picture, updatedPost.Summary, updatedPost.Author, updatedPost.EditorData, postID,
	)

	if err != nil {
		log.Println("Failed to update data in the database:", err)
		http.Error(w, "Failed to update data in the database", http.StatusInternalServerError)
		return
	}

	// Optionally, you can return the updated post as JSON response
	json.NewEncoder(w).Encode(updatedPost)
}

func DeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the ID parameter from the URL
	params := mux.Vars(r)
	postID := params["id"]

	// Delete the blog post from the database
	_, err := db.Exec("DELETE FROM posts WHERE id = $1", postID)
	if err != nil {
		log.Println("Failed to delete data from the database:", err)
		http.Error(w, "Failed to delete data from the database", http.StatusInternalServerError)
		return
	}

	// Return success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Blog post deleted successfully"))
}

func FilterBlogPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the category parameter from the URL query string
	category := r.URL.Query().Get("category")

	// Log the category value
	log.Println("Category:", category)

	// Build the SQL query based on the presence of the category parameter
	var query string
	var args []interface{}

	if category != "" {
		query = "SELECT id, category, title, picture, summary, author, editor_data, created_at FROM posts WHERE category = $1"
		args = append(args, category)
	} else {
		query = "SELECT id, category, title, picture, summary, author, editor_data, created_at FROM posts WHERE category = '' OR category IS NULL"
	}

	// Log the SQL query and arguments
	log.Println("Executing query:", query, "with args:", args)

	// Execute the SQL query
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println("Failed to execute query:", err)
		http.Error(w, "Failed to retrieve data from the database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the retrieved blog posts
	var retrievedPosts []BlogPost

	// Iterate through the result rows and scan them into BlogPost structs
	for rows.Next() {
		var post BlogPost
		if err := rows.Scan(
			&post.ID,
			&post.Category,
			&post.Title,
			&post.Picture,
			&post.Summary,
			&post.Author,
			&post.EditorData,
			&post.CreatedAt,
		); err != nil {
			log.Println("Failed to scan row:", err)
			// Skip this row and continue to the next one
			continue
		}
		retrievedPosts = append(retrievedPosts, post)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Println("Error while iterating over database rows:", err)
		http.Error(w, "Error while iterating over database rows", http.StatusInternalServerError)
		return
	}

	// Encode and send the retrieved blog posts as JSON response
	if len(retrievedPosts) == 0 {
		// If no blog posts were retrieved, return an empty array
		w.WriteHeader(http.StatusOK) // Set HTTP status code to 200
		w.Write([]byte("[]"))        // Write an empty array to the response
		return
	}

	json.NewEncoder(w).Encode(retrievedPosts)
}

/*
///////////////////////Curl Commands////////////////////
GetBlogPosts:
curl -X GET http://localhost:8080/blogposts

CreateBlogPost:
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "category": "Braking-System",
    "title": "Braking System Maintenance",
    "author": "Gustavo Ruiz",
    "editorData": "<p>For making a basic maintenance to the braking system, we're gonna follow these steps:</p><ul><li>Use Personal Protection Equipment</li><li>Check the manufacturer workshop manual</li><li>Follow the processes established in the workshop</li></ul><h3>Use PPE</h3><figure class=\"image\"><img src=\"https://www.safetyandhealthmagazine.com/ext/resources/images/news/PPE/work-PPE.jpg?t=1662648636&amp;width=768\" alt=\"NIOSH to host Equitable PPE Protections Workshop in November | 2022-09-08 | Safety+Health\"></figure><p>In the automotive service, we need to use Personal Protection Equipment (PPE) to ensure safety.</p><h3>Check the workshop manual</h3><p>Refer to the manufacturer's workshop manual for detailed instructions.</p>"
  }' \
  http://localhost:8080/blogposts

curl -X POST -H "Content-Type: application/json" -d @data.json http://localhost:8080/blogposts


UpdateBlogPost:
curl -X PUT http://localhost:8080/blogposts/2 \
-H "Content-Type: application/json" \
-d '{"category": "Updated Category", "title": "Updated Title", "author": "Updated Author", "editor_data": "<p>This is an <strong>example</strong> of HTML content.</p>"}'


DeleteBlogPost:
curl -X DELETE http://localhost:8080/blogposts/{{post_id}}

FilterBlogPost:
curl -X GET "http://localhost:8080/blogposts/filter?category=<category_name>
curl -X GET "http://localhost:8080/blogposts/filter?category=fuel-system
///////////////////////Postgres Configuration///////////
# Init Postgres in bash
psql
# List databases
\l
# Create database
CREATE DATABASE blog_automotive;
# Switch to orders database
\c blog_automotive
# Check you path in UNIX bash
pwd
# Execute sql script
\i /Users/tavito/Documents/go/blogging-automotive/blog_automotive.sql
# Delete database in case you need
DROP DATABASE blog_automotive;


# List the all tables/table within a database
\dt
\d table_name

# Query
SELECT * FROM blog_elements;
SELECT * FROM posts;

//////////////New Features to be updated//////////////////
- Make a view blogpost counter to give insight on readers likes data.

*/
